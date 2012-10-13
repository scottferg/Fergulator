package main

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
	RomBanks  [][]Word
	VromBanks [][]Word

	PrgBankCount int
	ChrRomCount  int
	Battery      bool
	Data         []byte

	Latch0     int
	Latch1     int
	Latch0High int
	Latch1High int
	Latch0Low  int
	Latch1Low  int
}

func NewMmc2(r *Rom) *Mmc2 {
	m := &Mmc2{
		RomBanks:     r.RomBanks,
		VromBanks:    r.VromBanks,
		PrgBankCount: r.PrgBankCount,
		ChrRomCount:  r.ChrRomCount,
		Battery:      r.Battery,
		Data:         r.Data,
	}

	m.Latch0 = 0xFE
	m.Latch1 = 0xFE
	m.Latch0Low = 0
	m.Latch1Low = 4
	m.Latch0High = 0
	m.Latch1High = 0

	m.LoadRom()

	return m
}

func (m *Mmc2) LoadRom() {
	// The PRG banks are 8192 bytes in size, half the size of an 
	// iNES PRG bank. If your emulator or copier handles PRG data 
	// in 16384 byte chunks, you can think of the lower bit as 
	// selecting the first or second half of the bank
	//
	// http://forums.nesdev.com/viewtopic.php?p=38182#p38182

	// Write swappable PRG banks (0x8000 and 0xA000)
	m.Write8kRamBank(0, 0x8000)

	// Write hardwired PRG banks (0xC000 and 0xE000) 
	m.Write8kRamBank((len(m.RomBanks)*2)-2, 0xA000)
	m.Write8kRamBank((len(m.RomBanks)*2)-1, 0xC000)
	m.Write8kRamBank((len(m.RomBanks) * 2), 0xE000)

	WriteVramBank(m.VromBanks, 4, 0x0000, Size4k)
	WriteVramBank(m.VromBanks, 0, 0x1000, Size4k)
}

func (m *Mmc2) Write(v Word, a int) {
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

func (m *Mmc2) Hook(a int) {
	a &= 0x1FF0

	switch {
	case a == 0xFD0 && m.Latch0 != 0xFD:
        fmt.Printf("Latch A: 0x%X\n", a)
		m.Latch0 = 0xFD
		WriteVramBank(m.VromBanks, m.Latch0Low, 0x0000, Size4k)
	case a == 0xFE0 && m.Latch0 != 0xFE:
        fmt.Printf("Latch A: 0x%X\n", a)
		m.Latch0 = 0xFE
		WriteVramBank(m.VromBanks, m.Latch1Low, 0x0000, Size4k)
	case a == 0x1F00 && m.Latch1 != 0xFD:
        fmt.Printf("Latch A: 0x%X\n", a)
		m.Latch1 = 0xFD
		WriteVramBank(m.VromBanks, m.Latch0High, 0x1000, Size4k)
	case a == 0x1FE0 && m.Latch1 != 0xFE:
        fmt.Printf("Latch A: 0x%X\n", a)
		m.Latch1 = 0xFE
		WriteVramBank(m.VromBanks, m.Latch1High, 0x1000, Size4k)
	}
}

func (m *Mmc2) BatteryBacked() bool {
	return m.Battery
}

func (m *Mmc2) RegisterNumber(a int) int {
	switch {
	case a >= 0xA000 && a <= 0xAFFF:
		return RegisterPrgBankSelect
	case a >= 0xB000 && a <= 0xBFFF:
		return RegisterChrBank1Select
	case a >= 0xC000 && a <= 0xCFFF:
		return RegisterChrBank2Select
	case a >= 0xD000 && a <= 0xDFFF:
		return RegisterChrBank3Select
	case a >= 0xE000 && a <= 0xEFFF:
		return RegisterChrBank4Select
	case a >= 0xF000 && a <= 0xFFFF:
		return RegisterMirroringSelect
	}

	return -1
}

func (m *Mmc2) PrgBankSelect(v Word) {
	bank := int(v & 0xF)
	WriteRamBank(m.RomBanks, bank, 0x8000, Size8k)
}

func (m *Mmc2) ChrBankSelect(v Word, b int) {
	switch b {
	case 1:
		m.Latch0Low = int(v)

		if m.Latch0 == 0xFD {
			WriteVramBank(m.VromBanks, m.Latch0Low, 0x0000, Size4k)
		}
	case 2:
		m.Latch1Low = int(v)

		if m.Latch0 == 0xFE {
			WriteVramBank(m.VromBanks, m.Latch1Low, 0x0000, Size4k)
		}
	case 3:
		m.Latch0High = int(v)

		if m.Latch1 == 0xFD {
			WriteVramBank(m.VromBanks, m.Latch0High, 0x1000, Size4k)
		}
	case 4:
		m.Latch1High = int(v)

		if m.Latch1 == 0xFE {
			WriteVramBank(m.VromBanks, m.Latch1High, 0x1000, Size4k)
		}
	}
}

func (m *Mmc2) MirroringSelect(v Word) {
	if v&0x1 == 0x1 {
		ppu.Nametables.SetMirroring(MirroringHorizontal)
	} else {
		ppu.Nametables.SetMirroring(MirroringVertical)
	}
}

func (m *Mmc2) Write8kRamBank(bank, dest int) {
	b := (bank >> 1) % len(m.RomBanks)
	offset := (bank % 2) * 0x2000

	WriteOffsetRamBank(m.RomBanks, b, dest, Size8k, offset)
}
