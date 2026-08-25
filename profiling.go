//go:build v8go_profiling

// Copyright 2024 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

import "runtime/pprof"

var (
	isolateProfile = pprof.NewProfile("v8go.isolate")
	contextProfile = pprof.NewProfile("v8go.context")
)

func profileIsolateCreated(iso *Isolate) {
	isolateProfile.Add(iso, 1)
}

func profileIsolateDisposed(iso *Isolate) {
	isolateProfile.Remove(iso)
}

func profileContextCreated(ctx *Context) {
	contextProfile.Add(ctx, 1)
}

func profileContextClosed(ctx *Context) {
	contextProfile.Remove(ctx)
}
