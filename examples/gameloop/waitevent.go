package main

import "github.com/veandco/go-sdl2/sdl"

const (
	TITLE  = "README"
	WIDTH  = 800
	HEIGHT = 600
)

var (
	playerX, playerY = int32(WIDTH / 2), int32(HEIGHT / 2)
	running          = true
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(TITLE, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WIDTH, HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	// Draw initial frame
	draw(window, surface)

	// Event loop
	for event := sdl.WaitEvent(); event != nil; event = sdl.WaitEvent() {
		switch t := event.(type) {
		case sdl.QuitEvent: // NOTE: Please use `*sdl.QuitEvent` for `v0.4.x` (current version).
			println("Quitting..")
			return
		case sdl.KeyboardEvent:
			if t.State == sdl.RELEASED {
				if t.Keysym.Sym == sdl.K_LEFT {
					playerX -= 4
				} else if t.Keysym.Sym == sdl.K_RIGHT {
					playerX += 4
				}
				if t.Keysym.Sym == sdl.K_UP {
					playerY -= 4
				} else if t.Keysym.Sym == sdl.K_DOWN {
					playerY += 4
				}
				if playerX < 0 {
					playerX = WIDTH
				} else if playerX > WIDTH {
					playerX = 0
				}
				if playerY < 0 {
					playerY = HEIGHT
				} else if playerY > HEIGHT {
					playerY = 0
				}
				draw(window, surface)
			}
			break
		}
	}
}

func draw(window *sdl.Window, surface *sdl.Surface) {
	// Clear surface
	surface.FillRect(nil, 0)

	// Draw on the surface
	rect := sdl.Rect{playerX, playerY, 4, 4}
	colour := sdl.Color{R: 255, G: 0, B: 255, A: 255} // purple
	pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
	surface.FillRect(&rect, pixel)

	window.UpdateSurface()
}
