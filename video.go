package main

import (
	"fmt"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"log"
	"math"
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
	for {
		select {
        case d := <-v.debug:
			copy((*[512 * 480]int)(v.screen.Pixels)[:], d)
			v.screen.Flip()
			// 60hz
			// time.Sleep(16000000 * time.Nanosecond)
			// time.Sleep(12000000 * time.Nanosecond)
		case val := <-v.tick:
			bigscreen := make([]int, 512*480)
			for i := len(val) - 1; i >= 0; i-- {
				y := int(math.Floor(float64(i / 256)))
				x := i - (y * 256)

				y *= 2
				x *= 2

				bigscreen[(y*512)+x] = val[i]
				bigscreen[((y+1)*512)+x] = val[i]
				bigscreen[(y*512)+(x+1)] = val[i]
				bigscreen[((y+1)*512)+(x+1)] = val[i]
			}

			copy((*[512 * 480]int)(v.screen.Pixels)[:], bigscreen)
			v.screen.Flip()
			// 60hz
			// time.Sleep(16000000 * time.Nanosecond)
			time.Sleep(8000000 * time.Nanosecond)
		}
	}
}

func (v *Video) Close() {
	sdl.Quit()
}
