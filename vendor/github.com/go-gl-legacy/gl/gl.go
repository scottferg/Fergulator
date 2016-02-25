// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

// #cgo darwin LDFLAGS: -framework OpenGL -lGLEW
// #cgo windows LDFLAGS: -lglew32 -lopengl32
// #cgo linux LDFLAGS: -lGLEW -lGL
// #cgo freebsd  CFLAGS: -I/usr/local/include
// #cgo freebsd LDFLAGS: -L/usr/local/lib -lglfw
// #include "gl.h"
// void SetGlewExperimental(GLboolean v) {  glewExperimental = v;  }
import "C"
import "unsafe"
import "reflect"

type GLenum C.GLenum
type GLbitfield C.GLbitfield
type GLclampf C.GLclampf
type GLclampd C.GLclampd

type Pointer unsafe.Pointer

// those types are left for compatibility reasons
type GLboolean C.GLboolean
type GLbyte C.GLbyte
type GLshort C.GLshort
type GLint C.GLint
type GLsizei C.GLsizei
type GLubyte C.GLubyte
type GLushort C.GLushort
type GLuint C.GLuint
type GLfloat C.GLfloat
type GLdouble C.GLdouble

// helpers

func glBool(v bool) C.GLboolean {
	if v {
		return 1
	}

	return 0
}

func goBool(v C.GLboolean) bool {
	return v != 0
}

func glString(s string) *C.GLchar { return (*C.GLchar)(C.CString(s)) }

func freeString(ptr *C.GLchar) { C.free(unsafe.Pointer(ptr)) }

func ptr(v interface{}) unsafe.Pointer {

	if v == nil {
		return unsafe.Pointer(nil)
	}

	rv := reflect.ValueOf(v)
	var et reflect.Value
	switch rv.Type().Kind() {
	case reflect.Uintptr:
		offset, _ := v.(uintptr)
		return unsafe.Pointer(offset)
	case reflect.Ptr:
		if rv.IsNil() {
			return unsafe.Pointer(nil)
		}
		et = rv.Elem()
	case reflect.Slice:
		if rv.IsNil() || rv.Len() == 0 {
			return unsafe.Pointer(nil)
		}
		et = rv.Index(0)
	default:
		panic("type must be a pointer, a slice, uintptr or nil")
	}

	return unsafe.Pointer(et.UnsafeAddr())
}

/*
uniformMatrix2fv
uniformMatrix2fv
uniformMatrix3fv
uniformMatrix3fv
uniformMatrix4fv
uniformMatrix4fv
*/

// Main

func BlendColor(red GLclampf, green GLclampf, blue GLclampf, alpha GLclampf) {
	C.glBlendColor(C.GLclampf(red), C.GLclampf(green), C.GLclampf(blue), C.GLclampf(alpha))
}

func BlendEquation(mode GLenum) { C.glBlendEquation(C.GLenum(mode)) }

func BlendEquationSeparate(modeRGB GLenum, modeAlpha GLenum) {
	C.glBlendEquationSeparate(C.GLenum(modeRGB), C.GLenum(modeAlpha))
}

func BlendFuncSeparate(srcRGB GLenum, dstRGB GLenum, srcAlpha GLenum, dstAlpha GLenum) {
	C.glBlendFuncSeparate(C.GLenum(srcRGB), C.GLenum(dstRGB), C.GLenum(srcAlpha), C.GLenum(dstAlpha))
}

func SampleCoverage(value GLclampf, invert bool) {
	C.glSampleCoverage(C.GLclampf(value), glBool(invert))
}

func StencilFuncSeparate(face GLenum, func_ GLenum, ref int, mask uint) {
	C.glStencilFuncSeparate(C.GLenum(face), C.GLenum(func_), C.GLint(ref), C.GLuint(mask))
}

func StencilMaskSeparate(face GLenum, mask uint) {
	C.glStencilMaskSeparate(C.GLenum(face), C.GLuint(mask))
}

func StencilOpSeparate(face GLenum, fail GLenum, zfail GLenum, zpass GLenum) {
	C.glStencilOpSeparate(C.GLenum(face), C.GLenum(fail), C.GLenum(zfail), C.GLenum(zpass))
}

//void glAccum (GLenum op, float32 value)
func Accum(op GLenum, value float32) {
	C.glAccum(C.GLenum(op), C.GLfloat(value))
}

//void glAlphaFunc (GLenum func, GLclampf ref)
func AlphaFunc(func_ GLenum, ref GLclampf) {
	C.glAlphaFunc(C.GLenum(func_), C.GLclampf(ref))
}

//void glArrayElement (int i)
func ArrayElement(i int) {
	C.glArrayElement(C.GLint(i))
}

//void glBegin (GLenum mode)
func Begin(mode GLenum) {
	C.glBegin(C.GLenum(mode))
}

//void glBitmap (GLsizei width, int height, float32 xorig, float32 yorig, float32 xmove, float32 ymove, const uint8 *bitmap)
func Bitmap(width int, height int, xorig float32, yorig float32, xmove float32, ymove float32, bitmap *uint8) {
	C.glBitmap(C.GLsizei(width), C.GLsizei(height), C.GLfloat(xorig), C.GLfloat(yorig), C.GLfloat(xmove), C.GLfloat(ymove), (*C.GLubyte)(bitmap))
}

//void glBlendFunc (GLenum sfactor, GLenum dfactor)
func BlendFunc(sfactor GLenum, dfactor GLenum) {
	C.glBlendFunc(C.GLenum(sfactor), C.GLenum(dfactor))
}

//void glCallList (uint list)
func CallList(list uint) {
	C.glCallList(C.GLuint(list))
}

//void glCallLists (GLsizei n, GLenum type, const GLvoid *lists)
func CallLists(n int, typ GLenum, lists interface{}) {
	C.glCallLists(C.GLsizei(n), C.GLenum(typ), ptr(lists))
}

