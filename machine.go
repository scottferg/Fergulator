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
	"time"
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
	audioOut Audio
	pads     [2]*Controller

	totalCpuCycles int

	gamename       string
	saveStateFile  string
	batteryRamFile string

	memprofile = flag.String("memprofile", "", "write memory profile to this file")
	cpuprofile = flag.String("cprof", "", "write cpu profile to file")
)

func setResetVector() {
	high, _ := Ram.Read(0xFFFD)
	low, _ := Ram.Read(0xFFFC)

	ProgramCounter = (uint16(high) << 8) + uint16(low)
}

func RunSystem(c <-chan int) {
	var lastApuTick int
	var cycles int
	var flip int

	for running {
		select {
		case s := <-c:
			switch s {
			case LoadState:
				LoadGameState()
			case SaveState:
				SaveGameState()
			}
		default:
			cycles = cpu.Step()
			totalCpuCycles += cycles

			if cpu.Timestamp >= 13000 {
				ppu.Step(cpu.Timestamp)
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
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please specify a ROM file")
		return
	}

	var shutdown chan bool

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()

		// We need to kill the server gracefully (i.e. not with Ctrl-C)
		// if we're profiling. 2 minutes is a nice amount of time to
		// gather samples.
		go func() {
			time.Sleep(120 * time.Second)
			shutdown <- true
		}()
	} else {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	// Init the hardware, get communication channels
	// from the PPU and APU
	Ram.Init()
	cpu.Init()
	videoTick := ppu.Init()
	audioTick := apu.Init()

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

	r, shutdown := videoOut.Init(videoTick, gamename)

	controllerInterrupt := make(chan int)

	audioOut := NewAudio(audioTick)
	defer audioOut.Close()

	// Main runloop, in a separate goroutine so that
	// the video rendering can happen on this one
	go audioOut.Run()
	go ReadInput(r, controllerInterrupt, shutdown)
	go RunSystem(controllerInterrupt)

	// This needs to happen on the main thread for OSX
	runtime.LockOSThread()

	defer videoOut.Close()
	videoOut.Render()

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}

	return
}
