package main

import (
	"math"
)

const (
	StatusSpriteOverflow = iota
	StatusSprite0Hit
	StatusVblankStarted

	MirroringVertical
	MirroringHorizontal
	MirroringSingleLower
	MirroringSingleUpper
)

type SpriteData struct {
	Tiles        [256]Word
	YCoordinates [256]Word
	Attributes   [256]Word
	XCoordinates [256]Word
}

type Flags struct {
	BaseNametableAddress     Word
	VramAddressInc           Word
	SpritePatternAddress     Word
	BackgroundPatternAddress Word
	SpriteSize               Word
	MasterSlaveSel           Word
	NmiOnVblank              Word
}

type Masks struct {
	Grayscale            bool
	ShowBackgroundOnLeft bool
	ShowSpritesOnLeft    bool
	ShowBackground       bool
	ShowSprites          bool
	IntensifyReds        bool
	IntensifyGreens      bool
	IntensifyBlues       bool
}

type AddressRegisters struct {
	FV  Word
	V   Word
	H   Word
	VT  Word
	HT  Word
	FH  Word
	S   Word
	PAR Word
	AR  Word
}

type AddressCounters struct {
	FV Word
	V  Word
	H  Word
	VT Word
	HT Word
}

type Registers struct {
	Control    Word
	Mask       Word
	Status     Word
	OamAddress int
	OamData    Word
	Scroll     Word
	Address    int
	Data       Word
	FirstWrite bool
}

type Ppu struct {
	Registers
	Flags
	Masks
	SpriteData
	AddressRegister AddressRegisters
	AddressCounter  AddressCounters
	Vram            [0xFFFF]Word
	SpriteRam       [0x100]Word
	PaletteRam      [0x20]Word
	Mirroring       int

	Framebuffer []int

	Output chan []int
	Cycle  int
}

func (r *AddressRegisters) Address() int {
	high := (r.FV & 0x07) << 4
	high = high | ((r.V & 0x01) << 3)
	high = high | ((r.H & 0x01) << 2)
	high = high | ((r.VT >> 3) & 0x03)

	low := (r.VT & 0x07) << 5
	low = low | (r.HT & 0x1F)

	return ((int(high) << 8) | int(low)) & 0x7FFF
}

func (r *AddressRegisters) InitFromAddress(a int) {
	high := Word((a >> 8) & 0xFF)
	low := Word(a & 0xFF)

	r.FV = (high >> 4) & 0x07
	r.V = (high >> 3) & 0x01
	r.H = (high >> 2) & 0x01
	r.VT = (r.VT & 7) | ((high & 0x03) << 3)

	r.VT = (r.VT & 24) | ((low >> 0x05) & 7)
	r.HT = low & 0x1F
}

func (c *AddressCounters) Address() int {
	high := (c.FV & 0x07) << 4
	high = high | ((c.V & 0x01) << 3)
	high = high | ((c.H & 0x01) << 2)
	high = high | ((c.VT >> 3) & 0x03)

	low := (c.VT & 0x07) << 5
	low = low | (c.HT & 0x1F)

	return ((int(high) << 8) | int(low)) & 0x7FFF
}

func (c *AddressCounters) InitFromAddress(a int) {
	high := Word((a >> 8) & 0xFF)
	low := Word(a & 0xFF)

	c.FV = (high >> 4) & 0x03
	c.V = (high >> 3) & 0x01
	c.H = (high >> 2) & 0x01
	c.VT = (c.VT & 0x07) | ((high & 0x03) << 3)

	c.VT = (c.VT & 0x18) | ((low >> 5) & 0x07)
	c.HT = low & 0x1F
}

func (p *Ppu) Init() chan []int {
	p.FirstWrite = true
	p.Output = make(chan []int)

	p.Cycle = 0

	for i, _ := range p.Vram {
		p.Vram[i] = 0x00
	}

	for i, _ := range p.SpriteRam {
		p.SpriteRam[i] = 0x00
	}

	p.Framebuffer = make([]int, 0xF000)

	return p.Output
}

