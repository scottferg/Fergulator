package main

import (
	"fmt"
	"github.com/banthar/gl"
	"github.com/jteeuwen/glfw"
	"os"
)

type Video struct {
	tick  <-chan []int32
	debug <-chan []int32
	textures []gl.Texture
}

func reshape(width int, height int) {

	h := float64(height) / float64(width)

	gl.Viewport(0, 0, width, height)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Frustum(-1.0, 1.0, -h, h, 5.0, 60.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translatef(0.0, 0.0, -40.0)
}

func (v *Video) Init(t <-chan []int32, d <-chan []int32, n string) {
	v.tick = t
	v.debug = d

	if err := glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}

	if err := glfw.OpenWindow(300, 300, 0, 0, 0, 0, 0, 0, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}

	if gl.Init() != 0 {
		panic("gl error")
	}

	glfw.SetWindowTitle("gears")
	reshape(300, 300)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

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

			slice := val[:]

			tex := gl.GenTexture()

			tex.Bind(gl.TEXTURE_2D);
			gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, 256, 240, 0, gl.RGB, gl.INT, slice)

			gl.Begin(gl.QUADS)
			gl.TexCoord2f(0.0, 0.0)
			gl.Vertex3f(-1.0, -1.0,  1.0)
			gl.TexCoord2f(1.0, 0.0)
			gl.Vertex3f( 1.0, -1.0,  1.0)
			gl.TexCoord2f(1.0, 1.0)
			gl.Vertex3f( 1.0,  1.0,  1.0)
			gl.TexCoord2f(0.0, 1.0)
			gl.Vertex3f(-1.0,  1.0,  1.0)			
			gl.End()

			glfw.SwapBuffers()
		}
	}
}

func (v *Video) Close() {
	glfw.CloseWindow()
	glfw.Terminate()

}
