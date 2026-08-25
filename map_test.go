package v8go_test

import (
	"testing"

	v8 "github.com/iquirino/v8go"
)

func TestNewMap(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	m, err := v8.NewMap(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if m.MapSize() != 0 {
		t.Errorf("expected size 0, got %d", m.MapSize())
	}
}

func TestMapSetGetHasDelete(t *testing.T) {
	t.Parallel()
	iso := v8.NewIsolate()
	defer iso.Dispose()
	ctx := v8.NewContext(iso)
	defer ctx.Close()

	m, _ := v8.NewMap(ctx)
	key, _ := v8.NewValue(iso, "name")
	val, _ := v8.NewValue(iso, "Alice")

	m.MapSet(key, val)
	if m.MapSize() != 1 {
		t.Errorf("expected size 1, got %d", m.MapSize())
	}

	if !m.MapHas(key) {
		t.Error("expected map to have key 'name'")
	}

	got, err := m.MapGet(key)
	if err != nil {
		t.Fatal(err)
	}
	if got.String() != "Alice" {
		t.Errorf("expected 'Alice', got %q", got.String())
	}

	m.MapDelete(key)
	if m.MapSize() != 0 {
		t.Errorf("expected size 0 after delete, got %d", m.MapSize())
	}
}

func TestValueAsMap(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	val, _ := ctx.RunScript(`new Map([["a", 1], ["b", 2]])`, "")
	m, err := val.AsMap()
	if err != nil {
		t.Fatal(err)
	}
	if m.MapSize() != 2 {
		t.Errorf("expected size 2, got %d", m.MapSize())
	}
}
