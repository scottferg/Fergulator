// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

// #include "gl.h"
import "C"
import "unsafe"

//bool glAreTexturesResident (GLsizei n, const uint *textures, bool *residences)
func AreTexturesResident(textures []uint, residences []bool) bool {
	sz := len(textures)
	if sz == 0 {
		return false
	}

	if sz != len(residences) {
		panic("Residences slice must be equal in length to textures slice.")
	}

	ret := C.glAreTexturesResident(
		C.GLsizei(sz),
		(*C.GLuint)(unsafe.Pointer(&textures[0])),
		(*C.GLboolean)(unsafe.Pointer(&residences[0])),
	)

	if ret == TRUE {
		return true
	}

	return false
}

func ActiveTexture(texture GLenum) { C.glActiveTexture(C.GLenum(texture)) }

// Texture

type Texture Object

// Create single texture object
func GenTexture() Texture {
	var b C.GLuint
	C.glGenTextures(1, &b)
	return Texture(b)
}

// Fill slice with new textures
func GenTextures(textures []Texture) {
	if len(textures) > 0 {
		C.glGenTextures(C.GLsizei(len(textures)), (*C.GLuint)(&textures[0]))
	}
}

// Delete texture object
func (texture Texture) Delete() {
	b := C.GLuint(texture)
	C.glDeleteTextures(1, &b)
}

// Delete all textures in slice
func DeleteTextures(textures []Texture) {
	if len(textures) > 0 {
		C.glDeleteTextures(C.GLsizei(len(textures)), (*C.GLuint)(&textures[0]))
	}
}

// Bind this texture as target
func (texture Texture) Bind(target GLenum) {
	C.glBindTexture(C.GLenum(target), C.GLuint(texture))
}

// Unbind this texture
func (texture Texture) Unbind(target GLenum) {
	C.glBindTexture(C.GLenum(target), 0)
}

//void glTexImage1D (GLenum target, int level, int internalformat, int width, int border, GLenum format, GLenum type, const GLvoid *pixels)
func TexImage1D(target GLenum, level int, internalformat int, width int, border int, format, typ GLenum, pixels interface{}) {
	C.glTexImage1D(C.GLenum(target), C.GLint(level), C.GLint(internalformat),
		C.GLsizei(width), C.GLint(border), C.GLenum(format), C.GLenum(typ),
		ptr(pixels))
}

//void glTexImage2D (GLenum target, int level, int internalformat, int width, int height, int border, GLenum format, GLenum type, const GLvoid *pixels)
func TexImage2D(target GLenum, level int, internalformat int, width int, height int, border int, format, typ GLenum, pixels interface{}) {
	C.glTexImage2D(C.GLenum(target), C.GLint(level), C.GLint(internalformat),
		C.GLsizei(width), C.GLsizei(height), C.GLint(border), C.GLenum(format),
		C.GLenum(typ), ptr(pixels))
}

//void glCompressedTexImage2D(	GLenum  target, GLint  level, GLenum internalformat, GLsizei width,
// GLsizei height, GLint border, GLsizei imagesize, const GLvoid * data )
func CompressedTexImage2D(target GLenum, level int, internalformat GLenum, width int, height int, border int, imagesize int, data interface{}) {
	C.glCompressedTexImage2D(C.GLenum(target), C.GLint(level), C.GLenum(internalformat),
		C.GLsizei(width), C.GLsizei(height), C.GLint(border), C.GLsizei(imagesize), ptr(data))
}

//void glGetCompressedTexImage( GLenum target, GLint lod, GLvoid *img )
func GetCompressedTexImage(target GLenum, lod int, data interface{}) {
	C.glGetCompressedTexImage(C.GLenum(target), C.GLint(lod), ptr(data))
}

