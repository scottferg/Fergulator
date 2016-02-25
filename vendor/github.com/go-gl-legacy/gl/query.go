// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

// #include "gl.h"
import "C"

type Query Object

func GenQuery() (q Query) {
	C.glGenQueries(1, (*C.GLuint)(&q))
	return
}

func GenQueries(queries []Query) {
	if len(queries) > 0 {
		C.glGenQueries(C.GLsizei(len(queries)), (*C.GLuint)(&queries[0]))
	}
}

func (query Query) Begin(target GLenum) {
	C.glBeginQuery(C.GLenum(target), C.GLuint(query))
}

func (query Query) BeginIndexed(target GLenum, index uint) {
	C.glBeginQueryIndexed(C.GLenum(target), C.GLuint(index), C.GLuint(query))
}

func (query Query) Delete() {
	C.glDeleteQueries(1, (*C.GLuint)(&query))
}

func (query Query) GetObjecti(pname GLenum) (param int32) {
	C.glGetQueryObjectiv(C.GLuint(query), C.GLenum(pname), (*C.GLint)(&param))
	return
}

func (query Query) GetObjectui(pname GLenum) (param uint32) {
	C.glGetQueryObjectuiv(C.GLuint(query), C.GLenum(pname), (*C.GLuint)(&param))
	return
}

func (query Query) GetObjecti64(pname GLenum) (param int64) {
	C.glGetQueryObjecti64v(C.GLuint(query), C.GLenum(pname), (*C.GLint64)(&param))
	return
}

func (query Query) GetObjectui64(pname GLenum) (param uint64) {
	C.glGetQueryObjectui64v(C.GLuint(query), C.GLenum(pname), (*C.GLuint64)(&param))
	return
}

func (query Query) Counter(target GLenum) {
	C.glQueryCounter(C.GLuint(query), C.GLenum(target))
}

// Returns whether the passed samples counter is immediately available. If a delay
// would not occur waiting for the query result, true is returned, which also indicates
// that the results of all previous queries are available as well.
func (query Query) ResultAvailable() bool {
	return query.GetObjectui(QUERY_RESULT_AVAILABLE) == TRUE
}

func (query Query) BeginConditionalRender(mode GLenum) {
	C.glBeginConditionalRender(C.GLuint(query), C.GLenum(mode))
}

func DeleteQueries(queries []Query) {
	if len(queries) > 0 {
		C.glDeleteQueries(C.GLsizei(len(queries)), (*C.GLuint)(&queries[0]))
	}
}

func EndQuery(target GLenum) {
	C.glEndQuery(C.GLenum(target))
}

func GetQuery(target GLenum, pname GLenum) (param int32) {
	C.glGetQueryiv(C.GLenum(target), C.GLenum(pname), (*C.GLint)(&param))
	return
}

func GetQueryIndexed(target GLenum, index uint, pname GLenum) (param int32) {
	C.glGetQueryIndexediv(C.GLenum(target), C.GLuint(index), C.GLenum(pname), (*C.GLint)(&param))
	return
}

func (query Query) EndQueryIndexed(target GLenum, index uint) {
	C.glEndQueryIndexed(C.GLenum(target), C.GLuint(index))
}

func (query Query) EndConditionalRender() {
	C.glEndConditionalRender()
}
