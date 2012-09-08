package main

import (
	"fmt"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
)

var (
	Sync chan sdl.KeyboardEvent
)

type Controller struct {
	ButtonState [8]Word
	StrobeState int
	LastWrite   Word
}

func (c *Controller) SetButtonState(k sdl.KeyboardEvent, v Word) {
	switch k.Keysym.Sym {
	case sdl.K_z: // A
		c.ButtonState[0] = v
	case sdl.K_x: // B
		c.ButtonState[1] = v
	case sdl.K_RSHIFT: // Select
		c.ButtonState[2] = v
	case sdl.K_RETURN: // Start
		c.ButtonState[4] = v
	case sdl.K_UP: // Up
		c.ButtonState[4] = v
	case sdl.K_DOWN: // Down
		c.ButtonState[5] = v
	case sdl.K_LEFT: // Left
		c.ButtonState[6] = v
	case sdl.K_RIGHT: // Right
		c.ButtonState[7] = v
	}
}

func (c *Controller) KeyDown(e sdl.KeyboardEvent) {
	c.SetButtonState(e, 0x41)
}

func (c *Controller) KeyUp(e sdl.KeyboardEvent) {
	c.SetButtonState(e, 0x40)
}

func (c *Controller) Write(v Word) {
	if v == 0 && c.LastWrite == 1 {
		c.StrobeState = 0
	}

	c.LastWrite = v
}

func (c *Controller) Read() (r Word) {
	if c.StrobeState < 8 {
		r = c.ButtonState[c.StrobeState]
	} else {
		r = 0x0
	}

	c.StrobeState++

	if c.StrobeState == 24 {
		c.StrobeState = 0x0
	}

	return r
}

func (c *Controller) Init() chan sdl.KeyboardEvent {
	Sync = make(chan sdl.KeyboardEvent)

	for i := 0; i < len(c.ButtonState); i++ {
		c.ButtonState[i] = 0x40
	}

	return Sync
}

func JoypadListen() {
	for {
		select {
		case ev := <-sdl.Events:
			switch e := ev.(type) {
			case sdl.KeyboardEvent:
                fmt.Println("Key!")
				Sync <- e
			}
		}
	}
}
