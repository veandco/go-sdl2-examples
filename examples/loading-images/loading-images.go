package main

import (
	"os"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	bmpImagePath = "../../assets/test.bmp"
	pngImagePath = "../../assets/test.png"
)

func run() (err error) {
	var window *sdl.Window
	var surface *sdl.Surface
	var bmpImage *sdl.Surface
	var pngImage *sdl.Surface

	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		return
	}
	defer sdl.Quit()

	// Create a window for us to draw the images on
	if window, err = sdl.CreateWindow("Loading images", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN); err != nil {
		return
	}
	defer window.Destroy()

	if surface, err = window.GetSurface(); err != nil {
		return
	}

	// Load a BMP image
	if bmpImage, err = sdl.LoadBMP(bmpImagePath); err != nil {
		return err
	}
	defer bmpImage.Free()

	// Load a PNG image
	if pngImage, err = img.Load(pngImagePath); err != nil {
		return err
	}
	defer pngImage.Free()

	// Draw the BMP image on the first half of the window
	bmpImage.BlitScaled(nil, surface, &sdl.Rect{X: 0, Y: 0, W: 400, H: 400})

	// Draw the PNG image on the first half of the window
	pngImage.BlitScaled(nil, surface, &sdl.Rect{X: 400, Y: 0, W: 400, H: 400})

	// Update the window surface with what we have drawn
	window.UpdateSurface()

	// Run infinite loop until user closes the window
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case sdl.QuitEvent:
				running = false
			}
		}

		sdl.Delay(16)
	}

	return
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}
