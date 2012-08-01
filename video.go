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
			registers := ttf.RenderText_Solid(v.font, val.DumpRegisterState(), white)
			opcode := ttf.RenderText_Solid(v.font, fmt.Sprintf("Opcode: 0x%X", val.Opcode), white)
			pc := ttf.RenderText_Solid(v.font, fmt.Sprintf("PC: 0x%X", programCounter), white)
			status := ttf.RenderText_Solid(v.font, fmt.Sprintf("P: 0x%X", val.P), white)
			negative := ttf.RenderText_Solid(v.font, fmt.Sprintf("Neg Flag: %t", val.getNegative()), white)
			test := ttf.RenderText_Solid(v.font, fmt.Sprintf("Test Output: %s", *Ram[0x6004]), white)

			ppuctl := make([]*sdl.Surface, 8)

			ppuctl[0] = ttf.RenderText_Solid(v.font, fmt.Sprintf("PPUCTL: 0x%X", *Ram[0x2000]), white)
			ppuctl[1] = ttf.RenderText_Solid(v.font, fmt.Sprintf("PPUMASK: 0x%X", *Ram[0x2001]), white)
			ppuctl[2] = ttf.RenderText_Solid(v.font, fmt.Sprintf("PPUSTATUS: 0x%X", *Ram[0x2002]), white)
			ppuctl[3] = ttf.RenderText_Solid(v.font, fmt.Sprintf("OAMADDR: 0x%X", *Ram[0x2003]), white)
			ppuctl[4] = ttf.RenderText_Solid(v.font, fmt.Sprintf("OAMDATA: 0x%X", *Ram[0x2004]), white)
			ppuctl[5] = ttf.RenderText_Solid(v.font, fmt.Sprintf("PPUSCROLL: 0x%X", *Ram[0x2005]), white)
			ppuctl[6] = ttf.RenderText_Solid(v.font, fmt.Sprintf("PPUADDR: 0x%X", *Ram[0x2006]), white)
			ppuctl[7] = ttf.RenderText_Solid(v.font, fmt.Sprintf("PPUDATA: 0x%X", *Ram[0x2007]), white)

			v.screen.FillRect(nil, 0x000000)

			ppuY := int16(320)
			for _, val := range ppuctl {
				v.screen.Blit(&sdl.Rect{460, ppuY, 0, 0}, val, nil)
				ppuY = ppuY + 20
			}

			v.screen.Blit(&sdl.Rect{2, 360, 0, 0}, test, nil)
			v.screen.Blit(&sdl.Rect{2, 380, 0, 0}, negative, nil)
			v.screen.Blit(&sdl.Rect{2, 400, 0, 0}, status, nil)
			v.screen.Blit(&sdl.Rect{2, 420, 0, 0}, opcode, nil)
			v.screen.Blit(&sdl.Rect{2, 440, 0, 0}, pc, nil)
			v.screen.Blit(&sdl.Rect{2, 460, 0, 0}, registers, nil)

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
