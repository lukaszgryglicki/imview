package imview

import (
	"fmt"
	"image"
	"os"

	"github.com/go-gl/gl/v3.2-compatibility/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// Window - holds window data
type Window struct {
	*glfw.Window
	Image      image.Image
	Texture    *Texture
	Fullscreen bool
}

// NewWindow - creates a new window
func NewWindow(im image.Image) (*Window, error) {
	const maxSize = 1200
	w := im.Bounds().Size().X
	h := im.Bounds().Size().Y
	a := float64(w) / float64(h)
	if a >= 1 {
		if w > maxSize {
			w = maxSize
			h = int(maxSize / a)
		}
	} else {
		if h > maxSize {
			h = maxSize
			w = int(maxSize * a)
		}
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(w, h, os.Args[0], nil, nil)
	if err != nil {
		fmt.Printf("glfw.CreateWindow: %+v\n", err)
		return nil, err
	}

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	texture := NewTexture()
	texture.SetImage(im)
	result := &Window{window, im, texture, false}
	result.SetRefreshCallback(result.onRefresh)
	return result, nil
}

// ToggleFullscreen - toggles fullscreen mode
func (window *Window) ToggleFullscreen() {
	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()
	if !window.Fullscreen {
		window.SetMonitor(monitor, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
		window.Fullscreen = true
	} else {
		window.SetMonitor(nil, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
		window.Fullscreen = false
	}
}

// SetImage - sets window image
func (window *Window) SetImage(im image.Image) {
	window.Image = im
	window.Texture.SetImage(im)
	window.Draw()
}

func (window *Window) onRefresh(x *glfw.Window) {
	window.Draw()
}

// Draw - draws window
func (window *Window) Draw() {
	window.MakeContextCurrent()
	gl.Clear(gl.COLOR_BUFFER_BIT)
	window.DrawImage()
	window.SwapBuffers()
}

// DrawImage - draws window image
func (window *Window) DrawImage() {
	const padding = 0
	iw := window.Image.Bounds().Size().X
	ih := window.Image.Bounds().Size().Y
	w, h := window.GetFramebufferSize()
	s1 := float32(w) / float32(iw)
	s2 := float32(h) / float32(ih)
	f := float32(1 - padding)
	var x, y float32
	if s1 >= s2 {
		x = f * s2 / s1
		y = f
	} else {
		x = f
		y = f * s1 / s2
	}
	gl.Enable(gl.TEXTURE_2D)
	window.Texture.Bind()
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(-x, -y)
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(x, -y)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(x, y)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(-x, y)
	gl.End()
}
