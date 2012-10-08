package main

import (
	"fmt"
	"github.com/banthar/gl"
	"github.com/jteeuwen/glfw"
	"os"
)

type Video struct {
	tick  <-chan []uint32
	debug <-chan []uint32
	tex gl.Texture
}

func reshape(width int, height int) {

	gl.Viewport(0, 0, width, height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(-1, 1, -1, 1, -1, 1);
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Disable(gl.DEPTH_TEST)
}

func quit_event() int {
	running = false
	return 0
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

	glfw.SetWindowTitle("FergulatorGL")
	glfw.SetWindowSizeCallback(reshape)
	glfw.SetWindowCloseCallback(quit_event)	
	glfw.SetKeyCallback(KeyListener)
	reshape(512, 480)

	v.tex = gl.GenTexture()

}

func (v *Video) Render() {
	for {
		select {
		/*case d := <-v.debug:*/
		// 60hz
		// time.Sleep(16000000 * time.Nanosecond)
		// time.Sleep(12000000 * time.Nanosecond)
		case val := <-v.tick:
			// 60hz
			// time.Sleep(16000000 * time.Nanosecond)
			// time.Sleep(4000000 * time.Nanosecond)

			slice := make([]uint8, len(val) * 3)
			for i := 0; i < len(val); i = i+1 {
				slice[i * 3 + 0] = (uint8)((val[i] >> 16) & 0xff)
				slice[i * 3 + 1] = (uint8)((val[i] >> 8) & 0xff)
				slice[i * 3 + 2] = (uint8)((val[i]) & 0xff)
			} 

			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			v.tex.Bind(gl.TEXTURE_2D);
			gl.TexImage2D(gl.TEXTURE_2D, 0, 3, 256, 240, 0, gl.RGB, gl.UNSIGNED_BYTE, slice)
			gl.TexParameteri(gl.TEXTURE_2D,gl.TEXTURE_MIN_FILTER,gl.NEAREST);
			gl.TexParameteri(gl.TEXTURE_2D,gl.TEXTURE_MAG_FILTER,gl.NEAREST)

			gl.Begin(gl.QUADS)
			gl.TexCoord2f(0.0, 1.0)
			gl.Vertex3f(-1.0, -1.0,  0.0)
			gl.TexCoord2f(1.0, 1.0)
			gl.Vertex3f( 1.0, -1.0,  0.0)
			gl.TexCoord2f(1.0, 0.0)
			gl.Vertex3f( 1.0,  1.0,  0.0)
			gl.TexCoord2f(0.0, 0.0)
			gl.Vertex3f(-1.0,  1.0,  0.0)			
			gl.End()

			glfw.SwapBuffers()
		}
	}
}

func (v *Video) Close() {
	glfw.CloseWindow()
	glfw.Terminate()

}
