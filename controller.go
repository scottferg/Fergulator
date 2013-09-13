package main

import (
	"github.com/scottferg/Fergulator/nes"
	"github.com/scottferg/Go-SDL/sdl"
)

func GetKey(ev interface{}) int {
	if k, ok := ev.(sdl.KeyboardEvent); ok {
		switch k.Keysym.Sym {
		case sdl.K_z: // A
			return nes.ButtonA
		case sdl.K_x: // B
			return nes.ButtonB
		case sdl.K_RSHIFT: // Select
			return nes.ButtonSelect
		case sdl.K_RETURN: // Start
			return nes.ButtonStart
		case sdl.K_UP: // Up
			return nes.ButtonUp
		case sdl.K_DOWN: // Down
			return nes.ButtonDown
		case sdl.K_LEFT: // Left
			return nes.ButtonLeft
		case sdl.K_RIGHT: // Right
			return nes.ButtonRight
		}
	}

	return -1
}
