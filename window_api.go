package syscore

import (
	"fmt"
	"syscore/internal"
)

func SYSCORE_Window_WindowTest() {
	window, err := internal.WindowCreate()
	if err != nil {
		panic(err)
	}

	switch window.Type {
	case internal.WindowTypeWayland:
		fmt.Println("Detected windowing API: Wayland")
	case internal.WindowTypeX11:
		fmt.Println("Detected windowing API: X11")
	}
}
