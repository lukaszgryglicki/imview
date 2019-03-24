package imview

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"  // to support GIFs
	_ "image/jpeg" // to support JPGs
	_ "image/png"  // to support PNGs
	"os"

	_ "golang.org/x/image/bmp" // to support BMPs
)

// LoadImage - loads image from file
func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("LoadImage: %s: %+v\n", path, err)
		return nil, err
	}
	defer func() { _ = file.Close() }()
	im, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("image.Decode: %+v\n", err)
	}
	return im, err
}

// ImageToRGBA - converts image to RGBA
func ImageToRGBA(src image.Image) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Rect, src, image.ZP, draw.Src)
	return dst
}
