package main

import (
	"fmt"
	"os"
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

	ir := &xprotoIR{Requests: make(map[string]xprotoRequestIR)}
	root.Walk(func(node *specxml.Node) {
		if node == nil || node.Name != "request" {
			return
		}
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
	})
	return ir, nil
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
	if strings.HasSuffix(rest, "_reply") {
		rest = strings.TrimSuffix(rest, "_reply")
	}
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
