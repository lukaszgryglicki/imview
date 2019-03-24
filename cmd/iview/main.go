package main

import (
	"fmt"
	"image"
	"os"

	"github.com/lukaszgryglicki/imview"
)

func main() {
	image, err := loadImage(os.Args[1])
	if err != nil {
		fmt.Printf("loadImage: %+v\n", err)
		return
	}
	err = imview.ShowSingle(image)
	if err != nil {
		fmt.Printf("Show: %+v\n", err)
	}
}

// loadImages loads all images given by their filenames
func loadImage(path string) (*image.RGBA, error) {
	im, err := imview.LoadImage(path)
	if err != nil {
		fmt.Printf("LoadImage: %s: %v\n", path, err)
		return nil, err
	}
	return imview.ImageToRGBA(im), nil
}
