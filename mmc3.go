package main

import (
	"fmt"
)

const (
	ChrBank2k0000 = 0
	ChrBank2k0800 = 1
	ChrBank1k1000 = 2
	ChrBank1k1400 = 3
	ChrBank1k1800 = 4
	ChrBank1k1C00 = 5
	PrgBank8k8000 = 6
	PrgBank8kA000 = 7

	PrgBankSwapModeLow  = 0
	PrgBankSwapModeHigh = 1

	ChrA12InversionModeLow  = 0
	ChrA12InversionModeHigh = 1

	RegisterBankSelect = iota
	RegisterBankData
	RegisterMirroring
	RegisterPrgRamProtect
	RegisterIrqLatch
	RegisterIrqReload
	RegisterIrqDisable
	RegisterIrqEnable
)

type Mmc3 struct {
	RomBanks  [][]Word
	VromBanks [][]Word

	PrgBankCount int
	ChrRomCount  int
	Battery      bool
	Data         []byte

	BankSelection   int
	PrgBankMode     int
	ChrA12Inversion int
	AddressChanged  bool
	IrqEnabled      bool
	IrqLatchValue   Word
	IrqCounter      Word
	IrqReset        bool
	IrqResetVbl     bool

	RamProtectDest [16]int
}

func NewMmc3(r *Rom) *Mmc3 {
	m := &Mmc3{
		RomBanks:     r.RomBanks,
		VromBanks:    r.VromBanks,
		PrgBankCount: r.PrgBankCount,
		ChrRomCount:  r.ChrRomCount,
		Battery:      r.Battery,
		Data:         r.Data,
	}

	// This just needs to be non-zero and not a 1
	// so that it can be caught when it's changed
	m.PrgBankMode = 10

	m.LoadRom()

	return m
}

func (m *Mmc3) LoadRom() {
	// The PRG banks are 8192 bytes in size, half the size of an
	// iNES PRG bank. If your emulator or copier handles PRG data
	// in 16384 byte chunks, you can think of the lower bit as
	// selecting the first or second half of the bank
	//
	// http://forums.nesdev.com/viewtopic.php?p=38182#p38182

	// Write hardwired PRG banks (0xC000 and 0xE000)
	// Second to last bank
	m.Write8kRamBank(0, 0xC000)
	// Last bank
	m.Write8kRamBank(((len(m.RomBanks)-1)*2)+1, 0xE000)

	// Write swappable PRG banks (0x8000 and 0xA000)
	WriteRamBank(m.RomBanks, 0, 0x8000, Size16k)
}

func (m *Mmc3) BatteryBacked() bool {
	return m.Battery
}

func (m *Mmc3) Write(v Word, a int) {
	switch m.RegisterNumber(a) {
	case RegisterBankSelect:
		m.BankSelect(int(v))
	case RegisterBankData:
		m.BankData(int(v))
	case RegisterMirroring:
		m.SetMirroring(int(v))
	case RegisterPrgRamProtect:
		m.RamProtection(int(v))
	case RegisterIrqLatch:
		m.IrqLatch(int(v))
	case RegisterIrqReload:
		m.IrqReload(int(v))
	case RegisterIrqDisable:
		m.IrqDisable(int(v))
	case RegisterIrqEnable:
		m.IrqEnable(int(v))
	}
}

func (m *Mmc3) RegisterNumber(a int) int {
	switch {
	case a >= 0x8000 && a <= 0x9FFF:
		if a%2 == 0 {
			return RegisterBankSelect
		} else {
			return RegisterBankData
		}
	case a >= 0xA000 && a <= 0xBFFF:
		if a%2 == 0 {
			return RegisterMirroring
		} else {
			return RegisterPrgRamProtect
		}
	case a >= 0xC000 && a <= 0xDFFF:
		if a%2 == 0 {
			return RegisterIrqLatch
		} else {
			return RegisterIrqReload
		}
	}

	if a%2 == 0 {
		return RegisterIrqDisable
	}

	return RegisterIrqEnable
}

func (m *Mmc3) BankSelect(v int) {
	// Next bank register to update
	m.BankSelection = v & 0x7

	address := (v >> 6) & 0x1
	if m.PrgBankMode != address {
		m.AddressChanged = true
	}

	m.PrgBankMode = address
	m.ChrA12Inversion = (v >> 7) & 0x1
}

