package imview

import (
	"fmt"
	"image"
	"runtime"

	"github.com/go-gl/gl/v3.2-compatibility/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var gwindow *Window

func init() {
	runtime.LockOSThread()
}

func glInit() error {
	if err := gl.Init(); err != nil {
		fmt.Printf("gl.Init: %+v\n", err)
		return err
	}

	if err := glfw.Init(); err != nil {
		fmt.Printf("glfw.Init: %+v\n", err)
		return err
	}
	return nil
}

func keyboardCallbackFunc(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		if key == glfw.KeyEscape || key == glfw.KeyQ {
			w.SetShouldClose(true)
		} else if key == glfw.KeyF {
			gwindow.ToggleFullscreen()
		}
	}
}

// ShowSingle - shows single image
func ShowSingle(image *image.RGBA) error {
	err := glInit()
	if err != nil {
		fmt.Printf("glInit: %+v\n", err)
		return err
	}
	defer glfw.Terminate()

	window, err := NewWindow(image)
	if err != nil {
		fmt.Printf("NewWindow: %+v\n", err)
		return err
	}

	// Keyboard
	gwindow = window
	window.SetKeyCallback(keyboardCallbackFunc)

	// Fullscreen
	window.ToggleFullscreen()

	for {
		if !window.ShouldClose() {
			window.Draw()
			glfw.WaitEvents()
		} else {
			window.Destroy()
			break
		}
	}

	return err
}

// Show - shows images
func Show(images ...*image.RGBA) error {
	err := glInit()
	if err != nil {
		fmt.Printf("glInit: %+v\n", err)
		return err
	}
	defer glfw.Terminate()

	var (
		windows    []*Window
		macOSHacks []bool
	)
	for _, im := range images {
		window, err := NewWindow(im)
		if err != nil {
			fmt.Printf("NewWindow: %+v\n", err)
			continue
		}
		windows = append(windows, window)
		macOSHacks = append(macOSHacks, false)
	}

	n := len(windows)
	for n > 0 {
		for i := 0; i < n; i++ {
			window := windows[i]
			if window.ShouldClose() {
				window.Destroy()
				windows[i] = windows[n-1]
				windows = windows[:n-1]
				n--
				i--
				continue
			}
			window.Draw()
		}
		glfw.PollEvents()
		for i := 0; i < n; i++ {
			window := windows[i]
			if !macOSHacks[i] {
				x, y := window.GetPos()
				window.SetPos(x+1, y)
				macOSHacks[i] = true
			}
		}
	}

	return err
}
