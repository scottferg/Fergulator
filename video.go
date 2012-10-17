package main

import (
	"fmt"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/gfx"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"github.com/banthar/gl"
	"log"
	"math"
)

type Video struct {
	tick       <-chan []uint32
	debug      <-chan []uint32
	fpsmanager *gfx.FPSmanager
	screen     *sdl.Surface
	tex        gl.Texture
	joy        []*sdl.Joystick
	Fullscreen bool
}

func (v *Video) Init(t <-chan []uint32, d <-chan []uint32, n string) {
	v.tick = t
	v.debug = d

	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
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

	v.joy = make([]*sdl.Joystick, sdl.NumJoysticks())

	for i := 0; i < sdl.NumJoysticks(); i++ {
		v.joy[i] = sdl.JoystickOpen(i)

		fmt.Println("-----------------")
		if v.joy[i] != nil {
			fmt.Printf("Joystick %d\n", i)
			fmt.Println("  Name: ", sdl.JoystickName(0))
			fmt.Println("  Number of Axes: ", v.joy[i].NumAxes())
			fmt.Println("  Number of Buttons: ", v.joy[i].NumButtons())
			fmt.Println("  Number of Balls: ", v.joy[i].NumBalls())
		} else {
			fmt.Println("  Couldn't open Joystick!")
		}
	}
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
		case ev := <-sdl.Events:
			// TODO: Should see if there's a way to do this
			// from another goroutine. Had to move it here for
			// the ResizeEvent
			switch e := ev.(type) {
			case sdl.ResizeEvent:
				v.ResizeEvent(int(e.W), int(e.H))
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
				joy := int(e.Which)

				if joy > 0 {
					joy = 1
				}

				switch v.joy[int(e.Which)].GetButton(int(e.Button)) {
				case 1:
					pads[joy].ButtonDown(int(e.Button))
				case 0:
					pads[joy].ButtonUp(int(e.Button))
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
						// Trigger reset interrupt
						LoadState()
					}
				case sdl.K_s:
					if e.Type == sdl.KEYDOWN {
						// Trigger reset interrupt
						SaveState()
					}
				case sdl.K_o:
					if e.Type == sdl.KEYDOWN {
						ppu.OverscanEnabled = !ppu.OverscanEnabled
					}
				case sdl.K_1:
					if e.Type == sdl.KEYDOWN {
						v.ResizeEvent(256, 240)
					}
				case sdl.K_2:
					if e.Type == sdl.KEYDOWN {
						v.ResizeEvent(512, 480)
					}
				case sdl.K_3:
					if e.Type == sdl.KEYDOWN {
						v.ResizeEvent(768, 720)
					}
				case sdl.K_4:
					if e.Type == sdl.KEYDOWN {
						v.ResizeEvent(1024, 960)
					}
				case sdl.K_f:
					if e.Type == sdl.KEYDOWN {
						if v.Fullscreen {
							v.ResizeEvent(512, 480)
						} else {
							v.FullscreenEvent()
						}

						v.Fullscreen = !v.Fullscreen
					}
				}

				switch e.Type {
				case sdl.KEYDOWN:
					pads[1].KeyDown(e)
				case sdl.KEYUP:
					pads[1].KeyUp(e)
				}
			}
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