func (p *Ppu) writeNametableData(a int, v Word) {
	switch p.Mirroring {
	case MirroringVertical:
		if a >= 0x2000 && a < 0x2400 {
			p.Vram[a] = v
			p.Vram[0x2800+(a-0x2000)] = v
		} else if a >= 0x2800 && a < 0x2C00 {
			p.Vram[a] = v
			p.Vram[0x2000+(a-0x2800)] = v
		} else if a >= 0x2400 && a < 0x2800 {
			p.Vram[a] = v
			p.Vram[0x2C00+(a-0x2400)] = v
		} else if a >= 0x2C00 && a < 0x3000 {
			p.Vram[a] = v
			p.Vram[0x2400+(a-0x2C00)] = v
		}
	case MirroringHorizontal:
		if a >= 0x2000 && a < 0x2400 {
			p.Vram[a] = v
			p.Vram[0x2400+(a-0x2000)] = v
		} else if a >= 0x2400 && a < 0x2800 {
			p.Vram[a] = v
			p.Vram[0x2000+(a-0x2400)] = v
		} else if a >= 0x2800 && a < 0x2C00 {
			p.Vram[a] = v
			p.Vram[0x2C00+(a-0x2800)] = v
		} else if a >= 0x2C00 && a < 0x3000 {
			p.Vram[a] = v
			p.Vram[0x2800+(a-0x2C00)] = v
		}
	}
}

// Writes to mirrored regions of VRAM
func (p *Ppu) writeMirroredVram(a int, v Word) {
	if a >= 0x3F00 && a < 0x3F20 {
		// Palette table entries
		p.PaletteRam[a-0x3F00] = v
	} else if a >= 0x3F20 && a < 0x3F40 {
		// Palette table entries
		p.PaletteRam[a-0x3F20] = v
	} else if a >= 0x3F40 && a < 0x3F80 {
		// Palette table entries
		p.PaletteRam[a-0x3F40] = v
	} else if a >= 0x3F80 && a < 0x3FC0 {
		// Palette table entries
		p.PaletteRam[a-0x3F80] = v
	} else {
		p.Vram[a-0x1000] = v
	}
}

func (p *Ppu) Step() {
	switch p.Cycle {
	case 89001:
		// We're in VBlank
		p.setStatus(StatusVblankStarted)
		// Request NMI
		cpu.RequestInterrupt(InterruptNmi)

		p.RenderNametable(0)
		p.RenderSprites()

		p.Output <- p.Framebuffer
		p.Cycle = 0
	}

	p.Cycle++
}

// $2000
func (p *Ppu) WriteControl(v Word) {
	p.Control = v

	// Control flag
	// 7654 3210
	// |||| ||||
	// |||| ||++- Base nametable address
	// |||| ||    (0 = $2000; 1 = $2400; 2 = $2800; 3 = $2C00)
	// |||| |+--- VRAM address increment per CPU read/write of PPUDATA
	// |||| |     (0: increment by 1, going across; 1: increment by 32, going down)
	// |||| +---- Sprite pattern table address for 8x8 sprites
	// ||||       (0: $0000; 1: $1000; ignored in 8x16 mode)
	// |||+------ Background pattern table address (0: $0000; 1: $1000)
	// ||+------- Sprite size (0: 8x8; 1: 8x16)
	// |+-------- PPU master/slave select (has no effect on the NES)
	// +--------- Generate an NMI at the start of the
	//            vertical blanking interval (0: off; 1: on)
	p.BaseNametableAddress = v & 0x03
	p.VramAddressInc = (v >> 2) & 0x01
	p.SpritePatternAddress = (v >> 3) & 0x01
	p.BackgroundPatternAddress = (v >> 4) & 0x01
	p.SpriteSize = (v >> 5) & 0x01
	p.NmiOnVblank = (v >> 7) & 0x01
}