//void glClear (GLbitfield mask)
func Clear(mask GLbitfield) {
	C.glClear(C.GLbitfield(mask))
}

//void glClearAccum (float32 red, float32 green, float32 blue, float32 alpha)
func ClearAccum(red float32, green float32, blue float32, alpha float32) {
	C.glClearAccum(C.GLfloat(red), C.GLfloat(green), C.GLfloat(blue), C.GLfloat(alpha))
}

//void glClearColor (GLclampf red, GLclampf green, GLclampf blue, GLclampf alpha)
func ClearColor(red GLclampf, green GLclampf, blue GLclampf, alpha GLclampf) {
	C.glClearColor(C.GLclampf(red), C.GLclampf(green), C.GLclampf(blue), C.GLclampf(alpha))
}

//void glClearDepth (GLclampd depth)
func ClearDepth(depth GLclampd) {
	C.glClearDepth(C.GLclampd(depth))
}

//void glClearIndex (float32 c)
func ClearIndex(c float32) {
	C.glClearIndex(C.GLfloat(c))
}

//void glClearStencil (int s)
func ClearStencil(s int) {
	C.glClearStencil(C.GLint(s))
}

//void glClipPlane (GLenum plane, const float64 *equation)
func ClipPlane(plane GLenum, equation []float64) {
	C.glClipPlane(C.GLenum(plane), (*C.GLdouble)(&equation[0]))
}

//void glCopyPixels (int x, int y, int width, int height, GLenum type)
func CopyPixels(x int, y int, width int, height int, type_ GLenum) {
	C.glCopyPixels(C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height), C.GLenum(type_))
}

//void glCullFace (GLenum mode)
func CullFace(mode GLenum) {
	C.glCullFace(C.GLenum(mode))
}

//void glDeleteLists (uint list, int range)
func DeleteLists(list uint, range_ int) {
	C.glDeleteLists(C.GLuint(list), C.GLsizei(range_))
}

//void glDepthFunc (GLenum func)
func DepthFunc(func_ GLenum) {
	C.glDepthFunc(C.GLenum(func_))
}

//void glDepthMask (bool flag)
func DepthMask(flag bool) {
	C.glDepthMask(glBool(flag))
}

//void glDepthRange (GLclampd zNear, GLclampd zFar)
func DepthRange(zNear GLclampd, zFar GLclampd) {
	C.glDepthRange(C.GLclampd(zNear), C.GLclampd(zFar))
}

//void glDisable (GLenum cap)
func Disable(cap GLenum) {
	C.glDisable(C.GLenum(cap))
}

//void glDisableClientState (GLenum array)
func DisableClientState(array GLenum) {
	C.glDisableClientState(C.GLenum(array))
}

//void glDrawArrays (GLenum mode, int first, int count)
func DrawArrays(mode GLenum, first int, count int) {
	C.glDrawArrays(C.GLenum(mode), C.GLint(first), C.GLsizei(count))
}

//void glDrawArraysInstanced(GLenum mode,  GLint first,  GLsizei count,  GLsizei primcount)
func DrawArraysInstanced(mode GLenum, first int, count, primcount int) {
	C.glDrawArraysInstanced(C.GLenum(mode), C.GLint(first), C.GLsizei(count), C.GLsizei(primcount))
}

//void glDrawBuffer (GLenum mode)
func DrawBuffer(mode GLenum) {
	C.glDrawBuffer(C.GLenum(mode))
}

// //void glDrawBuffers(GLsizei n, const GLenum *bufs)
func DrawBuffers(n int, bufs []GLenum) {
	C.glDrawBuffers(C.GLsizei(n), (*C.GLenum)(&bufs[0]))
}

//void glDrawElements (GLenum mode, int count, GLenum type, const GLvoid *indices)
func DrawElements(mode GLenum, count int, typ GLenum, indices interface{}) {
	C.glDrawElements(C.GLenum(mode), C.GLsizei(count), C.GLenum(typ),
		ptr(indices))
}

//void glDrawRangeElements (GLenum mode, int start, int end, int count, GLenum type, const GLvoid *indices)
func DrawRangeElements(mode GLenum, start, end uint, count int, typ GLenum, indices interface{}) {
	C.glDrawRangeElements(C.GLenum(mode), C.GLuint(start), C.GLuint(end), C.GLsizei(count), C.GLenum(typ),
		ptr(indices))
}

//void glDrawElementsInstanced(GLenum  mode,  GLsizei  count,  GLenum  type,  const void *  indices,  GLsizei  primcount)
func DrawElementsInstanced(mode GLenum, count int, typ GLenum, indices interface{}, primcount int) {
	C.glDrawElementsInstanced(C.GLenum(mode), C.GLsizei(count), C.GLenum(typ),
		ptr(indices), C.GLsizei(primcount))
}

//void glDrawElementsBaseVertex(GLenum mode, int count, GLenum type, GLvoid *indices, int basevertex)
func DrawElementsBaseVertex(mode GLenum, count int, typ GLenum, indices interface{}, basevertex int) {
	C.glDrawElementsBaseVertex(C.GLenum(mode), C.GLsizei(count),
		C.GLenum(typ), ptr(indices), C.GLint(basevertex))
}

//void glDrawPixels (GLsizei width, int height, GLenum format, GLenum type, const GLvoid *pixels)
func DrawPixels(width int, height int, format, typ GLenum, pixels interface{}) {
	C.glDrawPixels(C.GLsizei(width), C.GLsizei(height), C.GLenum(format),
		C.GLenum(typ), ptr(pixels))
}

//void glEdgeFlag (bool flag)
func EdgeFlag(flag bool) {
	C.glEdgeFlag(glBool(flag))
}

//void glEdgeFlagPointer (GLsizei stride, const GLvoid *pointer)
func EdgeFlagPointer(stride int, pointer unsafe.Pointer) {
	C.glEdgeFlagPointer(C.GLsizei(stride), pointer)
}

