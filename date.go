package v8go

// #include "object.h"
// #include "context.h"
import "C"
import (
	"errors"
	"time"
)

// Date is a JavaScript Date object.
type Date struct {
	*Object
}

// NewDate creates a new Date from a Go time.Time.
func NewDate(ctx *Context, t time.Time) (*Date, error) {
	if ctx == nil {
		return nil, errors.New("v8go: Context is required")
	}
	ms := C.double(float64(t.UnixMilli()))
	ptr := C.NewDate(ctx.ptr, ms)
	if ptr == nil {
		return nil, errors.New("v8go: failed to create Date")
	}
	val := &Value{ptr, ctx}
	return &Date{&Object{val}}, nil
}

// Time returns the Go time.Time equivalent of this Date.
func (d *Date) Time() time.Time {
	ms := float64(C.DateValueOf(d.ptr))
	return time.UnixMilli(int64(ms))
}

// ToISOString returns the date as an ISO 8601 string (e.g., "2024-01-15T10:30:00.000Z").
func (d *Date) ToISOString() (string, error) {
	val, err := d.MethodCall("toISOString")
	if err != nil {
		return "", err
	}
	return val.String(), nil
}

// GetFullYear returns the year (4 digits for dates between 1000 and 9999).
func (d *Date) GetFullYear() (int, error) {
	val, err := d.MethodCall("getFullYear")
	if err != nil {
		return 0, err
	}
	return int(val.Int32()), nil
}

// GetMonth returns the month (0-11).
func (d *Date) GetMonth() (int, error) {
	val, err := d.MethodCall("getMonth")
	if err != nil {
		return 0, err
	}
	return int(val.Int32()), nil
}

// GetDate returns the day of the month (1-31).
func (d *Date) GetDate() (int, error) {
	val, err := d.MethodCall("getDate")
	if err != nil {
		return 0, err
	}
	return int(val.Int32()), nil
}

// GetHours returns the hour (0-23).
func (d *Date) GetHours() (int, error) {
	val, err := d.MethodCall("getHours")
	if err != nil {
		return 0, err
	}
	return int(val.Int32()), nil
}

// GetTime returns the number of milliseconds since Unix epoch.
func (d *Date) GetTime() (int64, error) {
	val, err := d.MethodCall("getTime")
	if err != nil {
		return 0, err
	}
	return val.Integer(), nil
}
