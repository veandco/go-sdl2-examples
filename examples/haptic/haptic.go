package main

import (
	"fmt"

	"github.com/malashin/go-sdl2/sdl"
)

func main() {
	// Initialize SDL.
	err := sdl.Init(sdl.INIT_HAPTIC)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// Find the number of available haptic devices.
	fmt.Println("Looking for haptic devices...")
	nh, err := sdl.NumHaptics()
	if err != nil {
		panic(err)
	}
	if nh == 0 {
		panic("No haptic devices found")
	}

	// Iterate over all found haptic devices.
	for i := 0; i < nh; i++ {
		// Print out the name of the haptic device.
		n, err := sdl.HapticName(i)
		if err != nil {
			panic(err)
		}
		fmt.Printf("\nHaptic device %v: %v\n", i, n)

		// Open haptic device.
		h, err := sdl.HapticOpen(i)
		if err != nil {
			panic(err)
		}

		// Print out effects and capabilities that are supported by the haptic device.
		printSupportedEffects(h)

		supported, err := h.Query()
		if err != nil {
			panic(err)
		}

		var ids []int
		idn := 0

		fmt.Println("  Creating new haptic effects on the specified device...")

		if supported&sdl.HAPTIC_CONSTANT != 0 {
			fmt.Printf("    effect %v: HAPTIC_CONSTANT\n", idn)
			e := &sdl.HapticConstant{
				Type:         sdl.HAPTIC_CONSTANT,
				Direction:    sdl.HapticDirection{Type: sdl.HAPTIC_CARTESIAN, Dir: [3]int32{-1, 0}},
				Length:       5000,
				Level:        0x6000,
				AttackLength: 1000,
				FadeLength:   1000,
			}
			id, err := h.NewEffect(e)
			if err != nil {
				fmt.Println("   ", err)
			} else {
				ids = append(ids, id)
				idn++
			}
		}

		if supported&sdl.HAPTIC_SINE != 0 {
			fmt.Printf("    effect %v: HAPTIC_SINE\n", idn)
			e := &sdl.HapticPeriodic{
				Type:         sdl.HAPTIC_SINE,
				Direction:    sdl.HapticDirection{Type: sdl.HAPTIC_CARTESIAN, Dir: [3]int32{0, -1}},
				Period:       1000,
				Magnitude:    -0x2000,
				Phase:        18000,
				Length:       5000,
				AttackLength: 1000,
				FadeLength:   1000,
			}
			id, err := h.NewEffect(e)
			if err != nil {
				fmt.Println("   ", err)
			} else {
				ids = append(ids, id)
				idn++
			}
		}

		if supported&sdl.HAPTIC_TRIANGLE != 0 {
			fmt.Printf("    effect %v: HAPTIC_TRIANGLE\n", idn)
			e := &sdl.HapticPeriodic{
				Type:         sdl.HAPTIC_TRIANGLE,
				Direction:    sdl.HapticDirection{Type: sdl.HAPTIC_CARTESIAN, Dir: [3]int32{1, 0}},
				Period:       1000,
				Magnitude:    0x4000,
				Length:       5000,
				AttackLength: 1000,
				FadeLength:   1000,
			}
			id, err := h.NewEffect(e)
			if err != nil {
				fmt.Println("   ", err)
			} else {
				ids = append(ids, id)
				idn++
			}
		}

		if supported&sdl.HAPTIC_SAWTOOTHUP != 0 {
			fmt.Printf("    effect %v: HAPTIC_SAWTOOTHUP\n", idn)
			e := &sdl.HapticPeriodic{
				Type:         sdl.HAPTIC_SAWTOOTHUP,
				Direction:    sdl.HapticDirection{Type: sdl.HAPTIC_CARTESIAN, Dir: [3]int32{0, 1}},
				Period:       500,
				Magnitude:    0x5000,
				Length:       5000,
				AttackLength: 1000,
				FadeLength:   1000,
			}
			id, err := h.NewEffect(e)
			if err != nil {
				fmt.Println("   ", err)
			} else {
				ids = append(ids, id)
				idn++
			}
		}

		if supported&sdl.HAPTIC_SAWTOOTHDOWN != 0 {
			fmt.Printf("    effect %v: HAPTIC_SAWTOOTHDOWN\n", idn)
			e := &sdl.HapticPeriodic{
				Type:         sdl.HAPTIC_SAWTOOTHDOWN,
				Direction:    sdl.HapticDirection{Type: sdl.HAPTIC_CARTESIAN, Dir: [3]int32{-1, 0}},
				Period:       1000,
				Magnitude:    0x4000,
				Length:       5000,
				AttackLength: 1000,
				FadeLength:   1000,
			}
			id, err := h.NewEffect(e)
			if err != nil {
				fmt.Println("   ", err)
			} else {
				ids = append(ids, id)
				idn++
			}
		}

		if supported&sdl.HAPTIC_SPRING != 0 {
			fmt.Printf("    effect %v: HAPTIC_SPRING\n", idn)
			e := &sdl.HapticCondition{
				Type:       sdl.HAPTIC_SPRING,
				Length:     5000,
				RightSat:   [3]uint16{0xFFFF},
				LeftSat:    [3]uint16{0xFFFF},
				RightCoeff: [3]int16{0x2000},
				LeftCoeff:  [3]int16{0x2000},
				Center:     [3]int16{0x1000},
			}
			id, err := h.NewEffect(e)
			if err != nil {
				fmt.Println("   ", err)
			} else {
				ids = append(ids, id)
				idn++
			}
		}

		if supported&sdl.HAPTIC_RAMP != 0 {
			fmt.Printf("    effect %v: HAPTIC_RAMP\n", idn)
			e := &sdl.HapticRamp{
				Type:         sdl.HAPTIC_RAMP,
				Direction:    sdl.HapticDirection{Type: sdl.HAPTIC_CARTESIAN, Dir: [3]int32{0, -1}},
				Start:        0x4000,
				End:          -0x4000,
				Length:       5000,
				AttackLength: 1000,
				FadeLength:   1000,
			}
			id, err := h.NewEffect(e)
			if err != nil {
				fmt.Println("    ", err)
			} else {
				ids = append(ids, id)
				idn++
			}
		}

		if supported&sdl.HAPTIC_LEFTRIGHT != 0 {
			fmt.Printf("    effect %v: HAPTIC_LEFTRIGHT\n", idn)
			e := &sdl.HapticLeftRight{
				Type:           sdl.HAPTIC_LEFTRIGHT,
				Length:         5000,
				LargeMagnitude: 0x3000,
				SmallMagnitude: 0xFFFF,
			}
			id, err := h.NewEffect(e)
			if err != nil {
				fmt.Println("    ", err)
			} else {
				ids = append(ids, id)
				idn++
			}
		}

		// Play all created effects.
		fmt.Println("  Now playing effects for 5 seconds each with 1 second delay between")
		for i, id := range ids {
			h.RunEffect(id, 1)
			fmt.Printf("    playing effect %v\n", i)
			sdl.Delay(6000)
			// Destroy the effect.
			h.DestroyEffect(id)
		}
		// Close the haptic device.
		h.Close()
	}
}

