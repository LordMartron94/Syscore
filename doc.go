/*
Package syscore provides low-level FFI and OS interaction for C-compatible libraries.

Use the SYSCORE_C_* helpers for string conversion, SYSCORE_Pure_* for generic dynamic
library loading, and SYSCORE_Window_* for platform windowing bindings and display API
detection. Implementation details live in syscore/internal and syscore/window; import
those only when extending syscore itself.
*/
package syscore
