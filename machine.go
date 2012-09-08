package main

import (
	"fmt"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"io/ioutil"
	"os"
	"time"
)

var (
	cycle         = "559ns"
	clockspeed, _ = time.ParseDuration(cycle)
	running       = true
	breakpoint    = false

	cpu   Cpu
	ppu   Ppu
	rom   Mapper
	video Video
	joy   Controller

	io chan sdl.KeyboardEvent
)

func setResetVector() {
	high, _ := Ram.Read(0xFFFD)
	low, _ := Ram.Read(0xFFFC)

	ProgramCounter = (int(high) << 8) + int(low)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify a ROM file")
		return
	}

	Ram.Init()
	cpu.Init()
    v := ppu.Init()
	io = joy.Init()

	video.Init(v)
	defer video.Close()

	if contents, err := ioutil.ReadFile(os.Args[1]); err == nil {
		if rom, err = LoadRom(contents); err != nil {
			fmt.Println(err.Error())
			return
		}

		rom.Init(contents)
		setResetVector()
	} else {
		fmt.Println(err.Error())
		return
	}

	go JoypadListen()
	go RunCycles()
    video.Render()

	return
}

func RunCycles() {
	for running {
		cpu.Step()

		// 3 PPU cycles for each CPU cycle
		for i := 0; i < 3; i++ {
			ppu.Step()
		}

		// time.Sleep(clockspeed)
	}
}
