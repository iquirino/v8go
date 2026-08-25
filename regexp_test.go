package v8go_test

import (
	"testing"

	v8 "github.com/iquirino/v8go"
)

func TestNewRegExp(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	re, err := v8.NewRegExp(ctx, `\d+`, v8.RegExpGlobal)
	if err != nil {
		t.Fatal(err)
	}

	src, _ := re.Source()
	if src != `\d+` {
		t.Errorf("expected source '\\d+', got %q", src)
	}

	flags, _ := re.Flags()
	if flags != "g" {
		t.Errorf("expected flags 'g', got %q", flags)
	}
}

func TestRegExpTest(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	re, _ := v8.NewRegExp(ctx, `^hello`, v8.RegExpIgnoreCase)

	match, _ := v8.NewValue(ctx.Isolate(), "Hello World")
	result, err := re.Test(match)
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Error("expected regex to match 'Hello World'")
	}

	noMatch, _ := v8.NewValue(ctx.Isolate(), "World Hello")
	result, _ = re.Test(noMatch)
	if result {
		t.Error("expected regex NOT to match 'World Hello'")
	}
}
