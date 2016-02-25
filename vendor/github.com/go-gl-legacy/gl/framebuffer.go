// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

// #include "gl.h"
import "C"

// Framebuffer Objects
// TODO: implement GLsync stuff
type Framebuffer Object

// void glBindFramebuffer(GLenum target, GLuint framebuffer);
//
// Binds fb to target FRAMEBUFFER. To bind to a specific target, see BindTarget.
func (fb Framebuffer) Bind() {
	C.glBindFramebuffer(C.GLenum(FRAMEBUFFER), C.GLuint(fb))
}

// Binds fb to the specified target.
//
// See issue at github for why this function exists:
// http://github.com/go-gl/gl/issues/113
func (fb Framebuffer) BindTarget(target GLenum) {
	C.glBindFramebuffer(C.GLenum(target), C.GLuint(fb))
}

// Unbinds target FRAMEBUFFER. To unbind a a specific target, see UnbindTarget.
func (fb Framebuffer) Unbind() {
	C.glBindFramebuffer(C.GLenum(FRAMEBUFFER), 0)
}

// Unbinds the specified target.
//
// See issue at github for why this function exists:
// http://github.com/go-gl/gl/issues/113
func (fb Framebuffer) UnbindTarget(target GLenum) {
	C.glBindFramebuffer(C.GLenum(target), 0)
}

// void glBlitFramebuffer(GLint srcX0, GLint srcY0, GLint srcX1, GLint srcY1, GLint dstX0, GLint dstY0, GLint dstX1, GLint dstY1, GLbitfield mask, GLenum filter);
func BlitFramebuffer(srcX0, srcY0, srcX1, srcY1, dstX0, dstY0, dstX1, dstY1 int, mask GLbitfield, filter GLenum) {
	C.glBlitFramebuffer(C.GLint(srcX0), C.GLint(srcY0), C.GLint(srcX1), C.GLint(srcY1), C.GLint(dstX0), C.GLint(dstY0), C.GLint(dstX1), C.GLint(dstY1), C.GLbitfield(mask), C.GLenum(filter))
}

// GLenum glCheckFramebufferStatus(GLenum target);
func CheckFramebufferStatus(target GLenum) GLenum {
	return (GLenum)(C.glCheckFramebufferStatus(C.GLenum(target)))
}

// void glDeleteFramebuffers(GLsizei n, GLuint* framebuffers);
func (fb Framebuffer) Delete() {
	C.glDeleteFramebuffers(1, (*C.GLuint)(&fb))
}

func DeleteFramebuffers(bufs []Framebuffer) {
	if len(bufs) > 0 {
		C.glDeleteFramebuffers(C.GLsizei(len(bufs)), (*C.GLuint)(&bufs[0]))
	}
}

// void glFramebufferTexture1D(GLenum target, GLenum attachment, GLenum textarget, GLuint texture, GLint level);
func FramebufferTexture1D(target, attachment, textarget GLenum, texture Texture, level int) {
	C.glFramebufferTexture1D(C.GLenum(target), C.GLenum(attachment), C.GLenum(textarget), C.GLuint(texture), C.GLint(level))
}

// void glFramebufferTexture2D(GLenum target, GLenum attachment, GLenum textarget, GLuint texture, GLint level);
func FramebufferTexture2D(target, attachment, textarget GLenum, texture Texture, level int) {
	C.glFramebufferTexture2D(C.GLenum(target), C.GLenum(attachment), C.GLenum(textarget), C.GLuint(texture), C.GLint(level))
}

// void glFramebufferTexture3D(GLenum target, GLenum attachment, GLenum textarget, GLuint texture, GLint level, GLint layer);
func FramebufferTexture3D(target, attachment, textarget GLenum, texture Texture, level int, layer int) {
	C.glFramebufferTexture3D(C.GLenum(target), C.GLenum(attachment), C.GLenum(textarget), C.GLuint(texture), C.GLint(level), C.GLint(layer))
}

// void glFramebufferTextureLayer(GLenum target, GLenum attachment, GLuint texture, GLint level, GLint layer);
func FramebufferTextureLayer(target, attachment GLenum, texture Texture, level, layer int) {
	C.glFramebufferTextureLayer(C.GLenum(target), C.GLenum(attachment), C.GLuint(texture), C.GLint(level), C.GLint(layer))
}

// void glGenFramebuffers(GLsizei n, GLuint* ids);
func GenFramebuffer() Framebuffer {
	var b C.GLuint
	C.glGenFramebuffers(1, &b)
	return Framebuffer(b)
}

func GenFramebuffers(bufs []Framebuffer) {
	if len(bufs) > 0 {
		C.glGenFramebuffers(C.GLsizei(len(bufs)), (*C.GLuint)(&bufs[0]))
	}
}
