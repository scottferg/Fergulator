package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"log"
	"time"
)

type Video struct {
	screen *sdl.Surface
	tick   <-chan []int
}

func (v *Video) Init(t <-chan []int) {
	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		log.Fatal(sdl.GetError())
	}

	v.screen = sdl.SetVideoMode(256, 240, 32, sdl.RESIZABLE)

	if v.screen == nil {
		log.Fatal(sdl.GetError())
	}

	sdl.WM_SetCaption("GoNES Emulator", "")

	v.tick = t
}

func (v *Video) Render() {
	for {
		select {
		case val := <-v.tick:
            copy((*[256 * 240]int)(v.screen.Pixels)[:], val)
            v.screen.Flip()
            // 60hz
            time.Sleep(16000000 * time.Nanosecond)
		}
	}
}

func (v *Video) Close() {
	sdl.Quit()
}
