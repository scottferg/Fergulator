package nes

import (
	"errors"
	"fmt"
)

const (
	Size1k  = 0x0400
	Size2k  = 0x0800
	Size4k  = 0x1000
	Size8k  = 0x2000
	Size16k = 0x4000
	Size32k = 0x8000
)

type Mapper interface {
	Write(v Word, a int)
	Read(a int) Word
	WriteVram(v Word, a int)
	ReadVram(a int) Word
	ReadTile(a int) []Word
	BatteryBacked() bool
}

func WriteRamBank(rom [][]Word, bank, dest, size int) {
	bank %= len(rom)

	for i := 0; i < size; i++ {
		Ram[i+dest] = rom[bank][i]
	}
}

// Used by MMC3 for selecting 8kb chunks of a PRG-ROM bank
func WriteOffsetRamBank(rom [][]Word, bank, dest, size, offset int) {
	bank %= len(rom)

	for i := 0; i < size; i++ {
		Ram[i+dest] = rom[bank][i+offset]
	}
}

func WriteVramBank(rom [][]Word, bank, dest, size int) {
	bank %= len(rom)

	for i := 0; i < size; i++ {
		ppu.Vram[i+dest] = rom[bank][i]
	}
}

func WriteOffsetVramBank(rom [][]Word, bank, dest, size, offset int) {
	bank %= len(rom)

	for i := 0; i < size; i++ {
		ppu.Vram[i+dest] = rom[bank][i+offset]
	}
}

func LoadRom(rom []byte) (m Mapper, e error) {
	r := new(Nrom)

	if string(rom[0:3]) != "NES" {
		return m, errors.New("Invalid ROM file")

		if rom[3] != 0x1a {
			return m, errors.New("Invalid ROM file")
		}
	}

	r.PrgBankCount = int(rom[4])
	r.ChrRomCount = int(rom[5])

	fmt.Printf("-----------------\nROM:\n  ")

	fmt.Printf("PRG-ROM banks: %d (%d real)\n  ", r.PrgBankCount, r.PrgBankCount)
	fmt.Printf("CHR-ROM banks: %d (%d real)\n  ", 2*r.ChrRomCount, r.ChrRomCount)

	fmt.Printf("Mirroring: ")
	switch rom[6] & 0x1 {
	case 0x0:
		fmt.Printf("Horizontal\n  ")
		ppu.Nametables.SetMirroring(MirroringHorizontal)
	case 0x1:
		fmt.Printf("Vertical\n  ")
		ppu.Nametables.SetMirroring(MirroringVertical)
	}

	if (rom[6]>>0x1)&0x1 == 0x1 {
		r.Battery = true
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

	if r.ChrRomCount > 0 {
		r.VromBanks = make([][]Word, r.ChrRomCount*2)
	} else {
		r.VromBanks = make([][]Word, 2)
	}

	for i := 0; i < cap(r.VromBanks); i++ {
		// Move 16kb chunk to 16kb bank
		r.VromBanks[i] = make([]Word, 0x1000)

		// If the game doesn't have CHR banks we
		// just need to allocate VRAM
		if r.ChrRomCount == 0 {
			continue
		}

		for x := 0; x < 0x1000; x++ {
			r.VromBanks[i][x] = Word(chrRom[(0x1000*i)+x])
		}
	}

	// Check mapper, get the proper type
	mapper := (Word(rom[6])>>4 | (Word(rom[7]) & 0xF0))
	fmt.Printf("Mapper: 0x%X -> ", mapper)
	switch mapper {
	case 0x00:
		fallthrough
	case 0x40:
		fallthrough
	case 0x41:
		// NROM
		fmt.Printf("NROM\n")
		return r, nil
	case 0x01:
		// MMC1
		fmt.Printf("MMC1\n")
        m = NewMmc1(r)
	case 0x42:
		fallthrough
	case 0x02:
		// Unrom
		fmt.Printf("UNROM\n")
		m = &Unrom{
			RomBanks:     r.RomBanks,
			VromBanks:    r.VromBanks,
			PrgBankCount: r.PrgBankCount,
			ChrRomCount:  r.ChrRomCount,
			Battery:      r.Battery,
			Data:         r.Data,
		}
	case 0x43:
		fallthrough
	case 0x03:
		// Cnrom
		fmt.Printf("CNROM\n")
		m = &Cnrom{
			RomBanks:     r.RomBanks,
			VromBanks:    r.VromBanks,
			PrgBankCount: r.PrgBankCount,
			ChrRomCount:  r.ChrRomCount,
			Battery:      r.Battery,
			Data:         r.Data,
		}
	case 0x07:
		// Anrom
		fmt.Printf("ANROM\n")
		m = &Anrom{
			RomBanks:     r.RomBanks,
			VromBanks:    r.VromBanks,
			PrgBankCount: r.PrgBankCount,
			ChrRomCount:  r.ChrRomCount,
			Battery:      r.Battery,
			Data:         r.Data,
			PrgUpperBank: len(r.RomBanks) - 1,
		}
		/*
			case 0x09:
				// MMC2
				fmt.Printf("MMC2\n")
				m = NewMmc2(r)
			case 0x44:
				fallthrough
			case 0x04:
				// MMC3
				fmt.Printf("MMC3\n")
				m = NewMmc3(r)
			case 0x05:
				// MMC5
				fmt.Printf("MMC5\n")
				m = NewMmc5(r)
		*/
	default:
		// Unsupported
		fmt.Printf("Unsupported\n")
		return m, errors.New(fmt.Sprintf("Unsupported memory mapper: 0x%X", mapper))
	}

	fmt.Printf("-----------------\n")

	return
}
