package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/ttf"
    "log"
)

type Video struct {
    screen *sdl.Surface
    font *ttf.Font
    tick <-chan string
}

func (v *Video) Init(t <-chan string) {
	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		log.Fatal(sdl.GetError())
	}

	if ttf.Init() != 0 {
		log.Fatal(sdl.GetError())
	}

	v.screen = sdl.SetVideoMode(640, 480, 32, sdl.RESIZABLE)

	if v.screen == nil {
		log.Fatal(sdl.GetError())
	}

	sdl.WM_SetCaption("GoNES Emulator", "")

	v.font = ttf.OpenFont("./fonts/arial.ttf", 16)

	if v.font == nil {
		log.Fatal(sdl.GetError())
	}

	v.font.SetStyle(ttf.STYLE_BOLD)

    v.tick = t
}

func (v *Video) Render() {
	white := sdl.Color{255, 255, 255, 0}

    for {
        select {
        case val := <-v.tick:
            text := ttf.RenderText_Blended(v.font, val, white)

            v.screen.FillRect(nil, 0x000000)
            v.screen.Blit(&sdl.Rect{2, 460, 0, 0}, text, nil)
            v.screen.Flip()

        case ev := <-sdl.Events:
            switch e := ev.(type) {
            case sdl.KeyboardEvent:
                if e.Keysym.Sym == sdl.K_ESCAPE {
                    running = false
                }
            }
        }
    }
}

func (v *Video) Close() {
    v.font.Close()

    ttf.Quit()
    sdl.Quit()
}
