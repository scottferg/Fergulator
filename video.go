package main

import (
	"fmt"
	"github.com/banthar/gl"
	"github.com/scottferg/Go-SDL/gfx"
	"github.com/scottferg/Go-SDL/sdl"
	"log"
	"math"
)

type Video struct {
	tick       <-chan []uint32
	resize     chan [2]int
	fpsmanager *gfx.FPSmanager
	screen     *sdl.Surface
	tex        gl.Texture
	Fullscreen bool
}

func (v *Video) Init(t <-chan []uint32, n string) chan [2]int {
	v.tick = t
	v.resize = make(chan [2]int)

	if sdl.Init(sdl.INIT_VIDEO|sdl.INIT_JOYSTICK|sdl.INIT_AUDIO) != 0 {
		log.Fatal(sdl.GetError())
	}

	v.screen = sdl.SetVideoMode(512, 480, 32, sdl.OPENGL|sdl.RESIZABLE)

	if v.screen == nil {
		log.Fatal(sdl.GetError())
	}

	sdl.WM_SetCaption(fmt.Sprintf("Fergulator - %s", n), "")

	if gl.Init() != 0 {
		panic("gl error")
	}

	gl.Enable(gl.TEXTURE_2D)
	v.Reshape(int(v.screen.W), int(v.screen.H))

	v.tex = gl.GenTexture()

	v.fpsmanager = gfx.NewFramerate()
	v.fpsmanager.SetFramerate(70)

	joy = make([]*sdl.Joystick, sdl.NumJoysticks())

	for i := 0; i < sdl.NumJoysticks(); i++ {
		joy[i] = sdl.JoystickOpen(i)

		fmt.Println("-----------------")
		if joy[i] != nil {
			fmt.Printf("Joystick %d\n", i)
			fmt.Println("  Name: ", sdl.JoystickName(0))
			fmt.Println("  Number of Axes: ", joy[i].NumAxes())
			fmt.Println("  Number of Buttons: ", joy[i].NumButtons())
			fmt.Println("  Number of Balls: ", joy[i].NumBalls())
		} else {
			fmt.Println("  Couldn't open Joystick!")
		}
	}

	return v.resize
}

func (v *Video) ResizeEvent(w, h int) {
	v.screen = sdl.SetVideoMode(w, h, 32, sdl.OPENGL|sdl.RESIZABLE)
	v.Reshape(w, h)
}

func (v *Video) FullscreenEvent() {
	v.screen = sdl.SetVideoMode(1440, 900, 32, sdl.OPENGL|sdl.FULLSCREEN)
	v.Reshape(1440, 900)
}

func (v *Video) Reshape(width int, height int) {
	x_offset := 0
	y_offset := 0

	r := ((float64)(height)) / ((float64)(width))

	if r > 0.9375 { // Height taller than ratio
		h := (int)(math.Floor((float64)(0.9375 * (float64)(width))))
		y_offset = (height - h) / 2
		height = h
	} else if r < 0.9375 { // Width wider
		var scrW, scrH float64
		if ppu.OverscanEnabled {
			scrW = 240.0
			scrH = 224.0
		} else {
			scrW = 256.0
			scrH = 240.0
		}

		w := (int)(math.Floor((float64)((scrH / scrW) * (float64)(height))))
		x_offset = (width - w) / 2
		width = w
	}

	gl.Viewport(x_offset, y_offset, width, height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(-1, 1, -1, 1, -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Disable(gl.DEPTH_TEST)
}

func quit_event() int {
	running = false
	return 0
}

func (v *Video) Render() {
	for running {
		select {
		case dimensions := <-v.resize:
			v.ResizeEvent(dimensions[0], dimensions[1])
		case val := <-v.tick:
			slice := make([]uint8, len(val)*3)
			for i := 0; i < len(val); i = i + 1 {
				slice[i*3+0] = (uint8)((val[i] >> 16) & 0xff)
				slice[i*3+1] = (uint8)((val[i] >> 8) & 0xff)
				slice[i*3+2] = (uint8)((val[i]) & 0xff)
			}

			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			v.tex.Bind(gl.TEXTURE_2D)

			if ppu.OverscanEnabled {
				gl.TexImage2D(gl.TEXTURE_2D, 0, 3, 240, 224, 0, gl.RGB, gl.UNSIGNED_BYTE, slice)
			} else {
				gl.TexImage2D(gl.TEXTURE_2D, 0, 3, 256, 240, 0, gl.RGB, gl.UNSIGNED_BYTE, slice)
			}

			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

			gl.Begin(gl.QUADS)
			gl.TexCoord2f(0.0, 1.0)
			gl.Vertex3f(-1.0, -1.0, 0.0)
			gl.TexCoord2f(1.0, 1.0)
			gl.Vertex3f(1.0, -1.0, 0.0)
			gl.TexCoord2f(1.0, 0.0)
			gl.Vertex3f(1.0, 1.0, 0.0)
			gl.TexCoord2f(0.0, 0.0)
			gl.Vertex3f(-1.0, 1.0, 0.0)
			gl.End()

			if v.screen != nil {
				sdl.GL_SwapBuffers()
				v.fpsmanager.FramerateDelay()
			}
		}
	}
}

func (v *Video) Close() {
	sdl.Quit()
}
