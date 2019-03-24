package main

import (
	"fmt"
	"image"
	"os"

	"github.com/lukaszgryglicki/imview"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%s: required at least one argument\n", os.Args[0])
		return
	}
	err := imview.InitProcessNextImages(os.Args[1:])
	if err != nil {
		fmt.Printf("InitProcessNextImages: %+v, args: %+v\n", err, os.Args[1:])
		return
	}
	image, imageRGBA, err := loadImage(os.Args[1])
	if err != nil {
		fmt.Printf("loadImage: %+v\n", err)
		return
	}
	err = imview.ShowSingle(image, imageRGBA)
	if err != nil {
		fmt.Printf("Show: %+v\n", err)
	}
}

func loadImage(path string) (image.Image, *image.RGBA, error) {
	im, err := imview.LoadImage(path)
	if err != nil {
		fmt.Printf("LoadImage: %s: %v\n", path, err)
		return nil, nil, err
	}
	return im, imview.ImageToRGBA(im), nil
}