//void glTexImage3D (GLenum target, int level, int internalformat, int width, int height, int depth, int border, GLenum format, GLenum type, const GLvoid *pixels)
func TexImage3D(target GLenum, level int, internalformat int, width, height, depth int, border int, format, typ GLenum, pixels interface{}) {
	C.glTexImage3D(C.GLenum(target), C.GLint(level), C.GLint(internalformat),
		C.GLsizei(width), C.GLsizei(height), C.GLsizei(depth), C.GLint(border),
		C.GLenum(format), C.GLenum(typ), ptr(pixels))
}

//void glTexBuffer (GLenum target, GLenum internalformat, GLuint buffer)
func TexBuffer(target, internalformat GLenum, buffer Buffer) {
	C.glTexBuffer(C.GLenum(target), C.GLenum(internalformat), C.GLuint(buffer))
}

//void glPixelMapfv (GLenum map, int mapsize, const float *values)
func PixelMapfv(map_ GLenum, mapsize int, values *float32) {
	C.glPixelMapfv(C.GLenum(map_), C.GLsizei(mapsize), (*C.GLfloat)(values))
}

//void glPixelMapuiv (GLenum map, int mapsize, const uint *values)
func PixelMapuiv(map_ GLenum, mapsize int, values *uint32) {
	C.glPixelMapuiv(C.GLenum(map_), C.GLsizei(mapsize), (*C.GLuint)(values))
}

//void glPixelMapusv (GLenum map, int mapsize, const uint16 *values)
func PixelMapusv(map_ GLenum, mapsize int, values *uint16) {
	C.glPixelMapusv(C.GLenum(map_), C.GLsizei(mapsize), (*C.GLushort)(values))
}

//void glTexSubImage1D (GLenum target, int level, int xoffset, int width, GLenum format, GLenum type, const GLvoid *pixels)
func TexSubImage1D(target GLenum, level int, xoffset int, width int, format, typ GLenum, pixels interface{}) {
	C.glTexSubImage1D(C.GLenum(target), C.GLint(level), C.GLint(xoffset),
		C.GLsizei(width), C.GLenum(format), C.GLenum(typ), ptr(pixels))
}

//void glTexSubImage2D (GLenum target, int level, int xoffset, int yoffset, int width, int height, GLenum format, GLenum type, const GLvoid *pixels)
func TexSubImage2D(target GLenum, level int, xoffset int, yoffset int, width int, height int, format, typ GLenum, pixels interface{}) {
	C.glTexSubImage2D(C.GLenum(target), C.GLint(level), C.GLint(xoffset),
		C.GLint(yoffset), C.GLsizei(width), C.GLsizei(height), C.GLenum(format),
		C.GLenum(typ), ptr(pixels))
}

//void glTexImage3D (GLenum target, int level, int xoffset, int yoffset, int zoffset, int width, int height, int depth, GLenum format, GLenum type, const GLvoid *pixels)
func TexSubImage3D(target GLenum, level int, xoffset, yoffset, zoffset, width, height, depth int, format, typ GLenum, pixels interface{}) {
	C.glTexSubImage3D(C.GLenum(target), C.GLint(level),
		C.GLint(xoffset), C.GLint(yoffset), C.GLint(zoffset),
		C.GLsizei(width), C.GLsizei(height), C.GLsizei(depth),
		C.GLenum(format), C.GLenum(typ), ptr(pixels))
}

//void glCopyTexImage1D (GLenum target, int level, GLenum internalFormat, int x, int y, int width, int border)
func CopyTexImage1D(target GLenum, level int, internalFormat GLenum, x int, y int, width int, border int) {
	C.glCopyTexImage1D(C.GLenum(target), C.GLint(level), C.GLenum(internalFormat), C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLint(border))
}

//void glCopyTexImage2D (GLenum target, int level, GLenum internalFormat, int x, int y, int width, int height, int border)
func CopyTexImage2D(target GLenum, level int, internalFormat GLenum, x int, y int, width int, height int, border int) {
	C.glCopyTexImage2D(C.GLenum(target), C.GLint(level), C.GLenum(internalFormat), C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height), C.GLint(border))
}

