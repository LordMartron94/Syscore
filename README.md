# SysCore

SysCore provides low-level **FFI** and **OS interaction** helpers: dynamic library loading, C string conversion, and platform-native **windowing** bindings. It does not implement rendering, input loops, or window lifecycle policy beyond what you build on top of the loaded APIs.

Consumers include `gpuarch` (Vulkan ICD load) and any engine that needs direct calls into `libxcb`, `libwayland-client`, Win32, or AppKit without cgo.

## Requirements

- Go 1.25+
- Platform libraries present at runtime (e.g. `libxcb.so.1` on Linux with X11)

## Module layout

```
syscore/
├── c_api.go                 # SYSCORE_C_* string helpers
├── pure_api.go              # SYSCORE_Pure_* dlopen / bind / dlsym
├── window_*_api.go          # SYSCORE_Window_* facades per backend
└── window/
    ├── platform/            # display API detection (WAYLAND_DISPLAY, DISPLAY, GOOS)
    ├── loadutil/            # shared loader bind helpers
    ├── xcb/                 # bindings + loader (linux)
    ├── wayland/             # bindings + loader (linux)
    ├── win32/               # bindings + loader (windows)
    └── appkit/              # bindings + loader (darwin)
```

| Import path | Role |
|-------------|------|
| `syscore` | Public `SYSCORE_*` facades |
| `syscore/window/<backend>/bindings` | 1:1 C types and `PFN_*` typedefs |
| `syscore/window/<backend>/loader` | `*ModuleLoad`, `*CommandsLoad`, manifest |
| `syscore/window/platform` | `DisplayAPIDetect`, `DisplayAPIListAvailable` |

Do not import `syscore/internal` from application code.

## C string utilities (`SYSCORE_C_*`)

Helpers for passing data across the Go/C boundary without cgo:

| Function | Purpose |
|----------|---------|
| `SYSCORE_C_StringToCString` | UTF-8 Go string → NUL-terminated `[]byte` + `*byte` |
| `SYSCORE_C_StringToCStringFirstByte` | Same, pointer only |
| `SYSCORE_C_StringToCStringSlice` | Same, slice only |
| `SYSCORE_C_CStringToString` | NUL-terminated `[]byte` → Go string |
| `SYSCORE_C_CStringPointerToString` | C `char*` address → Go string |
| `SYSCORE_C_StringToUTF16` | Go string → NUL-terminated UTF-16 (Win32) |

## Pure FFI (`SYSCORE_Pure_*`)

