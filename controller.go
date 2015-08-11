package main

import (
	"github.com/scottferg/Fergulator/nes"
	"github.com/scottferg/Go-SDL/sdl"
)

func GetKey(ev interface{}) int {
	if k, ok := ev.(sdl.KeyboardEvent); ok {
		switch k.Keysym.Sym {
		case sdl.K_z: // A
			fallthrough
		case sdl.K_g: // A
			return nes.ButtonA
		case sdl.K_x: // B
			fallthrough
		case sdl.K_i: // B
			return nes.ButtonB
		case sdl.K_RSHIFT: // Select
			fallthrough
		case sdl.K_n: // B
			return nes.ButtonSelect
		case sdl.K_RETURN: // Start
			fallthrough
		case sdl.K_o: // B
			return nes.ButtonStart
		case sdl.K_UP: // Up
			fallthrough
		case sdl.K_c: // B
			return nes.ButtonUp
		case sdl.K_DOWN: // Down
			fallthrough
		case sdl.K_d: // B
			return nes.ButtonDown
		case sdl.K_LEFT: // Left
			fallthrough
		case sdl.K_e: // B
			return nes.ButtonLeft
		case sdl.K_RIGHT: // Right
			fallthrough
		case sdl.K_f: // B
			return nes.ButtonRight
		}
	}

	return -1
}