//void glEdgeFlagv (const bool *flag)
func EdgeFlagv(flag []bool) {
	if len(flag) > 0 {
		C.glEdgeFlagv((*C.GLboolean)(unsafe.Pointer(&flag[0])))
	}
}

//void glEnable (GLenum cap)
func Enable(cap GLenum) {
	C.glEnable(C.GLenum(cap))
}

//void glEnableClientState (GLenum array)
func EnableClientState(array GLenum) {
	C.glEnableClientState(C.GLenum(array))
}

//void glEnd (void)
func End() {
	C.glEnd()
}

//void glEndList (void)
func EndList() {
	C.glEndList()
}

//void glEvalCoord1d (float64 u)
func EvalCoord1d(u float64) {
	C.glEvalCoord1d(C.GLdouble(u))
}

//void glEvalCoord1dv (const float64 *u)
func EvalCoord1dv(u *float64) {
	C.glEvalCoord1dv((*C.GLdouble)(u))
}

//void glEvalCoord1f (float32 u)
func EvalCoord1f(u float32) {
	C.glEvalCoord1f(C.GLfloat(u))
}

//void glEvalCoord1fv (const float *u)
func EvalCoord1fv(u *[1]float32) {
	C.glEvalCoord1fv((*C.GLfloat)(&u[0]))
}

//void glEvalCoord2d (float64 u, float64 v)
func EvalCoord2d(u float64, v float64) {
	C.glEvalCoord2d(C.GLdouble(u), C.GLdouble(v))
}

//void glEvalCoord2dv (const float64 *u)
func EvalCoord2dv(u *float64) {
	C.glEvalCoord2dv((*C.GLdouble)(u))
}

//void glEvalCoord2f (float32 u, float32 v)
func EvalCoord2f(u float32, v float32) {
	C.glEvalCoord2f(C.GLfloat(u), C.GLfloat(v))
}

//void glEvalCoord2fv (const float *u)
func EvalCoord2fv(u *[2]float32) {
	C.glEvalCoord2fv((*C.GLfloat)(&u[0]))
}

//void glEvalMesh1 (GLenum mode, int i1, int i2)
func EvalMesh1(mode GLenum, i1 int, i2 int) {
	C.glEvalMesh1(C.GLenum(mode), C.GLint(i1), C.GLint(i2))
}

//void glEvalMesh2 (GLenum mode, int i1, int i2, int j1, int j2)
func EvalMesh2(mode GLenum, i1 int, i2 int, j1 int, j2 int) {
	C.glEvalMesh2(C.GLenum(mode), C.GLint(i1), C.GLint(i2), C.GLint(j1), C.GLint(j2))
}

//void glEvalPoint1 (int i)
func EvalPoint1(i int) {
	C.glEvalPoint1(C.GLint(i))
}

//void glEvalPoint2 (int i, int j)
func EvalPoint2(i int, j int) {
	C.glEvalPoint2(C.GLint(i), C.GLint(j))
}

//void glFeedbackBuffer (GLsizei size, GLenum type, float32 *buffer)
func FeedbackBuffer(size int, type_ GLenum, buffer *float32) {
	C.glFeedbackBuffer(C.GLsizei(size), C.GLenum(type_), (*C.GLfloat)(buffer))
}

//void glFinish (void)
func Finish() {
	C.glFinish()
}

//void glFlush (void)
func Flush() {
	C.glFlush()
}

//void glFogf (GLenum pname, float32 param)
func Fogf(pname GLenum, param float32) {
	C.glFogf(C.GLenum(pname), C.GLfloat(param))
}

//void glFogfv (GLenum pname, const float *params)
func Fogfv(pname GLenum, params []float32) {
	if len(params) == 0 {
		panic("Invalid params length")
	}
	C.glFogfv(C.GLenum(pname), (*C.GLfloat)(&params[0]))
}

//void glFogi (GLenum pname, int param)
func Fogi(pname GLenum, param int) {
	C.glFogi(C.GLenum(pname), C.GLint(param))
}

//void glFogiv (GLenum pname, const int *params)
func Fogiv(pname GLenum, params []int32) {
	if len(params) == 0 {
		panic("Invalid params length")
	}
	C.glFogiv(C.GLenum(pname), (*C.GLint)(&params[0]))
}

//void glFrontFace (GLenum mode)
func FrontFace(mode GLenum) {
	C.glFrontFace(C.GLenum(mode))
}

//uint glGenLists (GLsizei range)
func GenLists(range_ int) uint {
	return uint(C.glGenLists(C.GLsizei(range_)))
}

//void glGetBooleanv (GLenum pname, bool *params)
func GetBooleanv(pname GLenum, params []bool) {
	if len(params) == 0 {
		panic("Invalid params length")
	}
	C.glGetBooleanv(C.GLenum(pname), (*C.GLboolean)(unsafe.Pointer(&params[0])))
}

// Convenience function for GetBooleanv
func GetBoolean4(pname GLenum) (v0, v1, v2, v3 bool) {
	var values [4]C.GLboolean
	C.glGetBooleanv(C.GLenum(pname), &values[0])
	v0 = values[0] != 0
	v1 = values[1] != 0
	v2 = values[2] != 0
	v3 = values[3] != 0
	return
}

//void glGetClipPlane (GLenum plane, float64 *equation)
func GetClipPlane(plane GLenum, equation *float64) {
	C.glGetClipPlane(C.GLenum(plane), (*C.GLdouble)(equation))
}

//void glGetDoublev (GLenum pname, float64 *params)
func GetDoublev(pname GLenum, params []float64) {
	if len(params) == 0 {
		panic("Invalid params length")
	}
	C.glGetDoublev(C.GLenum(pname), (*C.GLdouble)(&params[0]))
}

// Convenience function for GetDoublev
func GetDouble4(pname GLenum) (v0, v1, v2, v3 float64) {
	var values [4]C.GLdouble
	C.glGetDoublev(C.GLenum(pname), &values[0])
	v0 = float64(values[0])
	v1 = float64(values[1])
	v2 = float64(values[2])
	v3 = float64(values[3])
	return
}

