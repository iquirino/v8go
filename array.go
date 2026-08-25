// Copyright 2024 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include "object.h"
// #include "context.h"
import "C"
import (
	"errors"
)

// Array is a JavaScript Array object.
type Array struct {
	*Object
}

// NewArray creates a new JavaScript Array with the given length.
func NewArray(ctx *Context, length int) (*Array, error) {
	if ctx == nil {
		return nil, errors.New("v8go: Context is required")
	}
	ptr := C.NewArray(ctx.ptr, C.int(length))
	if ptr == nil {
		return nil, errors.New("v8go: failed to create Array")
	}
	val := &Value{ptr, ctx}
	return &Array{&Object{val}}, nil
}

// Length returns the length of the array.
func (a *Array) Length() int {
	return int(C.ArrayLength(a.ptr))
}

// Get returns the value at the given index.
func (a *Array) Get(idx uint32) (*Value, error) {
	return a.GetIdx(idx)
}

// Set sets the value at the given index.
func (a *Array) Set(idx uint32, val interface{}) error {
	return a.SetIdx(idx, val)
}

// Push appends one or more values to the end of the array and returns the new length.
func (a *Array) Push(args ...Valuer) (int, error) {
	result, err := a.MethodCall("push", args...)
	if err != nil {
		return 0, err
	}
	return int(result.Int32()), nil
}

// Pop removes and returns the last element of the array.
func (a *Array) Pop() (*Value, error) {
	return a.MethodCall("pop")
}

// Shift removes and returns the first element of the array.
func (a *Array) Shift() (*Value, error) {
	return a.MethodCall("shift")
}

// Unshift prepends one or more values to the beginning of the array and returns the new length.
func (a *Array) Unshift(args ...Valuer) (int, error) {
	result, err := a.MethodCall("unshift", args...)
	if err != nil {
		return 0, err
	}
	return int(result.Int32()), nil
}

// Includes returns true if the array contains the given value.
func (a *Array) Includes(val Valuer) (bool, error) {
	result, err := a.MethodCall("includes", val)
	if err != nil {
		return false, err
	}
	return result.Boolean(), nil
}

// IndexOf returns the first index of the given value, or -1 if not found.
func (a *Array) IndexOf(val Valuer) (int, error) {
	result, err := a.MethodCall("indexOf", val)
	if err != nil {
		return -1, err
	}
	return int(result.Int32()), nil
}
