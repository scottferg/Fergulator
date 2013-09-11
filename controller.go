package main

import (
	"github.com/scottferg/Go-SDL/sdl"
)

var (
	joy []*sdl.Joystick
)

const (
	JoypadButtonA      = 1
	JoypadButtonB      = 2
	JoypadButtonStart  = 9
	JoypadButtonSelect = 8
	JoypadAxisUp       = -32768
	JoypadAxisDown     = 32767
	JoypadAxisLeft     = -32768
	JoypadAxisRight    = 32767
)

type Controller struct {
	ButtonState [16]Word
	StrobeState int
	LastWrite   Word
	LastYAxis   [2]int
	LastXAxis   [2]int
}

func (c *Controller) SetJoypadAxisState(a, d int, v Word, offset int) {
	resetAxis := func(d, i int) {
		switch d {
		case 0:
			if c.LastYAxis[i] != -1 {
				c.ButtonState[c.LastYAxis[i]] = 0x40
			}
		case 1:
			if c.LastXAxis[i] != -1 {
				c.ButtonState[c.LastXAxis[i]] = 0x40
			}
		}
	}

	index := 0
	if offset > 0 {
		index = 1
	}

	if a == 4 || a == 1 {
		switch d {
		case JoypadAxisUp: // Up
			resetAxis(0, index)
			c.ButtonState[4+offset] = v
			c.LastYAxis[index] = 4 + offset
		case JoypadAxisDown: // Down
			resetAxis(0, index)
			c.ButtonState[5+offset] = v
			c.LastYAxis[index] = 5 + offset
		default:
			resetAxis(0, index)
			c.LastYAxis[index] = -1
		}
	} else if a == 3 || a == 0 {
		switch d {
		case JoypadAxisLeft: // Left
			resetAxis(1, index)
			c.ButtonState[6+offset] = v
			c.LastXAxis[index] = 6 + offset
		case JoypadAxisRight: // Right
			resetAxis(1, index)
			c.ButtonState[7+offset] = v
			c.LastXAxis[index] = 7 + offset
		default:
			resetAxis(1, index)
			c.LastXAxis[index] = -1
		}
	}
}

func (c *Controller) SetJoypadButtonState(k int, v Word, offset int) {
	switch k {
	case JoypadButtonA: // A
		c.ButtonState[0+offset] = v
	case JoypadButtonB: // B
		c.ButtonState[1+offset] = v
	case JoypadButtonSelect: // Select
		c.ButtonState[2+offset] = v
	case JoypadButtonStart: // Start
		c.ButtonState[3+offset] = v
	}
}

func (c *Controller) SetButtonState(k sdl.KeyboardEvent, v Word, offset int) {
	switch k.Keysym.Sym {
	case sdl.K_z: // A
		c.ButtonState[0+offset] = v
	case sdl.K_x: // B
		c.ButtonState[1+offset] = v
	case sdl.K_RSHIFT: // Select
		c.ButtonState[2+offset] = v
	case sdl.K_RETURN: // Start
		c.ButtonState[3+offset] = v
	case sdl.K_UP: // Up
		c.ButtonState[4+offset] = v
	case sdl.K_DOWN: // Down
		c.ButtonState[5+offset] = v
	case sdl.K_LEFT: // Left
		c.ButtonState[6+offset] = v
	case sdl.K_RIGHT: // Right
		c.ButtonState[7+offset] = v
	}
}

func (c *Controller) AxisDown(a, d int, offset int) {
	c.SetJoypadAxisState(a, d, 0x41, offset)
}

func (c *Controller) AxisUp(a, d int, offset int) {
	c.SetJoypadAxisState(a, d, 0x40, offset)
}

func (c *Controller) ButtonDown(b int, offset int) {
	c.SetJoypadButtonState(b, 0x41, offset)
}

func (c *Controller) ButtonUp(b int, offset int) {
	c.SetJoypadButtonState(b, 0x40, offset)
}

func (c *Controller) KeyDown(e sdl.KeyboardEvent, offset int) {
	c.SetButtonState(e, 0x41, offset)
}

func (c *Controller) KeyUp(e sdl.KeyboardEvent, offset int) {
	c.SetButtonState(e, 0x40, offset)
}

func (c *Controller) Write(v Word) {
	if v == 0 && c.LastWrite == 1 {
		// 0x4016 writes manage strobe state for
		// both controllers. 0x4017 is reserved for
		// APU
		pads[0].StrobeState = 0
		pads[1].StrobeState = 0
	}

	c.LastWrite = v
}

func (c *Controller) Read() (r Word) {
	if c.StrobeState < 8 {
		r = ((c.ButtonState[c.StrobeState+8] & 1) << 1) | c.ButtonState[c.StrobeState]
	} else if c.StrobeState == 18 {
		r = 0x0
	} else if c.StrobeState == 19 {
		r = 0x0
	} else {
		r = 0x0
	}

	c.StrobeState++

	if c.StrobeState == 24 {
		c.StrobeState = 0
	}

	return
}

func NewController() *Controller {
	c := &Controller{}

	for i := 0; i < len(c.ButtonState); i++ {
		c.ButtonState[i] = 0x40
	}

	return c
}