// $2001
func (p *Ppu) WriteMask(v Word) {
	p.Mask = v

	// 76543210
	// ||||||||
	// |||||||+- Grayscale (0: normal color; 1: produce a monochrome display)
	// ||||||+-- 1: Show background in leftmost 8 pixels of screen; 0: Hide
	// |||||+--- 1: Show sprites in leftmost 8 pixels of screen; 0: Hide
	// ||||+---- 1: Show background
	// |||+----- 1: Show sprites
	// ||+------ Intensify reds (and darken other colors)
	// |+------- Intensify greens (and darken other colors)
	// +-------- Intensify blues (and darken other colors)
	p.Grayscale = (v&0x01 == 0x01)
	p.ShowBackgroundOnLeft = (((v >> 1) & 0x01) == 0x01)
	p.ShowSpritesOnLeft = (((v >> 2) & 0x01) == 0x01)
	p.ShowBackground = (((v >> 3) & 0x01) == 0x01)
	p.ShowSprites = (((v >> 4) & 0x01) == 0x01)
	p.IntensifyReds = (((v >> 5) & 0x01) == 0x01)
	p.IntensifyGreens = (((v >> 6) & 0x01) == 0x01)
	p.IntensifyBlues = (((v >> 7) & 0x01) == 0x01)
}

func (p *Ppu) clearStatus(s Word) {
	current := Ram[0x2002]

	switch s {
	case StatusSpriteOverflow:
		current = current & 0xDF
	case StatusSprite0Hit:
		current = current & 0xBF
	case StatusVblankStarted:
		current = current & 0x7F
	}

	Ram[0x2002] = current
}

func (p *Ppu) setStatus(s Word) {
	current := Ram[0x2002]

	switch s {
	case StatusSpriteOverflow:
		current = current | 0x20
	case StatusSprite0Hit:
		current = current | 0x40
	case StatusVblankStarted:
		current = current | 0x80
	}

	Ram[0x2002] = current
}

// $2002
func (p *Ppu) ReadStatus() (s Word, e error) {
	p.FirstWrite = true
	s = Ram[0x2002]

	// Clear VBlank flag
	p.clearStatus(StatusVblankStarted)

	return
}

// $2003
func (p *Ppu) WriteOamAddress(v Word) {
	p.OamAddress = int(v)
}

// $2004
func (p *Ppu) WriteOamData(v Word) {
	p.OamData = v
}

// $4014
func (p *Ppu) WriteDma(v Word) {
	// Halt the CPU for 512 cycles
	cpu.CyclesToWait = 512

	addr := int(v) * 0x100
	for i := 0; i < 256; i++ {
		d, _ := Ram.Read(addr + i)
		p.SpriteRam[i] = d
		p.updateBufferedSpriteMem(i, d)
	}
}

func (p *Ppu) updateBufferedSpriteMem(a int, v Word) {
	i := int(math.Floor(float64(a / 4)))

	switch a % 4 {
	case 0x0:
		p.YCoordinates[i] = v
	case 0x1:
		p.Tiles[i] = v
	case 0x2:
		// Attribute
		p.Attributes[i] = v
	case 0x3:
		p.XCoordinates[i] = v
	}
}

// $2004
func (p *Ppu) ReadOamData() (Word, error) {
	return Ram[0x2004], nil
}

// $2005
func (p *Ppu) WriteScroll(v Word) {
	p.Scroll = v
}

