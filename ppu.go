package main

import (
	"fmt"
)

type Flags struct {
	BaseNametableAddress     Word
	VramAddressInc           Word
	SpritePatternAddress     Word
	BackgroundPatternAddress Word
	SpriteSize               Word
	MasterSlaveSel           Word
	GenerateNmi              Word
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
	AddressRegister AddressRegisters
	AddressCounter  AddressCounters
	ScanCycleCount  int
	OpCycleCount    int
	Vram            [0xFFFF]Word
	SpriteRam       [0x100]Word
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

func (p *Ppu) WriteControl(v Word) {
	p.Control = v

	// Control flags
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
	p.GenerateNmi = (v >> 7) & 0x01
}

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

func (p *Ppu) ReadStatus() {
    p.FirstWrite = true
}

func (p *Ppu) WriteOamAddress(v Word) {
	p.OamAddress = int(v)
	fmt.Printf("OAM Address: 0x%X\n", v)
}

func (p *Ppu) WriteOamData(v Word) {
	p.OamData = v
	fmt.Printf("OAM Data: 0x%X\n", v)
}

func (p *Ppu) ReadOamData() Word {
	fmt.Printf("Reading OAM")
	return 0
}

func (p *Ppu) WriteScroll(v Word) {
	p.Scroll = v
}

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

func (p *Ppu) WriteData(v Word) {
    p.Address = p.AddressCounter.Address()

    if p.Address > 0x2000 {
        fmt.Println("Writing to mirrored VRAM")
    } else {
        // fmt.Printf("Writing to VRAM[0x%X]: 0x%X\n", p.Address, v)
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

func (p *Ppu) ReadData() Word {
	fmt.Printf("Reading Data")
	return 0
}

func (p *Ppu) Init() {
	Ram.Write(0x2002, 0x80)
    p.FirstWrite = true
}
