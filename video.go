package main

import (
	"fmt"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/ttf"
	"log"
)

var (
    PaletteRgb = []uint32{
        0x7C7C7C, 0x0000FC, 0x0000BC, 0x4428BC, 0x940084,
        0xA80020, 0xA81000, 0x881400, 0x503000, 0x007800,
        0x006800, 0x005800, 0x004058, 0x000000, 0x000000,
        0x000000, 0xBCBCBC, 0x0078F8, 0x0058F8, 0x6844FC,
        0xD800CC, 0xE40058, 0xF83800, 0xE45C10, 0xAC7C00,
        0x00B800, 0x00A800, 0x00A844, 0x008888, 0x000000,
        0x000000, 0x000000, 0xF8F8F8, 0x3CBCFC, 0x6888FC,
        0x9878F8, 0xF878F8, 0xF85898, 0xF87858, 0xFCA044,
        0xF8B800, 0xB8F818, 0x58D854, 0x58F898, 0x00E8D8,
        0x787878, 0x000000, 0x000000, 0xFCFCFC, 0xA4E4FC,
        0xB8B8F8, 0xD8B8F8, 0xF8B8F8, 0xF8A4C0, 0xF0D0B0,
        0xFCE0A8, 0xF8D878, 0xD8F878, 0xB8F8B8, 0xB8F8D8,
        0x00FCFC, 0xF8D8F8, 0x000000, 0x000000,
    }
)

type Video struct {
	screen *sdl.Surface
	font   *ttf.Font
	tick   <-chan Frame
}

func (cpu *Cpu) DumpRegisterState() string {
	return fmt.Sprintf("A: 0x%X X: 0x%X Y: 0x%X SP: 0x%X", cpu.A, cpu.X, cpu.Y, cpu.StackPointer)
}

func (v *Video) Init(t <-chan Frame) {
	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		log.Fatal(sdl.GetError())
	}

	if ttf.Init() != 0 {
		log.Fatal(sdl.GetError())
	}

	v.screen = sdl.SetVideoMode(256 * int(sizeFactor), 240 * int(sizeFactor), 32, sdl.RESIZABLE)

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

func (v *Video) DrawTile(t *Tile, xoff, yoff int16) {
	for y, r := range t.Rows {
		for x, p := range r.Pixels {
			rect := sdl.Rect{
				X: int16(t.X+x) + int16(sizeFactor) + xoff,
				Y: int16(t.Y+y) + int16(sizeFactor) + yoff,
				W: 1 * sizeFactor,
				H: 1 * sizeFactor,
			}

			var color uint32
            
            if len(t.Palette) == 0 {
                color = 0xEEEEEE
            } else if int(p) < len(PaletteRgb) {
                switch p {
                case 0:
                    color = PaletteRgb[t.Palette[0]]
                case 1:
                    color = PaletteRgb[t.Palette[1]]
                case 2:
                    color = PaletteRgb[t.Palette[2]]
                case 3:
                    color = PaletteRgb[t.Palette[3]]
                }
            } else {
                color = 0
            }

			v.screen.FillRect(&rect, color)
		}
	}
}

func (v *Video) DrawFrame(tiles []*Tile, x, y int16) {
	for c := 0; c < 30; c++ {
		for r := 0; r < 32; r++ {
			tile := tiles[c*32+r]
			v.DrawTile(tile, x, y)
		}
	}
}

func (v *Video) DrawSprites(tiles []*Tile) {
    for _, tile := range tiles {
        v.DrawTile(tile, 0, 0)
    }
}

func (v *Video) Render() {
	for {
		select {
		case val := <-v.tick:
			v.DrawFrame(val.Background, 0, 0)
			v.DrawSprites(val.Sprites)
			v.screen.Flip()
		}
	}
}

func (v *Video) Close() {
	v.font.Close()

	ttf.Quit()
	sdl.Quit()
}
