//go:build windows

package internal

import "github.com/ebitengine/purego"

func platformNewCallback(fn any) uintptr {
	return purego.NewCallback(fn)
}
