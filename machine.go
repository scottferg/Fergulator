package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var (
	//cycle = "559ns"
	cycle          = "0ns"
	programCounter = 0x8000
	clockspeed, _  = time.ParseDuration(cycle)
	running        = true

	cpu   Cpu
	ppu   Ppu
	rom   Rom
	video Video

	breakpoint = 0xDF58
	terminate  = 0xDF69
)

func setResetVector() {
	high, _ := Ram.Read(0xFFFD)
	low, _ := Ram.Read(0xFFFC)

	programCounter = (int(high) << 8) + int(low)
}

func main() {
	Ram.Init()

	ppu.Init()
	cpu.Reset()

	cpu.P = 0x34

	v := make(chan Cpu)
	video.Init(v)

	cpu.Verbose = true

	defer video.Close()

	if contents, err := ioutil.ReadFile(os.Args[1]); err == nil {
		if err = rom.Init(contents); err != nil {
			fmt.Println(err.Error())
			return
		}

		setResetVector()

		go video.Render()

	loop:
		for running {
			cpu.Step()
			v <- cpu

			switch {
			case programCounter == terminate:
				break loop
			case programCounter == breakpoint:
				clockspeed, _ = time.ParseDuration("3000ms")
			}

			time.Sleep(clockspeed)
		}
	}

	fmt.Printf("Status was: 0x%X\n", cpu.P)

	return
}
