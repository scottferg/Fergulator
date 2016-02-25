// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

// #include "gl.h"
import "C"

//void glVertex2d (float64 x, float64 y)
func Vertex2d(x float64, y float64) {
	C.glVertex2d(C.GLdouble(x), C.GLdouble(y))
}

//void glVertex2dv (const float64 *v)
func Vertex2dv(v *[2]float64) {
	C.glVertex2dv((*C.GLdouble)(&v[0]))
}

//void glVertex2f (float32 x, float32 y)
func Vertex2f(x float32, y float32) {
	C.glVertex2f(C.GLfloat(x), C.GLfloat(y))
}

//void glVertex2fv (const float *v)
func Vertex2fv(v *[2]float32) {
	C.glVertex2fv((*C.GLfloat)(&v[0]))
}

//void glVertex2i (int x, int y)
func Vertex2i(x int, y int) {
	C.glVertex2i(C.GLint(x), C.GLint(y))
}

//void glVertex2iv (const int *v)
func Vertex2iv(v *[2]int32) {
	C.glVertex2iv((*C.GLint)(&v[0]))
}

//void glVertex2s (int16 x, int16 y)
func Vertex2s(x int16, y int16) {
	C.glVertex2s(C.GLshort(x), C.GLshort(y))
}

//void glVertex2sv (const int16 *v)
func Vertex2sv(v *[2]int16) {
	C.glVertex2sv((*C.GLshort)(&v[0]))
}

//void glVertex3d (float64 x, float64 y, float64 z)
func Vertex3d(x float64, y float64, z float64) {
	C.glVertex3d(C.GLdouble(x), C.GLdouble(y), C.GLdouble(z))
}

//void glVertex3dv (const float64 *v)
func Vertex3dv(v *[3]float64) {
	C.glVertex3dv((*C.GLdouble)(&v[0]))
}

//void glVertex3f (float32 x, float32 y, float32 z)
func Vertex3f(x float32, y float32, z float32) {
	C.glVertex3f(C.GLfloat(x), C.GLfloat(y), C.GLfloat(z))
}

//void glVertex3fv (const float *v)
func Vertex3fv(v *[3]float32) {
	C.glVertex3fv((*C.GLfloat)(&v[0]))
}

//void glVertex3i (int x, int y, int z)
func Vertex3i(x int, y int, z int) {
	C.glVertex3i(C.GLint(x), C.GLint(y), C.GLint(z))
}

//void glVertex3iv (const int *v)
func Vertex3iv(v *[3]int32) {
	C.glVertex3iv((*C.GLint)(&v[0]))
}

//void glVertex3s (int16 x, int16 y, int16 z)
func Vertex3s(x int16, y int16, z int16) {
	C.glVertex3s(C.GLshort(x), C.GLshort(y), C.GLshort(z))
}

//void glVertex3sv (const int16 *v)
func Vertex3sv(v *[3]int16) {
	C.glVertex3sv((*C.GLshort)(&v[0]))
}

//void glVertex4d (float64 x, float64 y, float64 z, float64 w)
func Vertex4d(x float64, y float64, z float64, w float64) {
	C.glVertex4d(C.GLdouble(x), C.GLdouble(y), C.GLdouble(z), C.GLdouble(w))
}

//void glVertex4dv (const float64 *v)
func Vertex4dv(v *[4]float64) {
	C.glVertex4dv((*C.GLdouble)(&v[0]))
}

//void glVertex4f (float32 x, float32 y, float32 z, float32 w)
func Vertex4f(x float32, y float32, z float32, w float32) {
	C.glVertex4f(C.GLfloat(x), C.GLfloat(y), C.GLfloat(z), C.GLfloat(w))
}

//void glVertex4fv (const float *v)
func Vertex4fv(v *[4]float32) {
	C.glVertex4fv((*C.GLfloat)(&v[0]))
}

//void glVertex4i (int x, int y, int z, int w)
func Vertex4i(x int, y int, z int, w int) {
	C.glVertex4i(C.GLint(x), C.GLint(y), C.GLint(z), C.GLint(w))
}

//void glVertex4iv (const int *v)
func Vertex4iv(v *[4]int32) {
	C.glVertex4iv((*C.GLint)(&v[0]))
}

//void glVertex4s (int16 x, int16 y, int16 z, int16 w)
func Vertex4s(x int16, y int16, z int16, w int16) {
	C.glVertex4s(C.GLshort(x), C.GLshort(y), C.GLshort(z), C.GLshort(w))
}

//void glVertex4sv (const int16 *v)
func Vertex4sv(v *[4]int16) {
	C.glVertex4sv((*C.GLshort)(&v[0]))
}

//void glVertexPointer (int size, GLenum type, int stride, const GLvoid *pointer)
func VertexPointer(size int, typ GLenum, stride int, pointer interface{}) {
	C.glVertexPointer(C.GLint(size), C.GLenum(typ), C.GLsizei(stride), ptr(pointer))
}
