package imview

import (
	"image"
	"image/draw"

	"github.com/go-gl/gl/v3.2-compatibility/gl"
)

// Texture - data struct
type Texture struct {
	Handle uint32
}

// NewTexture - returns new texture handle
func NewTexture() *Texture {
	var handle uint32
	gl.GenTextures(1, &handle)
	t := &Texture{handle}
	t.SetMinFilter(gl.LINEAR)
	t.SetMagFilter(gl.NEAREST)
	t.SetWrapS(gl.CLAMP_TO_EDGE)
	t.SetWrapT(gl.CLAMP_TO_EDGE)
	return t
}

// Bind - binds texture
func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.Handle)
}

// SetImage - sets texture image
func (t *Texture) SetImage(im image.Image) {
	rgba := image.NewRGBA(im.Bounds())
	draw.Draw(rgba, rgba.Rect, im, image.ZP, draw.Src)
	t.SetRGBA(rgba)
}

// SetRGBA - sets texture RGBA image
func (t *Texture) SetRGBA(im *image.RGBA) {
	t.Bind()
	size := im.Rect.Size()
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA, int32(size.X), int32(size.Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(im.Pix))
}

// SetMinFilter - sets texture mode
func (t *Texture) SetMinFilter(x int32) {
	t.Bind()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, x)
}

// SetMagFilter - sets texture mode
func (t *Texture) SetMagFilter(x int32) {
	t.Bind()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, x)
}

// SetWrapS - sets texture mode
func (t *Texture) SetWrapS(x int32) {
	t.Bind()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, x)
}

// SetWrapT - sets texture mode
func (t *Texture) SetWrapT(x int32) {
	t.Bind()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, x)
}
