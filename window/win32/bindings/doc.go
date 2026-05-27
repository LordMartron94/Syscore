/*
Package bindings holds generated Win32 metadata bindings split by namespace and DLL.

Types live under bindings/types/<namespace_slug> (one Go package per WinMD namespace).
Function pointer typedefs live under bindings/pfn/<module_slug> (one Go package per DLL).

The full WinMD spec is emitted without dropping variants; duplicate names within a package
receive numeric suffixes (_1, _2, ...).
*/
package bindings
