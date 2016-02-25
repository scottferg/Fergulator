package sdl

import (
	"os"
)

func init() {
	if os.Getenv("SDL_VIDEODRIVER") == "" {
		os.Setenv("SDL_VIDEODRIVER", "x11")
	}
}
