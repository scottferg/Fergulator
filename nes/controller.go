package nes

const (
	ButtonA = iota
	ButtonB
	ButtonSelect
	ButtonStart
	ButtonUp
	ButtonDown
	ButtonLeft
	ButtonRight
)

type getButtonFunc func(int) int

type Controller struct {
	ButtonState [16]word
	StrobeState int
	LastWrite   word
	LastYAxis   [2]int
	LastXAxis   [2]int
	getter      getButtonFunc
}

func (c *Controller) SetButtonState(button int, v word, offset int) {
	switch button {
	case ButtonA: // A
		c.ButtonState[0+offset] = v
	case ButtonB: // B
		c.ButtonState[1+offset] = v
	case ButtonSelect: // Select
		c.ButtonState[2+offset] = v
	case ButtonStart: // Start
		c.ButtonState[3+offset] = v
	case ButtonUp: // Up
		c.ButtonState[4+offset] = v
	case ButtonDown: // Down
		c.ButtonState[5+offset] = v
	case ButtonLeft: // Left
		c.ButtonState[6+offset] = v
	case ButtonRight: // Right
		c.ButtonState[7+offset] = v
	}
}

func (c *Controller) KeyDown(e int, offset int) {
	c.SetButtonState(c.getter(e), 0x41, offset)
}

func (c *Controller) KeyUp(e int, offset int) {
	c.SetButtonState(c.getter(e), 0x40, offset)
}

func (c *Controller) Write(v word) {
	if v == 0 && c.LastWrite == 1 {
		// 0x4016 writes manage strobe state for
		// both controllers. 0x4017 is reserved for
		// APU
		Pads[0].StrobeState = 0
		Pads[1].StrobeState = 0
	}

	c.LastWrite = v
}

func (c *Controller) Read() (r word) {
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

func NewController(getter getButtonFunc) *Controller {
	c := &Controller{
		getter: getter,
	}

	for i := 0; i < len(c.ButtonState); i++ {
		c.ButtonState[i] = 0x40
	}

	return c
}
