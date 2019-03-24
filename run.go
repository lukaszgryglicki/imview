package imview

import (
	"fmt"
	"image"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func init() {
	runtime.LockOSThread()
}

// ImagesData - holds all remaining images
type ImagesData struct {
	n      int
	l      int
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
	m := []int{}
	n := 0
	e := 0
	for idx := range gdata.status {
		i := (idx + gwindow.C) % gdata.n
		if gdata.status[i] == 0 {
			m = append(m, i)
		}
	}
	thrN := runtime.NumCPU()
	nT := 0
	ch := make(chan error)
	for _, idx := range m {
		go func(i int, c chan error) {
			_, err := gdata.Load(i)
			c <- err
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

// UnprocessNextImages - free N processed images
func UnprocessNextImages(ni int) int {
	if ni <= 0 {
		return -1
	}
	var (
		im   image.Image
		rgba *image.RGBA
		n    int
	)
	for i := range gdata.status {
		if gdata.status[i] == 1 {
			gdata.images[i] = im
			gdata.rgbas[i] = rgba
			gdata.status[i] = 0
			n++
			if n == ni {
				break
			}
		}
	}
	gdata.l -= n
	if n > 0 {
		fmt.Printf("Unloaded %d images\n", n)
		return 1
	}
	return 0
}

// Load image at given index
func (imd *ImagesData) Load(i int) (bool, error) {
	if imd.status[i] != 0 {
		if imd.status[i] == -1 {
			return false, fmt.Errorf("image %d/%s is marked as failed", i, imd.names[i])
		}
		return false, nil
	}
	im, err := LoadImage(imd.names[i])
	if err != nil {
		imd.status[i] = -1
		fmt.Printf("Load: %s: %+v\n", imd.names[i], err)
		return false, err
	}
	imd.images[i] = im
	imd.rgbas[i] = ImageToRGBA(im)
	imd.status[i] = 1
	gdata.l++
	fmt.Printf("Loaded %d image\n", i)
	return true, nil
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

func loadStats() {
	str := fmt.Sprintf("(%d/%d)|", gdata.l, gdata.n)
	for _, status := range gdata.status {
		if status == -1 {
			str += "-"
		} else if status == 1 {
			str += "#"
		} else {
			str += " "
		}
	}
	fmt.Printf("%s|\n", str)
}

func keyboardCallbackFunc(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press || action == glfw.Repeat {
		load := false
		if key == glfw.KeyEscape || key == glfw.KeyQ {
			w.SetShouldClose(true)
		} else if key == glfw.KeyL {
			loadStats()
		} else if key == glfw.KeyF {
			gwindow.ToggleFullscreen()
		} else if key == glfw.KeyRight {
			load = gwindow.Move(1)
		} else if key == glfw.KeyUp {
			load = gwindow.Move(10)
		} else if key == glfw.KeyPageUp {
			load = gwindow.Move(100)
		} else if key == glfw.KeyEnd {
			gwindow.Move(2000000000)
		} else if key == glfw.KeyLeft {
			gwindow.Move(-1)
		} else if key == glfw.KeyDown {
			gwindow.Move(-10)
		} else if key == glfw.KeyPageDown {
			gwindow.Move(-100)
		} else if key == glfw.KeyHome {
			gwindow.Move(-2000000000)
		} else if key == glfw.Key1 {
			ProcessNextImages(1)
		} else if key == glfw.Key2 {
			ProcessNextImages(5)
		} else if key == glfw.Key3 {
			ProcessNextImages(10)
		} else if key == glfw.Key4 {
			ProcessNextImages(30)
		} else if key == glfw.Key5 {
			ProcessNextImages(100)
		} else if key == glfw.Key6 {
			UnprocessNextImages(1)
		} else if key == glfw.Key7 {
			UnprocessNextImages(5)
		} else if key == glfw.Key8 {
			UnprocessNextImages(10)
		} else if key == glfw.Key9 {
			UnprocessNextImages(30)
		} else if key == glfw.Key0 {
			UnprocessNextImages(100)
		}
		if gdata.l > 300 {
			UnprocessNextImages(gdata.l - 300)
		} else if load {
			fmt.Printf("Auto preloading\n")
			ProcessNextImages(1)
		}
		fmt.Printf("Current: %d/%d: %s (%d cached)\n", gwindow.C, gdata.n, gdata.names[gwindow.C], gdata.l)
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
	gdata.status[0] = 1
	window.Data = &gdata
	title := fmt.Sprintf("%d: %s (cached %d/%d)", window.C, window.Data.names[window.C], window.Data.l, window.Data.n)
	window.SetTitle(title)

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
