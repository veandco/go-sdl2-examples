package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var joysticks [16]*sdl.Joystick

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

	sdl.JoystickEventState(sdl.ENABLE)

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case sdl.QuitEvent:
				running = false
			case sdl.JoyAxisEvent:
				fmt.Printf("[%d ms] JoyAxis\ttype:%d\twhich:%c\taxis:%d\tvalue:%d\n",
					t.Timestamp, t.Type, t.Which, t.Axis, t.Value)
			case sdl.JoyBallEvent:
				fmt.Println("Joystick", t.Which, "trackball moved by", t.XRel, t.YRel)
			case sdl.JoyButtonEvent:
				if t.State == sdl.PRESSED {
					fmt.Println("Joystick", t.Which, "button", t.Button, "pressed")
				} else {
					fmt.Println("Joystick", t.Which, "button", t.Button, "released")
				}
			case sdl.JoyHatEvent:
				position := ""

				switch t.Value {
				case sdl.HAT_LEFTUP:
					position = "top-left"
				case sdl.HAT_UP:
					position = "top"
				case sdl.HAT_RIGHTUP:
					position = "top-right"
				case sdl.HAT_RIGHT:
					position = "right"
				case sdl.HAT_RIGHTDOWN:
					position = "bottom-right"
				case sdl.HAT_DOWN:
					position = "bottom"
				case sdl.HAT_LEFTDOWN:
					position = "bottom-left"
				case sdl.HAT_LEFT:
					position = "left"
				case sdl.HAT_CENTERED:
					position = "center"
				}

				fmt.Println("Joystick", t.Which, "hat", t.Hat, "moved to", position, "position")
			case sdl.JoyDeviceAddedEvent:
				// Open joystick for use
				joysticks[int(t.Which)] = sdl.JoystickOpen(int(t.Which))
				if joysticks[int(t.Which)] != nil {
					fmt.Println("Joystick", t.Which, "connected")
				}
			case sdl.JoyDeviceRemovedEvent:
				if joystick := joysticks[int(t.Which)]; joystick != nil {
					joystick.Close()
				}
				fmt.Println("Joystick", t.Which, "disconnected")
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
