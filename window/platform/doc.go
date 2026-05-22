/*
Package platform detects which native display / windowing API is appropriate for the host.

Use DisplayAPIDetect for the recommended API, DisplayAPIListAvailable to enumerate options on
Linux, and DisplayAPIIsAvailable to test a specific backend. The syscore package re-exports
these as SYSCORE_Window_DisplayAPI* facades.
*/
package platform
