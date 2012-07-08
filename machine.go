package main

import (
    "io/ioutil"
    "time"
    "fmt"
    "os"
)

var (
    programCounter = 0x8000
    // clockspeed, _ = time.ParseDuration("559ns") // 1.79Mhz
    clockspeed, _ = time.ParseDuration("10ms")
    running = true
)

func main() {
    cpu := new(Cpu)
    rom := new(Rom)
    video := new(Video)

    v := make(chan string)
    video.Init(v)

    defer video.Close()

    go video.Render()

    cpu.Verbose = true

    if contents, err := ioutil.ReadFile(os.Args[1]); err == nil {
        if err = rom.Init(contents); err != nil {
            fmt.Println(err.Error())
        }
    }

    for running {
        cpu.Step()
        time.Sleep(clockspeed)

        v <- cpu.DumpState()
    }

    return
}