Wraps [purego](https://github.com/ebitengine/purego) for `dlopen` / `dlsym` style loading on Linux and macOS, and `LoadLibrary` / `GetProcAddress` on Windows.

| Function | Purpose |
|----------|---------|
| `SYSCORE_Pure_LibraryLoad` | Load a dynamic library by path or name |
| `SYSCORE_Pure_LibraryFunctionBind` | Bind an exported symbol name to a Go function pointer |
| `SYSCORE_Pure_LibrarySymbolResolve` | Resolve symbol address without binding |
| `SYSCORE_Pure_FunctionBindAddress` | Bind a known address to a Go function pointer |
| `SYSCORE_Pure_FunctionPointerAddressSet` | Write a raw native callback address into a function-pointer slot |

Use these when integrating libraries that are not covered by the window loaders (Vulkan ICD resolution in `gpuarch` follows this pattern).

## Windowing

Platform window APIs mirror **gpuarch** Vulkan: **bindings** (registry-faithful types and PFNs) and **loader** (module handle, command struct, optional manifest for selective binding). There is no high-level `CreateWindow` helper; you load commands and call through holder fields.

| Backend | GOOS | Typical library |
|---------|------|-----------------|
| X11/XCB | `linux` | `libxcb.so.1` |
| Wayland | `linux` | `libwayland-client.so.0`, embedded `libsyscore_wayland_xdg.so` (xdg-shell) |
| Win32 | `windows` | `user32.dll`, `kernel32.dll` |
| AppKit | `darwin` | `libobjc.A.dylib`, `AppKit.framework` |

### Detecting which display API to use

Linux session policy matches the original syscore behavior: prefer **Wayland** when `WAYLAND_DISPLAY` is set, otherwise **X11** when `DISPLAY` is set. Windows and macOS map to a single native API.

```go
api, err := syscore.SYSCORE_Window_DisplayAPIDetect()
if err != nil {
    // no suitable Linux display environment
}

switch api {
case syscore.SYSCORE_Window_DisplayAPIWayland:
    // load syscore/window/wayland/loader
case syscore.SYSCORE_Window_DisplayAPIX11:
    // load syscore/window/xcb/loader
case syscore.SYSCORE_Window_DisplayAPIWin32:
    // load syscore/window/win32/loader
case syscore.SYSCORE_Window_DisplayAPIAppkit:
    // load syscore/window/appkit/loader
}
```

| Function | Purpose |
|----------|---------|
| `SYSCORE_Window_DisplayAPIDetect` | Recommended API for this process |
| `SYSCORE_Window_DisplayAPIIsAvailable` | Whether a specific API is usable now |
| `SYSCORE_Window_DisplayAPIListAvailable` | All APIs currently available (e.g. both Wayland and X11 on Linux) |
| `SYSCORE_Window_DisplayAPIName` | Stable name (`"wayland"`, `"x11"`, …) |

Environment variables (Linux): `SYSCORE_Window_EnvWaylandDisplay` (`WAYLAND_DISPLAY`), `SYSCORE_Window_EnvX11Display` (`DISPLAY`).

Equivalent package API: `syscore/window/platform`.

### Loading and calling (XCB example)

```go
module, err := syscore.SYSCORE_Window_Xcb_ModuleLoad()
var commands syscore.SYSCORE_Window_Xcb_Commands
err = syscore.SYSCORE_Window_Xcb_CommandsLoad(module, &commands)

conn := commands.Connect(nil, nil)
setup := commands.GetSetup(conn)
// commands.CreateWindow, MapWindow, Flush, ...
```

Wayland globals (`wl_registry_interface`, etc.) are resolved with `SYSCORE_Window_Wayland_ModuleSymbolResolve`. xdg-shell `wl_interface` pointers (`xdg_wm_base_interface`, etc.) come from `SYSCORE_Window_Wayland_Xdg_ModuleLoad` / `SYSCORE_Window_Wayland_Xdg_ModuleSymbolResolve`; use `ProxyMarshal` for `xdg_toplevel.set_title`, `xdg_surface.ack_configure`, and related requests. XCB adds `InternAtom`, `InternAtomReply`, `ChangeProperty`, and `Free` for window titles (`WM_NAME`, `_NET_WM_NAME`). AppKit exposes `NSString` / `setTitle:` selector constants. Win32 uses separate user32/kernel32 module handles and symbol resolve helpers.

### Selective command loading (manifest)

```go
var manifest syscore.SYSCORE_Window_Xcb_CommandManifest
syscore.SYSCORE_Window_Xcb_CommandManifestReset(&manifest)
syscore.SYSCORE_Window_Xcb_CommandManifestAddConnect(&manifest, &commands.Connect)
err = syscore.SYSCORE_Window_Xcb_CommandsLoadManifest(module, &manifest, &commands)
```

Each backend exposes `SYSCORE_Window_<Backend>_CommandManifestAdd*` helpers matching its PFN set.

## Boundaries

**In scope**

- Dynamic library load and symbol bind (purego / Win32)
- C string and UTF-16 conversion for FFI
- 1:1 windowing bindings and loaders per platform
- Display API detection from environment and GOOS

**Out of scope**

- GLFW/SDL-style window managers
- Event loops, input, or GPU surface creation (use `gpuarch` Vulkan WSI on top of handles you obtain)
- Generated bindings for arbitrary C headers (window sets are hand-maintained subsets)
