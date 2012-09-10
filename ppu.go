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

type Registers struct {
	Control        Word
	Mask           Word
	Status         Word
	OamAddress     int
	OamData        Word
	VramDataBuffer Word
	VramAddress    int
	VramLatch      int
	FineX          int
	Data           Word
	FirstWrite     bool
}

type Ppu struct {
	Registers
	Flags
	Masks
	SpriteData
	Vram              [0xFFFF]Word
	SpriteRam         [0x100]Word
	PaletteRam        [0x20]Word
	AttributeLocation [0x400]uint
	AttributeShift    [0x400]uint
	Mirroring         int

	Framebuffer []int

	Output   chan []int
	Cycle    int
	Scanline int
}

func (p *Ppu) Init() chan []int {
	p.FirstWrite = true
	p.Output = make(chan []int)

	p.Cycle = 0
	p.Scanline = -1

	for i, _ := range p.Vram {
		p.Vram[i] = 0x00
	}

	for i, _ := range p.SpriteRam {
		p.SpriteRam[i] = 0x00
	}

	for i, _ := range p.AttributeShift {
		x := uint(i)
		p.AttributeShift[i] = ((x >> 4) & 0x04) | (x & 0x02)
		p.AttributeLocation[i] = ((x >> 2) & 0x07) | (((x >> 4) & 0x38) | 0x3C0)
	}

	p.Framebuffer = make([]int, 0xF000)

	return p.Output
}

func (p *Ppu) writeNametableData(a int, v Word) {
	switch p.Mirroring {
	case MirroringVertical:
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
	case MirroringHorizontal:
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
	}
}

// Writes to mirrored regions of VRAM
func (p *Ppu) writeMirroredVram(a int, v Word) {
    // Background Palettes
	if a > 0x3F00 && a < 0x3F10 {
		// Palette table entries
		p.PaletteRam[a-0x3F00] = v
	} else if a >= 0x3F20 && a < 0x3F30 {
		// Palette table entries
		p.PaletteRam[a-0x3F20] = v
	} else if a >= 0x3F40 && a < 0x3F50 {
		// Palette table entries
		p.PaletteRam[a-0x3F40] = v
	} else if a >= 0x3F60 && a < 0x3F70 {
		// Palette table entries
		p.PaletteRam[a-0x3F60] = v
	} else if a >= 0x3F80 && a < 0x3F90 {
		// Palette table entries
		p.PaletteRam[a-0x3F80] = v
	} else if a >= 0x3FA0 && a < 0x3FB0 {
		// Palette table entries
		p.PaletteRam[a-0x3FA0] = v
	} else if a >= 0x3FC0 && a < 0x3FD0 {
		// Palette table entries
		p.PaletteRam[a-0x3FC0] = v
	} else if a >= 0x3FE0 && a < 0x3FF0 {
		// Palette table entries
		p.PaletteRam[a-0x3FE0] = v
	} 

    // Mirrored entries
    if a == 0x3F00 {
        p.PaletteRam[0x10] = v
        p.PaletteRam[0x00] = v
    } else if a == 0x3F10 {
        p.PaletteRam[0x00] = v
        p.PaletteRam[0x10] = v
    } else if a == 0x3F14 {
        p.PaletteRam[0x04] = v
    } else if a == 0x3F18 {
        p.PaletteRam[0x08] = v
    } else if a == 0x3F1C {
        p.PaletteRam[0x0C] = v
    }

    // Sprite palettes
    if a > 0x3F10 && a < 0x3F20 {
		// Palette table entries
		p.PaletteRam[a-0x3F00] = v
	} else if a >= 0x3F30 && a < 0x3F40 {
		// Palette table entries
		p.PaletteRam[a-0x3F20] = v
	} else if a >= 0x3F50 && a < 0x3F60 {
		// Palette table entries
		p.PaletteRam[a-0x3F40] = v
	} else if a >= 0x3F70 && a < 0x3F80 {
		// Palette table entries
		p.PaletteRam[a-0x3F70] = v
	} else if a >= 0x3F90 && a < 0x3FA0 {
		// Palette table entries
		p.PaletteRam[a-0x3F90] = v
	} else if a >= 0x3FB0 && a < 0x3FC0 {
		// Palette table entries
		p.PaletteRam[a-0x3FB0] = v
	} else if a >= 0x3FD0 && a < 0x3FE0 {
		// Palette table entries
		p.PaletteRam[a-0x3FD0] = v
	} else {
		p.Vram[a-0x1000] = v
	}
}