//GLenum glGetError (void)
func GetError() GLenum {
	return GLenum(C.glGetError())
}

//void glGetFloatv (GLenum pname, float *params)
func GetFloatv(pname GLenum, params []float32) {
	if len(params) == 0 {
		panic("Invalid params length")
	}
	C.glGetFloatv(C.GLenum(pname), (*C.GLfloat)(&params[0]))
}

// Convenience function for GetFloatv
func GetFloat4(pname GLenum) (v0, v1, v2, v3 float32) {
	var values [4]C.GLfloat
	C.glGetFloatv(C.GLenum(pname), &values[0])
	v0 = float32(values[0])
	v1 = float32(values[1])
	v2 = float32(values[2])
	v3 = float32(values[3])
	return
}

//void glGetIntegerv (GLenum pname, int *params)
func GetIntegerv(pname GLenum, params []int32) {
	if len(params) == 0 {
		panic("Invalid params length")
	}
	C.glGetIntegerv(C.GLenum(pname), (*C.GLint)(&params[0]))
}

// Convenience function for glGetIntegerv
func GetInteger4(pname GLenum) (v0, v1, v2, v3 int) {
	var values [4]C.GLint
	C.glGetIntegerv(C.GLenum(pname), &values[0])
	v0 = int(values[0])
	v1 = int(values[1])
	v2 = int(values[2])
	v3 = int(values[3])
	return
}

//void glGetLightfv (GLenum light, GLenum pname, float *params)
func GetLightfv(light GLenum, pname GLenum, params []float32) {
	if len(params) == 0 {
		panic("Invalid params length")
	}
	C.glGetLightfv(C.GLenum(light), C.GLenum(pname), (*C.GLfloat)(&params[0]))
}

//void glGetLightiv (GLenum light, GLenum pname, int *params)
func GetLightiv(light GLenum, pname GLenum, params []int32) {
	if len(params) == 0 {
		panic("Invalid params length")
	}
	C.glGetLightiv(C.GLenum(light), C.GLenum(pname), (*C.GLint)(&params[0]))
}

//void glGetMapdv (GLenum target, GLenum query, float64 *v)
func GetMapdv(target GLenum, query GLenum, v []float64) {
	if len(v) == 0 {
		panic("Invalid slice length")
	}
	C.glGetMapdv(C.GLenum(target), C.GLenum(query), (*C.GLdouble)(&v[0]))
}

//void glGetMapfv (GLenum target, GLenum query, float *v)
func GetMapfv(target GLenum, query GLenum, v []float32) {
	if len(v) == 0 {
		panic("Invalid slice length")
	}
	C.glGetMapfv(C.GLenum(target), C.GLenum(query), (*C.GLfloat)(&v[0]))
}

//void glGetMapiv (GLenum target, GLenum query, int *v)
func GetMapiv(target GLenum, query GLenum, v []int32) {
	if len(v) == 0 {
		panic("Invalid slice length")
	}
	C.glGetMapiv(C.GLenum(target), C.GLenum(query), (*C.GLint)(&v[0]))
}

//void glGetMaterialfv (GLenum face, GLenum pname, float *params)
func GetMaterialfv(face GLenum, pname GLenum, params []float32) {
	if len(params) == 0 {
		panic("Invalid params length")
	}
	C.glGetMaterialfv(C.GLenum(face), C.GLenum(pname), (*C.GLfloat)(&params[0]))
}

//void glGetMaterialiv (GLenum face, GLenum pname, int *params)
func GetMaterialiv(face GLenum, pname GLenum, params []int32) {
	if len(params) == 0 {
		panic("Invalid params length")
	}
	C.glGetMaterialiv(C.GLenum(face), C.GLenum(pname), (*C.GLint)(&params[0]))
}

//void glGetPixelMapfv (GLenum map, float *values)
func GetPixelMapfv(map_ GLenum, values []float32) {
	if len(values) == 0 {
		panic("Invalid values length")
	}
	C.glGetPixelMapfv(C.GLenum(map_), (*C.GLfloat)(&values[0]))
}

//void glGetPixelMapuiv (GLenum map, uint *values)
func GetPixelMapuiv(map_ GLenum, values *uint32) {
	C.glGetPixelMapuiv(C.GLenum(map_), (*C.GLuint)(values))
}

//void glGetPixelMapusv (GLenum map, uint16 *values)
func GetPixelMapusv(map_ GLenum, values *uint16) {
	C.glGetPixelMapusv(C.GLenum(map_), (*C.GLushort)(values))
}

//void glGetPointerv (GLenum pname, GLvoid* *params)
func GetPointerv(pname GLenum, params []unsafe.Pointer) {
	if len(params) == 0 {
		panic("Invalid params length")
	}
	C.glGetPointerv(C.GLenum(pname), &params[0])
}

//void glGetPolygonStipple (uint8 *mask)
func GetPolygonStipple(mask *uint8) {
	C.glGetPolygonStipple((*C.GLubyte)(mask))
}

//const uint8 * glGetString (GLenum name)
func GetString(name GLenum) string {
	s := unsafe.Pointer(C.glGetString(C.GLenum(name)))
	return C.GoString((*C.char)(s))
}

//void glHint (GLenum target, GLenum mode)
func Hint(target GLenum, mode GLenum) {
	C.glHint(C.GLenum(target), C.GLenum(mode))
}

//void glIndexMask (uint mask)
func IndexMask(mask uint) {
	C.glIndexMask(C.GLuint(mask))
}

//void glIndexPointer (GLenum type, int stride, const GLvoid *pointer)
func IndexPointer(typ GLenum, stride int, pointer interface{}) {
	C.glIndexPointer(C.GLenum(typ), C.GLsizei(stride), ptr(pointer))
}

//void glIndexd (float64 c)
func Indexd(c float64) {
	C.glIndexd(C.GLdouble(c))
}

