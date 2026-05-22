package loadutil

import (
	"fmt"

	"syscore/internal"
)

type CommandMapping struct {
	Target any
	Name   string
}

func LibraryCommandsBind(library internal.DynamicLibrary, mappings []CommandMapping) error {
	for i := range mappings {
		mapping := mappings[i]
		if err := internal.LibraryFunctionBind(library, mapping.Target, mapping.Name); err != nil {
			return fmt.Errorf("bind %q: %w", mapping.Name, err)
		}
	}
	return nil
}

func LibrarySymbolResolve(library internal.DynamicLibrary, symbolName string) (uintptr, error) {
	symbol, err := internal.DynamicLibrarySymbolResolve(library, symbolName)
	if err != nil {
		return 0, fmt.Errorf("resolve %q: %w", symbolName, err)
	}
	if symbol == 0 {
		return 0, fmt.Errorf("resolve %q: symbol not found", symbolName)
	}
	return symbol, nil
}
