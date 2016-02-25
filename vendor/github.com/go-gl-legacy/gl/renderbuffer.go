// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

// #include "gl.h"
import "C"

// Renderbuffer Objects

type Renderbuffer Object

// void glGenRenderbuffers(GLsizei n, GLuint *renderbuffers)
func GenRenderbuffer() Renderbuffer {
	var b C.GLuint
	C.glGenRenderbuffers(1, &b)
	return Renderbuffer(b)
}

// Fill slice with new renderbuffers
func GenRenderbuffers(bufs []Renderbuffer) {
	if len(bufs) > 0 {
		C.glGenRenderbuffers(C.GLsizei(len(bufs)), (*C.GLuint)(&bufs[0]))
	}
}

// void glBindRenderbuffer(GLenum target, GLuint renderbuffer);
func (rb Renderbuffer) Bind() {
	C.glBindRenderbuffer(C.GLenum(RENDERBUFFER), C.GLuint(rb))
}

// Unbind this texture
func (rb Renderbuffer) Unbind() {
	C.glBindRenderbuffer(C.GLenum(RENDERBUFFER), 0)
}

// void glDeleteRenderbuffers(GLsizei n, GLuint* renderbuffers);
func (rb Renderbuffer) Delete() {
	C.glDeleteRenderbuffers(1, (*C.GLuint)(&rb))
}

func DeleteRenderbuffers(bufs []Renderbuffer) {
	if len(bufs) > 0 {
		C.glDeleteRenderbuffers(C.GLsizei(len(bufs)), (*C.GLuint)(&bufs[0]))
	}
}

// void glGetRenderbufferParameteriv(GLenum target, GLenum pname, GLint* params);
func GetRenderbufferParameteriv(target, pname GLenum, params []int32) {
	if len(params) == 0 {
		panic("Invalid params size")
	}

	C.glGetRenderbufferParameteriv(C.GLenum(target), C.GLenum(pname), (*C.GLint)(&params[0]))
}

// void glRenderbufferStorage(GLenum target, GLenum internalformat, GLsizei width, GLsizei height);
func RenderbufferStorage(target, internalformat GLenum, width int, height int) {
	C.glRenderbufferStorage(C.GLenum(target), C.GLenum(internalformat), C.GLsizei(width), C.GLsizei(height))
}

// void glRenderbufferStorageMultisample(GLenum target, GLsizei samples, GLenum internalformat, GLsizei width, GLsizei height);
func RenderbufferStorageMultisample(target GLenum, samples int, internalformat GLenum, width, height int) {
	C.glRenderbufferStorageMultisample(C.GLenum(target), C.GLsizei(samples), C.GLenum(internalformat), C.GLsizei(width), C.GLsizei(height))
}

// GLsync glFramebufferRenderbuffer(GLenum target, GLenum attachment, GLenum renderbuffertarget, GLuint renderbuffer);
func (rb Renderbuffer) FramebufferRenderbuffer(target, attachment, renderbuffertarget GLenum) /* GLsync */ {
	// TODO: sync stuff.  return (GLsync)(C.glFramebufferRenderbuffer (C.GLenum(target), C.GLenum(attachment), C.GLenum(renderbuffertarget), C.GLuint(rb)))
	C.glFramebufferRenderbuffer(C.GLenum(target), C.GLenum(attachment), C.GLenum(renderbuffertarget), C.GLuint(rb))
}
