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
	tick   <-chan Cpu
}

func (cpu *Cpu) DumpRegisterState() string {
	return fmt.Sprintf("A: 0x%X X: 0x%X Y: 0x%X SP: 0x%X", cpu.A, cpu.X, cpu.Y, cpu.StackPointer)
}

func (v *Video) Init(t <-chan Cpu) {
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

func (v *Video) Render() {
	for {
		select {
		case val := <-v.tick:
            fmt.Sprintf("P: 0x%X\n", val.P)
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
