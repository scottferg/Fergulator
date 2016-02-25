// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

// #include "gl.h"
import "C"

// Vertex Arrays
type VertexArray Object

func GenVertexArray() VertexArray {
	var a C.GLuint
	C.glGenVertexArrays(1, &a)
	return VertexArray(a)
}

func GenVertexArrays(arrays []VertexArray) {
	if len(arrays) > 0 {
		C.glGenVertexArrays(C.GLsizei(len(arrays)), (*C.GLuint)(&arrays[0]))
	}
}

func (array VertexArray) Delete() {
	C.glDeleteVertexArrays(1, (*C.GLuint)(&array))
}

func DeleteVertexArrays(arrays []VertexArray) {
	if len(arrays) > 0 {
		C.glDeleteVertexArrays(C.GLsizei(len(arrays)), (*C.GLuint)(&arrays[0]))
	}
}

func (array VertexArray) Bind() {
	C.glBindVertexArray(C.GLuint(array))
}

func (array VertexArray) Unbind() {
	C.glBindVertexArray(C.GLuint(0))
}