func (m *Mmc3) BankData(v int) {
	loadHardBanks := func() {
		if m.AddressChanged {
			if m.PrgBankMode == PrgBankSwapModeLow {
				m.Write8kRamBank((len(m.RomBanks)-1)*2, 0xC000)
			} else {
				m.Write8kRamBank((len(m.RomBanks)-1)*2, 0x8000)
			}

			m.AddressChanged = false
		}
	}

	switch m.BankSelection {
	case ChrBank2k0000:
		if m.ChrRomCount == 0 {
			break
		}

		if m.ChrA12Inversion == ChrA12InversionModeLow {
			m.Write1kVramBank(v, 0x0000)
			m.Write1kVramBank(v+1, 0x0400)
		} else {
			m.Write1kVramBank(v, 0x1000)
			m.Write1kVramBank(v+1, 0x1400)
		}
	case ChrBank2k0800:
		if m.ChrRomCount == 0 {
			break
		}

		if m.ChrA12Inversion == ChrA12InversionModeLow {
			m.Write1kVramBank(v, 0x0800)
			m.Write1kVramBank(v+1, 0x0C00)
		} else {
			m.Write1kVramBank(v, 0x1800)
			m.Write1kVramBank(v+1, 0x1C00)
		}
	case ChrBank1k1000:
		if m.ChrRomCount == 0 {
			break
		}

		if m.ChrA12Inversion == ChrA12InversionModeLow {
			m.Write1kVramBank(v, 0x1000)
		} else {
			m.Write1kVramBank(v, 0x0000)
		}
	case ChrBank1k1400:
		if m.ChrRomCount == 0 {
			break
		}

		if m.ChrA12Inversion == ChrA12InversionModeLow {
			m.Write1kVramBank(v, 0x1400)
		} else {
			m.Write1kVramBank(v, 0x0400)
		}
	case ChrBank1k1800:
		if m.ChrRomCount == 0 {
			break
		}

		if m.ChrA12Inversion == ChrA12InversionModeLow {
			m.Write1kVramBank(v, 0x1800)
		} else {
			m.Write1kVramBank(v, 0x0800)
		}
	case ChrBank1k1C00:
		if m.ChrRomCount == 0 {
			break
		}

		if m.ChrA12Inversion == ChrA12InversionModeLow {
			m.Write1kVramBank(v, 0x1C00)
		} else {
			m.Write1kVramBank(v, 0x0C00)
		}
	case PrgBank8k8000:
		loadHardBanks()

		if m.PrgBankMode == PrgBankSwapModeLow {
			m.Write8kRamBank(v, 0x8000)
		} else {
			m.Write8kRamBank(v, 0xC000)
		}
	case PrgBank8kA000:
		m.Write8kRamBank(v, 0xA000)

		loadHardBanks()
	}
}

func (m *Mmc3) SetMirroring(v int) {
	switch v & 0x1 {
	case 0x0:
		ppu.Nametables.SetMirroring(MirroringVertical)
	case 0x1:
		ppu.Nametables.SetMirroring(MirroringHorizontal)
	}
}

func (m *Mmc3) RamProtection(v int) {
	// TODO: WhAT IS THIS I DON'T EVEN
	fmt.Println("RamProtection register")
}

func (m *Mmc3) IrqLatch(v int) {
	// $C000
	m.IrqLatchValue = Word(v)
}

func (m *Mmc3) IrqReload(v int) {
	// $C001
	m.IrqCounter |= 0x80

	if ppu.Scanline < 240 {
		m.IrqReset = true
	} else {
		m.IrqResetVbl = true
		m.IrqReset = false
	}
}

func (m *Mmc3) IrqDisable(v int) {
	// $E000
	m.IrqEnabled = false
	m.IrqCounter = m.IrqLatchValue
}

func (m *Mmc3) IrqEnable(v int) {
	// $E001
	m.IrqEnabled = true
}

func (m *Mmc3) Write8kRamBank(bank, dest int) {
	if m.RamProtectDest[dest>>12] == bank {
		return
	}
	m.RamProtectDest[dest>>12] = bank

	b := (bank >> 1) % len(m.RomBanks)
	offset := (bank % 2) * 0x2000

	WriteOffsetRamBank(m.RomBanks, b, dest, Size8k, offset)
}

func (m *Mmc3) Write1kVramBank(bank, dest int) {
	b := (bank >> 2) % len(m.VromBanks)
	offset := (bank % 4) * 0x400

	WriteOffsetVramBank(m.VromBanks, b, dest, Size1k, offset)
}

func (m *Mmc3) Hook() {
	// A12 Rising Edge
	if (ppu.Scanline > -1 && ppu.Scanline < 240) && (ppu.ShowBackground || ppu.ShowSprites) {
		if m.IrqResetVbl {
			m.IrqCounter = m.IrqLatchValue
			m.IrqResetVbl = false
		}

		if m.IrqReset {
			m.IrqCounter = m.IrqLatchValue
			m.IrqReset = false
		} else if m.IrqCounter > 0 {
			m.IrqCounter--
		}
	}

	if m.IrqCounter == 0 {
		if m.IrqEnabled {
			cpu.RequestInterrupt(InterruptIrq)
		}

		m.IrqReset = true
	}
}
