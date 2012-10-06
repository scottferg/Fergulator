package main

import (
	"fmt"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"log"
	"time"
)

type Video struct {
	screen *sdl.Surface
	tick   <-chan []int
	debug  <-chan []int
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

	v.tick = t
	v.debug = d
}

func (v *Video) Render() {
	buf := (*[512 * 480]int)(v.screen.Pixels)[:]
	for {
		select {
		case val := <-v.tick:
			for i := len(val) - 1; i >= 0; i-- {
				y := i >> 8
				x := i - (y * 256)

				y *= 2
				x *= 2

				buf[(y*512)+x] = val[i]
				buf[((y+1)*512)+x] = val[i]
				buf[(y*512)+(x+1)] = val[i]
				buf[((y+1)*512)+(x+1)] = val[i]
			}

			v.screen.Flip()
			time.Sleep(4000000 * time.Nanosecond)
		}
	}
}

func (v *Video) Close() {
	sdl.Quit()
}
