//go:build !linux

package internal

import "fmt"

func WindowCreate() (Window, error) {
	return Window{}, fmt.Errorf("syscore: window creation is only supported on linux")
}
