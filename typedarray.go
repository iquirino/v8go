// Copyright 2024 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include "object.h"
// #include "context.h"
import "C"
import (
	"errors"
	"unsafe"
)

// NewArrayBufferFromBytes creates a new ArrayBuffer containing a copy of the given bytes.
// The data is copied into V8's heap — modifications to the original slice won't affect it.
func NewArrayBufferFromBytes(ctx *Context, data []byte) (*Value, error) {
	if ctx == nil {
		return nil, errors.New("v8go: Context is required")
	}
	var ptr C.ValuePtr
	if len(data) > 0 {
		ptr = C.NewArrayBufferFromBytes(ctx.ptr, unsafe.Pointer(&data[0]), C.int(len(data)))
	} else {
		ptr = C.NewArrayBufferFromBytes(ctx.ptr, nil, 0)
	}
	if ptr == nil {
		return nil, errors.New("v8go: failed to create ArrayBuffer")
	}
	return &Value{ptr, ctx}, nil
}

// NewUint8ArrayFromBytes creates a new Uint8Array containing a copy of the given bytes.
// This is the equivalent of `new Uint8Array([...data])` in JS.
func NewUint8ArrayFromBytes(ctx *Context, data []byte) (*Value, error) {
	if ctx == nil {
		return nil, errors.New("v8go: Context is required")
	}
	var ptr C.ValuePtr
	if len(data) > 0 {
		ptr = C.NewUint8ArrayFromBytes(ctx.ptr, unsafe.Pointer(&data[0]), C.int(len(data)))
	} else {
		ptr = C.NewUint8ArrayFromBytes(ctx.ptr, nil, 0)
	}
	if ptr == nil {
		return nil, errors.New("v8go: failed to create Uint8Array")
	}
	return &Value{ptr, ctx}, nil
}
