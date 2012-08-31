package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var (
	cycle = "559ns"
	//cycle = "50ms"
	programCounter = 0x8000
	clockspeed, _  = time.ParseDuration(cycle)
	running        = true

    // Donkey Kong breakpoint: 0xC7BE

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
	Ram.Init()

	ppu.Init()
	cpu.Reset()

	cpu.P = 0x34

	v := make(chan Cpu)
	video.Init(v)

	// cpu.Verbose = true

	defer video.Close()

	if contents, err := ioutil.ReadFile(os.Args[1]); err == nil {
		if rom, err = LoadRom(contents); err != nil {
			fmt.Println(err.Error())
			return
		}

        rom.Init(contents)

		setResetVector()

		go video.Render()

		for running {
            if programCounter == 0xF50D {
                fmt.Println("Breakpoint!")
                s, _ := time.ParseDuration("99999s")
                time.Sleep(s)
            }

			cpu.Step()
			v <- cpu

			time.Sleep(clockspeed)
		}
	}

	return
}
