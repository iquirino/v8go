package v8go_test

import (
	"errors"
	"testing"
)

func fatalIf(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func recoverPanic(f func()) (recovered any) {
	defer func() {
		recovered = recover()
	}()
	f()
	return nil
}

// errorsJoin wraps errors.Join from the standard library.
func errorsJoin(errs ...error) error {
	return errors.Join(errs...)
}
