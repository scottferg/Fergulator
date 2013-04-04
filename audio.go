package main

import (
	"fmt"
	"github.com/timshannon/go-openal/openal"
)

type Audio struct {
	samples <-chan []byte
	source  openal.Source
	device  *openal.Device
	context *openal.Context
}

func NewAudio(s <-chan []byte) *Audio {

	a := Audio{
		samples: s,
	}

	return &a
}

func (a *Audio) Run() {
	a.device = openal.OpenDevice("")
	a.context = a.device.CreateContext()
	a.context.Activate()

	a.source = openal.NewSource()
	a.source.SetLooping(false)

    bufferIndex := 0
    buffers := []openal.Buffer{
        openal.NewBuffer(), 
        openal.NewBuffer(),
        openal.NewBuffer(),
    }

	for {
		v := <-a.samples
        fmt.Println(v)
		buffers[bufferIndex].SetData(openal.FormatMono8, v, 44100)

		a.source.QueueBuffer(buffers[bufferIndex])

        bufferIndex = bufferIndex + 1
        if bufferIndex == 3 {
            bufferIndex = 0
        }

        a.source.Play()

        for a.source.State() == openal.Playing {
            //loop long enough to let the wave file finish
        }

        a.source.Pause()

        buffers[bufferIndex] = a.source.UnqueueBuffer()
	}
}

func (a *Audio) Close() {
	fmt.Println("Closing!")
	a.source.Pause()
	a.context.Destroy()
}