//void glCopyTexSubImage1D (GLenum target, int level, int xoffset, int x, int y, int width)
func CopyTexSubImage1D(target GLenum, level int, xoffset int, x int, y int, width int) {
	C.glCopyTexSubImage1D(C.GLenum(target), C.GLint(level), C.GLint(xoffset), C.GLint(x), C.GLint(y), C.GLsizei(width))
}

//void glCopyTexSubImage2D (GLenum target, int level, int xoffset, int yoffset, int x, int y, int width, int height)
func CopyTexSubImage2D(target GLenum, level int, xoffset int, yoffset int, x int, y int, width int, height int) {
	C.glCopyTexSubImage2D(C.GLenum(target), C.GLint(level), C.GLint(xoffset), C.GLint(yoffset), C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height))
}

// TODO 3D textures

//void glTexEnvf (GLenum target, GLenum pname, float32 param)
func TexEnvf(target GLenum, pname GLenum, param float32) {
	C.glTexEnvf(C.GLenum(target), C.GLenum(pname), C.GLfloat(param))
}

//void glTexEnvfv (GLenum target, GLenum pname, const float *params)
func TexEnvfv(target GLenum, pname GLenum, params []float32) {
	if len(params) != 1 && len(params) != 4 {
		panic("Invalid params slice length")
	}
	C.glTexEnvfv(C.GLenum(target), C.GLenum(pname), (*C.GLfloat)(&params[0]))
}

//void glTexEnvi (GLenum target, GLenum pname, int param)
func TexEnvi(target GLenum, pname GLenum, param int) {
	C.glTexEnvi(C.GLenum(target), C.GLenum(pname), C.GLint(param))
}

//void glTexEnviv (GLenum target, GLenum pname, const int *params)
func TexEnviv(target GLenum, pname GLenum, params []int32) {
	if len(params) != 1 && len(params) != 4 {
		panic("Invalid params slice length")
	}
	C.glTexEnviv(C.GLenum(target), C.GLenum(pname), (*C.GLint)(&params[0]))
}

//void glTexGend (GLenum coord, GLenum pname, float64 param)
func TexGend(coord GLenum, pname GLenum, param float64) {
	C.glTexGend(C.GLenum(coord), C.GLenum(pname), C.GLdouble(param))
}

//void glTexGendv (GLenum coord, GLenum pname, const float64 *params)
func TexGendv(coord GLenum, pname GLenum, params []float64) {
	if len(params) == 0 {
		panic("Invalid params slice length")
	}
	C.glTexGendv(C.GLenum(coord), C.GLenum(pname), (*C.GLdouble)(&params[0]))
}

//void glTexGenf (GLenum coord, GLenum pname, float32 param)
func TexGenf(coord GLenum, pname GLenum, param float32) {
	C.glTexGenf(C.GLenum(coord), C.GLenum(pname), C.GLfloat(param))
}

//void glTexGenfv (GLenum coord, GLenum pname, const float *params)
func TexGenfv(coord GLenum, pname GLenum, params []float32) {
	if len(params) == 0 {
		panic("Invalid params slice length")
	}
	C.glTexGenfv(C.GLenum(coord), C.GLenum(pname), (*C.GLfloat)(&params[0]))
}

//void glTexGeni (GLenum coord, GLenum pname, int param)
func TexGeni(coord GLenum, pname GLenum, param int) {
	C.glTexGeni(C.GLenum(coord), C.GLenum(pname), C.GLint(param))
}

//void glTexGeniv (GLenum coord, GLenum pname, const int *params)
func TexGeniv(coord GLenum, pname GLenum, params []int32) {
	if len(params) == 0 {
		panic("Invalid params slice length")
	}
	C.glTexGeniv(C.GLenum(coord), C.GLenum(pname), (*C.GLint)(&params[0]))
}

