package main

import (
	"flag"
	"fmt"
	"github.com/scottferg/Fergulator/nes"
	"io/ioutil"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
)

var (
	running = true

	videoOut Video
	audioOut *Audio

	cpuprofile = flag.String("cprof", "", "write cpu profile to file")
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify a ROM file")
		return
	}

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Println(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	} else if false {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	contents, err := ioutil.ReadFile(os.Args[len(os.Args)-1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var rom RomInfo
	rom.Data = contents

	path := strings.Split(os.Args[1], "/")
	rom.Name = strings.Split(path[len(path)-1], ".")[0]
	rom.FilePath = path[:len(path)-2]

	audioOut = NewAudio()
	defer audioOut.Close()

	videoTick, err := nes.Init(rom, audioOut.AppendSample, GetKey)
	if err != nil {
		fmt.Println(err)
	}

	videoOut.Init(videoTick, gamename)

	// Main runloop, in a separate goroutine so that
	// the video rendering can happen on this one
	go nes.RunSystem()

	// This needs to happen on the main thread for OSX
	runtime.LockOSThread()

	defer videoOut.Close()
	videoOut.Render()

	return
}
