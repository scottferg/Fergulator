// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

// #include "gl.h"
import "C"

//void glColor3b (int8 red, int8 green, int8 blue)
func Color3b(red int8, green int8, blue int8) {
	C.glColor3b(C.GLbyte(red), C.GLbyte(green), C.GLbyte(blue))
}

//void glColor3bv (const int8 *v)
func Color3bv(v *[3]int8) {
	C.glColor3bv((*C.GLbyte)(&v[0]))
}

//void glColor3d (float64 red, float64 green, float64 blue)
func Color3d(red float64, green float64, blue float64) {
	C.glColor3d(C.GLdouble(red), C.GLdouble(green), C.GLdouble(blue))
}

//void glColor3dv (const float64 *v)
func Color3dv(v *[3]float64) {
	C.glColor3dv((*C.GLdouble)(&v[0]))
}

//void glColor3f (float32 red, float32 green, float32 blue)
func Color3f(red float32, green float32, blue float32) {
	C.glColor3f(C.GLfloat(red), C.GLfloat(green), C.GLfloat(blue))
}

//void glColor3fv (const float *v)
func Color3fv(v *[3]float32) {
	C.glColor3fv((*C.GLfloat)(&v[0]))
}

//void glColor3i (int red, int green, int blue)
func Color3i(red int, green int, blue int) {
	C.glColor3i(C.GLint(red), C.GLint(green), C.GLint(blue))
}

//void glColor3iv (const int *v)
func Color3iv(v *[3]int32) {
	C.glColor3iv((*C.GLint)(&v[0]))
}

//void glColor3s (int16 red, int16 green, int16 blue)
func Color3s(red int16, green int16, blue int16) {
	C.glColor3s(C.GLshort(red), C.GLshort(green), C.GLshort(blue))
}

//void glColor3sv (const int16 *v)
func Color3sv(v *[3]int16) {
	C.glColor3sv((*C.GLshort)(&v[0]))
}

//void glColor3ub (uint8 red, uint8 green, uint8 blue)
func Color3ub(red uint8, green uint8, blue uint8) {
	C.glColor3ub(C.GLubyte(red), C.GLubyte(green), C.GLubyte(blue))
}

//void glColor3ubv (const uint8 *v)
func Color3ubv(v *[3]uint8) {
	C.glColor3ubv((*C.GLubyte)(&v[0]))
}

//void glColor3ui (uint red, uint green, uint blue)
func Color3ui(red uint, green uint, blue uint) {
	C.glColor3ui(C.GLuint(red), C.GLuint(green), C.GLuint(blue))
}

//void glColor3uiv (const uint *v)
func Color3uiv(v *[3]uint32) {
	C.glColor3uiv((*C.GLuint)(&v[0]))
}

//void glColor3us (uint16 red, uint16 green, uint16 blue)
func Color3us(red uint16, green uint16, blue uint16) {
	C.glColor3us(C.GLushort(red), C.GLushort(green), C.GLushort(blue))
}

//void glColor3usv (const uint16 *v)
func Color3usv(v *[3]uint16) {
	C.glColor3usv((*C.GLushort)(&v[0]))
}

//void glColor4b (int8 red, int8 green, int8 blue, int8 alpha)
func Color4b(red int8, green int8, blue int8, alpha int8) {
	C.glColor4b(C.GLbyte(red), C.GLbyte(green), C.GLbyte(blue), C.GLbyte(alpha))
}

//void glColor4bv (const int8 *v)
func Color4bv(v *[4]int8) {
	C.glColor4bv((*C.GLbyte)(&v[0]))
}

//void glColor4d (float64 red, float64 green, float64 blue, float64 alpha)
func Color4d(red float64, green float64, blue float64, alpha float64) {
	C.glColor4d(C.GLdouble(red), C.GLdouble(green), C.GLdouble(blue), C.GLdouble(alpha))
}

//void glColor4dv (const float64 *v)
func Color4dv(v *[4]float64) {
	C.glColor4dv((*C.GLdouble)(&v[0]))
}

//void glColor4f (float32 red, float32 green, float32 blue, float32 alpha)
func Color4f(red float32, green float32, blue float32, alpha float32) {
	C.glColor4f(C.GLfloat(red), C.GLfloat(green), C.GLfloat(blue), C.GLfloat(alpha))
}

//void glColor4fv (const float *v)
func Color4fv(v *[4]float32) {
	C.glColor4fv((*C.GLfloat)(&v[0]))
}

//void glColor4i (int red, int green, int blue, int alpha)
func Color4i(red int, green int, blue int, alpha int) {
	C.glColor4i(C.GLint(red), C.GLint(green), C.GLint(blue), C.GLint(alpha))
}

//void glColor4iv (const int *v)
func Color4iv(v *[4]int32) {
	C.glColor4iv((*C.GLint)(&v[0]))
}

//void glColor4s (int16 red, int16 green, int16 blue, int16 alpha)
func Color4s(red int16, green int16, blue int16, alpha int16) {
	C.glColor4s(C.GLshort(red), C.GLshort(green), C.GLshort(blue), C.GLshort(alpha))
}

//void glColor4sv (const int16 *v)
func Color4sv(v *[4]int16) {
	C.glColor4sv((*C.GLshort)(&v[0]))
}

//void glColor4ub (uint8 red, uint8 green, uint8 blue, uint8 alpha)
func Color4ub(red uint8, green uint8, blue uint8, alpha uint8) {
	C.glColor4ub(C.GLubyte(red), C.GLubyte(green), C.GLubyte(blue), C.GLubyte(alpha))
}

//void glColor4ubv (const uint8 *v)
func Color4ubv(v *[4]uint8) {
	C.glColor4ubv((*C.GLubyte)(&v[0]))
}

//void glColor4ui (uint red, uint green, uint blue, uint alpha)
func Color4ui(red uint, green uint, blue uint, alpha uint) {
	C.glColor4ui(C.GLuint(red), C.GLuint(green), C.GLuint(blue), C.GLuint(alpha))
}

//void glColor4uiv (const uint *v)
func Color4uiv(v *[4]uint32) {
	C.glColor4uiv((*C.GLuint)(&v[0]))
}

//void glColor4us (uint16 red, uint16 green, uint16 blue, uint16 alpha)
func Color4us(red uint16, green uint16, blue uint16, alpha uint16) {
	C.glColor4us(C.GLushort(red), C.GLushort(green), C.GLushort(blue), C.GLushort(alpha))
}

//void glColor4usv (const uint16 *v)
func Color4usv(v *[4]uint16) {
	C.glColor4usv((*C.GLushort)(&v[0]))
}

//void glColorMask (bool red, bool green, bool blue, bool alpha)
func ColorMask(red bool, green bool, blue bool, alpha bool) {
	C.glColorMask(glBool(red), glBool(green), glBool(blue), glBool(alpha))
}

//void glColorMaterial (GLenum face, GLenum mode)
func ColorMaterial(face GLenum, mode GLenum) {
	C.glColorMaterial(C.GLenum(face), C.GLenum(mode))
}

//void glColorPointer (int size, GLenum type, int stride, const GLvoid *pointer)
func ColorPointer(size int, typ GLenum, stride int, pointer interface{}) {
	C.glColorPointer(C.GLint(size), C.GLenum(typ), C.GLsizei(stride),
		ptr(pointer))
}
