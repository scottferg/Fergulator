package nes

import (
	"fmt"
)

type Mmc5 struct {
	RomBanks  [][]Word
	VromBanks [][]Word

	ExtendedRam [0x400]Word

	PrgBankCount int
	ChrRomCount  int
	Battery      bool
	Data         []byte

	PrgSwitchMode   Word
	ChrSwitchMode   Word
	ExtendedRamMode Word
	ChrUpperBits    Word

	FillModeTile  Word
	FillModeColor Word

	SelectedPrgRamChip Word

	IrqLatch   int
	IrqCounter int
	IrqEnabled bool

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

	SpriteSwapFunc [8]func()
	BgSwapFunc     [4]func()
}

func NewMmc5(r *Nrom) *Mmc5 {
	m := &Mmc5{
		RomBanks:     r.RomBanks,
		VromBanks:    r.VromBanks,
		PrgBankCount: r.PrgBankCount,
		ChrRomCount:  r.ChrRomCount,
		Battery:      r.Battery,
		Data:         r.Data,
	}

	m.PrgSwitchMode = 0x3

	m.Load()

	return m
}

func (m *Mmc5) Load() {
	// 2x the banks since we're storing 8k per bank
	// instead of 16k
	fmt.Printf("  Emulated PRG banks: %d\n", 2*m.PrgBankCount)
	m.RomBanks = make([][]Word, 2*m.PrgBankCount)
	for i := 0; i < 2*m.PrgBankCount; i++ {
		// Move 8kb chunk to 8kb bank
		bank := make([]Word, 0x2000)
		for x := 0; x < 0x2000; x++ {
			bank[x] = Word(m.Data[(0x2000*i)+x])
		}

		m.RomBanks[i] = bank
	}

	// Everything after PRG-ROM
	chrRom := m.Data[0x2000*len(m.RomBanks):]

	// CHR is stored in 1k banks
	if m.ChrRomCount > 0 {
		m.VromBanks = make([][]Word, m.ChrRomCount*8)
	} else {
		m.VromBanks = make([][]Word, 2)
	}

	for i := 0; i < cap(m.VromBanks); i++ {
		// Move 16kb chunk to 16kb bank
		m.VromBanks[i] = make([]Word, 0x0400)

		// If the game doesn't have CHR banks we
		// just need to allocate VRAM

		for x := 0; x < 0x0400; x++ {
			var val Word
			if m.ChrRomCount == 0 {
				val = 0
			} else {
				val = Word(chrRom[(0x0400*i)+x])
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

	// Write hardwired PRG banks (0xC000 and 0xE000)
	// Second to last bank
	m.PrgUpperHighBank = (((len(m.RomBanks) - 1) * 2) + 1) >> 1
	m.PrgUpperLowBank = m.PrgUpperHighBank - 1
	// Last bank

	// Write swappable PRG banks (0x8000 and 0xA000)
	m.PrgLowerLowBank = 0
	m.PrgLowerHighBank = 1
}

func (m *Mmc5) Write(v Word, a int) {
	switch a {
	case 0x5100:
		// PRG Switching mode
		m.PrgSwitchMode = v & 0x3
	case 0x5101:
		// CHR Switching mode
		m.ChrSwitchMode = v & 0x3
	case 0x5102:
		// PRG RAM protect 1
	case 0x5103:
		// PRG RAM protect 2
	case 0x5104:
		// Extended RAM mode
		fmt.Printf("Extended RAM mode: 0x%X\n", v&0x3)
		m.ExtendedRamMode = v & 0x3
	case 0x5105:
		// Nametable mapping
		m.SetNametableMapping(v)
	case 0x5106:
		// Fill-mode tile
		m.FillModeTile = v
	case 0x5107:
		// Fill-mode tile
		m.FillModeColor = v & 0x3
	case 0x5113:
		// PRG-RAM bank
		m.SelectedPrgRamChip = ((v >> 2) & 0x1)
		// TODO: Bank the RAM
	case 0x5114:
		// PRG bank 0
		if m.PrgSwitchMode != 3 {
			return
		}

		// TODO: (v >> 7) & 0x1 is the RAM/ROM toggle bit
		m.PrgLowerLowBank = int(v) & 0x7F
	case 0x5115:
		// PRG bank 1
		// TODO: (v >> 7) & 0x1 is the RAM/ROM toggle bit
		bank := int(v) & 0x7F

		switch m.PrgSwitchMode {
		case 1:
			// TODO: This was >> 1
			bank &= 0xFE

			m.PrgLowerLowBank = bank
			m.PrgLowerHighBank = bank + 1
		case 2:
			// TODO: This was >> 1
			bank &= 0xFE

			m.PrgLowerLowBank = bank
			m.PrgLowerHighBank = bank + 1
		case 3:
			m.PrgLowerHighBank = bank
		}
	case 0x5116:
		// PRG bank 2
		// TODO: (v >> 7) & 0x1 is the RAM/ROM toggle bit
		bank := (int(v) & 0x7F) % len(m.RomBanks)

		switch m.PrgSwitchMode {
		case 2:
			m.PrgUpperLowBank = bank
		case 3:
			m.PrgUpperLowBank = bank
		}
	case 0x5117:
		// PRG bank 3
		bank := (int(v) & 0x7F) % len(m.RomBanks)

		switch m.PrgSwitchMode {
		case 0:
			// TODO: This was >> 1
			bank >>= 2

			m.PrgLowerLowBank = bank
			m.PrgLowerHighBank = bank + 1
			m.PrgUpperLowBank = bank + 2
			m.PrgUpperHighBank = bank + 3
		case 1:
			// TODO: This was >> 1
			bank >>= 1

			m.PrgUpperLowBank = bank
			m.PrgUpperHighBank = bank + 1
		case 2:
			m.PrgUpperHighBank = bank
		case 3:
			m.PrgUpperHighBank = bank
		}
	// TODO: Registers $5120-$5127 apply to sprite graphics
	// and $5128-$512B for background graphics, but ONLY when
	// 8x16 sprites are enabled.
	case 0x5120:
		// Sprite CHR bank 0
		if m.ChrSwitchMode != 3 {
			m.SpriteSwapFunc[0] = func() {}
			return
		}

		m.SpriteSwapFunc[0] = func() {
			bank := (int(m.ChrUpperBits)<<8 | int(v)) % len(m.VromBanks)
			m.Chr000Bank = bank
		}
	case 0x5121:
		// Sprite CHR bank 1
		bank := (int(m.ChrUpperBits)<<8 | int(v)) % len(m.VromBanks)

		switch m.ChrSwitchMode {
		case 2:
			m.SpriteSwapFunc[1] = func() {
				m.Chr000Bank = bank
				m.Chr400Bank = bank + 1
			}
		case 3:
			m.SpriteSwapFunc[1] = func() {
				m.Chr400Bank = bank
			}
		}
	case 0x5122:
		// Sprite CHR bank 2
		if m.ChrSwitchMode != 3 {
			m.SpriteSwapFunc[2] = func() {}
			return
		}

		bank := (int(m.ChrUpperBits)<<8 | int(v)) % len(m.VromBanks)

		m.SpriteSwapFunc[2] = func() {
			m.Chr800Bank = bank
		}
	case 0x5123:
		// Sprite CHR bank 3
		bank := (int(m.ChrUpperBits)<<8 | int(v)) % len(m.VromBanks)

		switch m.ChrSwitchMode {
		case 1:
			m.SpriteSwapFunc[3] = func() {
				m.Chr000Bank = bank
				m.Chr400Bank = bank + 1
				m.Chr800Bank = bank + 2
				m.ChrC00Bank = bank + 3
			}
		case 2:
			m.SpriteSwapFunc[3] = func() {
				m.Chr800Bank = bank
				m.ChrC00Bank = bank + 1
			}
		case 3:
			m.SpriteSwapFunc[3] = func() {
				m.ChrC00Bank = bank
			}
		}
	case 0x5124:
		// Sprite CHR bank 4
		if m.ChrSwitchMode != 3 {
			m.SpriteSwapFunc[4] = func() {}
			return
		}

		bank := (int(m.ChrUpperBits)<<8 | int(v)) % len(m.VromBanks)

		m.SpriteSwapFunc[4] = func() {
			m.Chr1000Bank = bank
		}
	case 0x5125:
		// Sprite CHR bank 5
		bank := (int(m.ChrUpperBits)<<8 | int(v)) % len(m.VromBanks)

		switch m.ChrSwitchMode {
		case 2:
			m.SpriteSwapFunc[5] = func() {
				m.Chr1000Bank = bank
				m.Chr1400Bank = bank + 1
			}
		case 3:
			m.SpriteSwapFunc[5] = func() {
				m.Chr1400Bank = bank
			}
		}
	case 0x5126:
		// Sprite CHR bank 6
		if m.ChrSwitchMode != 3 {
			m.SpriteSwapFunc[6] = func() {}
			return
		}

		bank := (int(m.ChrUpperBits)<<8 | int(v)) % len(m.VromBanks)

		m.SpriteSwapFunc[6] = func() {
			m.Chr1800Bank = bank
		}
	case 0x5127:
		// Sprite CHR bank 7
		bank := (int(m.ChrUpperBits)<<8 | int(v)) % len(m.VromBanks)

		switch m.ChrSwitchMode {
		case 0:
			m.SpriteSwapFunc[7] = func() {
				m.Chr000Bank = bank
				m.Chr400Bank = bank + 1
				m.Chr800Bank = bank + 2
				m.ChrC00Bank = bank + 3
				m.Chr1000Bank = bank + 4
				m.Chr1400Bank = bank + 5
				m.Chr1800Bank = bank + 6
				m.Chr1C00Bank = bank + 7
			}
		case 1:
			m.SpriteSwapFunc[7] = func() {
				m.Chr1000Bank = bank
				m.Chr1400Bank = bank + 1
				m.Chr1800Bank = bank + 2
				m.Chr1C00Bank = bank + 3
			}
		case 2:
			m.SpriteSwapFunc[7] = func() {
				m.Chr1800Bank = bank
				m.Chr1C00Bank = bank + 1
			}
		case 3:
			m.SpriteSwapFunc[7] = func() {
				m.Chr1C00Bank = bank
			}
		}
	case 0x5128:
		// Background CHR bank 0
		if m.ChrSwitchMode != 3 {
			m.BgSwapFunc[0] = func() {}
			return
		}

		bank := (int(m.ChrUpperBits)<<8 | int(v)) % len(m.VromBanks)

		m.BgSwapFunc[0] = func() {
			m.Chr000Bank = bank
			m.Chr1000Bank = bank
		}
	case 0x5129:
		// Background CHR bank 1
		bank := (int(m.ChrUpperBits)<<8 | int(v)) % len(m.VromBanks)

		switch m.ChrSwitchMode {
		case 2:
			m.BgSwapFunc[1] = func() {
				m.Chr000Bank = bank
				m.Chr400Bank = bank + 1
				m.Chr1000Bank = bank
				m.Chr1400Bank = bank + 1
			}
		case 3:
			m.BgSwapFunc[1] = func() {
				m.Chr400Bank = bank
				m.Chr1400Bank = bank
			}
		}
	case 0x512A:
		// Background CHR bank 2
		if m.ChrSwitchMode != 3 {
			m.BgSwapFunc[2] = func() {}
			return
		}

		bank := (int(m.ChrUpperBits)<<8 | int(v)) % len(m.VromBanks)

		if m.ExtendedRamMode == 0x0 {
			m.BgSwapFunc[2] = func() {
				m.Chr800Bank = bank
			}
		} else {
			m.BgSwapFunc[2] = func() {
				m.Chr1800Bank = bank
			}
		}
	case 0x512B:
		// Background CHR bank 3
		bank := (int(m.ChrUpperBits)<<8 | int(v)) % len(m.VromBanks)

		switch m.ChrSwitchMode {
		case 0:
			m.BgSwapFunc[3] = func() {
				m.Chr000Bank = bank
				m.Chr400Bank = bank + 1
				m.Chr800Bank = bank + 2
				m.ChrC00Bank = bank + 3
				m.Chr1000Bank = bank + 4
				m.Chr1400Bank = bank + 5
				m.Chr1800Bank = bank + 6
				m.Chr1C00Bank = bank + 7
			}
		case 1:
			if m.ExtendedRamMode == 0x0 {
				m.BgSwapFunc[3] = func() {
					m.Chr000Bank = bank
					m.Chr400Bank = bank + 1
					m.Chr800Bank = bank + 2
					m.ChrC00Bank = bank + 3
				}
			} else {
				m.BgSwapFunc[3] = func() {
					m.Chr1000Bank = bank
					m.Chr1400Bank = bank + 1
					m.Chr1800Bank = bank + 2
					m.Chr1C00Bank = bank + 3
				}
			}
		case 2:
			if m.ExtendedRamMode == 0x0 {
				m.BgSwapFunc[3] = func() {
					m.Chr000Bank = bank
					m.Chr400Bank = bank + 1
				}
			} else {
				m.BgSwapFunc[3] = func() {
					m.Chr1000Bank = bank
					m.Chr1400Bank = bank + 1
				}
			}
		case 3:
			if m.ExtendedRamMode == 0x0 {
				m.BgSwapFunc[3] = func() {
					m.Chr400Bank = bank
				}
			} else {
				m.BgSwapFunc[3] = func() {
					m.Chr1400Bank = bank
				}
			}
		}
	case 0x5130:
		// Upper CHR bank bits
		fmt.Printf("Upper bits: 0x%X\n", v&0x3)
		m.ChrUpperBits = v & 0x3
	case 0x5203:
		// IRQ Counter
		m.IrqLatch = int(v)
		m.IrqCounter = 0
	case 0x5204:
		m.IrqEnabled = (v&0x80 == 0x80)
	default:
		// fmt.Printf("Unhandled write to: 0x%X -> 0x%X\n", a, v)
	}

	if a >= 0x5C00 && a <= 0x5FFF {
		if m.ExtendedRamMode != 0x3 {
			Ram[a] = v
			m.ExtendedRam[a-0x5C00] = v
		}
	}
}

func (m *Mmc5) WriteVram(v Word, a int) {
	switch {
	case a >= 0x1C00:
		m.VromBanks[m.Chr1C00Bank][a&0x3FF] = v
	case a >= 0x1800:
		m.VromBanks[m.Chr1800Bank][a&0x3FF] = v
	case a >= 0x1400:
		m.VromBanks[m.Chr1400Bank][a&0x3FF] = v
	case a >= 0x1000:
		m.VromBanks[m.Chr1000Bank][a&0x3FF] = v
	case a >= 0x0C00:
		m.VromBanks[m.ChrC00Bank][a&0x3FF] = v
	case a >= 0x0800:
		m.VromBanks[m.Chr800Bank][a&0x3FF] = v
	case a >= 0x0400:
		m.VromBanks[m.Chr400Bank][a&0x3FF] = v
	default:
		m.VromBanks[m.Chr000Bank][a&0x3FF] = v
	}
}

func (m *Mmc5) ReadVram(a int) Word {
	switch {
	case a >= 0x1C00:
		return m.VromBanks[m.Chr1C00Bank][a&0x3FF]
	case a >= 0x1800:
		return m.VromBanks[m.Chr1800Bank][a&0x3FF]
	case a >= 0x1400:
		return m.VromBanks[m.Chr1400Bank][a&0x3FF]
	case a >= 0x1000:
		return m.VromBanks[m.Chr1000Bank][a&0x3FF]
	case a >= 0x0C00:
		return m.VromBanks[m.ChrC00Bank][a&0x3FF]
	case a >= 0x0800:
		return m.VromBanks[m.Chr800Bank][a&0x3FF]
	case a >= 0x0400:
		return m.VromBanks[m.Chr400Bank][a&0x3FF]
	default:
		return m.VromBanks[m.Chr000Bank][a&0x3FF]
	}
}

func (m *Mmc5) ReadTile(a int) []Word {
	switch {
	case a >= 0x1C00:
		return m.VromBanks[m.Chr1C00Bank][a&0x3FF : a&0x3FF+16]
	case a >= 0x1800:
		return m.VromBanks[m.Chr1800Bank][a&0x3FF : a&0x3FF+16]
	case a >= 0x1400:
		return m.VromBanks[m.Chr1400Bank][a&0x3FF : a&0x3FF+16]
	case a >= 0x1000:
		return m.VromBanks[m.Chr1000Bank][a&0x3FF : a&0x3FF+16]
	case a >= 0x0C00:
		return m.VromBanks[m.ChrC00Bank][a&0x3FF : a&0x3FF+16]
	case a >= 0x0800:
		return m.VromBanks[m.Chr800Bank][a&0x3FF : a&0x3FF+16]
	case a >= 0x0400:
		return m.VromBanks[m.Chr400Bank][a&0x3FF : a&0x3FF+16]
	default:
		return m.VromBanks[m.Chr000Bank][a&0x3FF : a&0x3FF+16]
	}
}

func (m *Mmc5) Read(a int) Word {
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

func (m *Mmc5) BatteryBacked() bool {
	return m.Battery
}

func (m *Mmc5) SetNametableMapping(v Word) {
	var i Word
	for i = 0; i < 4; i++ {
		bits := (v >> (i * 2)) & 0x3
		switch bits {
		case 0:
			ppu.Nametables.LogicalTables[i] = &ppu.Nametables.Nametable0
		case 1:
			ppu.Nametables.LogicalTables[i] = &ppu.Nametables.Nametable1
		case 2:
			ppu.Nametables.LogicalTables[i] = &m.ExtendedRam
		case 3:
			var fillmode [0x400]Word
			for x := range fillmode {
				fillmode[x] = 0x0
			}
			ppu.Nametables.LogicalTables[i] = &fillmode
		}
	}
}

func (m *Mmc5) SwapSpriteVram() {
	for _, s := range m.SpriteSwapFunc {
		s()
	}
}

func (m *Mmc5) SwapBgVram() {
	for _, bg := range m.BgSwapFunc {
		bg()
	}
}

func (m *Mmc5) NotifyScanline() {
	if Ram[0x5204]&0x40 == 0x40 {
		// If In-Frame flag is set
		m.IrqCounter++
		if m.IrqEnabled && m.IrqCounter == m.IrqLatch {
			Ram[0x5204] |= 0x80
			cpu.RequestInterrupt(InterruptIrq)
		}
	} else {
		Ram[0x5204] = 0x40
		m.IrqCounter = 0
	}
}
