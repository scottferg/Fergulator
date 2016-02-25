// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

// #include "gl.h"
import "C"

// AttribLocation

type AttribLocation int

func (indx AttribLocation) Attrib1f(x float32) {
	C.glVertexAttrib1f(C.GLuint(indx), C.GLfloat(x))
}

func (indx AttribLocation) Attrib1fv(values *[1]float32) {
	C.glVertexAttrib1fv(C.GLuint(indx), (*C.GLfloat)(&values[0]))
}

func (indx AttribLocation) Attrib2f(x float32, y float32) {
	C.glVertexAttrib2f(C.GLuint(indx), C.GLfloat(x), C.GLfloat(y))
}

func (indx AttribLocation) Attrib2fv(values *[2]float32) {
	C.glVertexAttrib2fv(C.GLuint(indx), (*C.GLfloat)(&values[0]))
}

func (indx AttribLocation) Attrib3f(x float32, y float32, z float32) {
	C.glVertexAttrib3f(C.GLuint(indx), C.GLfloat(x), C.GLfloat(y), C.GLfloat(z))
}

func (indx AttribLocation) Attrib3fv(values *[3]float32) {
	C.glVertexAttrib3fv(C.GLuint(indx), (*C.GLfloat)(&values[0]))
}

func (indx AttribLocation) Attrib4f(x float32, y float32, z float32, w float32) {
	C.glVertexAttrib4f(C.GLuint(indx), C.GLfloat(x), C.GLfloat(y), C.GLfloat(z), C.GLfloat(w))
}

func (indx AttribLocation) Attrib4fv(values *[4]float32) {
	C.glVertexAttrib4fv(C.GLuint(indx), (*C.GLfloat)(&values[0]))
}

func (indx AttribLocation) AttribPointer(size uint, typ GLenum, normalized bool, stride int, pointer interface{}) {
	C.glVertexAttribPointer(C.GLuint(indx), C.GLint(size), C.GLenum(typ),
		glBool(normalized), C.GLsizei(stride), ptr(pointer))
}

func (indx AttribLocation) EnableArray() {
	C.glEnableVertexAttribArray(C.GLuint(indx))
}

func (indx AttribLocation) DisableArray() {
	C.glDisableVertexAttribArray(C.GLuint(indx))
}

func (indx AttribLocation) AttribDivisor(divisor int) {
	C.glVertexAttribDivisor(C.GLuint(indx), C.GLuint(divisor))
}
