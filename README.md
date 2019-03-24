# Original author

    go get -u github.com/fogleman/imview/cmd/imview

# iview

Simple image viewer written in Go + OpenGL.

## Installation

    go get -u github.com/lukaszgryglicki/imview/cmd/iview

## Usage

    iview /path/to/images/folder/*

# Keyboard

- ESC/Q - quit
- f - toggle fullscreen
- Left/Right - +/- 1 image
- Up/Down +/- 10 images
- PgUp/PgDown +/- 100 images
- Home/End - first/last image
- 1,2,3,4,5 - preload next 1,5,10,30,100 images (but not less than N=number of your CPUs)
- 5,6,7,8,9,0 - unload last 1,5,10,30,100 images
- Cache will kepp no more than 300 images in memory
- Every time when new image is loaded (not from the cache) preload(1) is called (it will preload N next images, N = number of CPUs).
- Every time when N > 300 images are fetched into the cache - (N - 300) last images will be deleted from the cache.
- l - output cache stats into stdout
