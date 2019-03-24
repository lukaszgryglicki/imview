package main

import (
	"fmt"
	"image"
	"os"

	"github.com/lukaszgryglicki/imview"
)

func main() {
	images := loadImages(os.Args[1:])
	err := imview.Show(images...)
	if err != nil {
		fmt.Printf("Show: %+v\n", err)
	}
}

// loadImages loads all images given by their filenames
func loadImages(paths []string) []*image.RGBA {
	var result []*image.RGBA
	for _, path := range paths {
		im, err := imview.LoadImage(path)
		if err != nil {
			fmt.Printf("LoadImage: %s: %v\n", path, err)
			continue
		}
		result = append(result, imview.ImageToRGBA(im))
	}
	return result
}
