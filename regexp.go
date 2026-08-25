// Copyright 2024 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "object.h"
// #include "context.h"
import "C"
import (
	"errors"
	"unsafe"
)

// RegExpFlags are flags that can be passed to NewRegExp.
type RegExpFlags int

const (
	RegExpNone       RegExpFlags = 0
	RegExpGlobal     RegExpFlags = 1
	RegExpIgnoreCase RegExpFlags = 2
	RegExpMultiline  RegExpFlags = 4
	RegExpSticky     RegExpFlags = 8
	RegExpUnicode    RegExpFlags = 16
	RegExpDotAll     RegExpFlags = 32
)

// RegExp is a JavaScript RegExp object.
type RegExp struct {
	*Object
}

// NewRegExp creates a new RegExp object with the given pattern and flags.
func NewRegExp(ctx *Context, pattern string, flags RegExpFlags) (*RegExp, error) {
	if ctx == nil {
		return nil, errors.New("v8go: Context is required")
	}
	cPattern := C.CString(pattern)
	defer C.free(unsafe.Pointer(cPattern))
	ptr := C.NewRegExp(ctx.ptr, cPattern, C.int(len(pattern)), C.int(flags))
	if ptr == nil {
		return nil, errors.New("v8go: failed to create RegExp")
	}
	val := &Value{ptr, ctx}
	return &RegExp{&Object{val}}, nil
}

// Test tests the given string against this RegExp. Returns true if it matches.
func (r *RegExp) Test(str Valuer) (bool, error) {
	val, err := r.MethodCall("test", str)
	if err != nil {
		return false, err
	}
	return val.Boolean(), nil
}

// Source returns the pattern source text of this RegExp.
func (r *RegExp) Source() (string, error) {
	val, err := r.Get("source")
	if err != nil {
		return "", err
	}
	return val.String(), nil
}

// Flags returns the flags string (e.g., "gi") of this RegExp.
func (r *RegExp) Flags() (string, error) {
	val, err := r.Get("flags")
	if err != nil {
		return "", err
	}
	return val.String(), nil
}
