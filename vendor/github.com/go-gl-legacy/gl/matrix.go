// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

// #include "gl.h"
import "C"

//void glFrustum (float64 left, float64 right, float64 bottom, float64 top, float64 zNear, float64 zFar)
func Frustum(left float64, right float64, bottom float64, top float64, zNear float64, zFar float64) {
	C.glFrustum(C.GLdouble(left), C.GLdouble(right), C.GLdouble(bottom), C.GLdouble(top), C.GLdouble(zNear), C.GLdouble(zFar))
}

//void glLoadIdentity (void)
func LoadIdentity() {
	C.glLoadIdentity()
}

//void glLoadMatrixd (const float64 *m)
func LoadMatrixd(m *[16]float64) {
	C.glLoadMatrixd((*C.GLdouble)(&m[0]))
}

//void glLoadMatrixf (const float32 *m)
func LoadMatrixf(m *[16]float32) {
	C.glLoadMatrixf((*C.GLfloat)(&m[0]))
}

//void glMatrixMode (GLenum mode)
func MatrixMode(mode GLenum) {
	C.glMatrixMode(C.GLenum(mode))
}

//void glMultMatrixd (const float64 *m)
func MultMatrixd(m *[16]float64) {
	C.glMultMatrixd((*C.GLdouble)(&m[0]))
}

//void glMultMatrixf (const float32 *m)
func MultMatrixf(m *[16]float32) {
	C.glMultMatrixf((*C.GLfloat)(&m[0]))
}

//void glOrtho (float64 left, float64 right, float64 bottom, float64 top, float64 zNear, float64 zFar)
func Ortho(left float64, right float64, bottom float64, top float64, zNear float64, zFar float64) {
	C.glOrtho(C.GLdouble(left), C.GLdouble(right), C.GLdouble(bottom), C.GLdouble(top), C.GLdouble(zNear), C.GLdouble(zFar))
}

//void glPopMatrix (void)
func PopMatrix() {
	C.glPopMatrix()
}

//void glPushMatrix (void)
func PushMatrix() {
	C.glPushMatrix()
}

//void glRotated (float64 angle, float64 x, float64 y, float64 z)
func Rotated(angle float64, x float64, y float64, z float64) {
	C.glRotated(C.GLdouble(angle), C.GLdouble(x), C.GLdouble(y), C.GLdouble(z))
}

//void glRotatef (float32 angle, float32 x, float32 y, float32 z)
func Rotatef(angle float32, x float32, y float32, z float32) {
	C.glRotatef(C.GLfloat(angle), C.GLfloat(x), C.GLfloat(y), C.GLfloat(z))
}

//void glScaled (float64 x, float64 y, float64 z)
func Scaled(x float64, y float64, z float64) {
	C.glScaled(C.GLdouble(x), C.GLdouble(y), C.GLdouble(z))
}

//void glScalef (float32 x, float32 y, float32 z)
func Scalef(x float32, y float32, z float32) {
	C.glScalef(C.GLfloat(x), C.GLfloat(y), C.GLfloat(z))
}

//void glTranslated (float64 x, float64 y, float64 z)
func Translated(x float64, y float64, z float64) {
	C.glTranslated(C.GLdouble(x), C.GLdouble(y), C.GLdouble(z))
}

//void glTranslatef (float32 x, float32 y, float32 z)
func Translatef(x float32, y float32, z float32) {
	C.glTranslatef(C.GLfloat(x), C.GLfloat(y), C.GLfloat(z))
}
