package imview

import (
	"fmt"
	"image"
	"runtime"

	"github.com/go-gl/gl/v3.2-compatibility/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// ImagesData - holds all remaining images
type ImagesData struct {
	n      int
	images []image.Image
	rgbas  []*image.RGBA
	ready  []bool
	names  []string
}

var gwindow *Window
var gdata ImagesData

// InitProcessNextImages - init structures for parallel processing
func InitProcessNextImages(fns []string) error {
	n := len(fns)
	var (
		emptyImage image.Image
		emptyRGBA  *image.RGBA
	)
	for _, fn := range fns {
		gdata.names = append(gdata.names, fn)
		gdata.images = append(gdata.images, emptyImage)
		gdata.rgbas = append(gdata.rgbas, emptyRGBA)
		gdata.ready = append(gdata.ready, false)
	}
	gdata.n = n
	return nil
}

// ProcessNextImages - process N next unloaded images
func ProcessNextImages(ni int) {
	m := make(map[int]struct{})
	n := 0
	for i := range gdata.ready {
		if !gdata.ready[i] {
			m[i] = struct{}{}
			n++
		}
		if n == ni {
			break
		}
	}
	thrN := runtime.NumCPU()
	nT := 0
	ch := make(chan struct{})
	for idx := range m {
		go func(i int, c chan (struct{})) {
			gdata.Load(i)
			c <- struct{}{}
		}(idx, ch)
		nT++
		if nT == thrN {
			<-ch
			nT--
		}
	}
	for nT > 0 {
		<-ch
		nT--
	}
}

// Load image at given index
func (imd *ImagesData) Load(i int) {
	if imd.ready[i] {
		return
	}
	im, err := LoadImage(imd.names[i])
	if err != nil {
		fmt.Printf("Load: %s: %+v\n", imd.names[i], err)
		return
	}
	imd.images[i] = im
	imd.rgbas[i] = ImageToRGBA(im)
	imd.ready[i] = true
}

func glInit() error {
	runtime.LockOSThread()

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
	if action == glfw.Press || action == glfw.Repeat {
		if key == glfw.KeyEscape || key == glfw.KeyQ {
			w.SetShouldClose(true)
		} else if key == glfw.KeyF {
			gwindow.ToggleFullscreen()
		} else if key == glfw.KeyRight {
			gwindow.Move(1)
			ProcessNextImages(1)
		} else if key == glfw.KeyUp {
			gwindow.Move(10)
			ProcessNextImages(1)
		} else if key == glfw.KeyPageUp {
			gwindow.Move(100)
			ProcessNextImages(1)
		} else if key == glfw.KeyEnd {
			gwindow.Move(2000000000)
			ProcessNextImages(1)
		} else if key == glfw.KeyLeft {
			gwindow.Move(-1)
			ProcessNextImages(1)
		} else if key == glfw.KeyDown {
			gwindow.Move(-10)
			ProcessNextImages(1)
		} else if key == glfw.KeyPageDown {
			gwindow.Move(-100)
			ProcessNextImages(1)
		} else if key == glfw.KeyHome {
			gwindow.Move(-2000000000)
			ProcessNextImages(1)
		} else if key == glfw.Key1 {
			ProcessNextImages(1)
		} else if key == glfw.Key2 {
			ProcessNextImages(2)
		} else if key == glfw.Key3 {
			ProcessNextImages(5)
		} else if key == glfw.Key4 {
			ProcessNextImages(10)
		} else if key == glfw.Key5 {
			ProcessNextImages(20)
		} else if key == glfw.Key6 {
			ProcessNextImages(50)
		} else if key == glfw.Key7 {
			ProcessNextImages(100)
		} else if key == glfw.Key8 {
			ProcessNextImages(200)
		} else if key == glfw.Key9 {
			ProcessNextImages(1000)
		} else if key == glfw.Key0 {
			ProcessNextImages(2000000000)
		}
	}
}

// ShowSingle - shows single image
func ShowSingle(image image.Image, rgba *image.RGBA) error {
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
	gdata.images[0] = image
	gdata.rgbas[0] = rgba
	gdata.ready[0] = true
	window.Data = &gdata

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
