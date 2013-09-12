package main

const (
	StatusSpriteOverflow = iota
	StatusSprite0Hit
	StatusVblankStarted
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
	Color  uint32
	Value  int
	Pindex int
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
	A12High           bool

	Palettebuffer []Pixel
	Framebuffer   []uint32

	Output      chan []uint32
	Cycle       int
	Scanline    int
	Timestamp   int
	VblankTime  int
	FrameCount  int
	FrameCycles int

	SuppressNmi        bool
	SuppressVbl        bool
	OverscanEnabled    bool
	SpriteLimitEnabled bool
}

func (p *Ppu) Init() chan []uint32 {
	p.WriteLatch = true
	p.Output = make(chan []uint32)

	p.OverscanEnabled = true
	p.SpriteLimitEnabled = true
	p.Cycle = 0
	p.Scanline = 241
	p.FrameCount = 0

	p.VblankTime = 20 * 341 * 5 // NTSC

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
	p.Framebuffer = make([]uint32, 0xEFE0)

	return p.Output
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
		if a&0xF == 0 {
			a = 0
		}
		p.PaletteRam[a&0x1F] = v
	} else {
		p.Nametables.writeNametableData(a-0x1000, v)
	}
}

func (p *Ppu) raster() {
	length := len(p.Palettebuffer)
	for i := length - 1; i >= 0; i-- {
		y := i / 256
		x := i - (y * 256)

		var color uint32
		color = p.Palettebuffer[i].Color

		width := 256

		if p.OverscanEnabled {
			if y < 8 || y > 231 || x < 8 || x > 247 {
				continue
			} else {
				y -= 8
				x -= 8
			}

			width = 240

			if len(p.Framebuffer) == 0xF000 {
				p.Framebuffer = make([]uint32, 0xEFE0)
			}
		} else {
			if len(p.Framebuffer) == 0xEFE0 {
				p.Framebuffer = make([]uint32, 0xF000)
			}
		}

		p.Framebuffer[(y*width)+x] = color << 8
		p.Palettebuffer[i].Value = 0
		p.Palettebuffer[i].Pindex = -1
	}

	p.Output <- p.Framebuffer
}

