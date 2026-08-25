package v8go_test

import (
	"testing"

	v8 "github.com/iquirino/v8go"
)

func TestNewArray(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	arr, err := v8.NewArray(ctx, 3)
	if err != nil {
		t.Fatal(err)
	}
	if arr.Length() != 3 {
		t.Errorf("expected length 3, got %d", arr.Length())
	}
}

func TestArrayGetSet(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	arr, _ := v8.NewArray(ctx, 0)
	arr.Set(0, "hello")
	arr.Set(1, "world")

	val, err := arr.Get(0)
	if err != nil {
		t.Fatal(err)
	}
	if val.String() != "hello" {
		t.Errorf("expected 'hello', got %q", val.String())
	}
}

func TestArrayPushPop(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	arr, _ := v8.NewArray(ctx, 0)
	v1, _ := v8.NewValue(ctx.Isolate(), "a")
	v2, _ := v8.NewValue(ctx.Isolate(), "b")

	newLen, err := arr.Push(v1, v2)
	if err != nil {
		t.Fatal(err)
	}
	if newLen != 2 {
		t.Errorf("expected length 2 after push, got %d", newLen)
	}

	popped, err := arr.Pop()
	if err != nil {
		t.Fatal(err)
	}
	if popped.String() != "b" {
		t.Errorf("expected 'b' from pop, got %q", popped.String())
	}
	if arr.Length() != 1 {
		t.Errorf("expected length 1 after pop, got %d", arr.Length())
	}
}

func TestArrayShiftUnshift(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	arr, _ := v8.NewArray(ctx, 0)
	v1, _ := v8.NewValue(ctx.Isolate(), "first")
	v2, _ := v8.NewValue(ctx.Isolate(), "second")
	arr.Push(v1, v2)

	shifted, err := arr.Shift()
	if err != nil {
		t.Fatal(err)
	}
	if shifted.String() != "first" {
		t.Errorf("expected 'first' from shift, got %q", shifted.String())
	}

	v3, _ := v8.NewValue(ctx.Isolate(), "prepended")
	newLen, err := arr.Unshift(v3)
	if err != nil {
		t.Fatal(err)
	}
	if newLen != 2 {
		t.Errorf("expected length 2 after unshift, got %d", newLen)
	}
}

func TestArrayIncludesIndexOf(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	val, _ := ctx.RunScript(`["apple", "banana", "cherry"]`, "")
	arr, _ := val.AsArray()

	search, _ := v8.NewValue(ctx.Isolate(), "banana")
	found, err := arr.Includes(search)
	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Error("expected banana to be found")
	}

	idx, err := arr.IndexOf(search)
	if err != nil {
		t.Fatal(err)
	}
	if idx != 1 {
		t.Errorf("expected index 1, got %d", idx)
	}

	missing, _ := v8.NewValue(ctx.Isolate(), "grape")
	idx, _ = arr.IndexOf(missing)
	if idx != -1 {
		t.Errorf("expected -1 for missing item, got %d", idx)
	}
}

func TestValueAsArray(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	val, _ := ctx.RunScript(`[1, 2, 3]`, "")
	arr, err := val.AsArray()
	if err != nil {
		t.Fatal(err)
	}
	if arr.Length() != 3 {
		t.Errorf("expected length 3, got %d", arr.Length())
	}

	// Non-array should error
	val2, _ := ctx.RunScript(`"not an array"`, "")
	_, err = val2.AsArray()
	if err == nil {
		t.Error("expected error for non-array value")
	}
}
