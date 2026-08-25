//go:build !v8go_profiling

// Copyright 2024 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

func profileIsolateCreated(iso *Isolate)  {}
func profileIsolateDisposed(iso *Isolate) {}
func profileContextCreated(ctx *Context)  {}
func profileContextClosed(ctx *Context)   {}