//void glTexParameterf (GLenum target, GLenum pname, float32 param)
func TexParameterf(target GLenum, pname GLenum, param float32) {
	C.glTexParameterf(C.GLenum(target), C.GLenum(pname), C.GLfloat(param))
}

//void glTexParameterfv (GLenum target, GLenum pname, const float *params)
func TexParameterfv(target GLenum, pname GLenum, params []float32) {
	if len(params) == 0 {
		panic("Invalid params slice length")
	}
	C.glTexParameterfv(C.GLenum(target), C.GLenum(pname), (*C.GLfloat)(&params[0]))
}

//void glTexParameteri (GLenum target, GLenum pname, int param)
func TexParameteri(target GLenum, pname GLenum, param int) {
	C.glTexParameteri(C.GLenum(target), C.GLenum(pname), C.GLint(param))
}

//void glTexParameteriv (GLenum target, GLenum pname, const int *params)
func TexParameteriv(target GLenum, pname GLenum, params []int32) {
	if len(params) == 0 {
		panic("Invalid params slice length")
	}
	C.glTexParameteriv(C.GLenum(target), C.GLenum(pname), (*C.GLint)(&params[0]))
}

//void glPrioritizeTextures (GLsizei n, const uint *textures, const GLclampf *priorities)
func PrioritizeTextures(n int, textures *uint32, priorities *GLclampf) {
	C.glPrioritizeTextures(C.GLsizei(n), (*C.GLuint)(textures), (*C.GLclampf)(priorities))
}

//void glGetTexEnvfv (GLenum target, GLenum pname, float *params)
func GetTexEnvfv(target GLenum, pname GLenum, params []float32) {
	if len(params) == 0 {
		panic("Invalid params slice length")
	}
	C.glGetTexEnvfv(C.GLenum(target), C.GLenum(pname), (*C.GLfloat)(&params[0]))
}

//void glGetTexEnviv (GLenum target, GLenum pname, int *params)
func GetTexEnviv(target GLenum, pname GLenum, params []int32) {
	if len(params) == 0 {
		panic("Invalid params slice length")
	}
	C.glGetTexEnviv(C.GLenum(target), C.GLenum(pname), (*C.GLint)(&params[0]))
}

//void glGetTexGendv (GLenum coord, GLenum pname, float64 *params)
func GetTexGendv(coord GLenum, pname GLenum, params []float64) {
	if len(params) == 0 {
		panic("Invalid params slice length")
	}
	C.glGetTexGendv(C.GLenum(coord), C.GLenum(pname), (*C.GLdouble)(&params[0]))
}

//void glGetTexGenfv (GLenum coord, GLenum pname, float *params)
func GetTexGenfv(coord GLenum, pname GLenum, params []float32) {
	if len(params) == 0 {
		panic("Invalid params slice length")
	}
	C.glGetTexGenfv(C.GLenum(coord), C.GLenum(pname), (*C.GLfloat)(&params[0]))
}

//void glGetTexGeniv (GLenum coord, GLenum pname, int *params)
func GetTexGeniv(coord GLenum, pname GLenum, params []int32) {
	if len(params) == 0 {
		panic("Invalid params slice length")
	}
	C.glGetTexGeniv(C.GLenum(coord), C.GLenum(pname), (*C.GLint)(&params[0]))
}

//void glGetTexImage (GLenum target, int level, GLenum format, GLenum type, GLvoid *pixels)
func GetTexImage(target GLenum, level int, format, typ GLenum, pixels interface{}) {
	C.glGetTexImage(C.GLenum(target), C.GLint(level), C.GLenum(format),
		C.GLenum(typ), ptr(pixels))
}

//void glGetTexLevelParameterfv (GLenum target, int level, GLenum pname, float *params)
func GetTexLevelParameterfv(target GLenum, level int, pname GLenum, params []float32) {
	if len(params) == 0 {
		panic("Invalid params slice length")
	}
	C.glGetTexLevelParameterfv(C.GLenum(target), C.GLint(level), C.GLenum(pname), (*C.GLfloat)(&params[0]))
}

