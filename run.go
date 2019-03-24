package imview

import (
	"fmt"
	"image"
	"runtime"

	"github.com/go-gl/gl/v3.2-compatibility/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func init() {
	runtime.LockOSThread()
}

// Show - shows images
func Show(images ...*image.RGBA) error {
	if err := gl.Init(); err != nil {
		fmt.Printf("gl.Init: %+v\n", err)
		return err
	}

	if err := glfw.Init(); err != nil {
		fmt.Printf("glfw.Init: %+v\n", err)
		return err
	}
	defer glfw.Terminate()

	var windows []*Window
	var err error
	for _, im := range images {
		window, err := NewWindow(im)
		if err != nil {
			fmt.Printf("NewWindow: %+v\n", err)
			continue
		}
		windows = append(windows, window)
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
	}

	return err
}