// $2006
func (p *Ppu) WriteAddress(v Word) {
	// http://nesdev.com/PPU%20addressing.txt
	// 
	// …ÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕ—ÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕª
	// ∫2000           ≥      1  0                     4               ∫
	// ∫2005/1         ≥                   76543  210                  ∫
	// ∫2005/2         ≥ 210        76543                              ∫
	// ∫2006/1         ≥ -54  3  2  10                                 ∫
	// ∫2006/2         ≥              765  43210                       ∫
	// ∫NT read        ≥                                  76543210     ∫
	// ∫AT read (4)    ≥                                            10 ∫
	// «ƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒ≈ƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒ∂
	// ∫               ≥…ÕÕÕª…Õª…Õª…ÕÕÕÕÕª…ÕÕÕÕÕª…ÕÕÕª…Õª…ÕÕÕÕÕÕÕÕª…ÕÕª∫
	// ∫PPU registers  ≥∫ FV∫∫V∫∫H∫∫   VT∫∫   HT∫∫ FH∫∫S∫∫     PAR∫∫AR∫∫
	// ∫PPU counters   ≥«ƒƒƒ∂«ƒ∂«ƒ∂«ƒƒƒƒƒ∂«ƒƒƒƒƒ∂»ÕÕÕº»Õº»ÕÕÕÕÕÕÕÕº»ÕÕº∫
	// ∫               ≥»ÕÕÕº»Õº»Õº»ÕÕÕÕÕº»ÕÕÕÕÕº                      ∫
	// «ƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒ≈ƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒƒ∂
	// ∫2007 access    ≥  DC  B  A  98765  43210                       ∫
	// ∫NT read (1)    ≥      B  A  98765  43210                       ∫
	// ∫AT read (1,2,4)≥      B  A  543c   210b                        ∫
	// ∫PT read (3)    ≥ 210                           C  BA987654     ∫
	// »ÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕœÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕÕº
	if p.FirstWrite {
		p.AddressRegister.FV = (v >> 4) & 0x03
		p.AddressRegister.V = (v >> 3) & 0x01
		p.AddressRegister.H = (v >> 2) & 0x01

		p.AddressRegister.VT = (p.AddressRegister.VT & 7) | ((v & 0x03) << 3)
	} else {
		p.AddressRegister.VT = (p.AddressRegister.VT & 0x18) | ((v >> 5) & 0x07)
		p.AddressRegister.HT = v & 0x1F

		p.AddressCounter.FV = p.AddressRegister.FV
		p.AddressCounter.V = p.AddressRegister.V
		p.AddressCounter.H = p.AddressRegister.H
		p.AddressCounter.VT = p.AddressRegister.VT
		p.AddressCounter.HT = p.AddressRegister.HT

		p.Address = p.AddressCounter.Address()
	}

	p.FirstWrite = !p.FirstWrite
}

// $2007
func (p *Ppu) WriteData(v Word) {
	p.Address = p.AddressCounter.Address()

	if p.Address > 0x3000 {
		p.writeMirroredVram(p.Address, v)
	} else if p.Address >= 0x2000 && p.Address < 0x3000 {
		// Nametable mirroring
		p.writeNametableData(p.Address, v)
	} else {
		p.Vram[p.Address] = v
	}

	switch p.VramAddressInc {
	case 0x01:
		p.Address = p.Address + 0x20
	default:
		p.Address = p.Address + 0x01
	}

	p.AddressCounter.InitFromAddress(p.Address)
}

// $2007
func (p *Ppu) ReadData() (Word, error) {
	return Ram[0x2007], nil
}

func (p *Ppu) SprPatternTableAddress(i Word) int {
	var a int
	if p.SpritePatternAddress == 0x01 {
		a = 0x1000
	} else {
		a = 0x0
	}

	return int(i)*0x10 + a
}

func (p *Ppu) BgPatternTableAddress(i Word) int {
	var a int
	if p.BackgroundPatternAddress == 0x01 {
		a = 0x1000
	} else {
		a = 0x0
	}

	return int(i)*0x10 + a
}

func (p *Ppu) RenderNametable(table Word) {
	var a int

	if p.Mirroring == MirroringVertical {
		switch table {
		case 0:
			a = 0x2000
		case 1:
			a = 0x2800
		case 2:
			a = 0x2400
		case 3:
			a = 0x2C00
		}
	} else if p.Mirroring == MirroringHorizontal {
		switch table {
		case 0:
			a = 0x2000
		case 1:
			a = 0x2400
		case 2:
			a = 0x2800
		case 3:
			a = 0x2C00
		}
	}

	x := 0
	y := 0

	leftTile := true
	topTile := true
	// Generates each tile and applies the palette
	for i := a; i < a+0x3C0; i++ {
		attrBase := (i & ^0x3FF) + 0x3C0
		attr := uint(p.Vram[attrBase+((i&0x1F)>>2)+((i&0x3E0)>>7)*8])

		var attrValue uint
		switch {
		case leftTile && topTile:
			// Top left
			attrValue = attr & 0x03
		case !leftTile && topTile:
			// Top right
			attrValue = (attr >> 2) & 0x03
		case leftTile && !topTile:
			// Bottom left
			attrValue = (attr >> 4) & 0x03
		case !leftTile && !topTile:
			// Bottom right
			attrValue = (attr >> 6) & 0x03
		}

		t := p.BgPatternTableAddress(p.Vram[i])
		p.DecodePatternTile(t, x, y, p.BgPaletteEntry(attrValue), false, false)

		leftTile = !leftTile

		x += 8

		if x > 255 {
			x = 0
			y += 8

			topTile = !topTile
			leftTile = true
		}
	}
}

