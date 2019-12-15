package main

/*
typedef unsigned char Uint8;
void OnAudio(void *userdata, Uint8 *stream, int length);
*/
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

const (
	DefaultFrequency = 16000
	DefaultFormat    = sdl.AUDIO_S16
	DefaultChannels  = 1
	DefaultSamples   = 512
)

var (
	audioC = make(chan []int16, 1)
)

//export OnAudio
func OnAudio(userdata unsafe.Pointer, _stream *C.Uint8, _length C.int) {
	// We need to cast the stream from C uint8 array into Go int16 slice
	length := int(_length) / 2                                                                      // Divide by 2 because a single int16 consists of two uint8
	header := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(_stream)), Len: length, Cap: length} // Build the slice header for our int16 slice
	buf := *(*[]int16)(unsafe.Pointer(&header))                                                     // Use the slice header as int16 slice

	// Copy the audio samples into temporary buffer
	audioSamples := make([]int16, length)
	copy(audioSamples, buf)

	// Send the temporary buffer to our main function via our Go channel
	audioC <- audioSamples
}

func main() {
	var dev sdl.AudioDeviceID
	var err error

	// Initialize SDL2
	if err = sdl.Init(sdl.INIT_AUDIO); err != nil {
		log.Fatal(err)
	}
	defer sdl.Quit()

	// Specify the configuration for our default recording device
	spec := sdl.AudioSpec{
		Freq:     DefaultFrequency,
		Format:   DefaultFormat,
		Channels: DefaultChannels,
		Samples:  DefaultSamples,
		Callback: sdl.AudioCallback(C.OnAudio),
	}

	// Open default recording device
	defaultRecordingDeviceName := sdl.GetAudioDeviceName(0, true)
	if dev, err = sdl.OpenAudioDevice(defaultRecordingDeviceName, true, &spec, nil, 0); err != nil {
		log.Fatal(err)
	}
	defer sdl.CloseAudioDevice(dev)

	// Start recording audio
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
		case audioSamples := <-audioC:
			log.Printf("Got %d audio samples!\n", len(audioSamples))
		}
	}
}
