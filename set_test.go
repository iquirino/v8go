package v8go_test

import (
	"testing"

	v8 "github.com/iquirino/v8go"
)

func TestNewSet(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	s, err := v8.NewSet(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if s.SetSize() != 0 {
		t.Errorf("expected size 0, got %d", s.SetSize())
	}
}

func TestSetAddHasDelete(t *testing.T) {
	t.Parallel()
	iso := v8.NewIsolate()
	defer iso.Dispose()
	ctx := v8.NewContext(iso)
	defer ctx.Close()

	s, _ := v8.NewSet(ctx)
	val, _ := v8.NewValue(iso, "hello")

	s.SetAdd(val)
	if s.SetSize() != 1 {
		t.Errorf("expected size 1, got %d", s.SetSize())
	}
	if !s.SetHas(val) {
		t.Error("expected set to have 'hello'")
	}

	s.SetDelete(val)
	if s.SetSize() != 0 {
		t.Errorf("expected size 0 after delete, got %d", s.SetSize())
	}
}

func TestValueAsSet(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	val, _ := ctx.RunScript(`new Set([1, 2, 3])`, "")
	s, err := val.AsSet()
	if err != nil {
		t.Fatal(err)
	}
	if s.SetSize() != 3 {
		t.Errorf("expected size 3, got %d", s.SetSize())
	}
}
