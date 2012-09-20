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
	WriteRamBank(bank, dest, size int)
	WriteVramBank(bank, dest, size int)
	Write(v Word, a int)
	Init(rom []byte) error
}

type Rom struct {
	PrgBankCount int
	ChrRomCount  int
	Data         []byte

	// TODO: MMC1
	Buffer            int
	BufferCounter     uint
	Mirroring         int
	PrgSwitchingArea  int
	PrgSwitchingSize  int
	VromSwitchingSize int
	RomSelectionReg0  int
	RomSelectionReg1  int
	RomBankSelect     int

	RomBanks  [][]Word
	VromBanks [][]Word
}

type Nrom Rom
type Mmc1 Rom
type Unrom Rom
type Cnrom Rom

// TODO: HOLY SHIT

func (r *Nrom) WriteRamBank(bank, dest, size int) {
	for i := 0; i < size; i++ {
		Ram[i+dest] = r.RomBanks[bank][i]
	}
}

func (r *Nrom) WriteVramBank(bank, dest, size int) {
	for i := 0; i < size; i++ {
		ppu.Vram[i+dest] = r.VromBanks[bank][i]
	}
}

func (r *Mmc1) WriteRamBank(bank, dest, size int) {
    limit := 1
    if size > Size16k {
        limit = size / Size16k
    }

    for x := 0; x < limit; x++ {
        for i := 0; i < Size16k; i++ {
            Ram[i+dest] = r.RomBanks[bank][i]
        }

        bank += 1
    }
}

func (r *Mmc1) WriteVramBank(bank, dest, size int) {
	for i := 0; i < size; i++ {
		ppu.Vram[i+dest] = r.VromBanks[bank][i]
	}
}

func (r *Unrom) WriteRamBank(bank, dest, size int) {
	for i := 0; i < size; i++ {
		Ram[i+dest] = r.RomBanks[bank][i]
	}
}

func (r *Unrom) WriteVramBank(bank, dest, size int) {
	for i := 0; i < size; i++ {
		ppu.Vram[i+dest] = r.VromBanks[bank][i]
	}
}

func (r *Cnrom) WriteRamBank(bank, dest, size int) {
	for i := 0; i < size; i++ {
		Ram[i+dest] = r.RomBanks[bank][i]
	}
}

func (r *Cnrom) WriteVramBank(bank, dest, size int) {
	for i := 0; i < size; i++ {
		ppu.Vram[i+dest] = r.VromBanks[bank][i]
	}
}

// ----------------------------------------

func (r *Nrom) Write(v Word, a int) {
	// Nothing to do
}

