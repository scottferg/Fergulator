package main

import (
	"fmt"
	"github.com/scottferg/Go-SDL/gfx"
	"github.com/go-gl/gl"
	"github.com/scottferg/Go-SDL/sdl"
	"log"
	"math"
	"os"
)

type Video struct {
	tick       <-chan []uint32
	screen     *sdl.Surface
	fpsmanager *gfx.FPSmanager
	tex        gl.Texture
	Fullscreen bool
}

func (v *Video) Init(t <-chan []uint32, n string) {
	v.tick = t

	if sdl.Init(sdl.INIT_VIDEO|sdl.INIT_JOYSTICK|sdl.INIT_AUDIO) != 0 {
		log.Fatal(sdl.GetError())
	}

	v.screen = sdl.SetVideoMode(512, 480, 32,
		sdl.OPENGL|sdl.RESIZABLE|sdl.GL_DOUBLEBUFFER)
	if v.screen == nil {
		log.Fatal(sdl.GetError())
	}

	sdl.WM_SetCaption(fmt.Sprintf("Fergulator - %s", n), "")

	if gl.Init() != 0 {
		panic(sdl.GetError())
	}

	gl.Enable(gl.TEXTURE_2D)
	v.Reshape(int(v.screen.W), int(v.screen.H))

	v.tex = gl.GenTexture()

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

	v.fpsmanager = gfx.NewFramerate()
	v.fpsmanager.SetFramerate(60)

	return
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
		case buf := <-v.tick:
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			v.tex.Bind(gl.TEXTURE_2D)

			if ppu.OverscanEnabled {
				gl.TexImage2D(gl.TEXTURE_2D, 0, 3, 240, 224, 0, gl.RGBA,
					gl.UNSIGNED_INT_8_8_8_8, buf)
			} else {
				gl.TexImage2D(gl.TEXTURE_2D, 0, 3, 256, 240, 0, gl.RGBA,
					gl.UNSIGNED_INT_8_8_8_8, buf)
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
		case ev := <-sdl.Events:
			switch e := ev.(type) {
			case sdl.ResizeEvent:
				v.ResizeEvent(int(e.W), int(e.H))
			case sdl.QuitEvent:
				os.Exit(0)
			case sdl.JoyAxisEvent:
				j := int(e.Which)

				index := j
				var offset int
				if j > 1 {
					offset = 8
					index = j % 2
				}

				switch e.Value {
				// Same values for left/right
				case JoypadAxisUp:
					fallthrough
				case JoypadAxisDown:
					pads[index].AxisDown(int(e.Axis), int(e.Value), offset)
				default:
					pads[index].AxisUp(int(e.Axis), int(e.Value), offset)
				}
			case sdl.JoyButtonEvent:
				j := int(e.Which)

				index := j
				var offset int
				if j > 1 {
					offset = 8
					index = j % 2
				}

				switch joy[j].GetButton(int(e.Button)) {
				case 1:
					pads[index].ButtonDown(int(e.Button), offset)
				case 0:
					pads[index].ButtonUp(int(e.Button), offset)
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
						LoadGameState()
					}
				case sdl.K_p:
					if e.Type == sdl.KEYDOWN {
						// Enable/disable scanline sprite limiter flag
						ppu.SpriteLimitEnabled = !ppu.SpriteLimitEnabled
					}
				case sdl.K_s:
					if e.Type == sdl.KEYDOWN {
						SaveGameState()
					}
				case sdl.K_o:
					if e.Type == sdl.KEYDOWN {
						ppu.OverscanEnabled = !ppu.OverscanEnabled
					}
				case sdl.K_i:
					if e.Type == sdl.KEYDOWN {
						audioEnabled = !audioEnabled
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
				}

				switch e.Type {
				case sdl.KEYDOWN:
					pads[0].KeyDown(e, 0)
				case sdl.KEYUP:
					pads[0].KeyUp(e, 0)
				}
			}
		}
	}
}

func (v *Video) Close() {
	sdl.Quit()
}
