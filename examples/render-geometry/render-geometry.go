package main

import (
	"log"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

// Whether to use RenderGeometryRaw or not
const raw = false

func main() {
	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Test RenderGeometryRaw", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Destroy()

	vertices := []sdl.Vertex{
		{sdl.FPoint{400, 150}, sdl.Color{255, 0, 0, 255}, sdl.FPoint{0, 0}},
		{sdl.FPoint{200, 450}, sdl.Color{0, 0, 255, 255}, sdl.FPoint{0, 0}},
		{sdl.FPoint{600, 450}, sdl.Color{0, 255, 0, 255}, sdl.FPoint{0, 0}},
	}

	if raw {
		log.Println("Using RenderGeometryRaw")
	} else {
		log.Println("Using RenderGeometry")
	}

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		renderer.SetDrawColor(0, 0, 0, sdl.ALPHA_OPAQUE)
		renderer.Clear()
		if raw {
			renderer.RenderGeometryRaw(nil, (*float32)(unsafe.Pointer(&vertices[0])), 20, (*sdl.Color)(unsafe.Add(unsafe.Pointer(&vertices[0]), 8)), 20, (*float32)(unsafe.Add(unsafe.Pointer(&vertices[0]), 12)), 20, len(vertices), nil, 0, 0)
		} else {
			renderer.RenderGeometry(nil, vertices, nil)
		}
		renderer.Present()

		sdl.Delay(33)
	}
}