func (p *Ppu) Step() {
	switch {
	case p.Scanline == 240:
		// fmt.Println("Scanline 240")
		// We're in VBlank
		p.setStatus(StatusVblankStarted)
		// Request NMI
		cpu.RequestInterrupt(InterruptNmi)

		// p.renderNametable(p.BaseNametableAddress)
        if p.ShowSprites {
            p.renderSprites()
        }

		p.Output <- p.Framebuffer
		p.Cycle = 0
		p.Scanline++
	case p.Scanline == 261:
		// End of vblank
		// fmt.Println("Scanline 261")
		p.Scanline = -1
		p.Cycle = 0
	case p.Scanline < 240 && p.Scanline > 0:
		// Render 1 row of 8x8 tiles
		if p.Cycle == 341 {
			// fmt.Println("End scanline")
			p.Cycle = 0
			if p.Scanline%8 == 0 {
				if p.ShowBackground {
					p.renderTileRow()
				}
			}
			p.Scanline++
		}
		// fmt.Println("Scanline 240")
	case p.Scanline == -1:
		// fmt.Println("Scanline -1")
		if p.Cycle == 304 {
			// Copy scroll latch into VRAMADDR register
			p.VramAddress = p.VramLatch
		}
	}

	if p.Cycle == 341 {
		p.Cycle = 0
		p.Scanline++
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

	p.VramLatch = p.VramLatch & 0x73FF
	p.VramLatch = p.VramLatch | ((int(v) & 0x03) << 10)
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

	// Fill sprite RAM
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
	if p.FirstWrite {
		p.VramLatch = p.VramLatch & 0x7FE0
		p.VramLatch = p.VramLatch | ((int(v) & 0xF8) >> 3)
		p.FineX = int(v) & 0x07
	} else {
		p.VramLatch = p.VramLatch & 0xC1F
		p.VramLatch = p.VramLatch | (((int(v) & 0xF8) << 2) | ((int(v) & 0x07) << 12))
	}

	p.FirstWrite = !p.FirstWrite
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
		p.VramLatch = p.VramLatch & 0xFF
		p.VramLatch = p.VramLatch | ((int(v) & 0x3F) << 8)
	} else {
		p.VramLatch = p.VramLatch & 0x7F00
		p.VramLatch = p.VramLatch | int(v)
		p.VramAddress = p.VramLatch
	}

	p.FirstWrite = !p.FirstWrite
}

// $2007
func (p *Ppu) WriteData(v Word) {
	if p.VramAddress > 0x3000 {
		p.writeMirroredVram(p.VramAddress, v)
	} else if p.VramAddress >= 0x2000 && p.VramAddress < 0x3000 {
		// Nametable mirroring
		p.writeNametableData(p.VramAddress, v)
	} else {
		p.Vram[p.VramAddress] = v
	}

	switch p.VramAddressInc {
	case 0x01:
		p.VramAddress = p.VramAddress + 0x20
	default:
		p.VramAddress = p.VramAddress + 0x01
	}
}

// $2007
func (p *Ppu) ReadData() (Word, error) {
    // Reads from $2007 are buffered with a
    // 1-byte delay
    tmp := p.VramDataBuffer
    p.VramDataBuffer = Ram[0x2007]

	return tmp, nil
}

func (p *Ppu) sprPatternTableAddress(i int) int {
	if p.SpriteSize&0x01 != 0x0 {
		// 8x16 Sprites
		var bank int
		if i&0x01 != 0 {
			bank = 0x1000
		} else {
			bank = 0x0000
		}

		return bank + ((int(i) >> 1) * 0x20)
	}

	// 8x8 Sprites
	var a int
	if p.SpritePatternAddress == 0x01 {
		a = 0x1000
	} else {
		a = 0x0
	}

	return int(i)*0x10 + a
}

func (p *Ppu) bgPatternTableAddress(i Word) int {
	var a int
	if p.BackgroundPatternAddress == 0x01 {
		a = 0x1000
	} else {
		a = 0x0
	}

	return int(i)*0x10 + a
}

func (p *Ppu) selectNametable(t int) (a int) {
	switch t {
	case 0:
		a = 0x2000
	case 1:
		a = 0x2800
	case 2:
		a = 0x2400
	case 3:
		a = 0x2C00
	}

    return
}

