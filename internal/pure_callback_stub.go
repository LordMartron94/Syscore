//go:build !linux && !windows && !darwin

package internal

func platformNewCallback(_ any) uintptr {
	return 0
}
