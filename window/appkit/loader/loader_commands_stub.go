//go:build !darwin

package loader

import "fmt"

func AppkitCommandsLoad(_ AppkitModule, _ *AppkitCommands) error {
	return fmt.Errorf("appkit loader: only supported on darwin")
}