//void glIndexdv (const float64 *c)
func Indexdv(c *[1]float64) {
	C.glIndexdv((*C.GLdouble)(&c[0]))
}

//void glIndexf (float32 c)
func Indexf(c float32) {
	C.glIndexf(C.GLfloat(c))
}

//void glIndexfv (const float32 *c)
func Indexfv(c *[1]float32) {
	C.glIndexfv((*C.GLfloat)(&c[0]))
}

//void glIndexi (int c)
func Indexi(c int) {
	C.glIndexi(C.GLint(c))
}

//void glIndexiv (const int *c)
func Indexiv(c *[1]int32) {
	C.glIndexiv((*C.GLint)(&c[0]))
}

//void glIndexs (int16 c)
func Indexs(c int16) {
	C.glIndexs(C.GLshort(c))
}

//void glIndexsv (const int16 *c)
func Indexsv(c *[1]int16) {
	C.glIndexsv((*C.GLshort)(&c[0]))
}

//void glIndexub (uint8 c)
func Indexub(c uint8) {
	C.glIndexub(C.GLubyte(c))
}

//void glIndexubv (const uint8 *c)
func Indexubv(c *[1]uint8) {
	C.glIndexubv((*C.GLubyte)(&c[0]))
}

//void glInitNames (void)
func InitNames() {
	C.glInitNames()
}

//void glInterleavedArrays (GLenum format, int stride, const GLvoid *pointer)
func InterleavedArrays(format GLenum, stride int, pointer unsafe.Pointer) {
	C.glInterleavedArrays(C.GLenum(format), C.GLsizei(stride), pointer)
}

//bool glIsEnabled (GLenum cap)
func IsEnabled(cap GLenum) bool {
	return goBool(C.glIsEnabled(C.GLenum(cap)))
}

//bool glIsList (uint list)
func IsList(list uint) bool {
	return goBool(C.glIsList(C.GLuint(list)))
}

//void glLightModelf (GLenum pname, float32 param)
func LightModelf(pname GLenum, param float32) {
	C.glLightModelf(C.GLenum(pname), C.GLfloat(param))
}

//void glLightModelfv (GLenum pname, const float *params)
func LightModelfv(pname GLenum, params []float32) {
	if len(params) == 0 {
		panic("Invalid params length")
	}
	C.glLightModelfv(C.GLenum(pname), (*C.GLfloat)(&params[0]))
}

//void glLightModeli (GLenum pname, int param)
func LightModeli(pname GLenum, param int) {
	C.glLightModeli(C.GLenum(pname), C.GLint(param))
}

//void glLightModeliv (GLenum pname, const int *params)
func LightModeliv(pname GLenum, params []int32) {
	if len(params) == 0 {
		panic("Invalid params length")
	}
	C.glLightModeliv(C.GLenum(pname), (*C.GLint)(&params[0]))
}

//void glLightf (GLenum light, GLenum pname, float32 param)
func Lightf(light GLenum, pname GLenum, param float32) {
	C.glLightf(C.GLenum(light), C.GLenum(pname), C.GLfloat(param))
}

//void glLightfv (GLenum light, GLenum pname, const float *params)
func Lightfv(light GLenum, pname GLenum, params []float32) {
	if len(params) == 0 {
		panic("Invalid params length")
	}
	C.glLightfv(C.GLenum(light), C.GLenum(pname), (*C.GLfloat)(&params[0]))
}

//void glLighti (GLenum light, GLenum pname, int param)
func Lighti(light GLenum, pname GLenum, param int) {
	C.glLighti(C.GLenum(light), C.GLenum(pname), C.GLint(param))
}

//void glLightiv (GLenum light, GLenum pname, const int *params)
func Lightiv(light GLenum, pname GLenum, params []int32) {
	if len(params) == 0 {
		panic("Invalid params length")
	}
	C.glLightiv(C.GLenum(light), C.GLenum(pname), (*C.GLint)(&params[0]))
}

//void glLineStipple (int factor, uint16 pattern)
func LineStipple(factor int, pattern uint16) {
	C.glLineStipple(C.GLint(factor), C.GLushort(pattern))
}

//void glLineWidth (float32 width)
func LineWidth(width float32) {
	C.glLineWidth(C.GLfloat(width))
}

//void glListBase (uint base)
func ListBase(base uint) {
	C.glListBase(C.GLuint(base))
}

//void glLoadName (uint name)
func LoadName(name uint) {
	C.glLoadName(C.GLuint(name))
}

//void glLogicOp (GLenum opcode)
func LogicOp(opcode GLenum) {
	C.glLogicOp(C.GLenum(opcode))
}

//void glMap1d (GLenum target, float64 u1, float64 u2, int stride, int order, const float64 *points)
func Map1d(target GLenum, u1 float64, u2 float64, stride int, order int, points []float64) {
	if len(points) == 0 {
		panic("Invalid points size")
	}
	C.glMap1d(C.GLenum(target), C.GLdouble(u1), C.GLdouble(u2),
		C.GLint(stride), C.GLint(order), (*C.GLdouble)(&points[0]))
}

//void glMap1f (GLenum target, float32 u1, float32 u2, int stride, int order, const float32 *points)
func Map1f(target GLenum, u1 float32, u2 float32, stride int, order int, points []float32) {
	if len(points) == 0 {
		panic("Invalid points size")
	}
	C.glMap1f(C.GLenum(target), C.GLfloat(u1), C.GLfloat(u2), C.GLint(stride),
		C.GLint(order), (*C.GLfloat)(&points[0]))
}

//void glMap2d (GLenum target, float64 u1, float64 u2, int ustride, int uorder, float64 v1, float64 v2, int vstride, int vorder, const float64 *points)
func Map2d(target GLenum, u1 float64, u2 float64, ustride int, uorder int, v1 float64, v2 float64, vstride int, vorder int, points []float64) {
	if len(points) == 0 {
		panic("Invalid points size")
	}
	C.glMap2d(C.GLenum(target), C.GLdouble(u1), C.GLdouble(u2), C.GLint(ustride),
		C.GLint(uorder), C.GLdouble(v1), C.GLdouble(v2), C.GLint(vstride),
		C.GLint(vorder), (*C.GLdouble)(&points[0]))
}

