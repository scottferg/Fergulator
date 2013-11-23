package nes

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

var (
	cpuClockSpeed = 1789773
	AudioEnabled  = true

	totalCpuCycles int

	GameName       string
	SaveStateFile  string
	BatteryRamFile string

	Running bool
)

const (
	SaveState = iota
	LoadState
)

func LoadGameState() {
	fmt.Println("Loading state")

	state, err := ioutil.ReadFile(SaveStateFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for i, v := range state[:0x2000] {
		Ram[i] = Word(v)
	}

	pchigh := uint16(state[0x2000])
	pclow := uint16(state[0x2001])

	cpu.ProgramCounter = (pchigh << 8) | pclow

	cpu.A = Word(state[0x2002])
	cpu.X = Word(state[0x2003])
	cpu.Y = Word(state[0x2004])
	cpu.P = Word(state[0x2005])
	cpu.StackPointer = Word(state[0x2006])

	// Sprite RAM
	for i, v := range state[0x2007:0x2107] {
		ppu.SpriteRam[i] = Word(v)
	}

	/*
		// Pattern VRAM
		for i, v := range state[0x2107:0x4107] {
			ppu.Vram[i] = Word(v)
		}
	*/

	// Nametable VRAM
	for i, v := range state[0x4107:0x4507] {
		ppu.Nametables.LogicalTables[0][i] = Word(v)
	}
	for i, v := range state[0x4507:0x4907] {
		ppu.Nametables.LogicalTables[1][i] = Word(v)
	}
	for i, v := range state[0x4907:0x4D07] {
		ppu.Nametables.LogicalTables[2][i] = Word(v)
	}
	for i, v := range state[0x4D07:0x5107] {
		ppu.Nametables.LogicalTables[3][i] = Word(v)
	}

	// Palette RAM
	for i, v := range state[0x5107:0x5126] {
		ppu.PaletteRam[i] = Word(v)
	}
}

func SaveGameState() {
	fmt.Println("Saving state")
	buf := new(bytes.Buffer)

	// RAM
	for _, v := range Ram[:0x2000] {
		buf.WriteByte(byte(v))
	}

	// Cpu.ProgramCounter
	// High then low
	buf.WriteByte(byte(cpu.ProgramCounter>>8) & 0xFF)
	buf.WriteByte(byte(cpu.ProgramCounter & 0xFF))

	// CPU Registers
	buf.WriteByte(byte(cpu.A))
	buf.WriteByte(byte(cpu.X))
	buf.WriteByte(byte(cpu.Y))
	buf.WriteByte(byte(cpu.P))
	buf.WriteByte(byte(cpu.StackPointer))

	// Sprite RAM
	for _, v := range ppu.SpriteRam {
		buf.WriteByte(byte(v))
	}

	// Pattern VRAM
	/*
		for _, v := range ppu.Vram[:0x2000] {
			buf.WriteByte(byte(v))
		}
	*/

	// Nametable VRAM
	for _, v := range ppu.Nametables.LogicalTables[0] {
		buf.WriteByte(byte(v))
	}
	for _, v := range ppu.Nametables.LogicalTables[1] {
		buf.WriteByte(byte(v))
	}
	for _, v := range ppu.Nametables.LogicalTables[2] {
		buf.WriteByte(byte(v))
	}
	for _, v := range ppu.Nametables.LogicalTables[3] {
		buf.WriteByte(byte(v))
	}

	// Palette RAM
	for _, v := range ppu.PaletteRam {
		buf.WriteByte(byte(v))
	}

	if err := ioutil.WriteFile(SaveStateFile, buf.Bytes(), 0644); err != nil {
		panic(err.Error())
	}
}

func loadBatteryRam() {
	fmt.Println("Loading battery RAM")

	batteryRam, err := ioutil.ReadFile(BatteryRamFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for i, v := range batteryRam[:0x2000] {
		Ram[0x6000+i] = Word(v)
	}
}

func saveBatteryFile() {
	buf := new(bytes.Buffer)

	// Battery/Work RAM
	for _, v := range Ram[0x6000:0x7FFF] {
		buf.WriteByte(byte(v))
	}

	if err := ioutil.WriteFile(BatteryRamFile, buf.Bytes(), 0644); err != nil {
		panic(err.Error())
	}

	fmt.Println("Battery RAM saved to disk")
}

func Init(contents []byte, audioBuf func(int16), getter GetButtonFunc) (chan []int32, error) {
	// Init the hardware, get communication channels
	// from the PPU and APU
	Ram = NewMemory()
	cpu.Init()
	apu.Init(audioBuf)
	videoTick := ppu.Init()

	Pads[0] = NewController(getter)
	Pads[1] = NewController(getter)

	var err error
	if rom, err = LoadRom(contents); err != nil {
		return nil, err
	}

	if rom.BatteryBacked() {
		loadBatteryRam()
		defer saveBatteryFile()
	}

	cpu.SetResetVector()

	return videoTick, nil
}
