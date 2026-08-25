// Copyright 2024 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "object.h"
// #include "context.h"
import "C"
import "errors"

// Map is a JavaScript Map object.
type Map struct {
	*Object
}

// NewMap creates a new empty Map.
func NewMap(ctx *Context) (*Map, error) {
	if ctx == nil {
		return nil, errors.New("v8go: Context is required")
	}
	ptr := C.NewMap(ctx.ptr)
	if ptr == nil {
		return nil, errors.New("v8go: failed to create Map")
	}
	val := &Value{ptr, ctx}
	return &Map{&Object{val}}, nil
}

// MapGet returns the value for the given key, or undefined if not present.
func (m *Map) MapGet(key Valuer) (*Value, error) {
	rtn := C.MapGet(m.ptr, key.value().ptr)
	return valueResult(m.ctx, rtn)
}

// MapSet sets a key-value pair and returns the Map (for chaining).
func (m *Map) MapSet(key, val Valuer) error {
	C.MapSet(m.ptr, key.value().ptr, val.value().ptr)
	return nil
}

// MapHas returns true if the key exists in the Map.
func (m *Map) MapHas(key Valuer) bool {
	return C.MapHas(m.ptr, key.value().ptr) != 0
}

// MapDelete removes a key from the Map. Returns true if the key was present.
func (m *Map) MapDelete(key Valuer) bool {
	return C.MapDelete(m.ptr, key.value().ptr) != 0
}

// MapSize returns the number of entries in the Map.
func (m *Map) MapSize() int {
	return int(C.MapSize(m.ptr))
}

// AsMap casts the value to a Map. Returns error if not a Map.
func (v *Value) AsMap() (*Map, error) {
	if !v.IsMap() {
		return nil, errors.New("v8go: value is not a Map")
	}
	return &Map{&Object{v}}, nil
}
