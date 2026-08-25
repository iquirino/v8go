package v8go_test

import (
	"testing"

	v8 "github.com/iquirino/v8go"
)

func TestNewArrayBufferFromBytes(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	data := []byte{1, 2, 3, 4, 5}
	val, err := v8.NewArrayBufferFromBytes(ctx, data)
	if err != nil {
		t.Fatal(err)
	}
	if !val.IsArrayBuffer() {
		t.Error("expected ArrayBuffer")
	}

	buf, release, err := val.ArrayBufferGetContents()
	if err != nil {
		t.Fatal(err)
	}
	defer release()

	if len(buf) != 5 {
		t.Fatalf("expected length 5, got %d", len(buf))
	}
	if buf[0] != 1 || buf[4] != 5 {
		t.Errorf("unexpected content: %v", buf)
	}
}

func TestNewUint8ArrayFromBytes(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	data := []byte{10, 20, 30}
	val, err := v8.NewUint8ArrayFromBytes(ctx, data)
	if err != nil {
		t.Fatal(err)
	}
	if !val.IsUint8Array() {
		t.Error("expected Uint8Array")
	}

	// Verify from JS side
	ctx.Global().Set("arr", val)
	result, err := ctx.RunScript(`arr[0] + arr[1] + arr[2]`, "")
	if err != nil {
		t.Fatal(err)
	}
	if result.Int32() != 60 {
		t.Errorf("expected 60, got %d", result.Int32())
	}
}

func TestNewArrayBufferFromBytesEmpty(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	val, err := v8.NewArrayBufferFromBytes(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !val.IsArrayBuffer() {
		t.Error("expected ArrayBuffer")
	}
}
