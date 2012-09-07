package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
)

func RespondToPress(e sdl.KeyboardEvent) {
    if e.Keysym.Sym == sdl.K_ESCAPE {
        running = false
    }
}

func Listen() {
	for {
		select {
		case ev := <-sdl.Events:
			switch e := ev.(type) {
			case sdl.KeyboardEvent:
                RespondToPress(e)
			}
		}
	}
}
