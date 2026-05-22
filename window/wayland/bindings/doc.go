/*
Package bindings provides 1:1 Wayland client types and function pointer typedefs for libwayland-client.

Interface globals (wl_registry_interface, etc.) are resolved by symbol name from the loader module.
Clients invoke protocol through bound WaylandCommands fields.
*/
package bindings
