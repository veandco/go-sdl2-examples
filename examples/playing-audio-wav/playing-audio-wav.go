package main

// typedef unsigned char Uint8;
// void OnAudioPlayback(void *userdata, Uint8 *stream, int len);
import "C"
import (
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	audio  []byte
	offset int // We use this to keep track of which part of audio to play
)

//export OnAudioPlayback
func OnAudioPlayback(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]byte)(unsafe.Pointer(&hdr))
	for i := 0; i < n; i++ {
		buf[i] = audio[offset]
		offset = (offset + 1) % len(audio) // Increase audio offset and loop when it reaches the end
	}
}

func main() {
	var dev sdl.AudioDeviceID
	var err error

	// Initialize SDL2
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		log.Println(err)
		return
	}
	defer sdl.Quit()

	// Load WAV audio file
	tmpaudio, spec := sdl.LoadWAV("../../assets/test.wav")
	if spec == nil {
		log.Println(sdl.GetError())
		return
	}
	audio = tmpaudio

	spec.Callback = sdl.AudioCallback(C.OnAudioPlayback)

	// Open default playback device
	if dev, err = sdl.OpenAudioDevice("", false, spec, nil, 0); err != nil {
		log.Println(err)
		return
	}
	defer sdl.CloseAudioDevice(dev)

	// Start playback audio
	sdl.PauseAudioDevice(dev, false)

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
