package main

// typedef unsigned char Uint8;
// void SineWave(void *userdata, Uint8 *stream, int len);
import "C"
import (
	"log"
	"math"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	DefaultFrequency = 16000
	DefaultFormat    = sdl.AUDIO_S16
	DefaultChannels  = 2
	DefaultSamples   = 512

	toneHz = 440
	dPhase = 2 * math.Pi * toneHz / DefaultSamples
)

//export SineWave
func SineWave(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length) / 2
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.ushort)(unsafe.Pointer(&hdr))

	var phase float64
	for i := 0; i < n; i++ {
		phase += dPhase
		sample := C.ushort((math.Sin(phase) + 0.999999) * 32768)
		buf[i] = sample
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

	// Specify the configuration for our default playback device
	spec := sdl.AudioSpec{
		Freq:     DefaultFrequency,
		Format:   DefaultFormat,
		Channels: DefaultChannels,
		Samples:  DefaultSamples,
		Callback: sdl.AudioCallback(C.SineWave),
	}

	// Open default playback device
	if dev, err = sdl.OpenAudioDevice("", false, &spec, nil, 0); err != nil {
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
