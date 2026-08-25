package v8go_test

import (
	"testing"

	v8 "github.com/iquirino/v8go"
)

func TestObjectPrivateProperties(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	val, _ := ctx.RunScript(`({name: "visible"})`, "")
	obj, _ := val.AsObject()

	// Set a private property
	err := obj.SetPrivate("secret", "hidden_value")
	if err != nil {
		t.Fatal(err)
	}

	// Has private
	if !obj.HasPrivate("secret") {
		t.Error("expected HasPrivate to be true")
	}

	// Get private
	priv, err := obj.GetPrivate("secret")
	if err != nil {
		t.Fatal(err)
	}
	if priv.String() != "hidden_value" {
		t.Errorf("expected 'hidden_value', got %q", priv.String())
	}

	// Private is invisible to JS
	ctx.Global().Set("testObj", obj)
	keys, _ := ctx.RunScript(`Object.keys(testObj)`, "")
	keysArr, _ := keys.AsArray()
	if keysArr.Length() != 1 {
		t.Errorf("expected 1 visible key, got %d", keysArr.Length())
	}

	// Symbols also can't see it
	symKeys, _ := ctx.RunScript(`Object.getOwnPropertySymbols(testObj).length`, "")
	if symKeys.Int32() != 0 {
		t.Errorf("expected 0 symbol keys, got %d", symKeys.Int32())
	}

	// Delete private
	obj.DeletePrivate("secret")
	if obj.HasPrivate("secret") {
		t.Error("expected HasPrivate to be false after delete")
	}
}
