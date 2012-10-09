package main

import (
	"fmt"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/gfx"
	"github.com/banthar/gl"
	"github.com/jteeuwen/glfw"
	"math"
	"os"
	"runtime"
)

type Video struct {
	tick       <-chan []uint32
	debug      <-chan []uint32
	fpsmanager *gfx.FPSmanager
	tex        gl.Texture
}

func (v *Video) Init(t <-chan []uint32, d <-chan []uint32, n string) {
	v.tick = t
	v.debug = d

	if err := glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}

	if err := glfw.OpenWindow(512, 480, 0, 0, 0, 0, 0, 0, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}

	if gl.Init() != 0 {
		panic("gl error")
	}

	gl.Enable(gl.TEXTURE_2D)

	glfw.SetWindowTitle(fmt.Sprintf("Fergulator - %s", n))
	glfw.SetWindowSizeCallback(reshape)
	glfw.SetWindowCloseCallback(quit_event)
	glfw.SetKeyCallback(KeyListener)
	reshape(512, 480)

	v.tex = gl.GenTexture()

	v.fpsmanager = gfx.NewFramerate()
	v.fpsmanager.SetFramerate(70)
}

func reshape(width int, height int) {
	x_offset := 0
	y_offset := 0

	r := ((float64)(height)) / ((float64)(width))

	if r > 0.9375 { // Height taller than ratio
		h := (int)(math.Floor((float64)(0.9375 * (float64)(width))))
		y_offset = (height - h) / 2
		height = h
	} else if r < 0.9375 { // Width wider
		w := (int)(math.Floor((float64)((256.0 / 240.0) * (float64)(height))))
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
	runtime.LockOSThread()

	for {
		select {
		case val := <-v.tick:
			slice := make([]uint8, len(val)*3)
			for i := 0; i < len(val); i = i + 1 {
				slice[i*3+0] = (uint8)((val[i] >> 16) & 0xff)
				slice[i*3+1] = (uint8)((val[i] >> 8) & 0xff)
				slice[i*3+2] = (uint8)((val[i]) & 0xff)
			}

			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			v.tex.Bind(gl.TEXTURE_2D)
			gl.TexImage2D(gl.TEXTURE_2D, 0, 3, 256, 240, 0, gl.RGB, gl.UNSIGNED_BYTE, slice)
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

			glfw.SwapBuffers()
			v.fpsmanager.FramerateDelay()
		}
	}
}

func (v *Video) Close() {
	glfw.CloseWindow()
	glfw.Terminate()
}
