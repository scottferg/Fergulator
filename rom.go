package main

import (
	"errors"
	"fmt"
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

// Nrom
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
	WriteRamBank(m.RomBanks, int(v&0x7), 0x8000, Size16k)
}

func (m *Cnrom) Write(v Word, a int) {
	WriteVramBank(m.VromBanks, int(v&0x3), 0x0, Size8k)
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

	fmt.Printf("PRG-ROM Count: %d\n", r.PrgBankCount)
	fmt.Printf("CHR-ROM Count: %d\n", r.ChrRomCount)

	switch rom[6] & 0x1 {
	case 0x0:
		fmt.Println("Horizontal mirroring")
		ppu.Nametables.SetMirroring(MirroringHorizontal)
	case 0x1:
		fmt.Println("Vertical mirroring")
		ppu.Nametables.SetMirroring(MirroringVertical)
	}

	if (rom[6]>>0x1)&0x1 == 0x1 {
		r.BatteryBacked = true
	}

	r.Data = rom[16:]

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

	r.VromBanks = make([][]Word, r.ChrRomCount)
	for i := 0; i < r.ChrRomCount; i++ {
		// Move 16kb chunk to 16kb bank
		bank := make([]Word, 0x1000)
		for x := 0; x < 0x1000; x++ {
			bank[x] = Word(chrRom[(0x1000*i)+x])
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
		if r.ChrRomCount == 1 {
			WriteVramBank(r.VromBanks, 0, 0x0000, Size8k)
		} else {
			WriteVramBank(r.VromBanks, 0, 0x0000, Size4k)
			WriteVramBank(r.VromBanks, 1, 0x1000, Size4k)
		}
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
