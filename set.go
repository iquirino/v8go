// Copyright 2024 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "object.h"
// #include "context.h"
import "C"
import "errors"

// Set is a JavaScript Set object.
type Set struct {
	*Object
}

// NewSet creates a new empty Set.
func NewSet(ctx *Context) (*Set, error) {
	if ctx == nil {
		return nil, errors.New("v8go: Context is required")
	}
	ptr := C.NewSet(ctx.ptr)
	if ptr == nil {
		return nil, errors.New("v8go: failed to create Set")
	}
	val := &Value{ptr, ctx}
	return &Set{&Object{val}}, nil
}

// SetAdd adds a value to the Set.
func (s *Set) SetAdd(val Valuer) {
	C.SetAdd(s.ptr, val.value().ptr)
}

// SetHas returns true if the value exists in the Set.
func (s *Set) SetHas(val Valuer) bool {
	return C.SetHas(s.ptr, val.value().ptr) != 0
}

// SetDelete removes a value from the Set. Returns true if it was present.
func (s *Set) SetDelete(val Valuer) bool {
	return C.SetDelete(s.ptr, val.value().ptr) != 0
}

// SetSize returns the number of elements in the Set.
func (s *Set) SetSize() int {
	return int(C.SetSize(s.ptr))
}

// AsSet casts the value to a Set. Returns error if not a Set.
func (v *Value) AsSet() (*Set, error) {
	if !v.IsSet() {
		return nil, errors.New("v8go: value is not a Set")
	}
	return &Set{&Object{v}}, nil
}
