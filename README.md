# v8go — Execute JavaScript from Go

[![Go Reference](https://pkg.go.dev/badge/github.com/iquirino/v8go.svg)](https://pkg.go.dev/github.com/iquirino/v8go)

A Go binding to the V8 JavaScript engine. Run JavaScript, expose Go functions to JS, and exchange data between Go and V8 with zero serialization overhead.

## Install

```bash
go get github.com/iquirino/v8go
```

Prebuilt V8 static libraries are included for Linux and macOS (amd64/arm64) and Android.

---

## Quick Start

```go
package main

import (
    "fmt"
    v8 "github.com/iquirino/v8go"
)

func main() {
    ctx := v8.NewContext()
    defer ctx.Isolate().Dispose()
    defer ctx.Close()

    val, err := ctx.RunScript(`"Hello from V8!"`, "hello.js")
    if err != nil {
        panic(err)
    }
    fmt.Println(val.String()) // Hello from V8!
}
```

---

## Core Concepts

### Isolate (VM instance)

> **Think of an Isolate as a completely independent JavaScript universe.** It has its own heap, garbage collector, and JIT compiler. Two Isolates share nothing — they can't see each other's variables, objects, or functions. This is the strongest isolation boundary V8 offers. One Isolate = one thread at a time (V8 is not thread-safe within a single Isolate).

```go
iso := v8.NewIsolate()
defer iso.Dispose()
```

### Context (execution environment)

> **A Context is a global scope within an Isolate.** Think of it as a separate "tab" — it has its own `global` object, its own variables, but shares the underlying VM (JIT cache, GC) with other Contexts in the same Isolate. Contexts are cheap to create compared to Isolates.

```go
ctx := v8.NewContext(iso)
defer ctx.Close()
```

### One Isolate, many Contexts

> Use this when you need multiple independent scopes but want to save memory. Each Context gets its own global object, so variables defined in one are invisible to another.

```go
iso := v8.NewIsolate()
defer iso.Dispose()

ctx1 := v8.NewContext(iso)
ctx1.RunScript("var x = 1", "")

ctx2 := v8.NewContext(iso) // separate global scope
_, err := ctx2.RunScript("x", "") // ReferenceError: x is not defined
```

---

## Lifecycle & Cleanup

> **Critical:** V8 resources are NOT garbage collected by Go. You MUST explicitly dispose them or they'll leak.

```go
// Always follow this pattern:
iso := v8.NewIsolate()
defer iso.Dispose()       // frees the entire VM + heap

ctx := v8.NewContext(iso)
defer ctx.Close()         // frees the context + all tracked values

// If you created a context without an explicit isolate:
ctx := v8.NewContext()
defer ctx.Isolate().Dispose() // don't forget the implicit isolate!
defer ctx.Close()
```

**Order matters:** Close contexts before disposing their isolate.

---

## Running Scripts

### Basic execution

```go
val, err := ctx.RunScript(`1 + 2`, "math.js")
fmt.Println(val.String()) // "3"
```

### 🆕 Execution with timeout

```go
val, err := ctx.RunScriptWithTimeout(`while(true){}`, "loop.js", 100*time.Millisecond)
if errors.Is(err, v8.ErrScriptTimeout) {
    fmt.Println("Script timed out")
}
```

### Terminate long-running scripts manually

```go
go func() {
    time.Sleep(200 * time.Millisecond)
    iso.TerminateExecution()
}()
val, err := ctx.RunScript(longScript, "slow.js")
// err: "ExecutionTerminated: script execution has been terminated"
```

---

## Values

### Creating values from Go

```go
iso := v8.NewIsolate()

strVal, _ := v8.NewValue(iso, "hello")         // string
intVal, _ := v8.NewValue(iso, int32(42))        // int32
numVal, _ := v8.NewValue(iso, float64(3.14))    // float64
boolVal, _ := v8.NewValue(iso, true)            // bool
bigVal, _ := v8.NewValue(iso, big.NewInt(9999)) // BigInt
```

### Type checking

```go
val.IsString()    // true/false
val.IsNumber()
val.IsObject()
val.IsArray()
val.IsPromise()
val.IsDate()
val.IsRegExp()
// ... 40+ type checks available
```

### 🆕 Error-aware string conversion

```go
s, err := val.StringErr() // returns error if conversion fails (e.g., Symbol coercion)
s := val.String()         // same but returns "" on failure (for fmt.Stringer compat)
```

---

## Objects

### Get and set properties

```go
obj := ctx.Global()
obj.Set("version", "2.0.0")

val, _ := obj.Get("version")
fmt.Println(val.String()) // "2.0.0"

obj.Has("version")    // true
obj.Delete("version") // removes it
```

### 🆕 Property enumeration

```go
val, _ := ctx.RunScript(`({name: "Alice", age: 30})`, "")
obj, _ := val.AsObject()

names, _ := obj.GetPropertyNames()       // includes prototype chain
ownNames, _ := obj.GetOwnPropertyNames() // own only

for i := 0; i < ownNames.Length(); i++ {
    key, _ := ownNames.GetIdx(uint32(i))
    fmt.Println(key.String())
}
```

### 🆕 Define properties with attributes

```go
propVal, _ := v8.NewValue(iso, "immutable")
obj.DefineOwnProperty("locked", propVal, v8.ReadOnly|v8.DontDelete)
```

### 🆕 Private properties (invisible to JS)

```go
obj.SetPrivate("internal_id", "abc123")
val, _ := obj.GetPrivate("internal_id") // "abc123"
obj.HasPrivate("internal_id")           // true

// JS cannot see it:
// Object.keys(obj)                    → doesn't include "internal_id"
// Object.getOwnPropertySymbols(obj)   → doesn't include it either
```

---

## 🆕 Arrays

```go
arr, _ := v8.NewArray(ctx, 0)

v1, _ := v8.NewValue(iso, "hello")
v2, _ := v8.NewValue(iso, "world")

arr.Push(v1, v2)          // returns new length: 2
arr.Length()              // 2
val, _ := arr.Get(0)     // "hello"
arr.Pop()                 // removes and returns "world"
arr.Shift()               // removes and returns "hello"
arr.Unshift(v1)           // prepend, returns new length
arr.Includes(v1)          // true
arr.IndexOf(v1)           // 0

// Cast from Value:
val, _ = ctx.RunScript(`[1, 2, 3]`, "")
arr, _ = val.AsArray()
```

---

## 🆕 Dates

```go
date, _ := v8.NewDate(ctx, time.Now())

t := date.Time()                // Go time.Time
iso, _ := date.ToISOString()    // "2024-06-15T10:30:00.000Z"
year, _ := date.GetFullYear()   // 2024
month, _ := date.GetMonth()     // 5 (0-indexed)
ms, _ := date.GetTime()         // Unix milliseconds

// Cast from Value:
val, _ = ctx.RunScript(`new Date()`, "")
d, _ := val.AsDate()
```

---

## 🆕 Map and Set

### Map

```go
m, _ := v8.NewMap(ctx)

key, _ := v8.NewValue(iso, "name")
val, _ := v8.NewValue(iso, "Alice")

m.MapSet(key, val)
m.MapSize()          // 1
m.MapHas(key)        // true
got, _ := m.MapGet(key) // "Alice"
m.MapDelete(key)
```

### Set

```go
s, _ := v8.NewSet(ctx)

val, _ := v8.NewValue(iso, "item")
s.SetAdd(val)
s.SetSize()    // 1
s.SetHas(val)  // true
s.SetDelete(val)
```

---

## 🆕 RegExp

```go
re, _ := v8.NewRegExp(ctx, `\d+`, v8.RegExpGlobal|v8.RegExpIgnoreCase)

src, _ := re.Source() // `\d+`
flags, _ := re.Flags() // "gi"

str, _ := v8.NewValue(iso, "abc 123 def")
matched, _ := re.Test(str) // true
```

---

## 🆕 Binary Data (ArrayBuffer / TypedArray)

> **When to use this instead of JSON?** If you're passing binary data (images, protobuf, crypto buffers) or large numeric arrays between Go and JS, ArrayBuffers avoid the serialize→parse round-trip entirely. The byte slice you get back points directly into V8's heap — zero copy. For structured objects (maps, nested structs), JSON is still the simplest path.

### Create from Go bytes

```go
data := []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F}

// ArrayBuffer
buf, _ := v8.NewArrayBufferFromBytes(ctx, data)

// Uint8Array (JS can index it directly: arr[0], arr[1], etc.)
arr, _ := v8.NewUint8ArrayFromBytes(ctx, data)
```

### Read V8 buffer into Go

```go
val, _ := ctx.RunScript(`new ArrayBuffer(1024)`, "")
bytes, release, _ := val.ArrayBufferGetContents()
defer release()
// bytes is a []byte backed by V8 memory — zero copy!

// SharedArrayBuffer works too:
val2, _ := ctx.RunScript(`new SharedArrayBuffer(1024)`, "")
bytes2, release2, _ := val2.SharedArrayBufferGetContents()
defer release2()
```

---

## Functions

### Go function exposed to JS

```go
printFn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
    fmt.Println(info.Args()[0].String())
    return nil
})
global := v8.NewObjectTemplate(iso)
global.Set("print", printFn)

ctx := v8.NewContext(iso, global)
ctx.RunScript(`print("Hello from JS!")`, "")
```

### Go function with error handling

```go
fn := v8.NewFunctionTemplateWithError(iso, func(info *v8.FunctionCallbackInfo) (*v8.Value, error) {
    if len(info.Args()) == 0 {
        return nil, fmt.Errorf("argument required")
    }
    return info.Args()[0], nil
})
// If error is returned, it's thrown as a JS exception
```

### 🆕 Get function from template (returns error instead of panicking)

```go
fn, err := tmpl.GetFunction(ctx)
if err != nil {
    // handle template instantiation error
}
result, err := fn.Call(v8.Undefined(iso), arg1, arg2)
```

---

## Promises

> **What are microtasks?** In a browser or Node.js, Promise callbacks (`.then`, `.catch`, `async/await` continuations) don't execute immediately — they go into a "microtask queue" that runs after the current script finishes. In v8go, there's no event loop running automatically. You must explicitly tell V8 to process the queue by calling `ctx.PerformMicrotaskCheckpoint()`. Without this call, your Promise callbacks will never fire.

```go
resolver, _ := v8.NewPromiseResolver(ctx)
promise := resolver.GetPromise()

// Resolve from Go:
val, _ := v8.NewValue(iso, "done")
resolver.Resolve(val)

// IMPORTANT: Without this, .Then callbacks won't run!
ctx.PerformMicrotaskCheckpoint()

fmt.Println(promise.State())           // Fulfilled
fmt.Println(promise.Result().String()) // "done"
```

### 🆕 Then/Catch (returns error instead of panicking)

```go
p, err := promise.Then(func(info *v8.FunctionCallbackInfo) *v8.Value {
    fmt.Println("Resolved:", info.Args()[0].String())
    return nil
})
if err != nil {
    // handle error
}
ctx.PerformMicrotaskCheckpoint()
```

### 🆕 Microtask policy control

> By default, V8 runs microtasks (Promise callbacks) automatically after each script completes. If you want full control — for example, to batch multiple operations before resolving promises — switch to explicit mode.

```go
// Explicit: YOU decide when promises resolve
iso.SetMicrotasksPolicy(v8.MicrotasksExplicit)

ctx.RunScript(`fetch('/api').then(r => console.log(r))`, "") // .then won't fire yet!
// ... do other work ...
ctx.PerformMicrotaskCheckpoint() // NOW all pending .then/.catch callbacks run

// Auto (default): promises resolve immediately after each RunScript
iso.SetMicrotasksPolicy(v8.MicrotasksAuto)
```

---

## Error Handling

### JavaScript errors

```go
_, err := ctx.RunScript(`throw new TypeError("oops")`, "err.js")
if err != nil {
    jsErr := err.(*v8.JSError)
    fmt.Println(jsErr.Message)    // "TypeError: oops"
    fmt.Println(jsErr.Location)   // "err.js:1:1"
    fmt.Println(jsErr.StackTrace) // full stack trace
}
```

### 🆕 Exception value propagation

```go
_, err := ctx.RunScript(`throw new TypeError("oops")`, "")
jsErr := err.(*v8.JSError)

// Access the original V8 error object:
fmt.Println(jsErr.Value.IsNativeError()) // true
fmt.Println(jsErr.Value.String())        // "TypeError: oops"

// Rethrow in a callback:
iso.ThrowException(jsErr.Value)
```

---

## Modules (ES Modules)

> **Important:** V8 is not Node.js. There's no `require()`, no `module.exports`, no CommonJS. If you're loading code that uses `exports` or `require`, it won't work — that's Node.js syntax, not JavaScript. V8 only supports ES Modules (`import`/`export`). If you need CommonJS compat, prepend a shim: `var exports = {}; var module = {exports};`

```go
mod, err := v8.CompileModule(iso, `export const x = 42;`, "module.js")
if err != nil {
    panic(err)
}
err = mod.InstantiateModule(ctx, resolver)
val, err := mod.Evaluate(ctx)
```

---

## Resource Limits

### Memory limits

> V8 can consume unbounded memory if scripts allocate without limit. `WithResourceConstraints` sets a hard ceiling. When the limit is approached, V8 calls a near-heap-limit callback that terminates execution — your `RunScript` call returns an error instead of the process being OOM-killed.

```go
iso := v8.NewIsolate(v8.WithResourceConstraints(0, 50*1024*1024)) // max 50MB
// V8 calls TerminateExecution when limit is hit
```

### 🆕 Security tokens (multi-context isolation)

> **When would you use this?** If you run multiple tenants' code in separate Contexts on the same Isolate (to save memory), security tokens prevent one Context from accessing another's globals through shared prototype chains. If you use one Isolate per tenant, you don't need this — Isolates are already fully isolated.

```go
token, _ := v8.NewValue(iso, "tenant-A")
ctx.SetSecurityToken(token) // prevents cross-context access within same isolate
```

---

## Pre-compiled Scripts (Code Cache)

> **Why use this?** Parsing and compiling JavaScript has a real cost (especially for large scripts). If you run the same source code repeatedly in different contexts, you can compile it once and reuse the compiled bytecode. This skips the parsing+compilation step on subsequent runs — typically saving 20-40% of the first execution time.

```go
source := "const add = (a, b) => a + b"
script, _ := iso.CompileUnboundScript(source, "math.js", v8.CompileOptions{Mode: v8.CompileModeEager})
cache := script.CreateCodeCache()

// Later, in a new isolate:
script2, _ := iso2.CompileUnboundScript(source, "math.js", v8.CompileOptions{CachedData: cache})
val, _ := script2.Run(ctx2)
```

---

## CPU Profiler

```go
profiler := v8.NewCPUProfiler(iso)
profiler.StartProfiling("my-profile")

ctx.RunScript(code, "app.js")

profile := profiler.StopProfiling("my-profile")
root := profile.GetTopDownRoot()
// Walk the call tree...
```

---

## 🆕 Leak Detection (build tag)

> **Why?** If you create Isolates or Contexts without properly calling `Dispose()`/`Close()`, they leak V8 heap memory (each Isolate reserves ~4GB of virtual address space). This profiling integration lets you catch leaks using Go's standard `pprof` tooling.

Build with `-tags v8go_profiling` to enable pprof-based tracking of Isolate and Context creation/disposal:

```go
// With the build tag active:
// pprof.Lookup("v8go.isolate") tracks live isolates
// pprof.Lookup("v8go.context") tracks live contexts

// Example: check for leaks in tests
import "runtime/pprof"

func TestNoLeaks(t *testing.T) {
    // ... create and dispose isolates/contexts ...

    if n := pprof.Lookup("v8go.isolate").Count(); n != 0 {
        t.Errorf("leaked %d isolates", n)
    }
}
```

---

## Inspector (Console API)

```go
type MyHandler struct{}
func (h *MyHandler) ConsoleAPIMessage(msg v8.ConsoleAPIMessage) {
    fmt.Printf("[%d] %s\n", msg.ErrorLevel, msg.Message)
}

client := v8.NewInspectorClient(&MyHandler{})
inspector := v8.NewInspector(iso, client)
inspector.ContextCreated(ctx)

ctx.RunScript(`console.log("hello")`, "") // triggers handler
```

---

## Build Configuration

V8 is built with these GN flags (see `deps/build.py`):

| Flag | Value | Purpose |
|------|-------|---------|
| `v8_enable_sandbox` | `false` | Disabled — requires libc++ hardening which conflicts with Go's CGo linking |
| `v8_enable_pointer_compression` | `true` | 🆕 ~50% heap memory reduction |
| `v8_enable_maglev` | `true` | 🆕 Mid-tier JIT for faster warmup |
| `v8_enable_short_builtin_calls` | `true` | 🆕 Shorter x64 call sequences |
| `v8_enable_webassembly` | `false` | 🆕 Reduced attack surface |
| `v8_monolithic` | `true` | Single static archive |
| `v8_enable_i18n_support` | `true` | Full Intl API support |

---

## v1 — Breaking Changes

These methods changed signatures to return errors instead of panicking:

| Method | Before | After |
|--------|--------|-------|
| `Value.Object()` | `*Object` | `(*Object, error)` |
| `FunctionTemplate.GetFunction(ctx)` | `*Function` | `(*Function, error)` |
| `Promise.Then(cbs...)` | `*Promise` | `(*Promise, error)` |
| `Promise.Catch(cb)` | `*Promise` | `(*Promise, error)` |
| `Promise.ThenWithError(cbs...)` | `*Promise` | `(*Promise, error)` |
| `Promise.CatchWithError(cb)` | `*Promise` | `(*Promise, error)` |

**Migration:** Add `, err` to the left side of these calls and handle the error.

```go
// Before:
fn := tmpl.GetFunction(ctx)
prom.Then(callback)

// After:
fn, err := tmpl.GetFunction(ctx)
_, err = prom.Then(callback)
```

---

## v1 — New Features

| Feature | Description |
|---------|-------------|
| `Context.RunScriptWithTimeout` | Execute JS with a deadline — returns `ErrScriptTimeout` on expiry |
| `Value.StringErr()` | String conversion that reports errors instead of returning empty |
| `Value.AsArray()`, `Value.AsDate()`, `Value.AsMap()`, `Value.AsSet()` | Type-safe casting |
| `NewArray`, `Array.Push/Pop/Shift/Unshift/Includes/IndexOf` | Full Array API |
| `NewDate`, `Date.Time()`, `Date.ToISOString()`, `Date.GetFullYear()` | Date ↔ `time.Time` |
| `NewMap`, `Map.MapGet/MapSet/MapHas/MapDelete/MapSize` | Native Map without eval |
| `NewSet`, `Set.SetAdd/SetHas/SetDelete/SetSize` | Native Set without eval |
| `NewRegExp`, `RegExp.Test()`, `RegExp.Source()`, `RegExp.Flags()` | Create and use regex from Go |
| `NewArrayBufferFromBytes`, `NewUint8ArrayFromBytes` | Inject binary data into V8 |
| `Value.ArrayBufferGetContents()` | Zero-copy read of ArrayBuffer bytes |
| `Object.GetPropertyNames()`, `Object.GetOwnPropertyNames()` | Enumerate object keys |
| `Object.DefineOwnProperty(key, val, attrs)` | Define properties with ReadOnly/DontDelete |
| `Object.SetPrivate/GetPrivate/HasPrivate/DeletePrivate` | Properties invisible to JS |
| `Isolate.SetMicrotasksPolicy()` | Control when Promise callbacks execute |
| `Context.SetSecurityToken/GetSecurityToken` | Cross-context access control |
| `JSError.Value` | Access original V8 exception object for rethrowing |
| `Value.DetailString()` | No longer panics — returns fallback string on failure |
| `Value.Release()` | Nil-safe, double-call safe |
| Build tag `v8go_profiling` | pprof profiles for tracking Isolate/Context leaks |
| V8 sandbox (`v8_enable_sandbox`) | Disabled — requires libc++ hardening incompatible with CGo |
| Maglev JIT (`v8_enable_maglev`) | Faster warmup for short-lived scripts |
| WebAssembly disabled | Reduced binary size and attack surface |
| Null-byte safe strings | `RunScript`, `CompileModule`, `JSONParse` handle `\x00` correctly |
| malloc safety | All C++ allocations check for NULL |

---

## Supported Platforms

| OS | Arch | Status |
|----|------|--------|
| Linux | amd64, arm64 | ✅ |
| macOS | amd64, arm64 | ✅ |
| Android | amd64, arm64 | ✅ |
| Windows | — | Community contribution needed |

---

## License

See [LICENSE](LICENSE).

V8 Gopher image based on original artwork from [Renee French](http://reneefrench.blogspot.com).
