package v8go_test

import (
	"errors"
	"testing"
	"time"

	v8 "github.com/iquirino/v8go"
)

func TestRunScriptWithTimeout(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	// Script that completes in time
	val, err := ctx.RunScriptWithTimeout(`1 + 2`, "fast.js", time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val.String() != "3" {
		t.Errorf("expected '3', got %q", val.String())
	}
}

func TestRunScriptWithTimeoutExpires(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	// Infinite loop should timeout
	_, err := ctx.RunScriptWithTimeout(`while(true){}`, "slow.js", 50*time.Millisecond)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !errors.Is(err, v8.ErrScriptTimeout) {
		t.Errorf("expected ErrScriptTimeout, got: %v", err)
	}
}
