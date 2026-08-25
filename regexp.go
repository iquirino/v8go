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
