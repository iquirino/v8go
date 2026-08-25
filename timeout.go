// Copyright 2024 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

import (
	"errors"
	"time"
)

// ErrScriptTimeout is returned when a script execution exceeds the given deadline.
var ErrScriptTimeout = errors.New("v8go: script execution timed out")

// RunScriptWithTimeout executes the source JavaScript with a timeout.
// If the script does not complete within the given duration, execution is
// terminated and ErrScriptTimeout is returned.
func (c *Context) RunScriptWithTimeout(source, origin string, timeout time.Duration) (*Value, error) {
	done := make(chan struct{})
	var val *Value
	var err error

	go func() {
		val, err = c.RunScript(source, origin)
		close(done)
	}()

	select {
	case <-done:
		return val, err
	case <-time.After(timeout):
		c.iso.TerminateExecution()
		<-done // wait for RunScript to return after termination
		return nil, ErrScriptTimeout
	}
}
