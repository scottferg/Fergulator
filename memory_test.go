package main

import (
    "testing"
)

func TestMirroring(test *testing.T) {
    ppu1 := Word(0x1)
    ppu2 := Word(0x2)
    ppu3 := Word(0x3)

    Ram.Init()

    Ram.WriteMirrorable(0x1001, &ppu1)
    Ram.WriteMirrorable(0x1002, &ppu2)
    Ram.WriteMirrorable(0x1003, &ppu3)

    Ram.WriteMirrorable(0x2001, &ppu1)
    Ram.WriteMirrorable(0x2002, &ppu2)
    Ram.WriteMirrorable(0x2003, &ppu3)

    Ram.WriteMirrorable(0x3001, &ppu1)
    Ram.WriteMirrorable(0x3002, &ppu2)
    Ram.WriteMirrorable(0x3003, &ppu3)

    if val, _ := Ram.Read(0x1001); val != ppu1 {
        test.Errorf("Mirroring: 0x%X is not 0x%X\n", val, ppu1)
    }

    if val, _ := Ram.Read(0x2002); val != ppu2 {
        test.Errorf("Mirroring: 0x%X is not 0x%X\n", val, ppu2)
    }

    if val, _ := Ram.Read(0x3003); val != ppu3 {
        test.Errorf("Mirroring: 0x%X is not 0x%X\n", val, ppu3)
    }

    ppuT1 := Word(0xA)
    /*
    ppuT2 := Word(0xB)
    ppuT3 := Word(0xC)
    */

    Ram.WriteMirrorable(0x1001, &ppuT1)

    /*
    if val, _ := Ram.Read(0x1001); val != ppuT1 {
        test.Errorf("Mirroring: 0x%X is not 0x%X\n", val, ppuT1)
    }

    if val, _ := Ram.Read(0x2001); val != ppuT1 {
        test.Errorf("Mirroring: 0x%X is not 0x%X\n", val, ppuT1)
    }

    if val, _ := Ram.Read(0x3001); val != ppuT1 {
        test.Errorf("Mirroring: 0x%X is not 0x%X\n", val, ppuT1)
    }
    */
}
