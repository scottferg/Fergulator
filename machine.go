package main

import (
	"fmt"
	"io/ioutil"
	"os"
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
	v, d := ppu.Init()
	controller.Init()

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

	video.Init(v, d, os.Args[1])
	defer video.Close()

	go JoypadListen()
	go video.Render()

	for running {
		cycles := cpu.Step()

		// 3 PPU cycles for each CPU cycle
		for i := 0; i < 3*cycles; i++ {
			ppu.Step()
		}
	}

	return
}
