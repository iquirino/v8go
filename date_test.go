package v8go_test

import (
	"testing"
	"time"

	v8 "github.com/iquirino/v8go"
)

func TestNewDate(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	now := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	date, err := v8.NewDate(ctx, now)
	if err != nil {
		t.Fatal(err)
	}

	got := date.Time()
	if got.UnixMilli() != now.UnixMilli() {
		t.Errorf("expected %v, got %v", now, got)
	}
}

func TestDateToISOString(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	date, _ := v8.NewDate(ctx, time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC))
	iso, err := date.ToISOString()
	if err != nil {
		t.Fatal(err)
	}
	expected := "2024-01-15T10:30:00.000Z"
	if iso != expected {
		t.Errorf("expected %q, got %q", expected, iso)
	}
}

func TestDateGetters(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	// Use UTC date to avoid timezone issues
	val, _ := ctx.RunScript(`new Date(Date.UTC(2024, 5, 15, 10, 30))`, "")
	date, _ := val.AsDate()

	ms, err := date.GetTime()
	if err != nil {
		t.Fatal(err)
	}
	if ms <= 0 {
		t.Errorf("expected positive timestamp, got %d", ms)
	}

	year, _ := date.GetFullYear()
	if year != 2024 {
		t.Errorf("expected year 2024, got %d", year)
	}
}

func TestValueAsDate(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	val, _ := ctx.RunScript(`new Date()`, "")
	_, err := val.AsDate()
	if err != nil {
		t.Fatal("expected AsDate to succeed on Date object")
	}

	val2, _ := ctx.RunScript(`"not a date"`, "")
	_, err = val2.AsDate()
	if err == nil {
		t.Error("expected error for non-Date value")
	}
}
