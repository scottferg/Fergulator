package nes

import (
	"fmt"
)

const (
	RegisterPrgBankSelect = iota
	RegisterChrBank1Select
	RegisterChrBank2Select
	RegisterChrBank3Select
	RegisterChrBank4Select
	RegisterMirroringSelect
)

type Mmc2 struct {
	RomBanks  [][]word
	VromBanks [][]word

	PrgBankCount int
	ChrRomCount  int
	Battery      bool
	Data         []byte

	LatchLow  int
	LatchHigh int
	LatchFE0  int
	LatchFE1  int
	LatchFD0  int
	LatchFD1  int

	PrgUpperHighBank int
	PrgUpperLowBank  int
	PrgLowerHighBank int
	PrgLowerLowBank  int

	ChrHighBank int
	ChrLowBank  int
}

func NewMmc2(r *Nrom) *Mmc2 {
	m := &Mmc2{
		RomBanks:     r.RomBanks,
		VromBanks:    r.VromBanks,
		PrgBankCount: r.PrgBankCount,
		ChrRomCount:  r.ChrRomCount,
		Battery:      r.Battery,
		Data:         r.Data,
	}

	m.LatchLow = 0xFE
	m.LatchHigh = 0xFE
	m.LatchFD0 = 0
	m.LatchFE0 = 4
	m.LatchFD1 = 0
	m.LatchFE1 = 0

	m.Load()

	return m
}

func (m *Mmc2) Load() {
	// 2x the banks since we're storing 8k per bank
	// instead of 16k
	fmt.Printf("  Emulated PRG banks: %d\n", 2*m.PrgBankCount)
	m.RomBanks = make([][]word, 2*m.PrgBankCount)
	for i := 0; i < 2*m.PrgBankCount; i++ {
		// Move 8kb chunk to 8kb bank
		bank := make([]word, 0x2000)
		for x := 0; x < 0x2000; x++ {
			bank[x] = word(m.Data[(0x2000*i)+x])
		}

		m.RomBanks[i] = bank
	}

	// Everything after PRG-ROM
	chrRom := m.Data[0x2000*len(m.RomBanks):]

	// CHR is stored in 4k banks
	m.VromBanks = make([][]word, m.ChrRomCount*2)

	for i := 0; i < cap(m.VromBanks); i++ {
		// Move 16kb chunk to 16kb bank
		m.VromBanks[i] = make([]word, 0x1000)

		// If the game doesn't have CHR banks we
		// just need to allocate VRAM

		for x := 0; x < 0x1000; x++ {
			var val word
			if m.ChrRomCount == 0 {
				val = 0
			} else {
				val = word(chrRom[(0x1000*i)+x])
			}
			m.VromBanks[i][x] = val
		}
	}

	// The PRG banks are 8192 bytes in size, half the size of an
	// iNES PRG bank. If your emulator or copier handles PRG data
	// in 16384 byte chunks, you can think of the lower bit as
	// selecting the first or second half of the bank
	//
	// http://forums.nesdev.com/viewtopic.php?p=38182#p38182

	m.PrgLowerLowBank = 0
	// Write hardwired PRG banks (0xC000 and 0xE000)
	// Second to last bank
	m.PrgUpperHighBank = (((len(m.RomBanks) - 1) * 2) + 1) % len(m.RomBanks)
	m.PrgUpperLowBank = ((len(m.RomBanks) - 1) * 2) % len(m.RomBanks)
	// Last bank

	m.PrgLowerHighBank = (((len(m.RomBanks) - 2) * 2) + 1) % len(m.RomBanks)
}

func (m *Mmc2) Write(v word, a int) {
	switch m.RegisterNumber(a) {
	case RegisterPrgBankSelect:
		m.PrgBankSelect(v)
	case RegisterChrBank1Select:
		m.ChrBankSelect(v, 1)
	case RegisterChrBank2Select:
		m.ChrBankSelect(v, 2)
	case RegisterChrBank3Select:
		m.ChrBankSelect(v, 3)
	case RegisterChrBank4Select:
		m.ChrBankSelect(v, 4)
	case RegisterMirroringSelect:
		m.MirroringSelect(v)
	}
}

