package main

import "C"

import (
	"path/filepath"
	"runtime"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	winTitle  = "Go SDL2"
	winWidth  = 480
	winHeight = 800
)

// Game states.
const (
	stateRun = iota
	stateFlap
	stateDead
)

// States text.
var stateText = map[int]string{
	stateRun:  "RUN",
	stateFlap: "FLAP",
	stateDead: "DEAD",
}

// Text represents state text.
type Text struct {
	Width   int32
	Height  int32
	Texture *sdl.Texture
}

// Engine represents SDL engine.
type Engine struct {
	State     int
	Window    *sdl.Window
	Renderer  *sdl.Renderer
	Sprite    *sdl.Texture
	Font      *ttf.Font
	Music     *mix.Music
	Sound     *mix.Chunk
	StateText map[int]*Text
	running   bool
}

// NewEngine returns new engine.
func NewEngine() (e *Engine) {
	e = &Engine{}
	e.running = true
	return
}

// Init initializes SDL.
func (e *Engine) Init() (err error) {
	err = sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return
	}

	img.Init(img.INIT_PNG)

	err = mix.Init(mix.INIT_MP3)
	if err != nil {
		return
	}

	err = ttf.Init()
	if err != nil {
		return
	}

	err = mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.DEFAULT_FORMAT, mix.DEFAULT_CHANNELS, 3072)
	if err != nil {
		return
	}

	e.Window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		return
	}

	e.Renderer, err = sdl.CreateRenderer(e.Window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return
	}

	return
}

// Destroy destroys SDL and releases the memory.
func (e *Engine) Destroy() {
	e.Renderer.Destroy()
	e.Window.Destroy()
	mix.CloseAudio()

	img.Quit()
	mix.Quit()
	ttf.Quit()
	sdl.Quit()
}

// Running checks if loop is running.
func (e *Engine) Running() bool {
	return e.running
}

// Quit exits main loop.
func (e *Engine) Quit() {
	e.running = false
}

// Load loads resources.
func (e *Engine) Load() {
	assetDir := ""
	if runtime.GOOS != "android" {
		assetDir = filepath.Join("android", "src", "main", "assets")
	}

	var err error
	e.Sprite, err = img.LoadTexture(e.Renderer, filepath.Join(assetDir, "images", "sprite.png"))
	if err != nil {
		sdl.LogError(sdl.LOG_CATEGORY_APPLICATION, "LoadTexture: %s\n", err)
	}

	e.Font, err = ttf.OpenFont(filepath.Join(assetDir, "fonts", "universalfruitcake.ttf"), 24)
	if err != nil {
		sdl.LogError(sdl.LOG_CATEGORY_APPLICATION, "OpenFont: %s\n", err)
	}

	e.Music, err = mix.LoadMUS(filepath.Join(assetDir, "music", "frantic-gameplay.mp3"))
	if err != nil {
		sdl.LogError(sdl.LOG_CATEGORY_APPLICATION, "LoadMUS: %s\n", err)
	}

	e.Sound, err = mix.LoadWAV(filepath.Join(assetDir, "sounds", "click.wav"))
	if err != nil {
		sdl.LogError(sdl.LOG_CATEGORY_APPLICATION, "LoadWAV: %s\n", err)
	}

	e.StateText = map[int]*Text{}
	for k, v := range stateText {
		t, _ := e.renderText(v, sdl.Color{0, 0, 0, 0})
		_, _, tW, tH, _ := t.Query()
		e.StateText[k] = &Text{tW, tH, t}
	}
}

// Unload unloads resources.
func (e *Engine) Unload() {
	for _, v := range e.StateText {
		v.Texture.Destroy()
	}

	e.Sprite.Destroy()
	e.Font.Close()
	e.Music.Free()
	e.Sound.Free()
}

// renderText renders texture from ttf font.
func (e *Engine) renderText(text string, color sdl.Color) (texture *sdl.Texture, err error) {
	surface, err := e.Font.RenderUTF8Blended(text, color)
	if err != nil {
		return
	}

	defer surface.Free()

	texture, err = e.Renderer.CreateTextureFromSurface(surface)
	return
}

//export SDL_main
func SDL_main() {
	runtime.LockOSThread()
	e := NewEngine()

	err := e.Init()
	if err != nil {
		sdl.LogError(sdl.LOG_CATEGORY_APPLICATION, "Init: %s\n", err)
	}
	defer e.Destroy()

	e.Load()
	defer e.Unload()

	// Sprite size
	const n = 128

	// Sprite rects
	var rects []*sdl.Rect
	for x := 0; x < 6; x++ {
		rect := &sdl.Rect{int32(n * x), 0, n, n}
		rects = append(rects, rect)
	}

	e.Music.Play(-1)

	var frame int = 0
	var alpha uint8 = 255
	var showText bool = true

	var text *Text = e.StateText[stateRun]

	for e.Running() {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				e.Quit()

			case *sdl.MouseButtonEvent:
				e.Sound.Play(2, 0)
				if t.Type == sdl.MOUSEBUTTONDOWN && t.Button == sdl.BUTTON_LEFT {
					alpha = 255
					showText = true

					if e.State == stateRun {
						text = e.StateText[stateFlap]
						e.State = stateFlap
					} else if e.State == stateFlap {
						text = e.StateText[stateDead]
						e.State = stateDead
					} else if e.State == stateDead {
						text = e.StateText[stateRun]
						e.State = stateRun
					}
				}

			case *sdl.KeyboardEvent:
				if t.Keysym.Scancode == sdl.SCANCODE_ESCAPE || t.Keysym.Scancode == sdl.SCANCODE_AC_BACK {
					e.Quit()
				}
			}
		}

		e.Renderer.Clear()

		var clips []*sdl.Rect

		w, h := e.Window.GetSize()
		x, y := int32(w/2), int32(h/2)

		switch e.State {
		case stateRun:
			e.Renderer.SetDrawColor(168, 235, 254, 255)
			clips = rects[0:2]

		case stateFlap:
			e.Renderer.SetDrawColor(251, 231, 240, 255)
			clips = rects[2:4]

		case stateDead:
			e.Renderer.SetDrawColor(255, 250, 205, 255)
			clips = rects[4:6]
		}

		clip := clips[frame/2]

		e.Renderer.FillRect(nil)
		e.Renderer.Copy(e.Sprite, clip, &sdl.Rect{x - (n / 2), y - (n / 2), n, n})

		if showText {
			text.Texture.SetAlphaMod(alpha)
			e.Renderer.Copy(text.Texture, nil, &sdl.Rect{x - (text.Width / 2), y - n*1.5, text.Width, text.Height})
		}

		e.Renderer.Present()
		sdl.Delay(50)

		frame += 1
		if frame/2 >= 2 {
			frame = 0
		}

		alpha -= 10
		if alpha <= 10 {
			alpha = 255
			showText = false
		}
	}
}

func main() {
	SDL_main()
}
