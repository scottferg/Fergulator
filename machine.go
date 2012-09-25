package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var (
	cycle         = "559ns"
	clockspeed, _ = time.ParseDuration(cycle)

	running = true

	cpu        Cpu
	ppu        Ppu
	rom        Mapper
	video      Video
	controller Controller

	statefile string
	gamename  string
)

func setResetVector() {
	high, _ := Ram.Read(0xFFFD)
	low, _ := Ram.Read(0xFFFC)

	ProgramCounter = (int(high) << 8) + int(low)
}

func LoadState() {
	fmt.Println("Loading state")

	state, err := ioutil.ReadFile(statefile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for i, v := range state[:0x2000] {
		Ram[i] = Word(v)
	}

	pchigh := int(state[0x2000])
	pclow := int(state[0x2001])

	ProgramCounter = (pchigh << 8) | pclow

	cpu.A = Word(state[0x2002])
	cpu.X = Word(state[0x2003])
	cpu.Y = Word(state[0x2004])
	cpu.P = Word(state[0x2005])
	cpu.StackPointer = Word(state[0x2006])

	// Sprite RAM
	for i, v := range state[0x2007:0x2107] {
		ppu.SpriteRam[i] = Word(v)
	}

	// Pattern VRAM
	for i, v := range state[0x2107:0x4107] {
		ppu.Vram[i] = Word(v)
	}

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

func SaveState() {
	fmt.Println("Saving state")
	buf := new(bytes.Buffer)

	// RAM
	for _, v := range Ram[:0x2000] {
		buf.WriteByte(byte(v))
	}

	// ProgramCounter
	// High then low
	buf.WriteByte(byte(ProgramCounter>>8) & 0xFF)
	buf.WriteByte(byte(ProgramCounter & 0xFF))

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
	for _, v := range ppu.Vram[:0x2000] {
		buf.WriteByte(byte(v))
	}

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

	if err := ioutil.WriteFile(statefile, buf.Bytes(), 0644); err != nil {
		panic(err.Error())
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify a ROM file")
		return
	}

	Ram.Init()
	cpu.Init()
	v, d := ppu.Init()
	controller.Init()

	if contents, err := ioutil.ReadFile(os.Args[1]); err == nil {
		if rom, err = LoadRom(contents); err != nil {
			fmt.Println(err.Error())
			return
		}

		// Set the game name for save states
		path := strings.Split(os.Args[1], "/")
		gamename = strings.Split(path[len(path)-1], ".")[0]
		statefile = fmt.Sprintf(".%s.state", gamename)

		setResetVector()
	} else {
		fmt.Println(err.Error())
		return
	}

	video.Init(v, d, gamename)
	defer video.Close()

	go JoypadListen()
	go video.Render()

	for running {
		cycles := cpu.Step()

		for i := 0; i < 3*cycles; i++ {
			ppu.Step()
		}
	}

	return
}