//void glMap2f (GLenum target, float32 u1, float32 u2, int ustride, int uorder, float32 v1, float32 v2, int vstride, int vorder, const float32 *points)
func Map2f(target GLenum, u1 float32, u2 float32, ustride int, uorder int, v1 float32, v2 float32, vstride int, vorder int, points []float32) {
	if len(points) == 0 {
		panic("Invalid points size")
	}
	C.glMap2f(C.GLenum(target), C.GLfloat(u1), C.GLfloat(u2), C.GLint(ustride),
		C.GLint(uorder), C.GLfloat(v1), C.GLfloat(v2), C.GLint(vstride),
		C.GLint(vorder), (*C.GLfloat)(&points[0]))
}

//void glMapGrid1d (int un, float64 u1, float64 u2)
func MapGrid1d(un int, u1 float64, u2 float64) {
	C.glMapGrid1d(C.GLint(un), C.GLdouble(u1), C.GLdouble(u2))
}

//void glMapGrid1f (int un, float32 u1, float32 u2)
func MapGrid1f(un int, u1 float32, u2 float32) {
	C.glMapGrid1f(C.GLint(un), C.GLfloat(u1), C.GLfloat(u2))
}

//void glMapGrid2d (int un, float64 u1, float64 u2, int vn, float64 v1, float64 v2)
func MapGrid2d(un int, u1 float64, u2 float64, vn int, v1 float64, v2 float64) {
	C.glMapGrid2d(C.GLint(un), C.GLdouble(u1), C.GLdouble(u2), C.GLint(vn), C.GLdouble(v1), C.GLdouble(v2))
}

//void glMapGrid2f (int un, float32 u1, float32 u2, int vn, float32 v1, float32 v2)
func MapGrid2f(un int, u1 float32, u2 float32, vn int, v1 float32, v2 float32) {
	C.glMapGrid2f(C.GLint(un), C.GLfloat(u1), C.GLfloat(u2), C.GLint(vn), C.GLfloat(v1), C.GLfloat(v2))
}

//void glMaterialf (GLenum face, GLenum pname, float32 param)
func Materialf(face GLenum, pname GLenum, param float32) {
	C.glMaterialf(C.GLenum(face), C.GLenum(pname), C.GLfloat(param))
}

//void glMaterialfv (GLenum face, GLenum pname, const float *params)
func Materialfv(face GLenum, pname GLenum, params []float32) {
	if len(params) == 0 {
		panic("Invalid params length")
	}
	C.glMaterialfv(C.GLenum(face), C.GLenum(pname), (*C.GLfloat)(&params[0]))
}

//void glMateriali (GLenum face, GLenum pname, int param)
func Materiali(face GLenum, pname GLenum, param int) {
	C.glMateriali(C.GLenum(face), C.GLenum(pname), C.GLint(param))
}

//void glMaterialiv (GLenum face, GLenum pname, const int *params)
func Materialiv(face GLenum, pname GLenum, params []int32) {
	if len(params) == 0 {
		panic("Invalid params length")
	}
	C.glMaterialiv(C.GLenum(face), C.GLenum(pname), (*C.GLint)(&params[0]))
}

//void glNewList (uint list, GLenum mode)
func NewList(list uint, mode GLenum) {
	C.glNewList(C.GLuint(list), C.GLenum(mode))
}

//void glNormal3b (int8 nx, int8 ny, int8 nz)
func Normal3b(nx int8, ny int8, nz int8) {
	C.glNormal3b(C.GLbyte(nx), C.GLbyte(ny), C.GLbyte(nz))
}

//void glNormal3bv (const int8 *v)
func Normal3bv(v *[3]int8) {
	C.glNormal3bv((*C.GLbyte)(&v[0]))
}

//void glNormal3d (float64 nx, float64 ny, float64 nz)
func Normal3d(nx float64, ny float64, nz float64) {
	C.glNormal3d(C.GLdouble(nx), C.GLdouble(ny), C.GLdouble(nz))
}

//void glNormal3dv (const float64 *v)
func Normal3dv(v *[3]float64) {
	C.glNormal3dv((*C.GLdouble)(&v[0]))
}

//void glNormal3f (float32 nx, float32 ny, float32 nz)
func Normal3f(nx float32, ny float32, nz float32) {
	C.glNormal3f(C.GLfloat(nx), C.GLfloat(ny), C.GLfloat(nz))
}

//void glNormal3fv (const float *v)
func Normal3fv(v *[3]float32) {
	C.glNormal3fv((*C.GLfloat)(&v[0]))
}

//void glNormal3i (int nx, int ny, int nz)
func Normal3i(nx int, ny int, nz int) {
	C.glNormal3i(C.GLint(nx), C.GLint(ny), C.GLint(nz))
}

//void glNormal3iv (const int *v)
func Normal3iv(v *[3]int32) {
	C.glNormal3iv((*C.GLint)(&v[0]))
}

//void glNormal3s (int16 nx, int16 ny, int16 nz)
func Normal3s(nx int16, ny int16, nz int16) {
	C.glNormal3s(C.GLshort(nx), C.GLshort(ny), C.GLshort(nz))
}

//void glNormal3sv (const int16 *v)
func Normal3sv(v *[3]int16) {
	C.glNormal3sv((*C.GLshort)(&v[0]))
}

//void glNormalPointer (GLenum type, int stride, const GLvoid *pointer)
func NormalPointer(typ GLenum, stride int, pointer interface{}) {
	C.glNormalPointer(C.GLenum(typ), C.GLsizei(stride), ptr(pointer))
}

//void glPassThrough (float32 token)
func PassThrough(token float32) {
	C.glPassThrough(C.GLfloat(token))
}

