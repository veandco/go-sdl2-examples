package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

func fatal(msgs ...interface{}) {
	fmt.Fprint(os.Stderr, msgs...)
	os.Exit(1)
}

func main() {
	reader, err := os.Open("../../assets/test.png")
	if err != nil {
		fatal(err)
	}
	defer reader.Close()

	img, err := png.Decode(reader)
	if err != nil {
		fatal(err)
	}

	imgNRGBA, ok := img.(*image.NRGBA)
	if !ok {
		fatal("failed to cast img to image.NRGBA")
	}

	bounds := imgNRGBA.Bounds()

	var width, height int32

	width = int32(bounds.Max.X - bounds.Min.X)
	height = int32(bounds.Max.Y - bounds.Min.Y)

	surf, err := sdl.CreateRGBSurfaceWithFormat(0, width, height, 32, uint32(sdl.PIXELFORMAT_ABGR8888))
	if err != nil {
		fatal(err)
	}
	defer surf.Free()

	surf.Lock()

	if len(surf.Pixels()) != len(imgNRGBA.Pix) {
		fatal("surface and image had different pixel size")
	}

	copy(surf.Pixels(), imgNRGBA.Pix)

	surf.Unlock()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fatal(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateShapedWindow(
		"gopher",
		uint32(sdl.WINDOWPOS_CENTERED), uint32(sdl.WINDOWPOS_CENTERED),
		uint32(width), uint32(height),
		sdl.WINDOW_SHOWN)
	if err != nil {
		fatal(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		fatal(err)
	}
	defer renderer.Destroy()

	shmp := sdl.ShapeModeBinarizeAlpha{1 << 6}

	window.SetShape(surf, shmp)

	texture, err := renderer.CreateTextureFromSurface(surf)
	if err != nil {
		fatal(err)
	}
	defer texture.Destroy()

	renderer.Copy(texture, nil, nil)

	renderer.Present()

eventLoop:
	for {
		event := sdl.WaitEvent()
		if event == nil {
			break eventLoop
		}

		switch e := event.(type) {
		case *sdl.QuitEvent:
			fmt.Println("Close")
			break eventLoop

		case *sdl.MouseButtonEvent:
			if e.Button == sdl.BUTTON_LEFT && e.State == sdl.PRESSED {
				fmt.Println("LMB Click")
				break eventLoop
			}

			if e.Button == sdl.BUTTON_RIGHT && e.State == sdl.PRESSED {
				fmt.Println("RMB Click")
				fmt.Println("window.IsShaped():", window.IsShaped())
				wsm, ierr := window.GetShapeMode()
				fmt.Println("window.GetShapeMode():", wsm, ierr)
			}
		}
	}

	fmt.Println("Quit")
}
