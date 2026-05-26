package main

import (
	"flag"
	"fmt"
	"path/filepath"

	gocode "codegen/go"
)

func main() {
	var specRoot string
	var waylandBindings string
	var xdgBindings string
	var xcbBindings string
	var win32Bindings string
	var appkitBindings string
	var waylandLoader string
	var xcbLoader string
	var win32Loader string
	var appkitLoader string
	var linuxOnly bool
	var win32Only bool
	var appkitOnly bool
	var win32FlushEvery int

	flag.StringVar(&specRoot, "specRoot", "", "directory containing spec_manifest.json and vendor submodules")
	flag.StringVar(&waylandBindings, "waylandBindings", "", "output dir for wayland/bindings")
	flag.StringVar(&xdgBindings, "xdgBindings", "", "output dir for wayland/xdg/bindings")
	flag.StringVar(&xcbBindings, "xcbBindings", "", "output dir for xcb/bindings")
	flag.StringVar(&win32Bindings, "win32Bindings", "", "output dir for win32/bindings")
	flag.StringVar(&appkitBindings, "appkitBindings", "", "output dir for appkit/bindings")
	flag.StringVar(&waylandLoader, "waylandLoader", "", "output dir for wayland/loader")
	flag.StringVar(&xcbLoader, "xcbLoader", "", "output dir for xcb/loader")
	flag.StringVar(&win32Loader, "win32Loader", "", "output dir for win32/loader")
	flag.StringVar(&appkitLoader, "appkitLoader", "", "output dir for appkit/loader")
	flag.BoolVar(&linuxOnly, "linux", false, "generate Linux bindings only")
	flag.BoolVar(&win32Only, "win32", false, "generate Win32 bindings only")
	flag.BoolVar(&appkitOnly, "appkit", false, "generate AppKit bindings only")
	flag.IntVar(&win32FlushEvery, "win32FlushEvery", 75, "flush Win32 binding output every N codegen elements")
	flag.Parse()

	if specRoot == "" {
		panic("syscore window generate: -specRoot is required")
	}

	manifest, err := windowSpecManifestLoad(specRoot)
	if err != nil {
		panic(err)
	}

	runAll := !linuxOnly && !win32Only && !appkitOnly
	if runAll || linuxOnly {
		if err := windowGenerateLinux(specRoot, manifest, waylandBindings, xdgBindings, xcbBindings, waylandLoader, xcbLoader); err != nil {
			panic(err)
		}
	}
	if runAll || win32Only {
		if err := windowGenerateWin32(specRoot, manifest, win32Bindings, win32Loader, win32FlushEvery); err != nil {
			panic(err)
		}
	}
	if runAll || appkitOnly {
		if err := windowGenerateAppkit(specRoot, manifest, appkitBindings, appkitLoader); err != nil {
			panic(err)
		}
	}
}

