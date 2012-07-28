package main

type Ppu struct {
    Vram   [0xFFFF]Word
    SprRam [0x100]Word
}

func (p *Ppu) Init() {
    Ram.Write(0x2002, 0x80)
}