//void glPixelStoref (GLenum pname, float param)
func PixelStoref(pname GLenum, param float32) {
	C.glPixelStoref(C.GLenum(pname), C.GLfloat(param))
}

//void glPixelStorei (GLenum pname, int param)
func PixelStorei(pname GLenum, param int) {
	C.glPixelStorei(C.GLenum(pname), C.GLint(param))
}

//void glPixelTransferf (GLenum pname, float32 param)
func PixelTransferf(pname GLenum, param float32) {
	C.glPixelTransferf(C.GLenum(pname), C.GLfloat(param))
}

//void glPixelTransferi (GLenum pname, int param)
func PixelTransferi(pname GLenum, param int) {
	C.glPixelTransferi(C.GLenum(pname), C.GLint(param))
}

//void glPixelZoom (float32 xfactor, float32 yfactor)
func PixelZoom(xfactor float32, yfactor float32) {
	C.glPixelZoom(C.GLfloat(xfactor), C.GLfloat(yfactor))
}

//void glPointSize (float32 size)
func PointSize(size float32) {
	C.glPointSize(C.GLfloat(size))
}

//void glPolygonMode (GLenum face, GLenum mode)
func PolygonMode(face GLenum, mode GLenum) {
	C.glPolygonMode(C.GLenum(face), C.GLenum(mode))
}

//void glPolygonOffset (float32 factor, float32 units)
func PolygonOffset(factor float32, units float32) {
	C.glPolygonOffset(C.GLfloat(factor), C.GLfloat(units))
}

//void glPolygonStipple (const uint8 *mask)
func PolygonStipple(mask *uint8) {
	C.glPolygonStipple((*C.GLubyte)(mask))
}

//void glPopAttrib (void)
func PopAttrib() {
	C.glPopAttrib()
}

//void glPopClientAttrib (void)
func PopClientAttrib() {
	C.glPopClientAttrib()
}

//void glPopName (void)
func PopName() {
	C.glPopName()
}

//void glPrimitiveRestartIndex(GLuint index)
func PrimitiveRestartIndex(index GLuint) {
	C.glPrimitiveRestartIndex(C.GLuint(index))
}

//void glPushAttrib (GLbitfield mask)
func PushAttrib(mask GLbitfield) {
	C.glPushAttrib(C.GLbitfield(mask))
}

//void glPushClientAttrib (GLbitfield mask)
func PushClientAttrib(mask GLbitfield) {
	C.glPushClientAttrib(C.GLbitfield(mask))
}

//void glPushName (uint name)
func PushName(name uint) {
	C.glPushName(C.GLuint(name))
}

//void glRasterPos2d (float64 x, float64 y)
func RasterPos2d(x float64, y float64) {
	C.glRasterPos2d(C.GLdouble(x), C.GLdouble(y))
}

//void glRasterPos2dv (const float64 *v)
func RasterPos2dv(v *[2]float64) {
	C.glRasterPos2dv((*C.GLdouble)(&v[0]))
}

//void glRasterPos2f (float32 x, float32 y)
func RasterPos2f(x float32, y float32) {
	C.glRasterPos2f(C.GLfloat(x), C.GLfloat(y))
}

//void glRasterPos2fv (const float *v)
func RasterPos2fv(v *[2]float32) {
	C.glRasterPos2fv((*C.GLfloat)(&v[0]))
}

//void glRasterPos2i (int x, int y)
func RasterPos2i(x int, y int) {
	C.glRasterPos2i(C.GLint(x), C.GLint(y))
}

//void glRasterPos2iv (const int *v)
func RasterPos2iv(v *[2]int32) {
	C.glRasterPos2iv((*C.GLint)(&v[0]))
}

//void glRasterPos2s (int16 x, int16 y)
func RasterPos2s(x int16, y int16) {
	C.glRasterPos2s(C.GLshort(x), C.GLshort(y))
}

//void glRasterPos2sv (const int16 *v)
func RasterPos2sv(v *[2]int16) {
	C.glRasterPos2sv((*C.GLshort)(&v[0]))
}

//void glRasterPos3d (float64 x, float64 y, float64 z)
func RasterPos3d(x float64, y float64, z float64) {
	C.glRasterPos3d(C.GLdouble(x), C.GLdouble(y), C.GLdouble(z))
}

//void glRasterPos3dv (const float64 *v)
func RasterPos3dv(v *[3]float64) {
	C.glRasterPos3dv((*C.GLdouble)(&v[0]))
}

//void glRasterPos3f (float32 x, float32 y, float32 z)
func RasterPos3f(x float32, y float32, z float32) {
	C.glRasterPos3f(C.GLfloat(x), C.GLfloat(y), C.GLfloat(z))
}

//void glRasterPos3fv (const float *v)
func RasterPos3fv(v *[3]float32) {
	C.glRasterPos3fv((*C.GLfloat)(&v[0]))
}

//void glRasterPos3i (int x, int y, int z)
func RasterPos3i(x int, y int, z int) {
	C.glRasterPos3i(C.GLint(x), C.GLint(y), C.GLint(z))
}

//void glRasterPos3iv (const int *v)
func RasterPos3iv(v *[3]int32) {
	C.glRasterPos3iv((*C.GLint)(&v[0]))
}

//void glRasterPos3s (int16 x, int16 y, int16 z)
func RasterPos3s(x int16, y int16, z int16) {
	C.glRasterPos3s(C.GLshort(x), C.GLshort(y), C.GLshort(z))
}

//void glRasterPos3sv (const int16 *v)
func RasterPos3sv(v *[3]int16) {
	C.glRasterPos3sv((*C.GLshort)(&v[0]))
}

//void glRasterPos4d (float64 x, float64 y, float64 z, float64 w)
func RasterPos4d(x float64, y float64, z float64, w float64) {
	C.glRasterPos4d(C.GLdouble(x), C.GLdouble(y), C.GLdouble(z), C.GLdouble(w))
}

//void glRasterPos4dv (const float64 *v)
func RasterPos4dv(v *[3]float64) {
	C.glRasterPos4dv((*C.GLdouble)(&v[0]))
}

