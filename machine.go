package main

import (
    "io/ioutil"
    "time"
    "fmt"
    "os"
)

var (
    cycle = "559ns"
    programCounter = 0xC000
    clockspeed, _ = time.ParseDuration(cycle)
    running = true

    cpu Cpu
    ppu Ppu
    rom Rom
    video Video

    breakpoint = 0xC7DC
    terminate  = 0xC7F3
)

func setResetVector() {
    high, _ := Ram.Read(0xFFFD)
    low, _ := Ram.Read(0xFFFC)

    fmt.Printf("Reset: 0x%X%X\n", high, low)

    programCounter = (int(high) << 8) + int(low)
}

func main() {
    Ram.Init()

    ppu.Init()
    cpu.Reset()

    cpu.P = 0x24

    v := make(chan Cpu)
    video.Init(v)

    cpu.Verbose = true

    defer video.Close()

    if contents, err := ioutil.ReadFile(os.Args[1]); err == nil {
        if err = rom.Init(contents); err != nil {
            fmt.Println(err.Error())
            return
        }

        //setResetVector()

        go video.Render()

loop:
        for running {
            cpu.Step()
            v <- cpu

            switch {
            case programCounter >= terminate:
                break loop;
            case programCounter >= breakpoint:
                clockspeed, _ = time.ParseDuration("3000ms")
            }

            time.Sleep(clockspeed)
        }
    }

    return
}
