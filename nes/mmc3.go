package nes

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
	RomBanks  []Word
	VromBanks []Word

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

	PrgUpperHighBank int
	PrgUpperLowBank  int
	PrgLowerHighBank int
	PrgLowerLowBank  int

	Chr000Bank  int
	Chr400Bank  int
	Chr800Bank  int
	ChrC00Bank  int
	Chr1000Bank int
	Chr1400Bank int
	Chr1800Bank int
	Chr1C00Bank int

	RamProtectDest [16]int
}

func NewMmc3(r *Nrom) *Mmc3 {
	m := &Mmc3{
		PrgBankCount: r.PrgBankCount,
		ChrRomCount:  r.ChrRomCount,
		Battery:      r.Battery,
		Data:         r.Data,
	}

	// This just needs to be non-zero and not a 1
	// so that it can be caught when it's changed
	m.PrgBankMode = 10

	m.Load()

	return m
}

func (m *Mmc3) Load() {
	// 2x the banks since we're storing 8k per bank
	// instead of 16k
	fmt.Printf("  Emulated PRG banks: %d\n", 2*m.PrgBankCount)
	m.RomBanks = make([]Word, (2*m.PrgBankCount)*0x2000)
	for i := 0; i < 2*m.PrgBankCount; i++ {
		// Move 8kb chunk to 8kb bank
		for x := 0; x < 0x2000; x++ {
			m.RomBanks[(i*0x2000)+x] = Word(m.Data[(0x2000*i)+x])
		}
	}

	// Everything after PRG-ROM
	chrRom := m.Data[0x2000*2*m.PrgBankCount:]

	// CHR is stored in 1k banks
	if m.ChrRomCount > 0 {
		m.VromBanks = make([]Word, (m.ChrRomCount*8)*0x0400)
	} else {
		m.VromBanks = make([]Word, 8*0x0400)
	}

	for i := 0; i < m.ChrRomCount*8; i++ {
		// Move 16kb chunk to 16kb bank
		for x := 0; x < 0x0400; x++ {
			var val Word
			if m.ChrRomCount == 0 {
				val = 0
			} else {
				val = Word(chrRom[(0x0400*i)+x])
			}
			m.VromBanks[(0x0400*i)+x] = val
		}
	}

	// The PRG banks are 8192 bytes in size, half the size of an
	// iNES PRG bank. If your emulator or copier handles PRG data
	// in 16384 byte chunks, you can think of the lower bit as
	// selecting the first or second half of the bank
	//
	// http://forums.nesdev.com/viewtopic.php?p=38182#p38182

	// Write hardwired PRG banks (0xC000 and 0xE000)
	// Second to last bank

	m.PrgUpperHighBank = ((((2 * m.PrgBankCount) - 1) * 2) + 1) >> 1
	m.PrgUpperLowBank = m.PrgUpperHighBank - 1
	// Last bank

	// Write swappable PRG banks (0x8000 and 0xA000)
	m.PrgLowerLowBank = 0
	m.PrgLowerHighBank = 1
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

func (m *Mmc3) WriteVram(v Word, a int) {
	var addr int

	switch {
	case a&0x1C00 == 0x1C00:
		addr = (m.Chr1C00Bank << 0xA) + a&0x3FF
	case a&0x1800 == 0x1800:
		addr = (m.Chr1800Bank << 0xA) + a&0x3FF
	case a&0x1400 == 0x1400:
		addr = (m.Chr1400Bank << 0xA) + a&0x3FF
	case a&0x1000 == 0x1000:
		addr = (m.Chr1000Bank << 0xA) + a&0x3FF
	case a&0x0C00 == 0x0C00:
		addr = (m.ChrC00Bank << 0xA) + a&0x3FF
	case a&0x0800 == 0x0800:
		addr = (m.Chr800Bank << 0xA) + a&0x3FF
	case a&0x0400 == 0x0400:
		addr = (m.Chr400Bank << 0xA) + a&0x3FF
	default:
		addr = (m.Chr000Bank << 0xA) + a&0x3FF
	}

	m.VromBanks[addr] = v
}

func (m *Mmc3) ReadVram(a int) Word {
	var addr int

	switch {
	case a&0x1C00 == 0x1C00:
		addr = (m.Chr1C00Bank << 0xA) + a&0x3FF
	case a&0x1800 == 0x1800:
		addr = (m.Chr1800Bank << 0xA) + a&0x3FF
	case a&0x1400 == 0x1400:
		addr = (m.Chr1400Bank << 0xA) + a&0x3FF
	case a&0x1000 == 0x1000:
		addr = (m.Chr1000Bank << 0xA) + a&0x3FF
	case a&0x0C00 == 0x0C00:
		addr = (m.ChrC00Bank << 0xA) + a&0x3FF
	case a&0x0800 == 0x0800:
		addr = (m.Chr800Bank << 0xA) + a&0x3FF
	case a&0x0400 == 0x0400:
		addr = (m.Chr400Bank << 0xA) + a&0x3FF
	default:
		addr = (m.Chr000Bank << 0xA) + a&0x3FF
	}

	return m.VromBanks[addr]
}

func (m *Mmc3) ReadTile(a int) []Word {
	var addr int

	switch {
	case a&0x1C00 == 0x1C00:
		addr = (m.Chr1C00Bank << 0xA) + a&0x3FF
	case a&0x1800 == 0x1800:
		addr = (m.Chr1800Bank << 0xA) + a&0x3FF
	case a&0x1400 == 0x1400:
		addr = (m.Chr1400Bank << 0xA) + a&0x3FF
	case a&0x1000 == 0x1000:
		addr = (m.Chr1000Bank << 0xA) + a&0x3FF
	case a&0x0C00 == 0x0C00:
		addr = (m.ChrC00Bank << 0xA) + a&0x3FF
	case a&0x0800 == 0x0800:
		addr = (m.Chr800Bank << 0xA) + a&0x3FF
	case a&0x0400 == 0x0400:
		addr = (m.Chr400Bank << 0xA) + a&0x3FF
	default:
		addr = (m.Chr000Bank << 0xA) + a&0x3FF
	}

	return m.VromBanks[addr : addr+16]
}

func (m *Mmc3) Read(a int) Word {
	var addr int

	switch {
	case a&0xE000 == 0xE000:
		addr = (m.PrgUpperHighBank << 0xD) + a&0x1FFF
	case a&0xC000 == 0xC000:
		addr = (m.PrgUpperLowBank << 0xD) + a&0x1FFF
	case a&0xA000 == 0xA000:
		addr = (m.PrgLowerHighBank << 0xD) + a&0x1FFF
	case a&0x8000 == 0x8000:
		addr = (m.PrgLowerLowBank << 0xD) + a&0x1FFF
	}

	return m.RomBanks[addr]
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
			// TODO: +1?
			b := ((2 * m.PrgBankCount) - 1) * 2
			b = ((b >> 1) % (2 * m.PrgBankCount)) - 1

			if m.PrgBankMode == PrgBankSwapModeLow {
				if m.RamProtectDest[0xC000>>12] == b {
					return
				}
				m.RamProtectDest[0xC000>>12] = b
				m.PrgUpperLowBank = b
			} else {
				if m.RamProtectDest[0x8000>>12] == b {
					return
				}
				m.RamProtectDest[0x8000>>12] = b
				m.PrgLowerLowBank = b
			}

			m.AddressChanged = false
		}
	}

	switch m.BankSelection {
	case ChrBank2k0000:
		if m.ChrRomCount == 0 {
			break
		}

		b := v % len(m.VromBanks)
		if m.ChrA12Inversion == ChrA12InversionModeLow {
			m.Chr000Bank = b
			m.Chr400Bank = b + 1
		} else {
			m.Chr1000Bank = b
			m.Chr1400Bank = b + 1
		}
	case ChrBank2k0800:
		if m.ChrRomCount == 0 {
			break
		}

		b := v % len(m.VromBanks)
		if m.ChrA12Inversion == ChrA12InversionModeLow {
			m.Chr800Bank = b
			m.ChrC00Bank = b + 1
		} else {
			m.Chr1800Bank = b
			m.Chr1C00Bank = b + 1
		}
	case ChrBank1k1000:
		if m.ChrRomCount == 0 {
			break
		}

		b := v % len(m.VromBanks)
		if m.ChrA12Inversion == ChrA12InversionModeLow {
			m.Chr1000Bank = b
		} else {
			m.Chr000Bank = b
		}
	case ChrBank1k1400:
		if m.ChrRomCount == 0 {
			break
		}

		b := v % len(m.VromBanks)
		if m.ChrA12Inversion == ChrA12InversionModeLow {
			m.Chr1400Bank = b
		} else {
			m.Chr400Bank = b
		}
	case ChrBank1k1800:
		if m.ChrRomCount == 0 {
			break
		}

		b := v % len(m.VromBanks)
		if m.ChrA12Inversion == ChrA12InversionModeLow {
			m.Chr1800Bank = b
		} else {
			m.Chr800Bank = b
		}
	case ChrBank1k1C00:
		if m.ChrRomCount == 0 {
			break
		}

		b := v % len(m.VromBanks)
		if m.ChrA12Inversion == ChrA12InversionModeLow {
			m.Chr1C00Bank = b
		} else {
			m.ChrC00Bank = b
		}
	case PrgBank8k8000:
		b := v % (2 * m.PrgBankCount)

		if m.PrgBankMode == PrgBankSwapModeLow {
			if m.RamProtectDest[0x8000>>12] == b {
				return
			}
			m.RamProtectDest[0x8000>>12] = b

			m.PrgLowerLowBank = b
		} else {
			if m.RamProtectDest[0xC000>>12] == b {
				return
			}
			m.RamProtectDest[0xC000>>12] = b
			m.PrgUpperLowBank = b
		}

		loadHardBanks()
	case PrgBank8kA000:
		b := v % (2 * m.PrgBankCount)

		if m.RamProtectDest[0xA000>>12] == b {
			return
		}
		m.RamProtectDest[0xA000>>12] = b
		m.PrgLowerHighBank = b

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
