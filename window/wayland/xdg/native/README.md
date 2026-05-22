# xdg-shell native sources

Regenerate `../loader/embed/libsyscore_wayland_xdg.so` after updating the protocol:

```bash
gcc -shared -fPIC -o ../loader/embed/libsyscore_wayland_xdg.so xdg-shell-protocol.c -lwayland-client
```

Protocol C sources are produced from `xdg-shell.xml` via `wayland-scanner` (see Qt or wayland-protocols package on the build host).
