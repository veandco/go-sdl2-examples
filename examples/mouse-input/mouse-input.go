package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

func run() (err error) {
	var window *sdl.Window

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return
	}
	defer sdl.Quit()

	window, err = sdl.CreateWindow("Input", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return
	}
	defer window.Destroy()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseMotionEvent:
				fmt.Println("Mouse", t.Which, "moved by", t.XRel, t.YRel, "at", t.X, t.Y)
			case *sdl.MouseButtonEvent:
				if t.State == sdl.PRESSED {
					fmt.Println("Mouse", t.Which, "button", t.Button, "pressed at", t.X, t.Y)
				} else {
					fmt.Println("Mouse", t.Which, "button", t.Button, "released at", t.X, t.Y)
				}
			case *sdl.MouseWheelEvent:
				if t.X != 0 {
					fmt.Println("Mouse", t.Which, "wheel scrolled horizontally by", t.X)
				} else {
					fmt.Println("Mouse", t.Which, "wheel scrolled vertically by", t.Y)
				}
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
