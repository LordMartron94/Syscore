//go:build !darwin

package loader

import "fmt"

func AppkitCommandsLoadManifest(_ AppkitModule, _ *AppkitCommandManifest, _ *AppkitCommands) error {
	return fmt.Errorf("appkit loader: only supported on darwin")
}
