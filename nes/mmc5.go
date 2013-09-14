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

	SpriteSwapFunc [8]func()
	BgSwapFunc     [4]func()
}

func NewMmc5(r *Rom) *Mmc5 {
	m := &Mmc5{
		RomBanks:     r.RomBanks,
		VromBanks:    r.VromBanks,
		PrgBankCount: r.PrgBankCount,
		ChrRomCount:  r.ChrRomCount,
		Battery:      r.Battery,
		Data:         r.Data,
	}

	m.PrgSwitchMode = 0x3

	m.Write8kRamBank(((len(m.RomBanks) - 1) * 2), 0xC000)
	m.Write8kRamBank(((len(m.RomBanks)-1)*2)+1, 0xE000)

	return m
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
		bank := int(v) & 0x7F
		m.Write8kRamBank(bank, 0x8000)
	case 0x5115:
		// PRG bank 1
		// TODO: (v >> 7) & 0x1 is the RAM/ROM toggle bit
		bank := int(v) & 0x7F

		switch m.PrgSwitchMode {
		case 1:
			WriteRamBank(m.RomBanks, bank>>1, 0x8000, Size16k)
		case 2:
			WriteRamBank(m.RomBanks, bank>>1, 0x8000, Size16k)
		case 3:
			m.Write8kRamBank(bank, 0xA000)
		}
	case 0x5116:
		// PRG bank 2
		// TODO: (v >> 7) & 0x1 is the RAM/ROM toggle bit
		bank := int(v) & 0x7F

		switch m.PrgSwitchMode {
		case 2:
			m.Write8kRamBank(bank, 0xC000)
		case 3:
			m.Write8kRamBank(bank, 0xC000)
		}
	case 0x5117:
		// PRG bank 3
		bank := int(v) & 0x7F

		switch m.PrgSwitchMode {
		case 0:
			WriteRamBank(m.RomBanks, bank>>2, 0x8000, Size32k)
		case 1:
			WriteRamBank(m.RomBanks, bank>>1, 0xC000, Size16k)
		case 2:
			m.Write8kRamBank(bank, 0xE000)
		case 3:
			m.Write8kRamBank(bank, 0xE000)
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
			bank := int(m.ChrUpperBits)<<8 | int(v)
			m.Write1kVramBank(bank, 0x0000)
		}
	case 0x5121:
		// Sprite CHR bank 1
		bank := int(m.ChrUpperBits)<<8 | int(v)

		switch m.ChrSwitchMode {
		case 2:
			m.SpriteSwapFunc[1] = func() {
				WriteVramBank(m.VromBanks, bank, 0x0000, Size2k)
			}
		case 3:
			m.SpriteSwapFunc[1] = func() {
				m.Write1kVramBank(bank, 0x0400)
			}
		}
	case 0x5122:
		// Sprite CHR bank 2
		if m.ChrSwitchMode != 3 {
			m.SpriteSwapFunc[2] = func() {}
			return
		}

		bank := int(m.ChrUpperBits)<<8 | int(v)

		m.SpriteSwapFunc[2] = func() {
			m.Write1kVramBank(bank, 0x0800)
		}
	case 0x5123:
		// Sprite CHR bank 3
		bank := int(m.ChrUpperBits)<<8 | int(v)

		switch m.ChrSwitchMode {
		case 1:
			m.SpriteSwapFunc[3] = func() {
				WriteVramBank(m.VromBanks, bank, 0x0000, Size4k)
			}
		case 2:
			m.SpriteSwapFunc[3] = func() {
				WriteVramBank(m.VromBanks, bank, 0x0800, Size2k)
			}
		case 3:
			m.SpriteSwapFunc[3] = func() {
				m.Write1kVramBank(bank, 0x0C00)
			}
		}
	case 0x5124:
		// Sprite CHR bank 4
		if m.ChrSwitchMode != 3 {
			m.SpriteSwapFunc[4] = func() {}
			return
		}

		bank := int(m.ChrUpperBits)<<8 | int(v)

		m.SpriteSwapFunc[4] = func() {
			m.Write1kVramBank(bank, 0x1000)
		}
	case 0x5125:
		// Sprite CHR bank 5
		bank := int(m.ChrUpperBits)<<8 | int(v)

		switch m.ChrSwitchMode {
		case 2:
			m.SpriteSwapFunc[5] = func() {
				WriteVramBank(m.VromBanks, bank, 0x1000, Size2k)
			}
		case 3:
			m.SpriteSwapFunc[5] = func() {
				m.Write1kVramBank(bank, 0x1400)
			}
		}
	case 0x5126:
		// Sprite CHR bank 6
		if m.ChrSwitchMode != 3 {
			m.SpriteSwapFunc[6] = func() {}
			return
		}

		bank := int(m.ChrUpperBits)<<8 | int(v)

		m.SpriteSwapFunc[6] = func() {
			m.Write1kVramBank(bank, 0x1800)
		}
	case 0x5127:
		// Sprite CHR bank 7
		bank := int(m.ChrUpperBits)<<8 | int(v)

		switch m.ChrSwitchMode {
		case 0:
			m.SpriteSwapFunc[7] = func() {
				WriteVramBank(m.VromBanks, bank, 0x0000, Size8k)
			}
		case 1:
			m.SpriteSwapFunc[7] = func() {
				WriteVramBank(m.VromBanks, bank, 0x1000, Size4k)
			}
		case 2:
			m.SpriteSwapFunc[7] = func() {
				WriteVramBank(m.VromBanks, bank, 0x1800, Size2k)
			}
		case 3:
			m.SpriteSwapFunc[7] = func() {
				m.Write1kVramBank(bank, 0x1C00)
			}
		}
	case 0x5128:
		// Background CHR bank 0
		if m.ChrSwitchMode != 3 {
			m.BgSwapFunc[0] = func() {}
			return
		}

		bank := int(m.ChrUpperBits)<<8 | int(v)
		m.BgSwapFunc[0] = func() {
			m.Write1kVramBank(bank, 0x0000)
			m.Write1kVramBank(bank, 0x1000)
		}
	case 0x5129:
		// Background CHR bank 1
		bank := int(m.ChrUpperBits)<<8 | int(v)

		switch m.ChrSwitchMode {
		case 2:
			m.BgSwapFunc[1] = func() {
				WriteVramBank(m.VromBanks, bank, 0x0000, Size2k)
				WriteVramBank(m.VromBanks, bank, 0x1000, Size2k)
			}
		case 3:
			m.BgSwapFunc[1] = func() {
				m.Write1kVramBank(bank, 0x0400)
				m.Write1kVramBank(bank, 0x1400)
			}
		}
	case 0x512A:
		// Background CHR bank 2
		if m.ChrSwitchMode != 3 {
			m.BgSwapFunc[2] = func() {}
			return
		}

		bank := int(m.ChrUpperBits)<<8 | int(v)
		if m.ExtendedRamMode == 0x0 {
			m.BgSwapFunc[2] = func() {
				m.Write1kVramBank(bank, 0x0800)
			}
		} else {
			m.BgSwapFunc[2] = func() {
				m.Write1kVramBank(bank, 0x1800)
			}
		}
	case 0x512B:
		// Background CHR bank 3
		bank := int(m.ChrUpperBits)<<8 | int(v)

		switch m.ChrSwitchMode {
		case 0:
			m.BgSwapFunc[3] = func() {
				WriteVramBank(m.VromBanks, bank, 0x0000, Size8k)
			}
		case 1:
			if m.ExtendedRamMode == 0x0 {
				m.BgSwapFunc[3] = func() {
					WriteVramBank(m.VromBanks, bank, 0x0000, Size4k)
				}
			} else {
				m.BgSwapFunc[3] = func() {
					WriteVramBank(m.VromBanks, bank, 0x1000, Size4k)
				}
			}
		case 2:
			if m.ExtendedRamMode == 0x0 {
				m.BgSwapFunc[3] = func() {
					WriteVramBank(m.VromBanks, bank, 0x0000, Size2k)
				}
			} else {
				m.BgSwapFunc[3] = func() {
					WriteVramBank(m.VromBanks, bank, 0x1000, Size2k)
				}
			}
		case 3:
			if m.ExtendedRamMode == 0x0 {
				m.BgSwapFunc[3] = func() {
					m.Write1kVramBank(bank, 0x0400)
				}
			} else {
				m.BgSwapFunc[3] = func() {
					m.Write1kVramBank(bank, 0x1400)
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

func (m *Mmc5) Read(a int) Word {
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

func (m *Mmc5) Write8kRamBank(bank, dest int) {
	b := (bank >> 1) % len(m.RomBanks)
	offset := (bank % 2) * 0x2000

	WriteOffsetRamBank(m.RomBanks, b, dest, Size8k, offset)
}

func (m *Mmc5) Write1kVramBank(bank, dest int) {
	b := (bank >> 2) % len(m.VromBanks)
	offset := (bank % 4) * 0x400

	WriteOffsetVramBank(m.VromBanks, b, dest, Size1k, offset)
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
