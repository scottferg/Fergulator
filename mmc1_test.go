package main

import (
	"testing"
)

var (
	m *Mmc1
)

func verifyMirroredValue(a int, v Word, test *testing.T) {
	if ppu.Nametables.readNametableData(a) != v {
		test.Errorf("0x%X was 0x%X, expected 0x%X\n", a, ppu.Vram[0x2000], v)
	}
}

func TestVerticalToHorizontal(test *testing.T) {
	rom = &Mmc1{
		RomBanks:     make([][]Word, 16),
		VromBanks:    make([][]Word, 16),
		PrgBankCount: 8,
		ChrRomCount:  8,
		Battery:      false,
		Data:         make([]byte, 32),
		PrgSwapBank:  BankLower,
	}

	ppu.Init()

	ppu.Nametables.SetMirroring(MirroringHorizontal)

	if ppu.Nametables.Mirroring != MirroringHorizontal {
		test.Errorf("Mirroring was not horizontal")
	}

	// Setup Vertical mirroring
	Ram.Write(0x8000, 0x0)
	Ram.Write(0x8000, 0x1)
	Ram.Write(0x8000, 0x0)
	Ram.Write(0x8000, 0x0)
	Ram.Write(0x8000, 0x0)

	if ppu.Nametables.Mirroring != MirroringVertical {
		test.Errorf("Mirroring was not vertical")
	}

	ppu.VramAddress = 0x2000
	ppu.WriteData(0x11)
	ppu.VramAddress = 0x2110
	ppu.WriteData(0x22)
	ppu.VramAddress = 0x2220
	ppu.WriteData(0x33)
	ppu.VramAddress = 0x2330
	ppu.WriteData(0x44)
	ppu.VramAddress = 0x2338
	ppu.WriteData(0x55)

	verifyMirroredValue(0x2000, 0x11, test)
	verifyMirroredValue(0x2800, 0x11, test)

	verifyMirroredValue(0x2110, 0x22, test)
	verifyMirroredValue(0x2910, 0x22, test)

	verifyMirroredValue(0x2220, 0x33, test)
	verifyMirroredValue(0x2A20, 0x33, test)

	verifyMirroredValue(0x2330, 0x44, test)
	verifyMirroredValue(0x2B30, 0x44, test)

	verifyMirroredValue(0x2338, 0x55, test)
	verifyMirroredValue(0x2B38, 0x55, test)

	ppu.VramAddress = 0x2400
	ppu.WriteData(0x11)
	ppu.VramAddress = 0x2510
	ppu.WriteData(0x22)
	ppu.VramAddress = 0x2620
	ppu.WriteData(0x33)
	ppu.VramAddress = 0x2730
	ppu.WriteData(0x44)
	ppu.VramAddress = 0x2738
	ppu.WriteData(0x55)

	verifyMirroredValue(0x2400, 0x11, test)
	verifyMirroredValue(0x2C00, 0x11, test)

	verifyMirroredValue(0x2510, 0x22, test)
	verifyMirroredValue(0x2D10, 0x22, test)

	verifyMirroredValue(0x2620, 0x33, test)
	verifyMirroredValue(0x2E20, 0x33, test)

	verifyMirroredValue(0x2730, 0x44, test)
	verifyMirroredValue(0x2F30, 0x44, test)

	verifyMirroredValue(0x2738, 0x55, test)
	verifyMirroredValue(0x2F38, 0x55, test)
}

func TestHorizontalToVertical(test *testing.T) {
	rom = &Mmc1{
		RomBanks:     make([][]Word, 16),
		VromBanks:    make([][]Word, 16),
		PrgBankCount: 8,
		ChrRomCount:  8,
		Battery:      false,
		Data:         make([]byte, 32),
		PrgSwapBank:  BankLower,
	}

	ppu.Init()

	ppu.Nametables.SetMirroring(MirroringVertical)

	if ppu.Nametables.Mirroring != MirroringVertical {
		test.Errorf("Mirroring was not vertical")
	}

	// Setup Vertical mirroring
	Ram.Write(0x8000, 0x1)
	Ram.Write(0x8000, 0x1)
	Ram.Write(0x8000, 0x0)
	Ram.Write(0x8000, 0x0)
	Ram.Write(0x8000, 0x0)

	if ppu.Nametables.Mirroring != MirroringHorizontal {
		test.Errorf("Mirroring was not horizontal")
	}

	ppu.VramAddress = 0x2000
	ppu.WriteData(0x11)
	ppu.VramAddress = 0x2110
	ppu.WriteData(0x22)
	ppu.VramAddress = 0x2220
	ppu.WriteData(0x33)
	ppu.VramAddress = 0x2330
	ppu.WriteData(0x44)
	ppu.VramAddress = 0x2338
	ppu.WriteData(0x55)

	verifyMirroredValue(0x2000, 0x11, test)
	verifyMirroredValue(0x2400, 0x11, test)

	verifyMirroredValue(0x2110, 0x22, test)
	verifyMirroredValue(0x2510, 0x22, test)

	verifyMirroredValue(0x2220, 0x33, test)
	verifyMirroredValue(0x2620, 0x33, test)

	verifyMirroredValue(0x2330, 0x44, test)
	verifyMirroredValue(0x2730, 0x44, test)

	verifyMirroredValue(0x2338, 0x55, test)
	verifyMirroredValue(0x2738, 0x55, test)

	ppu.VramAddress = 0x2800
	ppu.WriteData(0x11)
	ppu.VramAddress = 0x2910
	ppu.WriteData(0x22)
	ppu.VramAddress = 0x2A20
	ppu.WriteData(0x33)
	ppu.VramAddress = 0x2B30
	ppu.WriteData(0x44)
	ppu.VramAddress = 0x2B38
	ppu.WriteData(0x55)

	verifyMirroredValue(0x2800, 0x11, test)
	verifyMirroredValue(0x2C00, 0x11, test)

	verifyMirroredValue(0x2910, 0x22, test)
	verifyMirroredValue(0x2D10, 0x22, test)

	verifyMirroredValue(0x2A20, 0x33, test)
	verifyMirroredValue(0x2E20, 0x33, test)

	verifyMirroredValue(0x2B30, 0x44, test)
	verifyMirroredValue(0x2F30, 0x44, test)

	verifyMirroredValue(0x2B38, 0x55, test)
	verifyMirroredValue(0x2F38, 0x55, test)
}
