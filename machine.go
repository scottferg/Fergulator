package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var (
	cycle = "0ns"
	//cycle = "559ns"
	//cycle = "50ms"
	programCounter = 0xC000
	clockspeed, _  = time.ParseDuration(cycle)
	running        = true

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
    fmt.Printf("0xEB6D: 0x%X\n", *Ram[0xEB6D])
}

func main() {
	Ram.Init()

	ppu.Init()
	cpu.Reset()

	cpu.P = 0x34

	v := make(chan Cpu)
	video.Init(v)

	//cpu.Verbose = true

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
			cpu.Step()
			v <- cpu

			// time.Sleep(clockspeed)
		}
	}

	fmt.Printf("Status was: 0x%X\n", cpu.P)

	return
}
