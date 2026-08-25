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
