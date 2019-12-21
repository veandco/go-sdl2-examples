package main

import "C"
import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	// Initialize SDL2
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		log.Println(err)
		return
	}
	defer sdl.Quit()

	// Initialize SDL2 mixer
	if err := mix.Init(mix.INIT_MP3); err != nil {
		log.Println(err)
		return
	}
	defer mix.Quit()

	// Open default playback device
	if err := mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.DEFAULT_FORMAT, mix.DEFAULT_CHANNELS, mix.DEFAULT_CHUNKSIZE); err != nil {
		log.Println(err)
		return
	}
	defer mix.CloseAudio()

	// Load compressed audio file (like MP3) with long duration as *mix.Music
	music, err := mix.LoadMUS("../../assets/test.mp3")
	if err != nil {
		log.Println(err)
	}
	defer music.Free()

	// Load WAV file with short duration as *mix.Chunk
	chunk, err := mix.LoadWAV("../../assets/test.wav")
	if err != nil {
		log.Println(err)
	}
	defer chunk.Free()

	// Play the music once. Change 0 to -1 for infinite playback!.
	if err := music.Play(0); err != nil {
		log.Println(err)
	}

	// Play the chunk once. Change 0 to -1 for infinite playback!
	if _, err := chunk.Play(-1, 0); err != nil {
		log.Println(err)
	}

	// Listen to OS signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Run infinite loop until we receive SIGINT or SIGTERM!
	running := true
	for running {
		select {
		case sig := <-c:
			log.Printf("Received signal %v. Exiting.\n", sig)
			running = false
		}
	}
}