func (m *Mmc2) WriteVram(v word, a int) {
	switch {
	case a >= 0x1000:
		m.VromBanks[m.ChrHighBank][a&0xFFF] = v
	default:
		m.VromBanks[m.ChrLowBank][a&0xFFF] = v
	}
}

func (m *Mmc2) ReadVram(a int) word {
	// TODO: Causes some minor glitching in the
	// Punch Out! title screen. This used to happen
	// in ppu.fetchTileAttributes()
	m.LatchTrigger(ppu.bgPatternTableAddress(
		ppu.Nametables.readNametableData(ppu.VramAddress)))

	switch {
	case a >= 0x1000:
		return m.VromBanks[m.ChrHighBank][a&0xFFF]
	default:
		return m.VromBanks[m.ChrLowBank][a&0xFFF]
	}
}

func (m *Mmc2) ReadTile(a int) []word {
	switch {
	case a >= 0x1000:
		return m.VromBanks[m.ChrHighBank][a&0xFFF : a&0xFFF+16]
	default:
		return m.VromBanks[m.ChrLowBank][a&0xFFF : a&0xFFF+16]
	}
}

func (m *Mmc2) Read(a int) word {
	switch {
	case a >= 0xE000:
		return m.RomBanks[m.PrgUpperHighBank][a&0x1FFF]
	case a >= 0xC000:
		return m.RomBanks[m.PrgUpperLowBank][a&0x1FFF]
	case a >= 0xA000:
		return m.RomBanks[m.PrgLowerHighBank][a&0x1FFF]
	case a >= 0x8000:
		return m.RomBanks[m.PrgLowerLowBank][a&0x1FFF]
	}

	return 0
}

func (m *Mmc2) LatchTrigger(a int) {
	a &= 0x3FF0

	switch {
	case a == 0x0FD0 && m.LatchLow != 0xFD:
		m.LatchLow = 0xFD

		m.ChrLowBank = m.LatchFD0 % len(m.VromBanks)
	case a == 0x0FE0 && m.LatchLow != 0xFE:
		m.LatchLow = 0xFE

		m.ChrLowBank = m.LatchFE0 % len(m.VromBanks)
	case a == 0x1FD0 && m.LatchHigh != 0xFD:
		m.LatchHigh = 0xFD

		m.ChrHighBank = m.LatchFD1 % len(m.VromBanks)
	case a == 0x1FE0 && m.LatchHigh != 0xFE:
		m.LatchHigh = 0xFE

		m.ChrHighBank = m.LatchFE1 % len(m.VromBanks)
	}
}

func (m *Mmc2) BatteryBacked() bool {
	return m.Battery
}

func (m *Mmc2) RegisterNumber(a int) int {
	switch a >> 12 {
	case 0xA:
		return RegisterPrgBankSelect
	case 0xB:
		return RegisterChrBank1Select
	case 0xC:
		return RegisterChrBank2Select
	case 0xD:
		return RegisterChrBank3Select
	case 0xE:
		return RegisterChrBank4Select
	case 0xF:
		return RegisterMirroringSelect
	}

	return -1
}

func (m *Mmc2) PrgBankSelect(v word) {
	m.PrgLowerLowBank = int(v&0xF) % len(m.RomBanks)
}

func (m *Mmc2) ChrBankSelect(v word, b int) {
	v &= 0x1F

	switch b {
	case 1:
		m.LatchFD0 = int(v)

		if m.LatchLow == 0xFD {
			m.ChrLowBank = m.LatchFD0 % len(m.VromBanks)
		}
	case 2:
		m.LatchFE0 = int(v)

		if m.LatchLow == 0xFE {
			m.ChrLowBank = m.LatchFE0 % len(m.VromBanks)
		}
	case 3:
		m.LatchFD1 = int(v)

		if m.LatchHigh == 0xFD {
			m.ChrHighBank = m.LatchFD1 % len(m.VromBanks)
		}
	case 4:
		m.LatchFE1 = int(v)

		if m.LatchHigh == 0xFE {
			m.ChrHighBank = m.LatchFE1 % len(m.VromBanks)
		}
	}
}

func (m *Mmc2) MirroringSelect(v word) {
	if v&0x1 == 0x1 {
		ppu.Nametables.SetMirroring(MirroringHorizontal)
	} else {
		ppu.Nametables.SetMirroring(MirroringVertical)
	}
}
