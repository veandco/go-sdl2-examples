package main

// typedef unsigned char Uint8;
// void AudioCallback(void *userdata, Uint8 *stream, int len);
import "C"
import (
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	sampleHz = 48000
)

var (
	audioData []byte
)

//export AudioCallback
func AudioCallback(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint8)(unsafe.Pointer(&hdr))

	if len(audioData) == 0 {
		for i := 0; i < n; i++ {
			buf[i] = 0
		}
	} else {
		for i := 0; i < n; i++ {
			if i >= len(audioData) {
				break
			}

			buf[i] = C.Uint8(audioData[i])
		}

		if n <= len(audioData) {
			// Remove processed audio data
			audioData = audioData[n:]
		} else {
			// Clear audio data
			audioData = audioData[len(audioData):]
		}
	}
}

// Handle audio data that is sent through HTTP
// Must be stereo unsigned 8-bit audio data at 48000Hz to play correctly
func audioHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	audioData = append(audioData, data...)

	w.WriteHeader(http.StatusOK)
}

func main() {
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		log.Println(err)
		return
	}
	defer sdl.Quit()

	spec := &sdl.AudioSpec{
		Freq:     sampleHz,
		Format:   sdl.AUDIO_U8,
		Channels: 2,
		Samples:  sampleHz,
		Callback: sdl.AudioCallback(C.AudioCallback),
	}
	if err := sdl.OpenAudio(spec, nil); err != nil {
		log.Println(err)
		return
	}
	defer sdl.CloseAudio()

	sdl.PauseAudio(false)

	http.HandleFunc("/audio", audioHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

