package main

import (
	"fmt"
	"os"
	"strings"

	specxml "foundation/spec/xml"
)

type waylandInterfaceIR struct {
	Name     string
	Version  string
	Requests []waylandMessageIR
	Events   []waylandMessageIR
}

type waylandMessageIR struct {
	Name  string
	Index uint32
}

func waylandProtocolIRBuild(xmlPath string, interfaceNames []string) (map[string]waylandInterfaceIR, error) {
	content, err := os.ReadFile(xmlPath)
	if err != nil {
		return nil, err
	}
	root, err := specxml.TreeParse(content)
	if err != nil {
		return nil, err
	}

	wantAll := len(interfaceNames) == 0
	want := make(map[string]struct{}, len(interfaceNames))
	if !wantAll {
		for _, name := range interfaceNames {
			want[name] = struct{}{}
		}
	}

	result := make(map[string]waylandInterfaceIR, len(interfaceNames))
	root.Walk(func(node *specxml.Node) {
		if node == nil || node.Name != "interface" {
			return
		}
		name := node.Attr("name")
		if !wantAll {
			if _, ok := want[name]; !ok {
				return
			}
		}
		if name == "" {
			return
		}
		if _, exists := result[name]; exists {
			return
		}
		result[name] = waylandInterfaceIRFromNode(node)
	})

	if !wantAll {
		for _, name := range interfaceNames {
			if _, ok := result[name]; !ok {
				return nil, fmt.Errorf("interface %q not found in %s", name, xmlPath)
			}
		}
	}
	return result, nil
}

func waylandInterfaceIRFromNode(node *specxml.Node) waylandInterfaceIR {
	ir := waylandInterfaceIR{
		Name:    node.Attr("name"),
		Version: node.Attr("version"),
	}
	var requestIndex uint32
	var eventIndex uint32
	for _, child := range node.Children {
		switch child.Name {
		case "request":
			ir.Requests = append(ir.Requests, waylandMessageIR{
				Name:  child.Attr("name"),
				Index: requestIndex,
			})
			requestIndex++
		case "event":
			ir.Events = append(ir.Events, waylandMessageIR{
				Name:  child.Attr("name"),
				Index: eventIndex,
			})
			eventIndex++
		}
	}
	return ir
}

func waylandGoConstPrefix(interfaceName string) string {
	return strings.ToUpper(strings.ReplaceAll(interfaceName, ".", "_"))
}

func waylandGoSymbolConst(interfaceName string) string {
	return waylandGoConstPrefix(interfaceName) + "_INTERFACE_SYMBOL"
}

func waylandGoOpcodeConst(interfaceName string, messageType string, messageName string) string {
	prefix := waylandGoConstPrefix(interfaceName)
	msg := strings.ToUpper(strings.ReplaceAll(messageName, "-", "_"))
	switch messageType {
	case "request":
		return prefix + "_" + msg + "_OPCODE"
	case "event":
		return prefix + "_" + msg + "_EVENT_OPCODE"
	default:
		return prefix + "_" + msg
	}
}

func waylandGoOpaqueType(interfaceName string) string {
	parts := strings.Split(interfaceName, "_")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, "")
}

func waylandGoListenerType(interfaceName string) string {
	return waylandGoOpaqueType(interfaceName) + "Listener"
}

func waylandGoListenerFieldName(eventName string) string {
	parts := strings.Split(eventName, "_")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, "")
}
