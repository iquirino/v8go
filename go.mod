module github.com/iquirino/v8go

go 1.22

require (
	github.com/iquirino/v8go/deps/android_amd64 v0.0.0-20260821023824-8fae4e44f4be
	github.com/iquirino/v8go/deps/android_arm64 v0.0.0-20260821023824-8fae4e44f4be
	github.com/iquirino/v8go/deps/darwin_amd64 v0.0.0-20260821023824-8fae4e44f4be
	github.com/iquirino/v8go/deps/darwin_arm64 v0.0.0-20260821023824-8fae4e44f4be
	github.com/iquirino/v8go/deps/linux_amd64 v0.0.0-20260821023824-8fae4e44f4be
	github.com/iquirino/v8go/deps/linux_arm64 v0.0.0-20260821023824-8fae4e44f4be
)

// The deps/<os>_<arch> modules live in this repository. Resolve them from the
// working tree so in-repo builds/tests/fmt use the checked-in prebuilt
// libraries directly, instead of the module proxy. (replace is ignored by
// downstream consumers, which use the require versions above.)
replace (
	github.com/iquirino/v8go/deps/android_amd64 => ./deps/android_amd64
	github.com/iquirino/v8go/deps/android_arm64 => ./deps/android_arm64
	github.com/iquirino/v8go/deps/darwin_amd64 => ./deps/darwin_amd64
	github.com/iquirino/v8go/deps/darwin_arm64 => ./deps/darwin_arm64
	github.com/iquirino/v8go/deps/linux_amd64 => ./deps/linux_amd64
	github.com/iquirino/v8go/deps/linux_arm64 => ./deps/linux_arm64
)
