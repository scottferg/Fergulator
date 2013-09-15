package nes

import (
	"fmt"
)

const (
	BankUpper = iota
	BankLower
)

type Mmc1 struct {
	RomBanks  [][]Word
	VromBanks [][]Word

	PrgBankCount int
	ChrRomCount  int
	Battery      bool
	Data         []byte

	Buffer        int
	BufferCounter uint
	PrgLowerBank  int
	PrgUpperBank  int
	PrgSwapBank   int
	PrgBankSize   int
	ChrBankSize   int
	ChrLowerBank  int
	ChrUpperBank  int
	Mirroring     int
}

func (m *Mmc1) Write(v Word, a int) {
	// If reset bit is set
	if v&0x80 != 0 {
		m.BufferCounter = 0
		m.Buffer = 0x0

		m.PrgSwapBank = BankLower
		m.PrgBankSize = Size16k
	} else {
		// Buffer the write
		m.Buffer = (m.Buffer & (0xFF - (0x1 << m.BufferCounter))) | ((int(v) & 0x1) << m.BufferCounter)
		m.BufferCounter++

		// If the buffer is filled
		if m.BufferCounter == 0x5 {
			m.SetRegister(m.RegisterNumber(a), m.Buffer)

			// Reset buffer
			m.BufferCounter = 0
			m.Buffer = 0
		}
	}
}

func (m *Mmc1) Read(a int) Word {
	if a >= 0xC000 {
		return m.RomBanks[m.PrgUpperBank][a&0x3FFF]
	}

	return m.RomBanks[m.PrgLowerBank][a&0x3FFF]
}

func (m *Mmc1) WriteVram(v Word, a int) {
	if a >= 0x1000 {
		m.VromBanks[m.ChrUpperBank][a&0xFFF] = v
	}

	m.VromBanks[m.ChrLowerBank][a&0xFFF] = v
}

func (m *Mmc1) ReadVram(a int) Word {
	if a >= 0x1000 {
		return m.VromBanks[m.ChrUpperBank][a&0xFFF]
	}

	return m.VromBanks[m.ChrLowerBank][a&0xFFF]
}

func (m *Mmc1) ReadTile(a int) []Word {
	if a >= 0x1000 {
		return m.VromBanks[m.ChrUpperBank][a&0xFFF : a&0xFFF+16]
	}

	return m.VromBanks[m.ChrLowerBank][a&0xFFF : a&0xFFF+16]
}

func (m *Mmc1) BatteryBacked() bool {
	return m.Battery
}

func (m *Mmc1) SetRegister(reg int, v int) {
	switch reg {
	// Control register
	case 0:
		// Set mirroring
		m.Mirroring = v & 0x3

		switch m.Mirroring {
		case 0x0:
			ppu.Nametables.SetMirroring(MirroringSingleUpper)
		case 0x1:
			ppu.Nametables.SetMirroring(MirroringSingleLower)
		case 0x2:
			ppu.Nametables.SetMirroring(MirroringVertical)
		case 0x3:
			ppu.Nametables.SetMirroring(MirroringHorizontal)
		}

		switch (v >> 0x2) & 0x3 {
		case 0x0:
			fallthrough
		case 0x1:
			m.PrgBankSize = Size32k
			m.PrgSwapBank = BankLower
		case 0x2:
			m.PrgBankSize = Size16k
			m.PrgSwapBank = BankUpper
		case 0x3:
			m.PrgBankSize = Size16k
			m.PrgSwapBank = BankLower
		}

		// Set CHR bank size
		switch (v >> 0x4) & 0x1 {
		case 0x0:
			fmt.Println("CHR 8k bank size")
			m.ChrBankSize = Size8k
		case 0x1:
			fmt.Println("CHR 4k bank size")
			m.ChrBankSize = Size4k
		}
	case 1:
		// CHR Bank 0
		if m.ChrRomCount == 0 {
			return
		}

		// Select VROM at 0x0000
		switch m.ChrBankSize {
		case Size8k:
			// Swap 8k VROM (in 8k mode, ignore first bit D0)
			bank := v & 0xF
			bank %= len(m.VromBanks)

			if v&0x10 == 0x10 {
				bank = (len(m.VromBanks) / 2) + (v & 0xF)
			} else {
				bank = v & 0xF
			}

			m.ChrUpperBank = bank + 1
			m.ChrLowerBank = bank
		case Size4k:
			// Swap 4k VROM
			var bank int

			if v&0x10 == 0x10 {
				bank = (len(m.VromBanks) / 2) + (v & 0xF)
			} else {
				bank = v & 0xF
			}
			m.ChrLowerBank = bank
		}
	case 2:
		// CHR Bank 1
		if m.ChrRomCount == 0 {
			return
		}

		// Select VROM bank at 0x1000, ignored in
		// 8k switching mode
		if m.ChrBankSize == Size4k {
			var bank int

			if v&0x10 == 0x10 {
				bank = (len(m.VromBanks) / 2) + (v & 0xF)
			} else {
				bank = v & 0xF
			}

			m.ChrUpperBank = bank
		}
	case 3:
		// PRG Bank
		switch m.PrgBankSize {
		case Size32k:
			// Swap 32k ROM (in 32k mode, ignore first bit D0)
			m.PrgLowerBank = ((v >> 0x1) & 0x7) * 2
			m.PrgUpperBank = m.PrgLowerBank + 1
		case Size16k:
			// Swap 16k ROM
			if m.PrgSwapBank == BankUpper {
				m.PrgUpperBank = v & 0xF
			} else {
				m.PrgLowerBank = v & 0xF
			}
		}
	}
}

func (m *Mmc1) RegisterNumber(a int) int {
	switch {
	case a >= 0x8000 && a <= 0x9FFF:
		return 0
	case a >= 0xA000 && a <= 0xBFFF:
		return 1
	case a >= 0xC000 && a <= 0xDFFF:
		return 2
	}

	return 3
}
