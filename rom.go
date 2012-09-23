package main

import (
	"errors"
	"fmt"
	"math"
)

const (
	Size4k  = 0x1000
	Size8k  = 0x2000
	Size16k = 0x4000
	Size32k = 0x8000
)

type Mapper interface {
	Write(v Word, a int)
}

type Rom struct {
	RomBanks  [][]Word
	VromBanks [][]Word

	PrgBankCount  int
	ChrRomCount   int
	BatteryBacked bool
	Data          []byte
}

type Unrom Rom
type Cnrom Rom

type Mmc1 struct {
	RomBanks  [][]Word
	VromBanks [][]Word

	PrgBankCount  int
	ChrRomCount   int
	BatteryBacked bool
	Data          []byte

	Buffer            int
	BufferCounter     uint
	Mirroring         int
	PrgSwitchingArea  int
	PrgSwitchingSize  int
	VromSwitchingSize int
	RomSelectionReg0  int
	RomSelectionReg1  int
	RomBankSelect     int
}

func WriteRamBank(rom [][]Word, bank, dest, size int) {
	for i := 0; i < size; i++ {
		Ram[i+dest] = rom[bank][i]
	}
}

func WriteVramBank(rom [][]Word, bank, dest, size int) {
	for i := 0; i < size; i++ {
		ppu.Vram[i+dest] = rom[bank][i]
	}
}

func (m *Rom) Write(v Word, a int) {
	// Nothing to do
}

func (m *Unrom) Write(v Word, a int) {
	WriteRamBank(m.RomBanks, int(v), 0x8000, Size16k)
}

func (m *Cnrom) Write(v Word, a int) {
	WriteVramBank(m.VromBanks, int(v&0x3), 0x0, Size8k)
}

