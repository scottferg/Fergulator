// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

// #include "gl.h"
// GLuint workaroundGlGetUniformBlockIndex(GLuint program, const GLchar* uniformBlockName) {
//     return glGetUniformBlockIndex(program, uniformBlockName);
// }
import "C"
import "unsafe"

// Program

type Program Object

func CreateProgram() Program { return Program(C.glCreateProgram()) }

func (program Program) Delete() { C.glDeleteProgram(C.GLuint(program)) }

func (program Program) AttachShader(shader Shader) {
	C.glAttachShader(C.GLuint(program), C.GLuint(shader))
}

func (program Program) GetAttachedShaders() []Object {
	var len C.GLint
	C.glGetProgramiv(C.GLuint(program), C.GLenum(ACTIVE_UNIFORM_MAX_LENGTH), &len)

	objects := make([]Object, len)
	C.glGetAttachedShaders(C.GLuint(program), C.GLsizei(len), nil, *((**C.GLuint)(unsafe.Pointer(&objects))))
	return objects
}

func (program Program) DetachShader(shader Shader) {
	C.glDetachShader(C.GLuint(program), C.GLuint(shader))
}

func (program Program) TransformFeedbackVaryings(names []string, buffer_mode GLenum) {
	if len(names) == 0 {
		C.glTransformFeedbackVaryings(C.GLuint(program), 0, (**C.GLchar)(nil), C.GLenum(buffer_mode))
	} else {
		gl_names := make([]*C.GLchar, len(names))

		for i := range names {
			gl_names[i] = glString(names[i])
		}

		C.glTransformFeedbackVaryings(C.GLuint(program), C.GLsizei(len(gl_names)), &gl_names[0], C.GLenum(buffer_mode))

		for _, s := range gl_names {
			freeString(s)
		}
	}
}

func (program Program) Link() { C.glLinkProgram(C.GLuint(program)) }

func (program Program) Validate() { C.glValidateProgram(C.GLuint(program)) }

func (program Program) Use() { C.glUseProgram(C.GLuint(program)) }

func (program Program) Unuse() { C.glUseProgram(C.GLuint(0)) }

// Deprecated, please use program.Unuse()
func ProgramUnuse() { C.glUseProgram(C.GLuint(0)) }

func (program Program) GetInfoLog() string {
	var length C.GLint
	C.glGetProgramiv(C.GLuint(program), C.GLenum(INFO_LOG_LENGTH), &length)
	// length is buffer size including null character

	if length > 1 {
		log := C.malloc(C.size_t(length))
		defer C.free(log)
		C.glGetProgramInfoLog(C.GLuint(program), C.GLsizei(length), nil, (*C.GLchar)(log))
		return C.GoString((*C.char)(log))
	}
	return ""

}

func (program Program) Get(param GLenum) int {
	var rv C.GLint

	C.glGetProgramiv(C.GLuint(program), C.GLenum(param), &rv)
	return int(rv)
}

// glGetActiveUniform(GLuint program, GLuint index, GLsizei bufSize, GLsizei *length, GLint *size, GLenum *type, GLchar *name)
func (program Program) GetActiveUniform(index int) (
	Size int, Type GLenum, Name string) {
	// Maximum length of active uniform name in program
	bufSize := program.Get(ACTIVE_UNIFORM_MAX_LENGTH)
	nameBuf := C.malloc(C.size_t(bufSize))
	defer C.free(nameBuf)
	var size C.GLint
	C.glGetActiveUniform(
		C.GLuint(program),
		C.GLuint(index),
		C.GLsizei(bufSize),
		nil, // length == len(Name)
		&size,
		(*C.GLenum)(&Type),
		(*C.GLchar)(nameBuf))
	Name = C.GoString((*C.char)(nameBuf))
	Size = int(size)
	return
}

func (program Program) GetUniformiv(location UniformLocation, values []int32) {
	if len(values) == 0 {
		panic("Invalid values length")
	}
	// FIXME(jimt): This should really yield only one return value instead of using a slice.
	// http://www.opengl.org/sdk/docs/man/xhtml/glGetUniform.xml
	C.glGetUniformiv(C.GLuint(program), C.GLint(location), (*C.GLint)(&(values[0])))
}

func (program Program) GetUniformfv(location UniformLocation, values []float32) {
	if len(values) == 0 {
		panic("Invalid values length")
	}
	// FIXME(jimt): This should really yield only one return value instead of using a slice.
	// http://www.opengl.org/sdk/docs/man/xhtml/glGetUniform.xml
	C.glGetUniformfv(C.GLuint(program), C.GLint(location), (*C.GLfloat)(&(values[0])))
}

func (program Program) GetUniformLocation(name string) UniformLocation {

	cname := glString(name)
	defer freeString(cname)

	return UniformLocation(C.glGetUniformLocation(C.GLuint(program), cname))
}

func (program Program) GetUniformBlockIndex(name string) UniformBlockIndex {

	cname := glString(name)
	defer freeString(cname)

	// Workaround bug in GLEW < 1.8 where glGetUniformBlockIndex expects
	// a string of char instead of a string of GLchar.  We could ask everybody
	// to bump their version of GLEW, or we could add a bit of C code that
	// will silently cast GLchar into char.

	//return UniformBlockIndex(C.glGetUniformBlockIndex(C.GLuint(program), cname))
	return UniformBlockIndex(C.workaroundGlGetUniformBlockIndex(C.GLuint(program), cname))
}

func (program Program) UniformBlockBinding(index UniformBlockIndex, binding uint) {
	C.glUniformBlockBinding(C.GLuint(program), C.GLuint(index), C.GLuint(binding))
}

// glGetActiveAttrib(GLuint program, GLuint index, GLsizei bufSize, GLsizei *length, GLint *size, GLenum *type, GLchar *name)
func (program Program) GetActiveAttrib(index int) (
	Size int, Type GLenum, Name string) {
	// Maximum length of active uniform name in program
	bufSize := program.Get(ACTIVE_ATTRIBUTE_MAX_LENGTH)
	nameBuf := C.malloc(C.size_t(bufSize))
	defer C.free(nameBuf)
	var size C.GLint
	C.glGetActiveAttrib(
		C.GLuint(program),
		C.GLuint(index),
		C.GLsizei(bufSize),
		nil, // length == len(Name)
		&size,
		(*C.GLenum)(&Type),
		(*C.GLchar)(nameBuf))
	Name = C.GoString((*C.char)(nameBuf))
	Size = int(size)
	return
}

func (program Program) GetAttribLocation(name string) AttribLocation {

	cname := glString(name)
	defer freeString(cname)

	return AttribLocation(C.glGetAttribLocation(C.GLuint(program), cname))
}

func (program Program) BindAttribLocation(index AttribLocation, name string) {

	cname := glString(name)
	defer freeString(cname)

	C.glBindAttribLocation(C.GLuint(program), C.GLuint(index), cname)

}

func (program Program) BindFragDataLocation(colorNumber int, name string) {

	cname := glString(name)
	defer freeString(cname)

	C.glBindFragDataLocation(C.GLuint(program), C.GLuint(colorNumber), cname)
}

func (program Program) GetFragDataLocation(name string) int {

	cname := glString(name)
	defer freeString(cname)

	return int(C.glGetFragDataLocation(C.GLuint(program), cname))
}