func (r *Nrom) Init(rom []byte) error {
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

	// ROM data dests at byte 16
	r.Data = rom[16:]

	fmt.Printf("PRG-ROM Count: %d\n", r.PrgBankCount)
	r.RomBanks = make([][]Word, (len(r.Data) / 0x4000))

	bankCount := (len(r.Data) / 0x4000)
	for i := 0; i < bankCount; i++ {
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

	switch r.PrgBankCount {
	case 0x01:
		r.WriteRamBank(0, 0x8000, Size16k)
		r.WriteRamBank(0, 0xC000, Size16k)

		if r.ChrRomCount != 0 {
			r.WriteVramBank(0, 0x0, Size8k)
		}
	case 0x02:
		r.WriteRamBank(0, 0x8000, Size16k)
		r.WriteRamBank(1, 0xC000, Size16k)
		r.WriteVramBank(0, 0x0, Size8k)
	}

	return nil
}

func (r *Mmc1) Write(v Word, a int) {
	// If reset bit is set
	if v&0x80 != 0 {
		r.BufferCounter = 0
		r.Buffer = 0

		// Reset it
		if r.RegisterNumber(a) == 0 {
			r.PrgSwitchingArea = 1
			r.PrgSwitchingSize = 1
		}
	} else {
		// Buffer the write
		r.Buffer = (r.Buffer & (0xFF - (0x1 << r.BufferCounter))) | ((int(v) & 0x1) << r.BufferCounter)
		r.BufferCounter++

		// If the buffer is filled
		if r.BufferCounter == 0x5 {
			r.SetRegister(r.RegisterNumber(a), r.Buffer)

			// Reset buffer
			r.BufferCounter = 0
			r.Buffer = 0
		}
	}
}

func (r *Mmc1) SetRegister(reg int, v int) {
	switch reg {
	case 0:
		// Control register
		tmp := v & 0x3
		if tmp != r.Mirroring {
			// Set mirroring
			r.Mirroring = tmp
			if (r.Mirroring & 0x2) == 0 {
				// TODO: Single screen mirroring
				ppu.Mirroring = MirroringSingleScreen
				ppu.Nametables.Init()
				fmt.Println("Single screen mirroring!")
			} else if (r.Mirroring & 0x1) != 0 {
				ppu.Mirroring = MirroringHorizontal
				ppu.Nametables.Init()
			} else {
				ppu.Mirroring = MirroringVertical
				ppu.Nametables.Init()
			}
		}

		r.PrgSwitchingArea = (v >> 0x2) & 0x1
		r.PrgSwitchingSize = (v >> 0x3) & 0x1
		r.VromSwitchingSize = (v >> 0x4) & 0x1
	case 1:
		// CHR Bank 0
		r.RomSelectionReg0 = (v >> 0x4) & 0x1

		if r.ChrRomCount > 0 {
			// Select VROM at 0x0000
			if r.VromSwitchingSize == 0 {
				// Swap 8k VROM
				if r.RomSelectionReg0 == 0 {
					r.WriteVramBank(v&0xF, 0x0, Size8k)
				} else {
					r.WriteVramBank(int(math.Floor(float64(r.ChrRomCount/2)))+(v&0xF), 0x0000, Size8k)
				}
			} else {
				// Swap 4k VROM
				if r.RomSelectionReg0 == 0 {
					r.WriteVramBank(v&0xF, 0x0, Size4k)
				} else {
                    fmt.Printf("Bank: %d\n", int(math.Floor(float64(r.ChrRomCount/2)))+(v&0xF))
					r.WriteVramBank(int(math.Floor(float64(r.ChrRomCount/2)))+(v&0xF), 0x0, Size4k)
				}
			}
		}
	case 2:
		// CHR Bank 1
		r.RomSelectionReg1 = (v >> 0x4) & 0x1

		if r.ChrRomCount > 0 {
			// Select VROM bank at 0x1000
			if r.VromSwitchingSize == 1 {
				if r.RomSelectionReg1 == 0 {
                    r.WriteRamBank(v&0xF, 0x1000, Size4k)
				} else {
                    r.WriteRamBank(int(math.Floor(float64(r.ChrRomCount/2)))+(v&0xF), 0x1000, Size4k)
				}
			}
		}
	case 3:
		// PRG Bank
		baseBank := 0

		var bank int

		if r.PrgBankCount >= 32 {
			// 1024kb Cartridge
			if r.VromSwitchingSize == 0 {
				if r.RomSelectionReg0 == 1 {
					baseBank = 16
				}
			} else {
				baseBank = (r.RomSelectionReg0 | (r.RomSelectionReg1 << 0x1)) << 0x3
			}
		} else if r.PrgBankCount >= 16 {
			if r.RomSelectionReg0 == 1 {
				baseBank = 8
			}
		}

		if r.PrgSwitchingSize == 0 {
			// 32k 
			bank = baseBank + (v & 0xF)
			// Load bank
            r.WriteRamBank(bank * 2, 0x8000, Size32k)
		} else {
            // 16k
			bank = (baseBank * 2) + (v & 0xF)
			if r.PrgSwitchingArea == 0 {
                r.WriteRamBank(bank, 0xC000, Size16k)
			} else {
                r.WriteRamBank(bank, 0x8000, Size16k)
			}
		}
	}
}

func (r *Mmc1) RegisterNumber(a int) int {
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

func (r *Mmc1) Init(rom []byte) error {
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
	r.WriteRamBank(0, 0x8000, Size16k)
	// and the last ROM bank
	r.WriteRamBank(r.PrgBankCount - 1, 0xC000, Size16k)

    if r.ChrRomCount > 0 {
        r.WriteVramBank(0, 0x0, Size4k)
        r.WriteVramBank(1, 0x1000, Size4k)
    }

	return nil
}

func (r *Unrom) Init(rom []byte) error {
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

	// ROM data dests at byte 16
	r.Data = rom[16:]
	r.RomBanks = make([][]Word, (len(r.Data) / 0x4000))

	fmt.Printf("PRG-ROM Count: %d\n", r.PrgBankCount)

	bankCount := (len(r.Data) / 0x4000)
	for i := 0; i < bankCount; i++ {
		// Move 16kb chunk to 16kb bank
		bank := make([]Word, 0x4000)
		for x := 0; x < 0x4000; x++ {
			bank[x] = Word(r.Data[(0x4000*i)+x])
		}

		r.RomBanks[i] = bank
	}

	// Write the first ROM bank
	r.WriteRamBank(0, 0x8000, Size16k)
	// and the last ROM bank
	r.WriteRamBank(7, 0xC000, Size16k)

	fmt.Printf("VROM: %d\n", r.ChrRomCount)

	return nil
}

func (r *Unrom) Write(v Word, a int) {
	r.WriteRamBank(int(v), 0x8000, Size16k)
}

func (r *Cnrom) Init(rom []byte) error {
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

	// ROM data dests at byte 16
	r.Data = rom[16:]

	fmt.Printf("PRG-ROM Count: %d\n", r.PrgBankCount)

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
	r.WriteRamBank(0, 0x8000, Size16k)
	// and the last ROM bank

	if r.PrgBankCount > 1 {
		r.WriteRamBank(1, 0xC000, Size16k)
	} else {
		r.WriteRamBank(0, 0xC000, Size16k)
	}

	r.WriteVramBank(0, 0x0, Size8k)

	fmt.Printf("VROM: %d\n", r.ChrRomCount)

	return nil
}

func (r *Cnrom) Write(v Word, a int) {
	r.WriteVramBank(int(v&0x3), 0x0, Size8k)
}

func LoadRom(rom []byte) (r Mapper, e error) {
	if string(rom[0:3]) != "NES" {
		return r, errors.New("Invalid ROM file")

		if rom[3] != 0x1a {
			return r, errors.New("Invalid ROM file")
		}
	}

	mapper := (Word(rom[6])>>4 | (Word(rom[7]) & 0xF0))
	fmt.Printf("Mapper: 0x%X\n", mapper)
	switch mapper {
	case 0x00:
		fallthrough
	case 0x40:
		fallthrough
	case 0x41:
		// NROM
		r = new(Nrom)
	case 0x01:
		// MMC1
		r = new(Mmc1)
	case 0x42:
		fallthrough
	case 0x02:
		// Unrom
		r = new(Unrom)
	case 0x03:
		// Cnrom
		r = new(Cnrom)
	default:
		// Unsupported
		return r, errors.New(fmt.Sprintf("Unsupported memory mapper: 0x%X", mapper))
	}

	return
}
