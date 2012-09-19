package main

import (
	"fmt"
	"math"
)

const (
	StatusSpriteOverflow = iota
	StatusSprite0Hit
	StatusVblankStarted

	MirroringVertical
	MirroringHorizontal
	MirroringSingleScreen

	Pal  = 70
	Ntsc = 20
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

type Pixel struct {
	Color  int
	Value  int
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
	Control          Word
	Mask             Word
	Status           Word
	VramDataBuffer   Word
	VramAddress      int
	VramLatch        int
	SpriteRamAddress int
	FineX            Word
	Data             Word
	WriteLatch       bool
    HighBitShift     uint16
    LowBitShift      uint16
}

type Ppu struct {
	Registers
	Flags
	Masks
	SpriteData
	Vram              [0xFFFF]Word
	SpriteRam         [0x100]Word
	Nametables        Nametable
	PaletteRam        [0x20]Word
	AttributeLocation [0x400]uint
	AttributeShift    [0x400]uint
	Mirroring         int
	TilerowCounter    int

	Palettebuffer []Pixel
	Framebuffer   []int

	DebugMode  bool
	DebugCount int

	Output      chan []int
	Debug       chan []int
	Cycle       int
	Scanline    int
	Timestamp   int
	VblankTime  int
	FrameCount  int
	FrameCycles int
}

func (p *Ppu) Init() (chan []int, chan []int) {
	p.WriteLatch = true
	p.Output = make(chan []int)
	p.Debug = make(chan []int)

	p.Cycle = 0
	p.Scanline = -1
	p.FrameCount = 0

	p.VblankTime = 20 * 341 * 5 // NTSC
	p.TilerowCounter = 0

    p.Nametables.Init()

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

	p.Palettebuffer = make([]Pixel, 0xF000)
	p.Framebuffer = make([]int, 0xF000)

	return p.Output, p.Debug
}

func (p *Ppu) PpuRegRead(a int) (Word, error) {
	switch a & 0x7 {
	case 0x2:
		return p.ReadStatus()
	case 0x4:
		return p.ReadOamData()
	case 0x7:
		return p.ReadData()
	}

	return 0, nil
}

func (p *Ppu) PpuRegWrite(v Word, a int) {
	switch a & 0x7 {
	case 0x0:
		p.WriteControl(v)
	case 0x1:
		p.WriteMask(v)
	case 0x3:
		p.WriteOamAddress(v)
	case 0x4:
		p.WriteOamData(v)
	case 0x5:
		p.WriteScroll(v)
	case 0x6:
		p.WriteAddress(v)
	case 0x7:
		p.WriteData(v)
	}

	if a == 0x4014 {
		p.WriteDma(v)
	}
}

// Writes to mirrored regions of VRAM
func (p *Ppu) writeMirroredVram(a int, v Word) {
	if a >= 0x3F00 {
		base := a & 0x1F // 0b11111

		if base == 0x0 || base == 0x10 {
			p.PaletteRam[0x10] = v
			p.PaletteRam[0x00] = v
		} else {
			p.PaletteRam[base] = v
		}

		p.PaletteRam[0x10] = p.PaletteRam[0x0]
		p.PaletteRam[0x14] = p.PaletteRam[0x4]
		p.PaletteRam[0x18] = p.PaletteRam[0x8]
		p.PaletteRam[0x1C] = p.PaletteRam[0xC]
	} else {
		p.Vram[a-0x1000] = v
	}
}

func (p *Ppu) raster() {
    length := len(p.Palettebuffer)
	for i := length - 1; i >= 0; i-- {
		y := int(math.Floor(float64(i / 256)))
		x := i - (y * 256)

        var color int
        color = p.Palettebuffer[i].Color
		p.Framebuffer[(y*256)+x] = color
	}

	p.Output <- p.Framebuffer
}

func (p *Ppu) Step() {
	switch {
	case p.Scanline == 240:
		if p.Cycle == 1 {
			// We're in VBlank
			p.setStatus(StatusVblankStarted)

			// $2000.7 enables/disables NMIs
			if p.NmiOnVblank == 0x1 {
				// Request NMI
				cpu.RequestInterrupt(InterruptNmi)
			}

			if p.DebugMode {
				p.RenderNametables()
				p.Debug <- p.Framebuffer
			} else {
				p.raster()
			}

			p.Cycle++
		}
	case p.Scanline == 261: // End of vblank
		if p.Cycle == 341 {
			p.clearStatus(StatusSprite0Hit)
			p.clearStatus(StatusSpriteOverflow)

			p.Scanline = -1
			p.Cycle = 0
			p.FrameCount++
			return
		}
	case p.Scanline < 240 && p.Scanline > -1:
		if p.Cycle == 341 {
			p.Cycle = 0

			if p.ShowBackground {
				p.renderTileRow()
				p.updateEndScanlineRegisters()
			}

			if p.ShowSprites {
				p.evaluateScanlineSprites()
			}

			p.Scanline++
			return
		}
	case p.Scanline == -1:
		if p.Cycle == 1 {
			// Clear VBlank flag
			p.clearStatus(StatusVblankStarted)
		} else if p.Cycle == 304 {
			// Copy scroll latch into VRAMADDR register
			if p.ShowBackground || p.ShowSprites {
				// p.VramAddress = (p.VramAddress) | (p.VramLatch & 0x41F)
				p.VramAddress = p.VramLatch
			}
		}
	}

	if p.Cycle == 341 {
		p.Cycle = 0
		p.Scanline++
	}

	p.Cycle++
}

func (p *Ppu) updateEndScanlineRegisters() {
	// Scanline has ended
	if p.VramAddress&0x7000 == 0x7000 {
		tmp := p.VramAddress & 0x3E0
		p.VramAddress = p.VramAddress & 0xFFF
		switch tmp {
		case 0x3A0:
			p.VramAddress = p.VramAddress ^ 0xBA0
		case 0x3E0:
			p.VramAddress = p.VramAddress ^ 0x3E0
		default:
			p.VramAddress = p.VramAddress + 0x20
		}
	} else {
		// Increment the fine-Y
		p.VramAddress = p.VramAddress + 0x1000
	}

	if p.ShowBackground || p.ShowSprites {
		p.VramAddress = p.VramAddress & 0xFBE0
		p.VramAddress = p.VramAddress | (p.VramLatch & 0x41F)
	}
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

	p.VramLatch = (p.VramLatch & 0x73FF) | ((int(v) & 0x03) << 10)
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

	Ram.WriteMirroredRam(current, 0x2002)
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

	Ram.WriteMirroredRam(current, 0x2002)
}

// $2002
func (p *Ppu) ReadStatus() (s Word, e error) {
	p.WriteLatch = true
	s = Ram[0x2002]

	// Clear VBlank flag
	p.clearStatus(StatusVblankStarted)
	p.VramLatch = 0

	return
}

// $2003
func (p *Ppu) WriteOamAddress(v Word) {
	p.SpriteRamAddress = int(v)
}

// $2004
func (p *Ppu) WriteOamData(v Word) {
	p.SpriteRam[p.SpriteRamAddress] = v

	p.updateBufferedSpriteMem(p.SpriteRamAddress, v)

	p.SpriteRamAddress++
	p.SpriteRamAddress %= 0x100
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
	fmt.Println("Reading OAM")
	return p.SpriteRam[p.SpriteRamAddress], nil
}

// $2005
func (p *Ppu) WriteScroll(v Word) {
	if p.WriteLatch {
		p.VramLatch = p.VramLatch & 0x7FE0
		p.VramLatch = p.VramLatch | ((int(v) & 0xF8) >> 3)
		p.FineX = v & 0x07
	} else {
		p.VramLatch = p.VramLatch & 0xC1F
		p.VramLatch = p.VramLatch | (((int(v) & 0xF8) << 2) | ((int(v) & 0x07) << 12))
	}

	p.WriteLatch = !p.WriteLatch
}

// $2006
func (p *Ppu) WriteAddress(v Word) {
	if p.WriteLatch {
		p.VramLatch = p.VramLatch & 0xFF
		p.VramLatch = p.VramLatch | ((int(v) & 0x3F) << 8)
	} else {
		p.VramLatch = p.VramLatch & 0x7F00
		p.VramLatch = p.VramLatch | int(v)
		p.VramAddress = p.VramLatch
	}

	p.WriteLatch = !p.WriteLatch
}

// $2007
func (p *Ppu) WriteData(v Word) {
	if p.VramAddress > 0x3000 {
		p.writeMirroredVram(p.VramAddress, v)
	} else if p.VramAddress >= 0x2000 && p.VramAddress < 0x3000 {
		// Nametable mirroring
		p.Nametables.writeNametableData(p.VramAddress, v)
	} else {
		p.Vram[p.VramAddress & 0x3FFF] = v
	}

	p.incrementVramAddress()
}

// $2007
func (p *Ppu) ReadData() (Word, error) {
	// Reads from $2007 are buffered with a
	// 1-byte delay
	tmp := p.VramDataBuffer
	p.VramDataBuffer = p.Vram[p.VramAddress]

	p.incrementVramAddress()

	return tmp, nil
}

func (p *Ppu) incrementVramAddress() {
	switch p.VramAddressInc {
	case 0x01:
		p.VramAddress = p.VramAddress + 0x20
	default:
		p.VramAddress = p.VramAddress + 0x01
	}
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
	switch p.Mirroring {
	case MirroringHorizontal:
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
	case MirroringVertical:
		switch t {
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

	return
}

func (p *Ppu) RenderNametables() {
	/*
		p.renderNametable(0, 0, 0)
			p.renderNametable(2, 0, 240)
			p.renderNametable(3, 256, 240)
	*/
	p.renderNametable(1, 256, 0)

	p.Debug <- p.Framebuffer
}

func (p *Ppu) renderNametable(table, xoff, yoff int) {
	a := p.selectNametable(table)

	x := 0
	y := 0

	// Generates each tile and applies the palette
	for i := a; i < a+0x3C0; i++ {
		attrAddr := 0x23C0 | (i & 0xC00) | int(p.AttributeLocation[i&0x3FF])
		shift := p.AttributeShift[i&0x3FF]
		attr := ((p.Vram[attrAddr] >> shift) & 0x03) << 2

		t := p.bgPatternTableAddress(p.Vram[i])
		tile := p.Vram[t : t+16]

		for c := 0; c < 8; c++ {
			p.decodePatternTile([]Word{tile[c], tile[c+8]},
				x+xoff,
				y+yoff+c,
				p.bgPaletteEntry(attr),
				nil,
				false)
		}

		x += 8

		if x >= 255 {
			x = 0
			y += 8
		}
	}
}

func (p *Ppu) renderTileRow() {
	// Generates each tile, one scanline at a time
    // and applies the palette

	// 32 total for 32 tiles, 8 pixels of each
    
    // Load first two tiles into shift registers at start, then load
    // one per loop and shift the other back out
    // xcoord := p.VramAddress & 0x1F

    fetchTileAttributes := func() (int, Word) {
        attrAddr := 0x23C0 | (p.VramAddress & 0xC00) | int(p.AttributeLocation[p.VramAddress&0x3FF])
        shift := p.AttributeShift[p.VramAddress&0x3FF]
        attr := ((p.Nametables.readNametableData(attrAddr) >> shift) & 0x03) << 2

        index := p.Nametables.readNametableData(p.VramAddress)
        t := p.bgPatternTableAddress(index)

        // Flip bit 10 on wraparound
        if p.VramAddress&0x1F == 0x1F {
            // If rendering is enabled, at the end of a scanline
            // copy bits 10 and 4-0 from VRAM latch into VRAMADDR
            p.VramAddress = p.VramAddress ^ 0x41F
        } else {
            p.VramAddress++
        }

        return t, attr
    }

    t, attr := fetchTileAttributes()
    tile := p.Vram[t : t+16]
    // Move first tile into shift registers
    p.LowBitShift = uint16(tile[p.TilerowCounter])
    p.HighBitShift = uint16(tile[p.TilerowCounter+8])

    t, attrBuf := fetchTileAttributes()
    tile = p.Vram[t : t+16]
    // Get second tile, move the pixels into the right side of
    // shift registers
    // Current tile to render is attrBuf
    p.LowBitShift = (p.LowBitShift << 8) | uint16(tile[p.TilerowCounter])
    p.HighBitShift = (p.HighBitShift << 8) | uint16(tile[p.TilerowCounter+8])

	for x := 0; x < 32; x++ {
        var palette []Word

        var b uint
        for b = 0; b < 8; b++ {
            fbRow := p.Scanline*256 + ((x*8) + int(b))

            // If we're grabbing the pixel from the high
            // part of the shift register, use the buffered
            // palette, not the current one
            if (15 - b - uint(p.FineX)) < 8 {
                palette = p.bgPaletteEntry(attrBuf)
            } else {
                palette = p.bgPaletteEntry(attr)
            }

            pixel := (p.LowBitShift >> (15 - b - uint(p.FineX))) & 0x01
            pixel += ((p.HighBitShift >> (15 - b - uint(p.FineX)) & 0x01) << 1)

            p.Palettebuffer[fbRow] = Pixel{
                PaletteRgb[int(palette[pixel])],
                int(pixel),
            }
        }

		// xcoord = p.VramAddress & 0x1F
        attr = attrBuf

        t, attrBuf = fetchTileAttributes()
        tile = p.Vram[t : t+16]

        // Shift the first tile out, bring the new tile in
        p.LowBitShift = (p.LowBitShift << 8) | uint16(tile[p.TilerowCounter])
        p.HighBitShift = (p.HighBitShift << 8) | uint16(tile[p.TilerowCounter+8])
	}

	p.TilerowCounter++

	if p.TilerowCounter == 8 {
		p.TilerowCounter = 0
	}
}

func (p *Ppu) evaluateScanlineSprites() {
	spriteCount := 0

	for i, y := range p.SpriteData.YCoordinates {
		// if p.Scanline - int(y)+1  >= 0 && p.Scanline - int(y)+1 < 8 {
		if int(y) > (p.Scanline-1)-8 && int(y)+7 < (p.Scanline-1)+8 {
			attrValue := p.Attributes[i] & 0x3
			t := p.SpriteData.Tiles[i]

			c := (p.Scanline - 1) - int(y)
			if p.SpriteSize&0x01 != 0x0 {
				// 8x16 Sprite
				s := p.sprPatternTableAddress(int(t))
				tile := p.Vram[s : s+16]

				p.decodePatternTile([]Word{tile[c], tile[c+8]},
					int(p.XCoordinates[i]),
					int(p.YCoordinates[i])+c+1,
					p.sprPaletteEntry(uint(attrValue)),
					&p.Attributes[i], i == 0)
				s = p.sprPatternTableAddress(int(t+1))
				tile = p.Vram[s : s+16]

				p.decodePatternTile([]Word{tile[c], tile[c+8]},
					int(p.XCoordinates[i]),
					int(p.YCoordinates[i])+c+9,
					p.sprPaletteEntry(uint(attrValue)),
					&p.Attributes[i], i == 0)
			} else {
				// 8x8 Sprite
				s := p.sprPatternTableAddress(int(t))
				tile := p.Vram[s : s+16]

				p.decodePatternTile([]Word{tile[c], tile[c+8]},
					int(p.XCoordinates[i]),
					int(p.YCoordinates[i])+c+1,
					p.sprPaletteEntry(uint(attrValue)),
					&p.Attributes[i], i == 0)
			}

			spriteCount++

			if spriteCount == 8 {
				p.setStatus(StatusSpriteOverflow)
				break
			}
		}
	}
}

func (p *Ppu) decodePatternTile(t []Word, x, y int, pal []Word, attr *Word, spZero bool) {
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
			ycoord = y + int(8-b)
		} else {
			ycoord = y
		}

		fbRow := ycoord*256 + xcoord

		// Store the bit 0/1
		pixel := (t[0] >> b) & 0x01
		pixel += ((t[1] >> b & 0x01) << 1)

		trans := false
		if attr != nil && pixel == 0 {
			trans = true
		}

		// Set the color of the pixel in the buffer
		//
		if fbRow < 0xF000 && !trans {
			if attr != nil {
				priority := (*attr >> 5) & 0x1

				if p.Palettebuffer[fbRow].Value != 0 && priority == 1 {
					// Pixel is already rendered and priority
					// 1 means show behind background
					// continue
				}
			}

			p.Palettebuffer[fbRow] = Pixel{
				PaletteRgb[int(pal[pixel])],
				int(pixel),
			}

			// Since we render sprites first, if we're placing an opaque
			// pixel here and the existing pixel is opaque, we've hit
			// Sprite 0 
			if p.Palettebuffer[fbRow].Value != 0 && spZero {
				p.setStatus(StatusSprite0Hit)
				// p.Framebuffer[fbRow] = 0x00FF00
			}
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
