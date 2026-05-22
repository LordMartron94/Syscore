/*
Package loader loads libxcb.so.1 and binds xcb_* entry points onto XcbCommands.

Clients call through XcbCommands fields after XcbCommandsLoad or XcbCommandsLoadManifest.
*/
package loader