//void glRasterPos4f (float32 x, float32 y, float32 z, float32 w)
func RasterPos4f(x float32, y float32, z float32, w float32) {
	C.glRasterPos4f(C.GLfloat(x), C.GLfloat(y), C.GLfloat(z), C.GLfloat(w))
}

//void glRasterPos4fv (const float *v)
func RasterPos4fv(v *[4]float32) {
	C.glRasterPos4fv((*C.GLfloat)(&v[0]))
}

//void glRasterPos4i (int x, int y, int z, int w)
func RasterPos4i(x int, y int, z int, w int) {
	C.glRasterPos4i(C.GLint(x), C.GLint(y), C.GLint(z), C.GLint(w))
}

//void glRasterPos4iv (const int *v)
func RasterPos4iv(v *[4]int32) {
	C.glRasterPos4iv((*C.GLint)(&v[0]))
}

//void glRasterPos4s (int16 x, int16 y, int16 z, int16 w)
func RasterPos4s(x int16, y int16, z int16, w int16) {
	C.glRasterPos4s(C.GLshort(x), C.GLshort(y), C.GLshort(z), C.GLshort(w))
}

//void glRasterPos4sv (const int16 *v)
func RasterPos4sv(v *[4]int16) {
	C.glRasterPos4sv((*C.GLshort)(&v[0]))
}

//void glReadBuffer (GLenum mode)
func ReadBuffer(mode GLenum) {
	C.glReadBuffer(C.GLenum(mode))
}

//void glReadPixels (int x, int y, int width, int height, GLenum format, GLenum type, GLvoid *pixels)
func ReadPixels(x int, y int, width int, height int, format, typ GLenum, pixels interface{}) {
	C.glReadPixels(C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height),
		C.GLenum(format), C.GLenum(typ), ptr(pixels))
}

//void glRectd (float64 x1, float64 y1, float64 x2, float64 y2)
func Rectd(x1 float64, y1 float64, x2 float64, y2 float64) {
	C.glRectd(C.GLdouble(x1), C.GLdouble(y1), C.GLdouble(x2), C.GLdouble(y2))
}

//void glRectdv (const float64 *v1, const float64 *v2)
func Rectdv(a, b *[2]float64) {
	C.glRectdv((*C.GLdouble)(&a[0]), (*C.GLdouble)(&b[0]))
}

//void glRectf (float32 x1, float32 y1, float32 x2, float32 y2)
func Rectf(x1 float32, y1 float32, x2 float32, y2 float32) {
	C.glRectf(C.GLfloat(x1), C.GLfloat(y1), C.GLfloat(x2), C.GLfloat(y2))
}

//void glRectfv (const float *v1, const float *v2)
func Rectfv(a, b *[2]float32) {
	C.glRectfv((*C.GLfloat)(&a[0]), (*C.GLfloat)(&b[0]))
}

//void glRecti (int x1, int y1, int x2, int y2)
func Recti(x1 int, y1 int, x2 int, y2 int) {
	C.glRecti(C.GLint(x1), C.GLint(y1), C.GLint(x2), C.GLint(y2))
}

//void glRectiv (const int *v1, const int *v2)
func Rectiv(a, b *[2]int32) {
	C.glRectiv((*C.GLint)(&a[0]), (*C.GLint)(&b[0]))
}

//void glRects (int16 x1, int16 y1, int16 x2, int16 y2)
func Rects(x1 int16, y1 int16, x2 int16, y2 int16) {
	C.glRects(C.GLshort(x1), C.GLshort(y1), C.GLshort(x2), C.GLshort(y2))
}

//void glRectsv (const int16 *v1, const int16 *v2)
func Rectsv(a, b *[2]int16) {
	C.glRectsv((*C.GLshort)(&a[0]), (*C.GLshort)(&b[0]))
}

//int glRenderMode (GLenum mode)
func RenderMode(mode GLenum) int {
	return int(C.glRenderMode(C.GLenum(mode)))
}

//void glScissor (int x, int y, int width, int height)
func Scissor(x int, y int, width int, height int) {
	C.glScissor(C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height))
}

//void glSelectBuffer (GLsizei size, uint *buffer)
func SelectBuffer(buffer []uint32) {
	if len(buffer) > 0 {
		C.glSelectBuffer(C.GLsizei(len(buffer)), (*C.GLuint)(&buffer[0]))
	}
}

//void glShadeModel (GLenum mode)
func ShadeModel(mode GLenum) {
	C.glShadeModel(C.GLenum(mode))
}

//void glStencilFunc (GLenum func, int ref, uint mask)
func StencilFunc(func_ GLenum, ref int, mask uint) {
	C.glStencilFunc(C.GLenum(func_), C.GLint(ref), C.GLuint(mask))
}

//void glStencilMask (uint mask)
func StencilMask(mask uint) {
	C.glStencilMask(C.GLuint(mask))
}

//void glStencilOp (GLenum fail, GLenum zfail, GLenum zpass)
func StencilOp(fail GLenum, zfail GLenum, zpass GLenum) {
	C.glStencilOp(C.GLenum(fail), C.GLenum(zfail), C.GLenum(zpass))
}

//void glViewport (int x, int y, int width, int height)
func Viewport(x int, y int, width int, height int) {
	C.glViewport(C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height))
}

// void glGetFramebufferAttachmentParameter(GLenum target, GLenum attachment, GLenum pname, GLint* params);
//func GetFramebufferAttachmentParameter (target, attachment, pname GLenum, params []int32) {
//	if len(params) == 0 {
//		panic("Invalid params size")
//	}
//  C.glGetFramebufferAttachmentParameter (C.GLenum(target), C.GLenum(attachment),
//  	C.GLenum(pname), (*C.GLint)(&params[0]))
//}

func Init() GLenum {
	C.SetGlewExperimental(C.GLboolean(1))
	return GLenum(C.glewInit())
}