func (m *Mmc1) Write(v Word, a int) {
	// If reset bit is set
	if v&0x80 != 0 {
		m.BufferCounter = 0
		m.Buffer = 0

		// Reset it
		if m.RegisterNumber(a) == 0 {
			fmt.Println("Resetting MMC")
			m.PrgSwitchingArea = 1
			m.PrgSwitchingSize = 1
		}
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

func (m *Mmc1) SetRegister(reg int, v int) {
	fmt.Printf("Register: 0x%X Value: 0x%X\n", reg, v)
	switch reg {
	case 0:
		// Control register
		tmp := v & 0x3
		if tmp != m.Mirroring {
			// Set mirroring
			m.Mirroring = tmp
			if (m.Mirroring & 0x2) == 0 {
				// TODO: Single screen mirroring
				ppu.Mirroring = MirroringSingleScreen
				ppu.Nametables.Init()
				fmt.Println("Single screen mirroring!")
			} else if (m.Mirroring & 0x1) != 0 {
				ppu.Mirroring = MirroringHorizontal
				ppu.Nametables.Init()
			} else {
				ppu.Mirroring = MirroringVertical
				ppu.Nametables.Init()
			}
		}

		m.PrgSwitchingArea = (v >> 0x2) & 0x1
		m.PrgSwitchingSize = (v >> 0x3) & 0x1
		m.VromSwitchingSize = (v >> 0x4) & 0x1
	case 1:
		// CHR Bank 0
		m.RomSelectionReg0 = (v >> 0x4) & 0x1

		if m.ChrRomCount > 0 {
			// Select VROM at 0x0000
			if m.VromSwitchingSize == 0 {
				// Swap 8k VROM
				if m.RomSelectionReg0 == 0 {
					WriteVramBank(m.VromBanks, v&0xF, 0x0, Size8k)
				} else {
					WriteVramBank(m.VromBanks, int(math.Floor(float64(m.ChrRomCount/2)))+(v&0xF), 0x0000, Size8k)
				}
			} else {
				// Swap 4k VROM
				if m.RomSelectionReg0 == 0 {
					WriteVramBank(m.VromBanks, v&0xF, 0x0, Size4k)
				} else {
					fmt.Printf("Bank: %d\n", int(math.Floor(float64(m.ChrRomCount/2)))+(v&0xF))
					WriteVramBank(m.VromBanks, int(math.Floor(float64(m.ChrRomCount/2)))+(v&0xF), 0x0, Size4k)
				}
			}
		}
	case 2:
		// CHR Bank 1
		m.RomSelectionReg1 = (v >> 0x4) & 0x1

		if m.ChrRomCount > 0 {
			// Select VROM bank at 0x1000
			if m.VromSwitchingSize == 1 {
				if m.RomSelectionReg1 == 0 {
					WriteRamBank(m.RomBanks, v&0xF, 0x1000, Size4k)
				} else {
					WriteRamBank(m.RomBanks, int(math.Floor(float64(m.ChrRomCount/2)))+(v&0xF), 0x1000, Size4k)
				}
			}
		}
	case 3:
		// PRG Bank
		baseBank := 0

		var bank int

		if m.PrgBankCount >= 32 {
			// 1024kb Cartridge
			if m.VromSwitchingSize == 0 {
				if m.RomSelectionReg0 == 1 {
					baseBank = 16
				}
			} else {
				baseBank = (m.RomSelectionReg0 | (m.RomSelectionReg1 << 0x1)) << 0x3
			}
		} else if m.PrgBankCount >= 16 {
			if m.RomSelectionReg0 == 1 {
				baseBank = 8
			}
		}

		fmt.Printf("Base bank: %d\n", baseBank)
		if m.PrgSwitchingSize == 0 {
			fmt.Println("32k ROM load")
			// 32k 
			bank = baseBank + (v & 0xF)
			// Load bank
			WriteRamBank(m.RomBanks, bank, 0x8000, Size16k)
			WriteRamBank(m.RomBanks, bank+1, 0xC000, Size16k)
		} else {
			// 16k
			bank = baseBank*2 + (v & 0xF)
			if m.PrgSwitchingArea == 0 {
				fmt.Printf("Upper 16k ROM load: %d\n", bank)
				WriteRamBank(m.RomBanks, bank, 0xC000, Size16k)
			} else {
				fmt.Printf("Lower 16k ROM load: %d\n", bank)
				WriteRamBank(m.RomBanks, bank, 0x8000, Size16k)
				fmt.Printf("Value: 0x%X\n", m.RomBanks[bank][0])
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

func LoadRom(rom []byte) (m Mapper, e error) {
	r := new(Rom)

	if string(rom[0:3]) != "NES" {
		return m, errors.New("Invalid ROM file")

		if rom[3] != 0x1a {
			return m, errors.New("Invalid ROM file")
		}
	}

	r.PrgBankCount = int(rom[4])
	r.ChrRomCount = int(rom[5])

	switch rom[6] & 0x1 {
	case 0x0:
		fmt.Println("Horizontal mirroring")
		ppu.Mirroring = MirroringHorizontal
		ppu.Nametables.Init()
	case 0x1:
		fmt.Println("Vertical mirroring")
		ppu.Mirroring = MirroringVertical
		ppu.Nametables.Init()
	}

	if (rom[6]>>0x1)&0x1 == 0x1 {
		r.BatteryBacked = true
	}

	r.Data = rom[16:]

	fmt.Printf("PRG-ROM Count: %d\n", r.PrgBankCount)
	fmt.Printf("CHR-ROM Count: %d\n", r.ChrRomCount)
	r.RomBanks = make([][]Word, r.PrgBankCount)
	for i := 0; i < r.PrgBankCount; i++ {
		// Move 16kb chunk to 16kb bank
		bank := make([]Word, 0x4000)
		for x := 0; x < 0x4000; x++ {
			bank[x] = Word(r.Data[(0x4000*i)+x])
		}

		r.RomBanks[i] = bank
	}

	// Everything after PRG-ROM
	chrRom := r.Data[0x4000*len(r.RomBanks):]

	vramBankCount := (len(chrRom) / 0x2000)
	r.VromBanks = make([][]Word, vramBankCount)
	for i := 0; i < vramBankCount; i++ {
		// Move 16kb chunk to 16kb bank
		bank := make([]Word, 0x2000)
		for x := 0; x < 0x2000; x++ {
			bank[x] = Word(chrRom[(0x2000*i)+x])
		}

		r.VromBanks[i] = bank
	}

	// Write the first ROM bank
	WriteRamBank(r.RomBanks, 0, 0x8000, Size16k)

	if r.PrgBankCount > 1 {
		// and the last ROM bank
		WriteRamBank(r.RomBanks, r.PrgBankCount-1, 0xC000, Size16k)
	} else {
		// Or write the first ROM bank to the upper region
		WriteRamBank(r.RomBanks, 0, 0xC000, Size16k)
	}

    // If we have CHR-ROM, load the first two banks
    // into VRAM region 0x0000-0x1000
	if r.ChrRomCount > 0 {
        WriteVramBank(r.VromBanks, 0, 0x0, Size8k)
	}

	// Check mapper, get the proper type
	mapper := (Word(rom[6])>>4 | (Word(rom[7]) & 0xF0))
	fmt.Printf("Mapper: 0x%X\n", mapper)
	switch mapper {
	case 0x00:
		fallthrough
	case 0x40:
		fallthrough
	case 0x41:
		// NROM
		return r, nil
	case 0x01:
		// MMC1
		m = &Mmc1{
			RomBanks:      r.RomBanks,
			VromBanks:     r.VromBanks,
			PrgBankCount:  r.PrgBankCount,
			ChrRomCount:   r.ChrRomCount,
			BatteryBacked: r.BatteryBacked,
			Data:          r.Data,
		}
	case 0x42:
		fallthrough
	case 0x02:
		// Unrom
		m = &Unrom{
			RomBanks:      r.RomBanks,
			VromBanks:     r.VromBanks,
			PrgBankCount:  r.PrgBankCount,
			ChrRomCount:   r.ChrRomCount,
			BatteryBacked: r.BatteryBacked,
			Data:          r.Data,
		}
	case 0x03:
		// Cnrom
		m = &Cnrom{
			RomBanks:      r.RomBanks,
			VromBanks:     r.VromBanks,
			PrgBankCount:  r.PrgBankCount,
			ChrRomCount:   r.ChrRomCount,
			BatteryBacked: r.BatteryBacked,
			Data:          r.Data,
		}
	default:
		// Unsupported
		return m, errors.New(fmt.Sprintf("Unsupported memory mapper: 0x%X", mapper))
	}

	return
}