// printSupportedEffects prints out effects and capabilities that are supported by the haptic device.
func printSupportedEffects(h *sdl.Haptic) {
	ne, err := h.NumEffects()
	if err != nil {
		fmt.Println(err)
		return
	}
	nep, err := h.NumEffectsPlaying()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("  Can store %d effects\n  Can play %d effects at the same time\n", ne, nep)

	supported, err := h.Query()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("  Supported effects:")
	if supported&sdl.HAPTIC_CONSTANT != 0 {
		fmt.Println("    constant")
	}
	if supported&sdl.HAPTIC_SINE != 0 {
		fmt.Println("    sine")
	}
	if supported&sdl.HAPTIC_TRIANGLE != 0 {
		fmt.Println("    triangle")
	}
	if supported&sdl.HAPTIC_SAWTOOTHUP != 0 {
		fmt.Println("    sawtoothup")
	}
	if supported&sdl.HAPTIC_SAWTOOTHDOWN != 0 {
		fmt.Println("    sawtoothdown")
	}
	if supported&sdl.HAPTIC_RAMP != 0 {
		fmt.Println("    ramp")
	}
	if supported&sdl.HAPTIC_FRICTION != 0 {
		fmt.Println("    friction")
	}
	if supported&sdl.HAPTIC_SPRING != 0 {
		fmt.Println("    spring")
	}
	if supported&sdl.HAPTIC_DAMPER != 0 {
		fmt.Println("    damper")
	}
	if supported&sdl.HAPTIC_INERTIA != 0 {
		fmt.Println("    inertia")
	}
	if supported&sdl.HAPTIC_CUSTOM != 0 {
		fmt.Println("    custom")
	}
	if supported&sdl.HAPTIC_LEFTRIGHT != 0 {
		fmt.Println("    left/right")
	}

	fmt.Println("  Supported capabilities:")
	if supported&sdl.HAPTIC_GAIN != 0 {
		fmt.Println("    gain")
	}
	if supported&sdl.HAPTIC_AUTOCENTER != 0 {
		fmt.Println("    autocenter")
	}
	if supported&sdl.HAPTIC_STATUS != 0 {
		fmt.Println("    status")
	}
	if supported&sdl.HAPTIC_PAUSE != 0 {
		fmt.Println("    pause")
	}
}