func windowGenerateLinux(specRoot string, manifest windowSpecManifest, waylandBindings, xdgBindings, xcbBindings, waylandLoader, xcbLoader string) error {
	waylandXML, err := windowSpecInputPathResolve(specRoot, manifest, "wayland")
	if err != nil {
		return err
	}
	xdgXML, err := windowSpecInputPathResolve(specRoot, manifest, "xdg-shell")
	if err != nil {
		return err
	}
	xprotoXML, err := windowSpecInputPathResolve(specRoot, manifest, "xproto")
	if err != nil {
		return err
	}
	xcbEventHeader, err := windowSpecInputPathResolve(specRoot, manifest, "xcb-event-header")
	if err != nil {
		return err
	}

	waylandSubmodule, ok := manifest.Submodules["wayland"]
	if !ok {
		return fmt.Errorf("spec manifest: wayland submodule missing")
	}
	waylandClientHeader := filepath.Join(specRoot, waylandSubmodule.Path, "src/wayland-client-core.h")

	if waylandBindings != "" {
		if err := waylandBindingsGenerate(
			waylandXML,
			waylandClientHeader,
			waylandBindings,
			[]string{"wl_display", "wl_registry", "wl_compositor", "wl_surface"},
			map[string][]string{
				"wl_registry": {"global", "global_remove"},
			},
			nil,
			true,
		); err != nil {
			return err
		}
		fmt.Printf("Wrote Wayland bindings to %s\n", waylandBindings)
	}

	if xdgBindings != "" {
		extra := []gocode.ConstSpec{
			bindingConstSpecString("XDG_WM_BASE_GLOBAL_NAME", "xdg_wm_base", "XDG_WM_BASE_GLOBAL_NAME is the registry global name for xdg_wm_base."),
		}
		if err := waylandBindingsGenerate(
			xdgXML,
			waylandClientHeader,
			xdgBindings,
			[]string{"xdg_wm_base", "xdg_surface", "xdg_toplevel"},
			map[string][]string{
				"xdg_surface":  {"configure"},
				"xdg_toplevel": {"configure", "close"},
			},
			extra,
			false,
		); err != nil {
			return err
		}
		fmt.Printf("Wrote xdg-shell bindings to %s\n", xdgBindings)
	}

	if xcbBindings != "" {
		if err := xcbBindingsGenerate(xprotoXML, xcbEventHeader, xcbBindings); err != nil {
			return err
		}
		fmt.Printf("Wrote XCB bindings to %s\n", xcbBindings)
	}

	if waylandLoader != "" {
		if err := loaderCommandsGenerate(loaderGenSpec{
			OutputDir:   waylandLoader,
			BuildTag:    "linux",
			PackageName: "loader",
			FuncName:    "waylandCommandMappingsAll",
			Imports:     []string{"syscore/window/loadutil"},
			Commands:    waylandLoaderCommands(),
		}); err != nil {
			return err
		}
		fmt.Printf("Wrote Wayland loader to %s\n", waylandLoader)
	}

	if xcbLoader != "" {
		if err := loaderCommandsGenerate(loaderGenSpec{
			OutputDir:   xcbLoader,
			BuildTag:    "linux",
			PackageName: "loader",
			FuncName:    "xcbCommandMappingsAll",
			Imports:     []string{"syscore/window/loadutil"},
			Commands:    xcbLoaderCommands(),
		}); err != nil {
			return err
		}
		fmt.Printf("Wrote XCB loader to %s\n", xcbLoader)
	}

	return nil
}

func windowGenerateWin32(specRoot string, manifest windowSpecManifest, win32Bindings, win32Loader string, win32FlushEvery int) error {
	if win32Bindings != "" {
		winmdPath, err := windowSpecInputPathResolve(specRoot, manifest, "win32-winmd")
		if err != nil {
			return err
		}
		if err := win32BindingsGenerate(winmdPath, win32Bindings, win32FlushEvery); err != nil {
			return err
		}
		fmt.Printf("Wrote Win32 bindings to %s\n", win32Bindings)
	}
	if win32Loader != "" {
		if err := win32LoaderCommandsGenerate(loaderGenSpec{
			OutputDir:   win32Loader,
			BuildTag:    "windows",
			PackageName: "loader",
			Imports:     []string{"syscore/window/loadutil"},
		}); err != nil {
			return err
		}
		fmt.Printf("Wrote Win32 loader to %s\n", win32Loader)
	}
	return nil
}

func windowGenerateAppkit(specRoot string, manifest windowSpecManifest, appkitBindings, appkitLoader string) error {
	if appkitBindings != "" {
		enumsPath, err := windowSpecInputPathResolve(specRoot, manifest, "macos-headers")
		if err != nil {
			return err
		}
		if err := appkitBindingsGenerate(enumsPath, appkitBindings); err != nil {
			return err
		}
		fmt.Printf("Wrote AppKit bindings to %s\n", appkitBindings)
	}
	if appkitLoader != "" {
		if err := loaderCommandsGenerate(loaderGenSpec{
			OutputDir:   appkitLoader,
			BuildTag:    "darwin",
			PackageName: "loader",
			FuncName:    "appkitCommandMappingsAll",
			Imports:     []string{"syscore/window/loadutil"},
			Commands:    appkitLoaderCommands(),
		}); err != nil {
			return err
		}
		fmt.Printf("Wrote AppKit loader to %s\n", appkitLoader)
	}
	return nil
}
