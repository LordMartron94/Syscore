//go:build darwin

package syscore

import "unsafe"

/*
SYSCORE_Window_Appkit_ClassGet resolves an Objective-C class by name.

[Side Effects]
Allocates a temporary C string for the duration of the call.
*/
func SYSCORE_Window_Appkit_ClassGet(commands SYSCORE_Window_Appkit_Commands, className string) SYSCORE_Window_Appkit_Class {
	_, classNamePtr := SYSCORE_C_StringToCString(className)
	return commands.ObjcGetClass(uintptr(unsafe.Pointer(classNamePtr)))
}

/*
SYSCORE_Window_Appkit_SelGet registers an Objective-C selector by name.

[Side Effects]
Allocates a temporary C string for the duration of the call.
*/
func SYSCORE_Window_Appkit_SelGet(commands SYSCORE_Window_Appkit_Commands, selectorName string) SYSCORE_Window_Appkit_Sel {
	_, selectorPtr := SYSCORE_C_StringToCString(selectorName)
	return commands.SelRegisterName(uintptr(unsafe.Pointer(selectorPtr)))
}

/*
SYSCORE_Window_Appkit_Send dispatches objc_msgSend with uintptr arguments.

[Context]
Struct arguments (for example NSRect) must be passed as pointers in args; keep backing
storage alive for the duration of the call.
*/
func SYSCORE_Window_Appkit_Send(
	commands SYSCORE_Window_Appkit_Commands,
	receiver SYSCORE_Window_Appkit_Object,
	selector SYSCORE_Window_Appkit_Sel,
	args ...uintptr,
) SYSCORE_Window_Appkit_Object {
	return commands.ObjcMsgSend(receiver, selector, args)
}
