package main

import (
	"fmt"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/ttf"
	"log"
)

type Video struct {
	screen *sdl.Surface
	font   *ttf.Font
	tick   <-chan []*Tile
}

func (cpu *Cpu) DumpRegisterState() string {
	return fmt.Sprintf("A: 0x%X X: 0x%X Y: 0x%X SP: 0x%X", cpu.A, cpu.X, cpu.Y, cpu.StackPointer)
}

func (v *Video) Init(t <-chan []*Tile) {
	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		log.Fatal(sdl.GetError())
	}

	if ttf.Init() != 0 {
		log.Fatal(sdl.GetError())
	}

	v.screen = sdl.SetVideoMode(256, 240, 32, sdl.RESIZABLE)

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

func (v *Video) DrawTile(t *Tile) {
	for y, r := range t.Rows {
		for x, p := range r.Pixels {
			rect := sdl.Rect{
				X: int16(t.X + x),
				Y: int16(t.Y + y),
				W: 4,
				H: 4,
			}

			var color uint32
			switch p {
			case 0:
				color = 0x222222
			case 1:
				color = 0x555555
			case 2:
				color = 0xAAAAAA
			case 3:
				color = 0xFFFFFF
			}

			v.screen.FillRect(&rect, color)
		}
	}
}

func (v *Video) DrawFrame(tiles []*Tile) {
	for c := 0; c < 30; c++ {
		for r := 0; r < 32; r++ {
			tile := tiles[c*32+r]
            v.DrawTile(tile)
		}
	}
}

func (v *Video) Render() {
	for {
		select {
		case val := <-v.tick:
            v.DrawFrame(val)
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
