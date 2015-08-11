package nes

import (
	"errors"
	"fmt"
)

type Mapper interface {
	Write(v word, a int)
	Read(a int) word
	WriteVram(v word, a int)
	ReadVram(a int) word
	ReadTile(a int) []word
	BatteryBacked() bool
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

	// Check mapper, get the proper type
	mapper := (word(rom[6])>>4 | (word(rom[7]) & 0xF0))
	fmt.Printf("Mapper: 0x%X -> ", mapper)
	switch mapper {
	case 0x00:
		fallthrough
	case 0x40:
		fallthrough
	case 0x41:
		// NROM
		fmt.Printf("NROM\n")
		r.Load()
		return r, nil
	case 0x01:
		// MMC1
		fmt.Printf("MMC1\n")
		r.Load()
		m = NewMmc1(r)
	case 0x42:
		fallthrough
	case 0x02:
		// Unrom
		fmt.Printf("UNROM\n")
		r.Load()
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
		r.Load()
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
		r.Load()
		m = &Anrom{
			RomBanks:     r.RomBanks,
			VromBanks:    r.VromBanks,
			PrgBankCount: r.PrgBankCount,
			ChrRomCount:  r.ChrRomCount,
			Battery:      r.Battery,
			Data:         r.Data,
			PrgUpperBank: len(r.RomBanks) - 1,
		}
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
	case 0x09:
		// MMC2
		fmt.Printf("MMC2\n")
		m = NewMmc2(r)
	default:
		// Unsupported
		fmt.Printf("Unsupported\n")
		return m, errors.New(fmt.Sprintf("Unsupported memory mapper: 0x%X", mapper))
	}

	fmt.Printf("-----------------\n")

	return
}
