package main

import "github.com/veandco/go-sdl2/sdl"

const (
	TITLE     = "README"
	WIDTH     = 800
	HEIGHT    = 600
	FRAMERATE = 60
)

var (
	playerX, playerY   = int32(WIDTH / 2), int32(HEIGHT / 2)
	playerVX, playerVY = int32(0), int32(0)
	running            = true
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

	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			handleEvent(event)
		}

		loopTime := loop(surface)
		window.UpdateSurface()

		delay := (1000 / FRAMERATE) - loopTime
		sdl.Delay(delay)
	}
}

func handleEvent(event sdl.Event) {
	switch t := event.(type) {
	case sdl.QuitEvent: // NOTE: Please use `*sdl.QuitEvent` for `v0.4.x` (current version).
		println("Quitting..")
		running = false
		break
	case sdl.KeyboardEvent:
		if t.State == sdl.RELEASED {
			if t.Keysym.Sym == sdl.K_LEFT {
				playerVX -= 1
			} else if t.Keysym.Sym == sdl.K_RIGHT {
				playerVX += 1
			}
			if t.Keysym.Sym == sdl.K_UP {
				playerVY -= 1
			} else if t.Keysym.Sym == sdl.K_DOWN {
				playerVY += 1
			}
		}
		break
	}
}

func loop(surface *sdl.Surface) (loopTime uint32) {
	// Get time at the start of the function
	startTime := sdl.GetTicks()

	// Update player position
	playerX += playerVX
	playerY += playerVY
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

	// Clear surface
	surface.FillRect(nil, 0)

	// Draw on the surface
	rect := sdl.Rect{playerX, playerY, 4, 4}
	colour := sdl.Color{R: 255, G: 0, B: 255, A: 255} // purple
	pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
	surface.FillRect(&rect, pixel)

	// Calculate time passed since start of the function
	endTime := sdl.GetTicks()
	return endTime - startTime
}
