// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

//go:generate clang-format -i --verbose -style=Chromium v8go.h v8go.cc

// NOTE: The V8 config macros below MUST match the flags V8 was built with in
// deps/build.py. V8 is built with v8_enable_sandbox=false because the sandbox
// requires use_custom_libcxx=true (libc++ hardening) which conflicts with Go's
// runtime linking against the system libc++. Do NOT define V8_ENABLE_SANDBOX
// here unless the build system is updated to handle the libc++ conflict.

// #cgo CXXFLAGS: -fno-rtti -fPIC -std=c++20 -I${SRCDIR}/deps/include -Wall
// #cgo CXXFLAGS: -DV8_COMPRESS_POINTERS -DV8_31BIT_SMIS_ON_64BIT_ARCH
// #cgo CXXFLAGS: -DV8_DEPRECATION_WARNINGS -DV8_IMMINENT_DEPRECATION_WARNINGS
import "C"
