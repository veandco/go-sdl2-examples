// This program creates an SDL window and renderer and pushes UserEvents using PeepEvents
package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var winTitle string = "Go-SDL2 TestWaitEvent"
var winWidth, winHeight int32 = 800, 600

const pushTime uint32 = 1000 // number of milliseconds between event pushes

func run() int {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var event sdl.Event
	var running bool
	var err error

	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()

	window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 2
	}
	defer renderer.Destroy()

	var peepArray []sdl.Event = make([]sdl.Event, 2)
	id, err := window.GetID()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get window ID: %s\n", err)
		return 3
	}
	peepArray[0] = &sdl.UserEvent{sdl.USEREVENT, sdl.GetTicks(), id, 1331, nil, nil}
	peepArray[1] = &sdl.UserEvent{sdl.USEREVENT, sdl.GetTicks(), id, 10101, nil, nil}

	running = true
	lastPushTime := sdl.GetTicks()
	for running {
		if lastPushTime+pushTime < sdl.GetTicks() {
			lastPushTime = sdl.GetTicks()
			sdl.PumpEvents()
			numEventsHandled, err := sdl.PeepEvents(peepArray, sdl.ADDEVENT, sdl.FIRSTEVENT, sdl.LASTEVENT)
			if err != nil {
				fmt.Fprintf(os.Stderr, "PeepEvents error: %s\n", err)
				return 4
			} else {
				fmt.Printf("Successful push of %d events\n", numEventsHandled)
			}
		}

		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case sdl.QuitEvent:
				running = false
			case sdl.MouseMotionEvent:
				fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
			case sdl.MouseButtonEvent:
				fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
			case sdl.MouseWheelEvent:
				fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y)
			case sdl.KeyboardEvent:
				fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
					t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
			case sdl.UserEvent:
				fmt.Printf("[%d ms] UserEvent\tcode:%d\n", t.Timestamp, t.Code)
			}
		}
		sdl.Delay(1000 / 30)
	}

	return 0
}

func main() {
	os.Exit(run())
}