//void glGetTexLevelParameteriv (GLenum target, int level, GLenum pname, int *params)
func GetTexLevelParameteriv(target GLenum, level int, pname GLenum, params []int32) {
	if len(params) == 0 {
		panic("Invalid params slice length")
	}
	C.glGetTexLevelParameteriv(C.GLenum(target), C.GLint(level), C.GLenum(pname), (*C.GLint)(&params[0]))
}

//void glGetTexParameterfv (GLenum target, GLenum pname, float *params)
func GetTexParameterfv(target GLenum, pname GLenum, params []float32) {
	if len(params) == 0 {
		panic("Invalid params slice length")
	}
	C.glGetTexParameterfv(C.GLenum(target), C.GLenum(pname), (*C.GLfloat)(&params[0]))
}

//void glGetTexParameteriv (GLenum target, GLenum pname, int *params)
func GetTexParameteriv(target GLenum, pname GLenum, params []int32) {
	if len(params) == 0 {
		panic("Invalid params slice length")
	}
	C.glGetTexParameteriv(C.GLenum(target), C.GLenum(pname), (*C.GLint)(&params[0]))
}

func GenerateMipmap(target GLenum) {
	C.glGenerateMipmap(C.GLenum(target))
}

//void glTexCoord1d (float64 s)
func TexCoord1d(s float64) {
	C.glTexCoord1d(C.GLdouble(s))
}

//void glTexCoord1dv (const float64 *v)
func TexCoord1dv(v *[1]float64) {
	C.glTexCoord1dv((*C.GLdouble)(&v[0]))
}

//void glTexCoord1f (float32 s)
func TexCoord1f(s float32) {
	C.glTexCoord1f(C.GLfloat(s))
}

//void glTexCoord1fv (const float *v)
func TexCoord1fv(v *[1]float32) {
	C.glTexCoord1fv((*C.GLfloat)(&v[0]))
}

//void glTexCoord1i (int s)
func TexCoord1i(s int) {
	C.glTexCoord1i(C.GLint(s))
}

//void glTexCoord1iv (const int *v)
func TexCoord1iv(v *[1]int32) {
	C.glTexCoord1iv((*C.GLint)(&v[0]))
}

//void glTexCoord1s (int16 s)
func TexCoord1s(s int16) {
	C.glTexCoord1s(C.GLshort(s))
}

//void glTexCoord1sv (const int16 *v)
func TexCoord1sv(v *[1]int16) {
	C.glTexCoord1sv((*C.GLshort)(&v[0]))
}

//void glTexCoord2d (float64 s, float64 t)
func TexCoord2d(s float64, t float64) {
	C.glTexCoord2d(C.GLdouble(s), C.GLdouble(t))
}

//void glTexCoord2dv (const float64 *v)
func TexCoord2dv(v *[2]float64) {
	C.glTexCoord2dv((*C.GLdouble)(&v[0]))
}

//void glTexCoord2f (float32 s, float32 t)
func TexCoord2f(s float32, t float32) {
	C.glTexCoord2f(C.GLfloat(s), C.GLfloat(t))
}

//void glTexCoord2fv (const float *v)
func TexCoord2fv(v *[2]float32) {
	C.glTexCoord2fv((*C.GLfloat)(&v[0]))
}

//void glTexCoord2i (int s, int t)
func TexCoord2i(s int, t int) {
	C.glTexCoord2i(C.GLint(s), C.GLint(t))
}

//void glTexCoord2iv (const int *v)
func TexCoord2iv(v *[2]int32) {
	C.glTexCoord2iv((*C.GLint)(&v[0]))
}

//void glTexCoord2s (int16 s, int16 t)
func TexCoord2s(s int16, t int16) {
	C.glTexCoord2s(C.GLshort(s), C.GLshort(t))
}

