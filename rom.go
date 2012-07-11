package main

import (
    "errors"
    "fmt"
)

type Rom struct {
    PrgFlag Word
    ChrFlag Word
    Data    []byte
}

func (r *Rom) WriteRamBank(start int, length int, offset int) {
    for i := 0; i < length; i++ {
        Ram.Write(i + start, Word(r.Data[i + offset]))
    }
}

func (r *Rom) WriteVramBank(start int, length int, offset int) {
    for i := start; i < length; i++ {
        Vram[i] = Word(r.Data[i + offset])
    }
}

func (r *Rom) Init(rom []byte) error {
    if string(rom[0:3]) != "NES" {
        return errors.New("Invalid ROM file")

        if rom[3] != 0x1a {
            return errors.New("Invalid ROM file")
        }
    }

    r.PrgFlag = Word(rom[4])
    r.ChrFlag = Word(rom[5])

    fmt.Printf("PRG Flag: %X\n", r.PrgFlag)

    r.Data = rom[16:]

    // PRG ROM starts at byte 16
    r.WriteRamBank(0x8000, 0x4000, 0x0)

    switch r.PrgFlag {
    case 0x01:
        r.WriteRamBank(0xC000, 0x4000, 0x0)
    case 0x02:
        r.WriteRamBank(0xC000, 0x4000, 0x4000)
    }

    r.WriteVramBank(0x0000, 0x2000, 0x0)

    fmt.Println("ROM loaded!")

    return nil
}
