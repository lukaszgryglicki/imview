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
	status []int
	names  []string
}

var gwindow *Window
var gdata ImagesData

// InitProcessNextImages - init structures for parallel processing
func InitProcessNextImages(fns []string) {
	n := len(fns)
	var (
		emptyImage image.Image
		emptyRGBA  *image.RGBA
	)
	for _, fn := range fns {
		gdata.names = append(gdata.names, fn)
		gdata.images = append(gdata.images, emptyImage)
		gdata.rgbas = append(gdata.rgbas, emptyRGBA)
		gdata.status = append(gdata.status, 0)
	}
	gdata.n = n
	fmt.Printf("%d images\n", n)
}

// ProcessNextImages - process N next unloaded images
func ProcessNextImages(ni int) int {
	if ni <= 0 {
		return -1
	}
	runtime.UnlockOSThread()
	defer runtime.LockOSThread()
	m := []int{}
	n := 0
	e := 0
	for i := range gdata.status {
		if gdata.status[i] == 0 {
			m = append(m, i)
		}
	}
	thrN := runtime.NumCPU()
	nT := 0
	ch := make(chan error)
	for _, idx := range m {
		go func(i int, c chan error) {
			c <- gdata.Load(i)
		}(idx, ch)
		nT++
		if nT == thrN {
			r := <-ch
			nT--
			if r == nil {
				n++
				if n == ni {
					break
				}
			} else {
				e++
			}
		}
	}
	for nT > 0 {
		r := <-ch
		nT--
		if r == nil {
			n++
		} else {
			e++
		}
	}
	if n > 0 || e > 0 {
		fmt.Printf("Preloaded %d images (%d errors)\n", n, e)
	}
	if n > 0 {
		return 1
	}
	if e > 0 {
		return -1
	}
	return 0
}

// Load image at given index
func (imd *ImagesData) Load(i int) error {
	if imd.status[i] != 0 {
		if imd.status[i] == -1 {
			return fmt.Errorf("image %d/%s is marked as failed", i, imd.names[i])
		}
		return nil
	}
	im, err := LoadImage(imd.names[i])
	if err != nil {
		imd.status[i] = -1
		fmt.Printf("Load: %s: %+v\n", imd.names[i], err)
		return err
	}
	imd.images[i] = im
	imd.rgbas[i] = ImageToRGBA(im)
	imd.status[i] = 1
	fmt.Printf("Loaded %d image\n", i)
	return nil
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
		fmt.Printf("Current: %d/%d: %s\n", gwindow.C, gwindow.Data.n, gwindow.Data.names[gwindow.C])
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

	runtime.LockOSThread()
	window, err := NewWindow(image)
	if err != nil {
		fmt.Printf("NewWindow: %+v\n", err)
		return err
	}
	gdata.images[0] = image
	gdata.rgbas[0] = rgba
	gdata.status[0] = 1
	window.Data = &gdata
	window.SetTitle(window.Data.names[0])

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