func (p *Ppu) renderNametable(table int) {
    a := p.selectNametable(table)

	x := 0
	y := 0

	// Generates each tile and applies the palette
	for i := a; i < a+0x3C0; i++ {
		attrAddr := 0x23C0 | (p.VramAddress & 0xC00)
		shift := p.AttributeShift[p.VramAddress&0x3FF]
		attr := p.Vram[attrAddr+((i&0x1F)>>2)+((i&0x3E0)>>7)*8]
		attr = (attr >> shift) & 0x03

		t := p.bgPatternTableAddress(p.Vram[i])
		p.decodePatternTile(t, x, y, p.bgPaletteEntry(attr), nil)

		x += 8

		if x > 255 {
			x = 0
			y += 8
		}
	}
}

func (p *Ppu) renderTileRow() {
	// Generates each tile and applies the palette
	for x := 0; x < 32; x++ {
		// for i := a; i < a+0x20; i++ {
		attrAddr := 0x23C0 | (p.VramAddress & 0xC00) | int(p.AttributeLocation[p.VramAddress&0x3FF])
		shift := p.AttributeShift[p.VramAddress&0x3FF]
		attr := ((p.Vram[attrAddr] >> shift) & 0x03) << 2

        ntAddress := p.selectNametable((p.VramAddress & 0xC00) >> 10)
		t := p.bgPatternTableAddress(p.Vram[p.VramAddress+ntAddress])
		p.decodePatternTile(t, x*8, p.Scanline-8, p.bgPaletteEntry(attr), nil)

		// Flip bit 10 on wraparound
		p.VramAddress++

		// If rendering is enabled, at the end of a scanline
        // copy bits 10 and 4-0 from VRAM latch into VRAMADDR
		p.VramAddress = p.VramAddress ^ (p.VramLatch & 0x41F)
	}
}

func (p *Ppu) decodePatternTile(t, x, y int, pal []Word, attr *Word) {
	tile := p.Vram[t : t+16]

	l := len(tile)
	for i := 0; i < l/2; i++ {
		var b uint
		for b = 0; b < 8; b++ {
			var xcoord int
			if attr != nil && (*attr>>6)&0x1 != 0 {
				xcoord = x + int(b)
			} else {
				xcoord = x + int(7-b)
			}

			var ycoord int
			if attr != nil && (*attr>>7)&0x1 != 0 {
				ycoord = y + int(7-b)
			} else {
				ycoord = y + i
			}

			fbRow := ycoord*256 + xcoord

			// Store the bit 0/1
			pixel := (tile[i] >> b) & 0x01
			pixel += ((tile[i+8] >> b & 0x01) << 1)

			trans := false
			if attr != nil && pixel == 0 {
				trans = true
			}

			// Set the color of the pixel in the buffer
			if fbRow < 0xF000 && !trans {
				p.Framebuffer[fbRow] = PaletteRgb[int(pal[pixel])]
			}
		}
	}
}

func (p *Ppu) renderSprites() {
	for i, t := range p.SpriteData.Tiles {
		attrValue := p.Attributes[i] & 0x3

		if p.SpriteSize&0x01 != 0x0 {
			// 8x16 Sprite
			p.decodePatternTile(p.sprPatternTableAddress(int(t)),
				int(p.XCoordinates[i]),
				int(p.YCoordinates[i])+1,
				p.sprPaletteEntry(uint(attrValue)),
				&p.Attributes[i])
		} else {
			p.decodePatternTile(p.sprPatternTableAddress(int(t)),
				int(p.XCoordinates[i]),
				int(p.YCoordinates[i])+1,
				p.sprPaletteEntry(uint(attrValue)),
				&p.Attributes[i])
		}
	}
}

func (p *Ppu) bgPaletteEntry(a Word) (pal []Word) {
	switch a {
	case 0x0:
		pal = []Word{
			p.PaletteRam[0x00],
			p.PaletteRam[0x01],
			p.PaletteRam[0x02],
			p.PaletteRam[0x03],
		}
	case 0x4:
		pal = []Word{
			p.PaletteRam[0x00],
			p.PaletteRam[0x05],
			p.PaletteRam[0x06],
			p.PaletteRam[0x07],
		}
	case 0x8:
		pal = []Word{
			p.PaletteRam[0x00],
			p.PaletteRam[0x09],
			p.PaletteRam[0x0A],
			p.PaletteRam[0x0B],
		}
	case 0xC:
		pal = []Word{
			p.PaletteRam[0x00],
			p.PaletteRam[0x0D],
			p.PaletteRam[0x0E],
			p.PaletteRam[0x0F],
		}
	}

	return
}

func (p *Ppu) sprPaletteEntry(a uint) (pal []Word) {
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
