package main

import (
    "io/ioutil"
    "time"
    "fmt"
    "os"
)

var (
    programCounter = 0x8000
    clockspeed, _ = time.ParseDuration("559ns") // 1.79Mhz
    //clockspeed, _ = time.ParseDuration("200ms")
    running = true
)

func main() {
    cpu := new(Cpu)
    rom := new(Rom)
    video := new(Video)

    Ram.Init()
    cpu.Reset()

    v := make(chan Cpu)
    video.Init(v)

    defer video.Close()

    cpu.Verbose = false

    if contents, err := ioutil.ReadFile(os.Args[1]); err == nil {
        if err = rom.Init(contents); err != nil {
            fmt.Println(err.Error())
            return
        }

        go video.Render()

        for running {
            cpu.Step()
            v <- *cpu

            time.Sleep(clockspeed)
        }
    }

    return
}
