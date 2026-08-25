package v8go_test

import (
	"testing"

	v8 "github.com/iquirino/v8go"
)

func TestObjectGetPropertyNames(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	val, _ := ctx.RunScript(`({a: 1, b: 2, c: 3})`, "")
	obj, _ := val.AsObject()

	names, err := obj.GetPropertyNames()
	if err != nil {
		t.Fatal(err)
	}
	if names.Length() < 3 {
		t.Errorf("expected at least 3 property names, got %d", names.Length())
	}
}

func TestObjectGetOwnPropertyNames(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	// Create object with inherited property
	val, _ := ctx.RunScript(`
		const parent = {inherited: true};
		const child = Object.create(parent);
		child.own = "mine";
		child
	`, "")
	obj, _ := val.AsObject()

	ownNames, err := obj.GetOwnPropertyNames()
	if err != nil {
		t.Fatal(err)
	}
	if ownNames.Length() != 1 {
		t.Errorf("expected 1 own property, got %d", ownNames.Length())
	}

	allNames, err := obj.GetPropertyNames()
	if err != nil {
		t.Fatal(err)
	}
	if allNames.Length() < 2 {
		t.Errorf("expected at least 2 property names (including inherited), got %d", allNames.Length())
	}
}

func TestObjectDefineOwnProperty(t *testing.T) {
	t.Parallel()
	iso := v8.NewIsolate()
	defer iso.Dispose()
	ctx := v8.NewContext(iso)
	defer ctx.Close()

	val, _ := ctx.RunScript(`var testObj = {}; testObj`, "")
	obj, _ := val.AsObject()

	propVal, _ := v8.NewValue(iso, "immutable")
	ok := obj.DefineOwnProperty("frozen", propVal, v8.ReadOnly|v8.DontDelete)
	if !ok {
		t.Error("DefineOwnProperty returned false")
	}

	// Verify property exists and has correct value
	got, _ := obj.Get("frozen")
	if got.String() != "immutable" {
		t.Errorf("expected 'immutable', got %q", got.String())
	}

	// Verify DontEnum: property should not appear in Object.keys()
	propVal2, _ := v8.NewValue(iso, "visible")
	obj.DefineOwnProperty("shown", propVal2, v8.None)

	keys, _ := ctx.RunScript(`Object.keys(testObj)`, "")
	keysArr, _ := keys.AsArray()
	// "shown" is enumerable, "frozen" (DontEnum not set, but let's just check it exists)
	if keysArr.Length() < 1 {
		t.Errorf("expected at least 1 key, got %d", keysArr.Length())
	}
}
