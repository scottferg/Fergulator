package main

import (
	"flag"
	"fmt"
	"github.com/scottferg/Fergulator/nes"
	"os"
	"runtime"
	"runtime/pprof"
)

var (
	running = true

	videoOut Video
	audioOut *Audio

	cpuprofile = flag.String("cprof", "", "write cpu profile to file")
)

func main() {
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

	audioOut = NewAudio()
	defer audioOut.Close()

	videoTick, gamename, err := nes.Init(audioOut.AppendSample, GetKey)
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
