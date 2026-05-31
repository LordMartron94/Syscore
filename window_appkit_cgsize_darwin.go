//go:build darwin

package syscore

import (
	"fmt"
	"math"
	"sync"

	"syscore/window/appkit/bindings"
)

type (
	// SYSCORE_Window_Appkit_CGSize is the CoreGraphics size layout returned by CAMetalLayer.drawableSize.
	SYSCORE_Window_Appkit_CGSize = bindings.CGSize
)

var (
	windowAppkitObjcMsgSendCGSizeOnce sync.Once
	windowAppkitObjcMsgSendCGSize     func(receiver, selector uintptr) bindings.CGSize
	windowAppkitObjcMsgSendCGSizeErr  error
)

/*
SYSCORE_Window_Appkit_ObjcMsgSendCGSizeGet binds objc_msgSend for selectors that return CGSize.

[Context]
Used for -drawableSize on CAMetalLayer. The returned function must only be used for CGSize-returning selectors.
*/
func SYSCORE_Window_Appkit_ObjcMsgSendCGSizeGet() (func(receiver, selector uintptr) SYSCORE_Window_Appkit_CGSize, error) {
	windowAppkitObjcMsgSendCGSizeOnce.Do(func() {
		lib, err := SYSCORE_Pure_LibraryLoad("/usr/lib/libobjc.A.dylib")
		if err != nil {
			windowAppkitObjcMsgSendCGSizeErr = err
			return
		}

		windowAppkitObjcMsgSendCGSizeErr = SYSCORE_Pure_LibraryFunctionBind(lib, &windowAppkitObjcMsgSendCGSize, "objc_msgSend")
	})

	return windowAppkitObjcMsgSendCGSize, windowAppkitObjcMsgSendCGSizeErr
}

/*
SYSCORE_Window_Appkit_MetalLayerDrawableSizeGet returns the drawable pixel extent of a CAMetalLayer.

[Context]
Queries -drawableSize on the metal layer object. Width and height are the values Vulkan WSI expects
for surface creation on macOS.
*/
func SYSCORE_Window_Appkit_MetalLayerDrawableSizeGet(
	commands SYSCORE_Window_Appkit_Commands,
	metalLayer SYSCORE_Window_Appkit_Object,
) (widthPx uint32, heightPx uint32, err error) {
	sendCGSize, err := SYSCORE_Window_Appkit_ObjcMsgSendCGSizeGet()
	if err != nil {
		return 0, 0, fmt.Errorf("bind objc_msgSend for CGSize: %w", err)
	}

	drawableSizeSel := SYSCORE_Window_Appkit_SelGet(commands, "drawableSize")
	drawableSize := sendCGSize(uintptr(metalLayer), uintptr(drawableSizeSel))
	if drawableSize.Width <= 0 || drawableSize.Height <= 0 || math.IsNaN(drawableSize.Width) || math.IsNaN(drawableSize.Height) {
		return 0, 0, fmt.Errorf("CAMetalLayer drawableSize is invalid")
	}

	widthPx = uint32(drawableSize.Width)
	heightPx = uint32(drawableSize.Height)
	if widthPx == 0 || heightPx == 0 {
		return 0, 0, fmt.Errorf("CAMetalLayer drawableSize rounds to zero")
	}

	return widthPx, heightPx, nil
}
