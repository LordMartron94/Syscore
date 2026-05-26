/*
Package main implements the syscore window binding code generator.

[Context]
Reads pinned protocol XML from git submodules under window/internal/spec/vendor.
Emits *_gen.go bindings and loader command tables; does not fetch from the network
or read host /usr/share paths during generation.
*/
package main
