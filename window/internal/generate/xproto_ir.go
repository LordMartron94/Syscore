package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	specxml "foundation/spec/xml"
)

type xprotoFieldIR struct {
	Kind     string // field, list, switch
	Type     string
	Name     string
	ListType string
}

type xprotoRequestIR struct {
	Name     string
	Fields   []xprotoFieldIR
	HasReply bool
}

type xprotoIR struct {
	Requests map[string]xprotoRequestIR
	Events   map[string]uint32
	Enums    map[string]map[string]int64
}

func xprotoIRBuild(xmlPath string) (*xprotoIR, error) {
	content, err := os.ReadFile(xmlPath)
	if err != nil {
		return nil, err
	}
	root, err := specxml.TreeParse(content)
	if err != nil {
		return nil, err
	}
	if root == nil {
		return nil, fmt.Errorf("xproto: empty document")
	}

	ir := &xprotoIR{
		Requests: make(map[string]xprotoRequestIR),
		Events:   make(map[string]uint32),
		Enums:    make(map[string]map[string]int64),
	}
	root.Walk(func(node *specxml.Node) {
		if node == nil {
			return
		}
		switch node.Name {
		case "request":
			name := node.Attr("name")
			if name == "" {
				return
			}
			req := xprotoRequestIR{
				Name:     name,
				Fields:   xprotoRequestFieldsParse(node),
				HasReply: xprotoRequestHasReply(node),
			}
			ir.Requests[name] = req
		case "event":
			name := node.Attr("name")
			numberString := node.Attr("number")
			if name == "" || numberString == "" {
				return
			}
			number, err := strconv.ParseUint(numberString, 10, 32)
			if err != nil {
				return
			}
			ir.Events[name] = uint32(number)
		case "enum":
			enumName := node.Attr("name")
			if enumName == "" {
				return
			}
			items := make(map[string]int64)
			for _, item := range node.NodesNamed("item") {
				if item == nil {
					continue
				}
				itemName := item.Attr("name")
				if itemName == "" {
					continue
				}
				valueNode := item.Child("value")
				if valueNode != nil && strings.TrimSpace(valueNode.Text) != "" {
					value, err := strconv.ParseInt(strings.TrimSpace(valueNode.Text), 10, 64)
					if err != nil {
						continue
					}
					items[itemName] = value
					continue
				}
				bitNode := item.Child("bit")
				if bitNode == nil || strings.TrimSpace(bitNode.Text) == "" {
					continue
				}
				bitIndex, err := strconv.ParseInt(strings.TrimSpace(bitNode.Text), 10, 64)
				if err != nil || bitIndex < 0 {
					continue
				}
				items[itemName] = int64(1) << bitIndex
			}
			if len(items) > 0 {
				ir.Enums[enumName] = items
			}
		}
	})
	return ir, nil
}

func xprotoEnumValueGet(ir *xprotoIR, enumName string, itemName string) (int64, bool) {
	if ir == nil {
		return 0, false
	}
	enum, ok := ir.Enums[enumName]
	if !ok {
		return 0, false
	}
	value, ok := enum[itemName]
	return value, ok
}

func xprotoEventNumberGet(ir *xprotoIR, eventName string) (uint32, bool) {
	if ir == nil {
		return 0, false
	}
	value, ok := ir.Events[eventName]
	return value, ok
}

func xprotoRequestHasReply(node *specxml.Node) bool {
	for _, child := range node.Children {
		if child.Name == "reply" {
			return true
		}
	}
	return false
}

func xprotoRequestFieldsParse(node *specxml.Node) []xprotoFieldIR {
	fields := make([]xprotoFieldIR, 0, 8)
	for _, child := range node.Children {
		switch child.Name {
		case "field":
			fields = append(fields, xprotoFieldIR{
				Kind: "field",
				Type: child.Attr("type"),
				Name: child.Attr("name"),
			})
		case "list":
			fields = append(fields, xprotoFieldIR{
				Kind:     "list",
				Name:     child.Attr("name"),
				ListType: child.Attr("type"),
			})
		case "switch":
			fields = append(fields, xprotoFieldIR{
				Kind: "switch",
				Name: child.Attr("name"),
			})
		}
	}
	return fields
}

func xcbRequestNameFromSymbol(symbol string) (string, bool) {
	if !strings.HasPrefix(symbol, "xcb_") {
		return "", false
	}
	rest := strings.TrimPrefix(symbol, "xcb_")
	rest = strings.TrimSuffix(rest, "_reply")
	if rest == "" {
		return "", false
	}
	parts := strings.Split(rest, "_")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, ""), true
}

func xcbFieldGoName(snake string) string {
	parts := strings.Split(snake, "_")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, "")
}
