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
	ButtonState [8]Word
	StrobeState int
	LastWrite   Word
	LastYAxis   int
	LastXAxis   int
}

func (c *Controller) SetJoypadAxisState(a, d int, v Word) {
	resetAxis := func(d int) {
		switch d {
		case 0:
			if c.LastYAxis != -1 {
				c.ButtonState[c.LastYAxis] = 0x40
			}
		case 1:
			if c.LastXAxis != -1 {
				c.ButtonState[c.LastXAxis] = 0x40
			}
		}
	}

	if a == 4 || a == 1 {
		switch d {
		case JoypadAxisUp: // Up
			resetAxis(0)
			c.ButtonState[4] = v
			c.LastYAxis = 4
		case JoypadAxisDown: // Down
			resetAxis(0)
			c.ButtonState[5] = v
			c.LastYAxis = 5
		default:
			resetAxis(0)
			c.LastYAxis = -1
		}
	} else if a == 3 || a == 0 {
		switch d {
		case JoypadAxisLeft: // Left
			resetAxis(1)
			c.ButtonState[6] = v
			c.LastXAxis = 6
		case JoypadAxisRight: // Right
			resetAxis(1)
			c.ButtonState[7] = v
			c.LastXAxis = 7
		default:
			resetAxis(1)
			c.LastXAxis = -1
		}
	}
}

func (c *Controller) SetJoypadButtonState(k int, v Word) {
	switch k {
	case JoypadButtonA: // A
		c.ButtonState[0] = v
	case JoypadButtonB: // B
		c.ButtonState[1] = v
	case JoypadButtonSelect: // Select
		c.ButtonState[2] = v
	case JoypadButtonStart: // Start
		c.ButtonState[3] = v
	}
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

func (c *Controller) AxisDown(a, d int) {
	c.SetJoypadAxisState(a, d, 0x41)
}

func (c *Controller) AxisUp(a, d int) {
	c.SetJoypadAxisState(a, d, 0x40)
}

func (c *Controller) ButtonDown(b int) {
	c.SetJoypadButtonState(b, 0x41)
}

func (c *Controller) ButtonUp(b int) {
	c.SetJoypadButtonState(b, 0x40)
}

func (c *Controller) KeyDown(e sdl.KeyboardEvent) {
	c.SetButtonState(e, 0x41)
}

func (c *Controller) KeyUp(e sdl.KeyboardEvent) {
	c.SetButtonState(e, 0x40)
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

func ReadInput(r chan [2]int, i chan int) {
	for {
		select {
		case ev := <-sdl.Events:
			switch e := ev.(type) {
			case sdl.ResizeEvent:
				r <- [2]int{int(e.W), int(e.H)}
			case sdl.QuitEvent:
				running = false
			case sdl.JoyAxisEvent:
				joy := int(e.Which)

				if joy > 0 {
					joy = 1
				}

				switch e.Value {
				// Same values for left/right
				case JoypadAxisUp:
					fallthrough
				case JoypadAxisDown:
					pads[joy].AxisDown(int(e.Axis), int(e.Value))
				default:
					pads[joy].AxisUp(int(e.Axis), int(e.Value))
				}
			case sdl.JoyButtonEvent:
				j := int(e.Which)

				if j > 0 {
					j = 1
				}

				switch joy[int(e.Which)].GetButton(int(e.Button)) {
				case 1:
					pads[j].ButtonDown(int(e.Button))
				case 0:
					pads[j].ButtonUp(int(e.Button))
				}
			case sdl.KeyboardEvent:
				switch e.Keysym.Sym {
				case sdl.K_ESCAPE:
					running = false
				case sdl.K_r:
					// Trigger reset interrupt
					if e.Type == sdl.KEYDOWN {
						cpu.RequestInterrupt(InterruptReset)
					}
				case sdl.K_l:
					if e.Type == sdl.KEYDOWN {
						i <- LoadState
					}
				case sdl.K_p:
					if e.Type == sdl.KEYDOWN {
						// Enable/disable scanline sprite limiter flag
						ppu.SpriteLimitEnabled = !ppu.SpriteLimitEnabled
					}
				case sdl.K_s:
					if e.Type == sdl.KEYDOWN {
						i <- SaveState
					}
				case sdl.K_o:
					if e.Type == sdl.KEYDOWN {
						ppu.OverscanEnabled = !ppu.OverscanEnabled
					}
				case sdl.K_1:
					if e.Type == sdl.KEYDOWN {
						r <- [2]int{256, 240}
					}
				case sdl.K_2:
					if e.Type == sdl.KEYDOWN {
						r <- [2]int{512, 480}
					}
				case sdl.K_3:
					if e.Type == sdl.KEYDOWN {
						r <- [2]int{768, 720}
					}
				case sdl.K_4:
					if e.Type == sdl.KEYDOWN {
						r <- [2]int{1024, 960}
					}
				}

				switch e.Type {
				case sdl.KEYDOWN:
					pads[0].KeyDown(e)
				case sdl.KEYUP:
					pads[0].KeyUp(e)
				}
			}
		}
	}
}