func (p *Ppu) DecodePatternTile(t, x, y int, pal []Word, hrev, yrev bool) {
	tile := p.Vram[t : t+16]

	l := len(tile)
	for i := 0; i < l/2; i++ {
		var b uint
		for b = 0; b < 8; b++ {
			ycoord := y + i

            var xcoord int

            if hrev {
                xcoord = x + int(b)
            } else {
                xcoord = x + int(7-b)
            }

			fbRow := ycoord*256 + xcoord

			// Store the bit 0/1
			pixel := (tile[i] >> b) & 0x01
			pixel += ((tile[i+8] >> b & 0x01) << 1)

			// Set the color of the pixel in the buffer
            p.Framebuffer[fbRow] = PaletteRgb[int(pal[pixel])]
		}
	}
}

func (p *Ppu) RenderSprites() {
	for i, t := range p.SpriteData.Tiles {
		if t == 0x0 {
			break
		}

		priority := (p.Attributes[i] >> 5) & 0x1

		if priority != 0x0 {
			continue
		}

        if p.XCoordinates[i] >= 0xF0 || p.YCoordinates[i] >= 0xF0 {
            continue
        }

		attrValue := p.Attributes[i] & 0x3

		p.DecodePatternTile(p.SprPatternTableAddress(t),
			int(p.XCoordinates[i]),
			int(p.YCoordinates[i]) + 1,
			p.SprPaletteEntry(uint(attrValue)),
            (p.Attributes[i] >> 6) & 0x1 != 0,
            (p.Attributes[i] >> 7) & 0x1 != 0)
	}
}

func (p *Ppu) BgPaletteEntry(a uint) (pal []Word) {
	switch a {
	case 0x0:
		pal = []Word{
			p.PaletteRam[0x00],
			p.PaletteRam[0x01],
			p.PaletteRam[0x02],
			p.PaletteRam[0x03],
		}
	case 0x1:
		pal = []Word{
			p.PaletteRam[0x00],
			p.PaletteRam[0x05],
			p.PaletteRam[0x06],
			p.PaletteRam[0x07],
		}
	case 0x2:
		pal = []Word{
			p.PaletteRam[0x00],
			p.PaletteRam[0x09],
			p.PaletteRam[0x0A],
			p.PaletteRam[0x0B],
		}
	case 0x3:
		pal = []Word{
			p.PaletteRam[0x00],
			p.PaletteRam[0x0D],
			p.PaletteRam[0x0E],
			p.PaletteRam[0x0F],
		}
	}

	return
}

func (p *Ppu) SprPaletteEntry(a uint) (pal []Word) {
	switch a {
	case 0x0:
		pal = []Word{
			p.PaletteRam[0x10],
			p.PaletteRam[0x11],
			p.PaletteRam[0x12],
			p.PaletteRam[0x13],
		}
	case 0x1:
		pal = []Word{
			p.PaletteRam[0x10],
			p.PaletteRam[0x15],
			p.PaletteRam[0x16],
			p.PaletteRam[0x17],
		}
	case 0x2:
		pal = []Word{
			p.PaletteRam[0x10],
			p.PaletteRam[0x19],
			p.PaletteRam[0x1A],
			p.PaletteRam[0x1B],
		}
	case 0x3:
		pal = []Word{
			p.PaletteRam[0x10],
			p.PaletteRam[0x1D],
			p.PaletteRam[0x1E],
			p.PaletteRam[0x1F],
		}
	}

	return
}