func (p *Ppu) Step() {
	switch {
	case p.Scanline == 240:
		switch p.Cycle {
		case 1:
			if !p.SuppressVbl {
				// We're in VBlank
				p.setStatus(StatusVblankStarted)
			}

			// $2000.7 enables/disables NMIs
			if p.NmiOnVblank == 0x1 && !p.SuppressNmi {
				// Request NMI
				cpu.RequestInterrupt(InterruptNmi)
			}
			p.raster()
		}
	case p.Scanline == 260: // End of vblank
		switch p.Cycle {
		case 1:
			// Clear VBlank flag
			p.clearStatus(StatusVblankStarted)
		case 341:
			p.Scanline = -1
			p.Cycle = 1
			p.FrameCount++
			return
		}
	case p.Scanline < 240 && p.Scanline > -1:
		switch p.Cycle {
		case 254:
			if p.ShowBackground {
				p.renderTileRow()
			}

			if p.ShowSprites {
				p.evaluateScanlineSprites(p.Scanline)
			}
		case 256:
			if p.ShowBackground {
				p.updateEndScanlineRegisters()
			}
		case 260:
			if p.SpritePatternAddress == 0x1 && p.BackgroundPatternAddress == 0x0 {
				if v, ok := rom.(*Mmc3); ok {
					v.Hook()
				}
			}
		}
	case p.Scanline == -1:
		switch p.Cycle {
		case 1:
			p.clearStatus(StatusSprite0Hit)
			p.clearStatus(StatusSpriteOverflow)
		case 304:
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

func (p *Ppu) renderingEnabled() bool {
	return p.ShowBackground && p.ShowSprites
}

func (p *Ppu) updateEndScanlineRegisters() {

	/*******************************************************
	  TODO: Some documentatino implies that the X increment
	  should occur 34 times per scanline. These may not be
	  necessary.
	 ***********************************************************/

	// Flip bit 10 on wraparound
	if p.VramAddress&0x1F == 0x1F {
		// If rendering is enabled, at the end of a scanline
		// copy bits 10 and 4-0 from VRAM latch into VRAMADDR
		p.VramAddress ^= 0x41F
	} else {
		p.VramAddress++
	}

	// Flip bit 10 on wraparound
	if p.VramAddress&0x1F == 0x1F {
		// If rendering is enabled, at the end of a scanline
		// copy bits 10 and 4-0 from VRAM latch into VRAMADDR
		p.VramAddress ^= 0x41F
	} else {
		p.VramAddress++
	}

	if p.ShowBackground || p.ShowSprites {
		// Scanline has ended
		if p.VramAddress&0x7000 == 0x7000 {
			tmp := p.VramAddress & 0x3E0
			p.VramAddress &= 0xFFF

			switch tmp {
			case 0x3A0:
				p.VramAddress ^= 0xBA0
			case 0x3E0:
				p.VramAddress ^= 0x3E0
			default:
				p.VramAddress += 0x20
			}

		} else {
			// Increment the fine-Y
			p.VramAddress += 0x1000
		}

		p.VramAddress = (p.VramAddress & 0x7BE0) | (p.VramLatch & 0x41F)
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

	p.VramLatch = (p.VramLatch & 0xF3FF) | (int(p.BaseNametableAddress) << 10)
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

	if p.Cycle == 1 && p.Scanline == 240 {
		s &= 0x7F
		p.SuppressNmi = true
		p.SuppressVbl = true
	} else {
		p.SuppressNmi = false
		p.SuppressVbl = false
		// Clear VBlank flag
		p.clearStatus(StatusVblankStarted)
	}

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
	for i := 0; i < 0x100; i++ {
		d, _ := Ram.Read(addr + i)
		p.SpriteRam[i] = d
		p.updateBufferedSpriteMem(i, d)
	}
}

func (p *Ppu) updateBufferedSpriteMem(a int, v Word) {
	i := a / 4

	switch a % 4 {
	case 0x0:
		p.SpriteData.YCoordinates[i] = v
	case 0x1:
		p.SpriteData.Tiles[i] = v
	case 0x2:
		// Attribute
		p.SpriteData.Attributes[i] = v
	case 0x3:
		p.SpriteData.XCoordinates[i] = v
	}
}

// $2004
func (p *Ppu) ReadOamData() (Word, error) {
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
		p.Vram[p.VramAddress&0x3FFF] = v
		// MMC2 latch trigger
		t := p.bgPatternTableAddress(p.Nametables.readNametableData(p.VramAddress))
		triggerMapperLatch(t)
	}

	p.incrementVramAddress()
}

func triggerMapperLatch(i int) {
	if v, ok := rom.(*Mmc2); ok {
		v.LatchTrigger(i)
	}
}

// $2007
func (p *Ppu) ReadData() (r Word, err error) {
	// Reads from $2007 are buffered with a
	// 1-byte delay
	if p.VramAddress >= 0x2000 && p.VramAddress < 0x3000 {
		r = p.VramDataBuffer
		p.VramDataBuffer = p.Nametables.readNametableData(p.VramAddress)
	} else if p.VramAddress < 0x3F00 {
		r = p.VramDataBuffer
		p.VramDataBuffer = p.Vram[p.VramAddress]

		if p.VramAddress < 0x2000 {
			// MMC2 latch trigger
			t := p.bgPatternTableAddress(p.Nametables.readNametableData(p.VramAddress))
			triggerMapperLatch(t)
		}
	} else {
		bufferAddress := p.VramAddress - 0x1000
		switch {
		case bufferAddress >= 0x2000 && bufferAddress < 0x3000:
			p.VramDataBuffer = p.Nametables.readNametableData(bufferAddress)
		default:
			p.VramDataBuffer = p.Vram[bufferAddress]
		}

		a := p.VramAddress
		if a&0xF == 0 {
			a = 0
		}

		r = p.PaletteRam[a&0x1F]

		if p.VramAddress < 0x2000 {
			// MMC2 latch trigger
			t := p.bgPatternTableAddress(p.Nametables.readNametableData(p.VramAddress))
			triggerMapperLatch(t)
		}
	}

	p.incrementVramAddress()

	return
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
		if i&0x01 != 0 {
			return 0x1000 | ((int(i) >> 1) * 0x20)
		} else {
			return ((int(i) >> 1) * 0x20)
		}

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

	return (int(i) << 4) | (p.VramAddress >> 12) | a
}

func (p *Ppu) renderTileRow() {
	// Generates each tile, one scanline at a time
	// and applies the palette

	// Load first two tiles into shift registers at start, then load
	// one per loop and shift the other back out
	fetchTileAttributes := func() (uint16, uint16, Word) {
		attrAddr := 0x23C0 | (p.VramAddress & 0xC00) | int(p.AttributeLocation[p.VramAddress&0x3FF])
		shift := p.AttributeShift[p.VramAddress&0x3FF]
		attr := ((p.Nametables.readNametableData(attrAddr) >> shift) & 0x03) << 2

		index := p.Nametables.readNametableData(p.VramAddress)
		t := p.bgPatternTableAddress(index)

		// Flip bit 10 on wraparound
		if p.VramAddress&0x1F == 0x1F {
			// If rendering is enabled, at the end of a scanline
			// copy bits 10 and 4-0 from VRAM latch into VRAMADDR
			p.VramAddress ^= 0x41F
		} else {
			p.VramAddress++
		}

		// MMC2 latch trigger
		triggerMapperLatch(p.bgPatternTableAddress(index))

		return uint16(p.Vram[t]), uint16(p.Vram[t+8]), attr
	}

	// Move first tile into shift registers
	low, high, attr := fetchTileAttributes()
	p.LowBitShift, p.HighBitShift = low, high

	low, high, attrBuf := fetchTileAttributes()
	// Get second tile, move the pixels into the right side of
	// shift registers
	// Current tile to render is attrBuf
	p.LowBitShift = (p.LowBitShift << 8) | low
	p.HighBitShift = (p.HighBitShift << 8) | high

	for x := 0; x < 32; x++ {
		var palette int

		var b uint
		for b = 0; b < 8; b++ {
			fbRow := p.Scanline*256 + ((x * 8) + int(b))

			pixel := (p.LowBitShift >> (15 - b - uint(p.FineX))) & 0x01
			pixel += ((p.HighBitShift >> (15 - b - uint(p.FineX)) & 0x01) << 1)

			// If we're grabbing the pixel from the high
			// part of the shift register, use the buffered
			// palette, not the current one
			if (15 - b - uint(p.FineX)) < 8 {
				palette = p.bgPaletteEntry(attrBuf, pixel)
			} else {
				palette = p.bgPaletteEntry(attr, pixel)
			}

			if p.Palettebuffer[fbRow].Value != 0 {
				// Pixel is already rendered and priority
				// 1 means show behind background
				continue
			}

			p.Palettebuffer[fbRow] = Pixel{
				PaletteRgb[palette%64],
				int(pixel),
				-1,
			}
		}

		// xcoord = p.VramAddress & 0x1F
		attr = attrBuf

		// Shift the first tile out, bring the new tile in
		low, high, attrBuf = fetchTileAttributes()

		p.LowBitShift = (p.LowBitShift << 8) | low
		p.HighBitShift = (p.HighBitShift << 8) | high
	}
}

func (p *Ppu) evaluateScanlineSprites(line int) {
	spriteCount := 0

	for i, y := range p.SpriteData.YCoordinates {
		spriteHeight := 8
		if p.SpriteSize&0x1 == 0x1 {
			spriteHeight = 16
		}

		if int(y) > (line-1)-spriteHeight && int(y)+(spriteHeight-1) < (line-1)+spriteHeight {
			attrValue := p.Attributes[i] & 0x3
			t := p.SpriteData.Tiles[i]

			c := (line - 1) - int(y)

			// TODO: Hack to fix random sprite appearing in upper
			// left. It should be cropped by overscan.
			if p.XCoordinates[i] == 0 && p.YCoordinates[i] == 0 {
				continue
			}

			var ycoord int

			yflip := (p.Attributes[i]>>7)&0x1 == 0x1
			if yflip {
				ycoord = int(p.YCoordinates[i]) + ((spriteHeight - 1) - c)
			} else {
				ycoord = int(p.YCoordinates[i]) + c + 1
			}

			if p.SpriteSize&0x01 != 0x0 {
				// 8x16 Sprite
				s := p.sprPatternTableAddress(int(t))
				var tile []Word

				top := p.Vram[s : s+16]
				bottom := p.Vram[s+16 : s+32]

				if c > 7 && yflip {
					tile = top
					ycoord += 8
				} else if c < 8 && yflip {
					tile = bottom
					ycoord -= 8
				} else if c > 7 {
					tile = bottom
				} else {
					tile = top
				}

				sprite0 := i == 0

				p.decodePatternTile([]Word{tile[c%8], tile[(c%8)+8]},
					int(p.XCoordinates[i]),
					ycoord,
					uint(attrValue),
					&p.Attributes[i], sprite0, i)
			} else {
				// 8x8 Sprite
				s := p.sprPatternTableAddress(int(t))
				tile := p.Vram[s : s+16]

				p.decodePatternTile([]Word{tile[c], tile[c+8]},
					int(p.XCoordinates[i]),
					ycoord,
					uint(attrValue),
					&p.Attributes[i], i == 0, i)
			}

			spriteCount++

			if spriteCount == 9 {
				if p.SpriteLimitEnabled {
					p.setStatus(StatusSpriteOverflow)
					break
				}
			}
		}
	}
}

func (p *Ppu) decodePatternTile(t []Word, x, y int, palIndex uint, attr *Word, spZero bool, index int) {
	var b uint
	for b = 0; b < 8; b++ {
		var xcoord int
		if (*attr>>6)&0x1 != 0 {
			xcoord = x + int(b)
		} else {
			xcoord = x + int(7-b)
		}

		// Don't wrap around if we're past the edge of the
		// screen
		if xcoord > 255 {
			continue
		}

		fbRow := y*256 + xcoord

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
			priority := (*attr >> 5) & 0x1

			hit := (Ram[0x2002]&0x40 == 0x40)
			if p.Palettebuffer[fbRow].Value != 0 && spZero && !hit {
				// Since we render background first, if we're placing an opaque
				// pixel here and the existing pixel is opaque, we've hit
				// Sprite 0
				p.setStatus(StatusSprite0Hit)
			}

			if p.Palettebuffer[fbRow].Pindex > -1 && p.Palettebuffer[fbRow].Pindex < index {
				// Pixel with a higher sprite priority (lower index)
				// is already here, so don't render this pixel
				continue
			} else if p.Palettebuffer[fbRow].Value != 0 && priority == 1 {
				// Pixel is already rendered and priority
				// 1 means show behind background
				// unless background pixel is not transparent
				continue
			}

			pal := p.PaletteRam[0x10+(palIndex*0x4)+uint(pixel)]

			p.Palettebuffer[fbRow] = Pixel{
				PaletteRgb[int(pal)%64],
				int(pixel),
				index,
			}
		}
	}
}

func (p *Ppu) bgPaletteEntry(a Word, pix uint16) int {
	if pix == 0x0 {
		return int(p.PaletteRam[0x00])
	}

	return int(p.PaletteRam[uint16(a)+pix])
}
