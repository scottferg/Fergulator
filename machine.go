package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
)

var (
	cpuClockSpeed = 1789773
	running       = true
	audioEnabled  = true

	cpu      Cpu
	ppu      Ppu
	apu      Apu
	rom      Mapper
	videoOut Video
	audioOut *Audio
	pads     [2]*Controller

	totalCpuCycles int

	gamename       string
	saveStateFile  string
	batteryRamFile string

	cpuprofile = flag.String("cprof", "", "write cpu profile to file")
)

func setResetVector() {
	high, _ := Ram.Read(0xFFFD)
	low, _ := Ram.Read(0xFFFC)

	ProgramCounter = (uint16(high) << 8) + uint16(low)
}

func RunSystem() {
	var lastApuTick int
	var cycles int
	var flip int

	for running {
		cycles = cpu.Step()
		totalCpuCycles += cycles

		for i := 0; i < 3*cycles; i++ {
			ppu.Step()
		}

		for i := 0; i < cycles; i++ {
			apu.Step()
		}

		if audioEnabled {
			if totalCpuCycles-apu.LastFrameTick >= (cpuClockSpeed / 240) {
				apu.FrameSequencerStep()
				apu.LastFrameTick = totalCpuCycles
			}

			if totalCpuCycles-lastApuTick >= ((cpuClockSpeed / 44100) + flip) {
				apu.PushSample()
				lastApuTick = totalCpuCycles

				flip = (flip + 1) & 0x1
			}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify a ROM file")
		return
	}

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	} else if false {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	// Init the hardware, get communication channels
	// from the PPU and APU
	Ram.Init()
	cpu.Init()
	apu.Init()
	videoTick := ppu.Init()

	pads[0] = NewController()
	pads[1] = NewController()

	contents, err := ioutil.ReadFile(os.Args[len(os.Args)-1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if rom, err = LoadRom(contents); err != nil {
		fmt.Println(err.Error())
		return
	}

	// Set the game name for save states
	path := strings.Split(os.Args[1], "/")
	gamename = strings.Split(path[len(path)-1], ".")[0]
	saveStateFile = fmt.Sprintf(".%s.state", gamename)
	batteryRamFile = fmt.Sprintf(".%s.battery", gamename)

	if rom.BatteryBacked() {
		loadBatteryRam()
		defer saveBatteryFile()
	}

	setResetVector()

	videoOut.Init(videoTick, gamename)

	audioOut = NewAudio()
	defer audioOut.Close()

	// Main runloop, in a separate goroutine so that
	// the video rendering can happen on this one
	go RunSystem()

	// This needs to happen on the main thread for OSX
	runtime.LockOSThread()

	defer videoOut.Close()
	videoOut.Render()

	return
}