//void glTexCoord2sv (const int16 *v)
func TexCoord2sv(v *[2]int16) {
	C.glTexCoord2sv((*C.GLshort)(&v[0]))
}

//void glTexCoord3d (float64 s, float64 t, float64 r)
func TexCoord3d(s float64, t float64, r float64) {
	C.glTexCoord3d(C.GLdouble(s), C.GLdouble(t), C.GLdouble(r))
}

//void glTexCoord3dv (const float64 *v)
func TexCoord3dv(v *[3]float64) {
	C.glTexCoord3dv((*C.GLdouble)(&v[0]))
}

//void glTexCoord3f (float32 s, float32 t, float32 r)
func TexCoord3f(s float32, t float32, r float32) {
	C.glTexCoord3f(C.GLfloat(s), C.GLfloat(t), C.GLfloat(r))
}

//void glTexCoord3fv (const float *v)
func TexCoord3fv(v *[3]float32) {
	C.glTexCoord3fv((*C.GLfloat)(&v[0]))
}

//void glTexCoord3i (int s, int t, int r)
func TexCoord3i(s int, t int, r int) {
	C.glTexCoord3i(C.GLint(s), C.GLint(t), C.GLint(r))
}

//void glTexCoord3iv (const int *v)
func TexCoord3iv(v *[3]int32) {
	C.glTexCoord3iv((*C.GLint)(&v[0]))
}

//void glTexCoord3s (int16 s, int16 t, int16 r)
func TexCoord3s(s int16, t int16, r int16) {
	C.glTexCoord3s(C.GLshort(s), C.GLshort(t), C.GLshort(r))
}

//void glTexCoord3sv (const int16 *v)
func TexCoord3sv(v *[3]int16) {
	C.glTexCoord3sv((*C.GLshort)(&v[0]))
}

//void glTexCoord4d (float64 s, float64 t, float64 r, float64 q)
func TexCoord4d(s float64, t float64, r float64, q float64) {
	C.glTexCoord4d(C.GLdouble(s), C.GLdouble(t), C.GLdouble(r), C.GLdouble(q))
}

//void glTexCoord4dv (const float64 *v)
func TexCoord4dv(v *[4]float64) {
	C.glTexCoord4dv((*C.GLdouble)(&v[0]))
}

//void glTexCoord4f (float32 s, float32 t, float32 r, float32 q)
func TexCoord4f(s float32, t float32, r float32, q float32) {
	C.glTexCoord4f(C.GLfloat(s), C.GLfloat(t), C.GLfloat(r), C.GLfloat(q))
}

//void glTexCoord4fv (const float *v)
func TexCoord4fv(v *[4]float32) {
	C.glTexCoord4fv((*C.GLfloat)(&v[0]))
}

//void glTexCoord4i (int s, int t, int r, int q)
func TexCoord4i(s int, t int, r int, q int) {
	C.glTexCoord4i(C.GLint(s), C.GLint(t), C.GLint(r), C.GLint(q))
}

//void glTexCoord4iv (const int *v)
func TexCoord4iv(v *[4]int32) {
	C.glTexCoord4iv((*C.GLint)(&v[0]))
}

//void glTexCoord4s (int16 s, int16 t, int16 r, int16 q)
func TexCoord4s(s int16, t int16, r int16, q int16) {
	C.glTexCoord4s(C.GLshort(s), C.GLshort(t), C.GLshort(r), C.GLshort(q))
}

//void glTexCoord4sv (const int16 *v)
func TexCoord4sv(v *[4]int16) {
	C.glTexCoord4sv((*C.GLshort)(&v[0]))
}

//void glTexCoordPointer (int size, GLenum type, int stride, const GLvoid *pointer)
func TexCoordPointer(size int, typ GLenum, stride int, pointer interface{}) {
	C.glTexCoordPointer(C.GLint(size), C.GLenum(typ), C.GLsizei(stride),
		ptr(pointer))
}
