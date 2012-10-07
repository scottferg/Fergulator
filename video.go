package main

import (
	"fmt"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/gfx"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"log"
)

type Video struct {
	screen     *sdl.Surface
	fpsmanager *gfx.FPSmanager
	tick       <-chan []int
	debug      <-chan []int
}

func (v *Video) Init(t <-chan []int, d <-chan []int, n string) {
	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		log.Fatal(sdl.GetError())
	}

	v.screen = sdl.SetVideoMode(512, 480, 32, 0)

	if v.screen == nil {
		log.Fatal(sdl.GetError())
	}

	sdl.WM_SetCaption(fmt.Sprintf("Fergulator - %s", n), "")

	v.fpsmanager = gfx.NewFramerate()
	v.fpsmanager.SetFramerate(70)

	v.tick = t
	v.debug = d
}

func (v *Video) Render() {
	buf := (*[512 * 480]int32)(v.screen.Pixels)[:]
	for {
		select {
		case val := <-v.tick:
			for i := len(val) - 1; i >= 0; i-- {
				y := i >> 8
				x := i - (y * 256)

				y *= 2
				x *= 2

				buf[(y*512)+x] = int32(val[i])
				buf[((y+1)*512)+x] = int32(val[i])
				buf[(y*512)+(x+1)] = int32(val[i])
				buf[((y+1)*512)+(x+1)] = int32(val[i])
			}

			v.screen.Flip()
			v.fpsmanager.FramerateDelay()
		}
	}
}

func (v *Video) Close() {
	sdl.Quit()
}
