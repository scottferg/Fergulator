package main

import (
	"fmt"
	"github.com/go-gl/gl"
	"github.com/scottferg/Fergulator/nes"
	"github.com/scottferg/Go-SDL/gfx"
	"github.com/scottferg/Go-SDL/sdl"
	"log"
	"math"
	"os"
	"unsafe"
)

type Video struct {
	videoTick     <-chan []int16
	screen        *sdl.Surface
	fpsmanager    *gfx.FPSmanager
	prog          gl.Program
	texture       gl.Texture
	width, height int
	textureUni    gl.AttribLocation
	Fullscreen    bool
}

func createProgram(vertShaderSrc string, fragShaderSrc string) gl.Program {
	vertShader := loadShader(gl.VERTEX_SHADER, vertShaderSrc)
	fragShader := loadShader(gl.FRAGMENT_SHADER, fragShaderSrc)

	prog := gl.CreateProgram()

	prog.AttachShader(vertShader)
	prog.AttachShader(fragShader)
	prog.Link()

	if prog.Get(gl.LINK_STATUS) != gl.TRUE {
		log := prog.GetInfoLog()
		panic(fmt.Errorf("Failed to link program: %v", log))
	}

	return prog
}

func loadShader(shaderType gl.GLenum, source string) gl.Shader {
	shader := gl.CreateShader(shaderType)
	if err := gl.GetError(); err != gl.NO_ERROR {
		panic(fmt.Errorf("gl error: %v", err))
	}

	shader.Source(source)
	shader.Compile()

	if shader.Get(gl.COMPILE_STATUS) != gl.TRUE {
		log := shader.GetInfoLog()
		panic(fmt.Errorf("Failed to compile shader: %v, shader: %v", log, source))
	}

	return shader
}

func (v *Video) Init(t <-chan []int16, n string) {
	v.videoTick = t

	if sdl.Init(sdl.INIT_VIDEO|sdl.INIT_JOYSTICK|sdl.INIT_AUDIO) != 0 {
		log.Fatal(sdl.GetError())
	}

	v.screen = sdl.SetVideoMode(512, 512, 32,
		sdl.OPENGL|sdl.RESIZABLE|sdl.GL_DOUBLEBUFFER)
	if v.screen == nil {
		log.Fatal(sdl.GetError())
	}

	sdl.WM_SetCaption(fmt.Sprintf("Fergulator - %s", n), "")

	v.initGL()
	v.Reshape(int(v.screen.W), int(v.screen.H))

	v.fpsmanager = gfx.NewFramerate()
	v.fpsmanager.SetFramerate(60)

	return
}

func (v *Video) initGL() {
	if gl.Init() != 0 {
		panic(sdl.GetError())
	}

	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)

	v.prog = createProgram(vertShaderSrcDef, fragShaderSrcDef)
	posAttrib := v.prog.GetAttribLocation("vPosition")
	texCoordAttr := v.prog.GetAttribLocation("vTexCoord")
	paletteLoc := v.prog.GetUniformLocation("palette")
	v.textureUni = v.prog.GetAttribLocation("texture")

	v.texture = gl.GenTexture()
	gl.ActiveTexture(gl.TEXTURE0)
	v.texture.Bind(gl.TEXTURE_2D)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	v.prog.Use()
	posAttrib.EnableArray()
	texCoordAttr.EnableArray()

	paletteLoc.Uniform3iv(64, nes.ShaderPalette)

	vertVBO := gl.GenBuffer()
	vertVBO.Bind(gl.ARRAY_BUFFER)
	verts := []float32{-1.0, 1.0, -1.0, -1.0, 1.0, -1.0, 1.0, -1.0, 1.0, 1.0, -1.0, 1.0}
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*int(unsafe.Sizeof(verts[0])), &verts[0], gl.STATIC_DRAW)

	textCoorBuf := gl.GenBuffer()
	textCoorBuf.Bind(gl.ARRAY_BUFFER)
	texVerts := []float32{0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0}
	gl.BufferData(gl.ARRAY_BUFFER, len(texVerts)*int(unsafe.Sizeof(texVerts[0])), &texVerts[0], gl.STATIC_DRAW)

	posAttrib.AttribPointer(2, gl.FLOAT, false, 0, uintptr(0))
	texCoordAttr.AttribPointer(2, gl.FLOAT, false, 0, uintptr(0))
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
		w := (int)(math.Floor((float64)((256.0 / 240.0) * (float64)(height))))
		x_offset = (width - w) / 2
		width = w
	}

	v.width = width
	v.height = height

	gl.Viewport(x_offset, y_offset, width, height)
}

func quit_event() int {
	running = false
	return 0
}

func (v *Video) Render() {
	for running {
		select {
		case buf := <-v.videoTick:
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			v.prog.Use()

			gl.ActiveTexture(gl.TEXTURE0)
			v.texture.Bind(gl.TEXTURE_2D)

			gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 256, 256, 0, gl.RGBA,
				gl.UNSIGNED_SHORT_4_4_4_4, buf)

			gl.DrawArrays(gl.TRIANGLES, 0, 6)

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
			case sdl.KeyboardEvent:
				switch e.Keysym.Sym {
				case sdl.K_ESCAPE:
					running = false
				case sdl.K_r:
					// Trigger reset interrupt
					if e.Type == sdl.KEYDOWN {
						// cpu.RequestInterrupt(InterruptReset)
					}
				case sdl.K_l:
					if e.Type == sdl.KEYDOWN {
						nes.LoadGameState()
					}
				case sdl.K_s:
					if e.Type == sdl.KEYDOWN {
						nes.SaveGameState()
					}
				case sdl.K_i:
					if e.Type == sdl.KEYDOWN {
						nes.AudioEnabled = !nes.AudioEnabled
					}
				case sdl.K_1:
					if e.Type == sdl.KEYDOWN {
						v.ResizeEvent(256, 256)
					}
				case sdl.K_2:
					if e.Type == sdl.KEYDOWN {
						v.ResizeEvent(512, 512)
					}
				case sdl.K_3:
					if e.Type == sdl.KEYDOWN {
						v.ResizeEvent(768, 768)
					}
				case sdl.K_4:
					if e.Type == sdl.KEYDOWN {
						v.ResizeEvent(1024, 1024)
					}
				}

				switch e.Type {
				case sdl.KEYDOWN:
					nes.Pads[0].KeyDown(e, 0)
				case sdl.KEYUP:
					nes.Pads[0].KeyUp(e, 0)
				}
			}
		}
	}
}

func (v *Video) Close() {
	sdl.Quit()
}
