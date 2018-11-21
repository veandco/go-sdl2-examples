package main

/*
// NOTE: Only works on SDL2 2.0.5 and above!

extern void cOnAudio(void *userdata, unsigned char *stream, int len);
*/
import "C"
import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	want, have sdl.AudioSpec
)

func makeBytes(raw *C.uchar, len int) (out []byte) {
	in := asBytes(raw, len)
	out = make([]byte, len)

	for i := 0; i < len; i++ {
		out[i] = in[i]
	}

	return
}

//export onAudio
func onAudio(raw *C.uchar, sz int) {
	data := makeBytes(raw, sz)
	fmt.Println("Received audio:", len(data), "bytes")
}

func asBytes(in *C.uchar, len int) (p []byte) {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&p))
	sliceHeader.Cap = len
	sliceHeader.Len = len
	sliceHeader.Data = uintptr(unsafe.Pointer(in))
	return
}

func open() (dev sdl.AudioDeviceID, err error) {
	var want, have sdl.AudioSpec

	want.Callback = sdl.AudioCallback(C.cOnAudio)
	want.Channels = 1
	want.Format = sdl.AUDIO_S16SYS
	want.Freq = 16000
	want.Samples = 512

	dev, err = sdl.OpenAudioDevice("", true, &want, &have, 0)
	if err != nil {
		return
	}

	sdl.PauseAudioDevice(dev, false)

	return
}

func main() {
	sdl.Init(sdl.INIT_AUDIO)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	func() {
		dev, err := open()
		if err != nil {
			panic(err)
		}
		defer sdl.CloseAudioDevice(dev)

		<-sigchan

		fmt.Println("Exiting..")
	}()

	sdl.Quit()
}
