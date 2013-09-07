package main

import (
	"testing"
)

var (
	p *Ppu
)

func verifyValue(a int, v Word, test *testing.T) {
	if p.Nametables.readNametableData(a) != v {
		test.Errorf("0x%X was 0x%X, expected 0x%X\n", a, p.Vram[0x2000], v)
	}
}

func TestVerticalNametableMirroring(test *testing.T) {
	p = new(Ppu)
	p.Init()

	p.Nametables.SetMirroring(MirroringVertical)

	p.VramAddress = 0x2000
	p.WriteData(0x11)
	p.VramAddress = 0x2110
	p.WriteData(0x22)
	p.VramAddress = 0x2220
	p.WriteData(0x33)
	p.VramAddress = 0x2330
	p.WriteData(0x44)
	p.VramAddress = 0x2338
	p.WriteData(0x55)

	verifyValue(0x2000, 0x11, test)
	verifyValue(0x2800, 0x11, test)

	verifyValue(0x2110, 0x22, test)
	verifyValue(0x2910, 0x22, test)

	verifyValue(0x2220, 0x33, test)
	verifyValue(0x2A20, 0x33, test)

	verifyValue(0x2330, 0x44, test)
	verifyValue(0x2B30, 0x44, test)

	verifyValue(0x2338, 0x55, test)
	verifyValue(0x2B38, 0x55, test)

	p.VramAddress = 0x2400
	p.WriteData(0x11)
	p.VramAddress = 0x2510
	p.WriteData(0x22)
	p.VramAddress = 0x2620
	p.WriteData(0x33)
	p.VramAddress = 0x2730
	p.WriteData(0x44)
	p.VramAddress = 0x2738
	p.WriteData(0x55)

	verifyValue(0x2400, 0x11, test)
	verifyValue(0x2C00, 0x11, test)

	verifyValue(0x2510, 0x22, test)
	verifyValue(0x2D10, 0x22, test)

	verifyValue(0x2620, 0x33, test)
	verifyValue(0x2E20, 0x33, test)

	verifyValue(0x2730, 0x44, test)
	verifyValue(0x2F30, 0x44, test)

	verifyValue(0x2738, 0x55, test)
	verifyValue(0x2F38, 0x55, test)
}

func TestHorizontalNametableMirroring(test *testing.T) {
	p = new(Ppu)
	p.Init()

	p.Nametables.SetMirroring(MirroringHorizontal)

	p.VramAddress = 0x2000
	p.WriteData(0x11)
	p.VramAddress = 0x2110
	p.WriteData(0x22)
	p.VramAddress = 0x2220
	p.WriteData(0x33)
	p.VramAddress = 0x2330
	p.WriteData(0x44)
	p.VramAddress = 0x2338
	p.WriteData(0x55)

	verifyValue(0x2000, 0x11, test)
	verifyValue(0x2400, 0x11, test)

	verifyValue(0x2110, 0x22, test)
	verifyValue(0x2510, 0x22, test)

	verifyValue(0x2220, 0x33, test)
	verifyValue(0x2620, 0x33, test)

	verifyValue(0x2330, 0x44, test)
	verifyValue(0x2730, 0x44, test)

	verifyValue(0x2338, 0x55, test)
	verifyValue(0x2738, 0x55, test)

	p.VramAddress = 0x2800
	p.WriteData(0x11)
	p.VramAddress = 0x2910
	p.WriteData(0x22)
	p.VramAddress = 0x2A20
	p.WriteData(0x33)
	p.VramAddress = 0x2B30
	p.WriteData(0x44)
	p.VramAddress = 0x2B38
	p.WriteData(0x55)

	verifyValue(0x2800, 0x11, test)
	verifyValue(0x2C00, 0x11, test)

	verifyValue(0x2910, 0x22, test)
	verifyValue(0x2D10, 0x22, test)

	verifyValue(0x2A20, 0x33, test)
	verifyValue(0x2E20, 0x33, test)

	verifyValue(0x2B30, 0x44, test)
	verifyValue(0x2F30, 0x44, test)

	verifyValue(0x2B38, 0x55, test)
	verifyValue(0x2F38, 0x55, test)
}
