package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var (
	cycle = "559ns"
	programCounter = 0x8000
	clockspeed, _  = time.ParseDuration(cycle)
	running        = true
    breakpoint     = false

	cpu   Cpu
	ppu   Ppu
	rom   Mapper
	video Video
)

func setResetVector() {
	high, _ := Ram.Read(0xFFFD)
	low, _ := Ram.Read(0xFFFC)

	programCounter = (int(high) << 8) + int(low)

    fmt.Printf("Setting reset: 0x%X\n", programCounter)
}

func main() {
	v := make(chan Nametable, 1000)
	video.Init(v)

	Ram.Init()

	cpu.Init()
	ppu.Init(v)

	cpu.P = 0x34

	// cpu.Verbose = true

	defer video.Close()

	if contents, err := ioutil.ReadFile(os.Args[1]); err == nil {
		if rom, err = LoadRom(contents); err != nil {
			fmt.Println(err.Error())
            // ppu.RenderNametable(0)
			return
		}

        rom.Init(contents)
		setResetVector()

		go video.Render()

		for running {
            if programCounter == 0xC2E2 && false {
                fmt.Println("Breakpoint!")
                breakpoint = true
            }

            if breakpoint {
                clockspeed, _  = time.ParseDuration("50ms")
            }

			cpu.Step()

            // 3 PPU cycles for each CPU cycle
            for i := 0; i < 3; i++ {
                ppu.Step()
            }

			// time.Sleep(clockspeed)
		}
	}

	return
}
