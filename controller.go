package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
)

type Controller struct {
	ButtonState [8]Word
	StrobeState int
	LastWrite   Word
}

func (c *Controller) SetButtonState(k sdl.KeyboardEvent, v Word) {
	switch k.Keysym.Sym {
	case sdl.K_x: // B
		c.ButtonState[0] = v
	case sdl.K_z: // A
		c.ButtonState[1] = v
	case sdl.K_RSHIFT: // Select
		c.ButtonState[2] = v
	case sdl.K_RETURN: // Start
		c.ButtonState[3] = v
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
	} else if c.StrobeState == 19 {
		r = 0x1
	} else {
		r = 0x0
	}

	c.StrobeState++

	if c.StrobeState == 24 {
		c.StrobeState = 0x0
	}

	return r
}

func (c *Controller) Init() {
	for i := 0; i < len(c.ButtonState); i++ {
		c.ButtonState[i] = 0x40
	}
}

func JoypadListen() {
	for {
		select {
		case ev := <-sdl.Events:
			switch e := ev.(type) {
			case sdl.QuitEvent:
				running = false
			case sdl.KeyboardEvent:
				switch e.Keysym.Sym {
				case sdl.K_ESCAPE:
					running = false
				case sdl.K_r:
					// Trigger reset interrupt
					if e.Type == sdl.KEYDOWN {
						cpu.RequestInterrupt(InterruptReset)
					}
				case sdl.K_s:
					if e.Type == sdl.KEYDOWN {
						// Trigger reset interrupt
                        SaveState()
					}
				}

				switch e.Type {
				case sdl.KEYDOWN:
					controller.KeyDown(e)
				case sdl.KEYUP:
					controller.KeyUp(e)
				}
			}
		}
	}
}
