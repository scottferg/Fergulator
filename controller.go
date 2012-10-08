package main

import (
	"github.com/jteeuwen/glfw"
)

const (
	A = 90
	B = 88
	SELECT = glfw.KeyRshift
	START = glfw.KeyEnter
	UP = glfw.KeyUp
	LEFT = glfw.KeyLeft
	DOWN = glfw.KeyDown
	RIGHT = glfw.KeyRight
	RESET = 82
	SAVE = 83
	LOAD = 76
)

type Controller struct {
	ButtonState [8]Word
	StrobeState int
	LastWrite   Word
}

func (c *Controller) SetButtonState(k int, v Word) {
	switch (k) {
	case A: // A
		c.ButtonState[0] = v
	case B: // B
		c.ButtonState[1] = v
	case SELECT: // Select
		c.ButtonState[2] = v
	case START: // Start
		c.ButtonState[3] = v
	case UP: // Up
		c.ButtonState[4] = v
	case DOWN: // Down
		c.ButtonState[5] = v
	case LEFT: // Left
		c.ButtonState[6] = v
	case RIGHT: // Right
		c.ButtonState[7] = v
	}
}

func (c *Controller) KeyDown(e int) {
	c.SetButtonState(e, 0x41)
}

func (c *Controller) KeyUp(e int) {
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

func KeyListener(key, state int) {
	if (state == glfw.KeyPress) {
		switch key  {
			case glfw.KeyEsc:
				running = false
			case RESET:
				cpu.RequestInterrupt(InterruptReset)
			case LOAD:
				LoadState()
			case SAVE:
				SaveState()
			default:
				controller.KeyDown(key)
		}
	} else {
		controller.KeyUp(key)
	}
}
