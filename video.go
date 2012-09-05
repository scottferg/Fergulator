package main

import (
	"fmt"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/ttf"
	"log"
	"math"
)

var (
    PaletteRgb = []uint32{
        0x7C7C7C,
        0x0000FC,
        0x0000BC,
        0x4428BC,
        0x940084,
        0xA80020,
        0xA81000,
        0x881400,
        0x503000,
        0x007800,
        0x006800,
        0x005800,
        0x004058,
        0x000000,
        0x000000,
        0x000000,
        0xBCBCBC,
        0x0078F8,
        0x0058F8,
        0x6844FC,
        0xD800CC,
        0xE40058,
        0xF83800,
        0xE45C10,
        0xAC7C00,
        0x00B800,
        0x00A800,
        0x00A844,
        0x008888,
        0x000000,
        0x000000,
        0x000000,
        0xF8F8F8,
        0x3CBCFC,
        0x6888FC,
        0x9878F8,
        0xF878F8,
        0xF85898,
        0xF87858,
        0xFCA044,
        0xF8B800,
        0xB8F818,
        0x58D854,
        0x58F898,
        0x00E8D8,
        0x787878,
        0x000000,
        0x000000,
        0xFCFCFC,
        0xA4E4FC,
        0xB8B8F8,
        0xD8B8F8,
        0xF8B8F8,
        0xF8A4C0,
        0xF0D0B0,
        0xFCE0A8,
        0xF8D878,
        0xD8F878,
        0xB8F8B8,
        0xB8F8D8,
        0x00FCFC,
        0xF8D8F8,
        0x000000,
        0x000000,
    }
)

type Video struct {
	screen *sdl.Surface
	font   *ttf.Font
	tick   <-chan Nametable
}

func (cpu *Cpu) DumpRegisterState() string {
	return fmt.Sprintf("A: 0x%X X: 0x%X Y: 0x%X SP: 0x%X", cpu.A, cpu.X, cpu.Y, cpu.StackPointer)
}

func (v *Video) Init(t <-chan Nametable) {
	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		log.Fatal(sdl.GetError())
	}

	if ttf.Init() != 0 {
		log.Fatal(sdl.GetError())
	}

	// v.screen = sdl.SetVideoMode(256, 240, 32, sdl.RESIZABLE)
	v.screen = sdl.SetVideoMode(512, 480, 32, sdl.RESIZABLE)

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
				X: int16(t.X+x) + xoff,
				Y: int16(t.Y+y) + yoff,
				W: 1,
				H: 1,
			}

			var color uint32
            if int(p) < len(PaletteRgb) {
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

func (v *Video) Render() {
	for {
		select {
		case ev := <-sdl.Events:
			switch e := ev.(type) {
			case sdl.KeyboardEvent:
				if e.Keysym.Sym == sdl.K_ESCAPE {
					running = false
				}
			}
		case val := <-v.tick:
			v.DrawFrame(val.Table0, 0, 0)
			v.DrawFrame(val.Table1, 0, 240)
			v.DrawFrame(val.Table2, 256, 0)
			v.DrawFrame(val.Table3, 256, 240)
			v.screen.Flip()
		}
	}
}

func (v *Video) Close() {
	v.font.Close()

	ttf.Quit()
	sdl.Quit()
}

func HsvToRgb(val Word) (rgb uint32) {
	var h, s, v float64
	var fR, fG, fB float64

	h = float64(val & 0xF)
	s = 0.5
	v = float64(val & 0x30)

	i := math.Floor(h * 6)
	f := h*6 - i
	p := v * (1.0 - s)
	q := v * (1.0 - f*s)
	t := v * (1.0 - (1.0-f)*s)
	switch int(i) % 6 {
	case 0:
		fR, fG, fB = v, t, p
	case 1:
		fR, fG, fB = q, v, p
	case 2:
		fR, fG, fB = p, v, t
	case 3:
		fR, fG, fB = p, q, v
	case 4:
		fR, fG, fB = t, p, v
	case 5:
		fR, fG, fB = v, p, q
	}

	r := uint32((fR * 255) + 0.5)
	g := uint32((fG * 255) + 0.5)
	b := uint32((fB * 255) + 0.5)

	rgb = (r << 16) | (g << 8) | b
	return
}
