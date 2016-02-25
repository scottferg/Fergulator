// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

// #include "gl.h"
import "C"

// Transform Feedback Objects

type TransformFeedback Object

// Create a single transform feedback object
func GenTransformFeedback() TransformFeedback {
	var t C.GLuint
	C.glGenTransformFeedbacks(1, &t)
	return TransformFeedback(t)
}

// Fill slice with new transform feedbacks
func GenTransformFeedbacks(feedbacks []TransformFeedback) {
	if len(feedbacks) > 0 {
		C.glGenTransformFeedbacks(C.GLsizei(len(feedbacks)), (*C.GLuint)(&feedbacks[0]))
	}
}

// Delete a transform feedback object
func (feedback TransformFeedback) Delete() {
	C.glDeleteTransformFeedbacks(1, (*C.GLuint)(&feedback))
}

// Draw the results of the last Begin/End cycle from this transform feedback using primitive type 'mode'
func (feedback TransformFeedback) Draw(mode GLenum) {
	C.glDrawTransformFeedback(C.GLenum(mode), C.GLuint(feedback))
}

// Delete all transform feedbacks in a slice
func DeleteTransformFeedbacks(feedbacks []TransformFeedback) {
	if len(feedbacks) > 0 {
		C.glDeleteTransformFeedbacks(C.GLsizei(len(feedbacks)), (*C.GLuint)(&feedbacks[0]))
	}
}

// Bind this transform feedback as target
func (feedback TransformFeedback) Bind(target GLenum) {
	C.glBindTransformFeedback(C.GLenum(target), C.GLuint(feedback))
}

// Begin transform feedback with primitive type 'mode'
func BeginTransformFeedback(mode GLenum) {
	C.glBeginTransformFeedback(C.GLenum(mode))
}

// Pause transform feedback
func PauseTransformFeedback() {
	C.glPauseTransformFeedback()
}

// End transform feedback
func EndTransformFeedback() {
	C.glEndTransformFeedback()
}
