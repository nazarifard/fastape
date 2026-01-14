# Fastape

Fastape is a Go serializer designed for maximum throughput and minimal allocations when you control both ends of the bytes.

It’s a great fit for in-memory databases, caches, and hot paths where you want:
- Precise buffer sizing (`Sizeof`) to avoid reallocations.
- Very fast encode/decode (`Roll`/`Unroll`).
- Generated, type-specific codecs (no reflection at runtime).

This is not intended as a general-purpose, portable wire format. See Safety & Portability.

## Quickstart (code generation)

Install the generator binary:

```sh
go install github.com/nazarifard/fastape/cmd/fastape@latest
```

Annotate the types you want to generate tapes for. The marker must be in its own comment group with a blank line before the type:

```go
package example

import "time"

type MyString string
type MyTime time.Time

// fastape:generate

type Info = []map[MyString][3]struct {
	*MyTime
	Count int
}

//go:generate fastape .
```

Generate and use the tape:

```sh
go generate ./...
```

```go
var src, dst Info
var tape InfoTape // generated into fastape_gen.go

buf := make([]byte, tape.Sizeof(src))
n, err := tape.Roll(src, buf)
if err != nil {
	// handle
}

_, err = tape.Unroll(buf[:n], &dst)
if err != nil {
	// handle
}
```

## Core API

Every tape implements:

```go
type Tape[V any] interface {
	Sizeof(v V) int
	Roll(v V, bs []byte) (n int, err error)    // encode
	Unroll(bs []byte, v *V) (n int, err error) // decode
}
```

`Sizeof` is the key: it enables exact preallocation so encode paths can be 0-alloc.

## What gets generated

The generator emits a `...Tape` type for each annotated type, built from a small set of composable primitives:
- `StringTape` for strings
- `TimeTape` for `time.Time`
- `PtrTape`, `SliceTape`, `MapTape` for containers
- `UnitTape` for fixed-size values (fast, but has portability tradeoffs)
- `NamedTape` wrappers for named/alias types

Complex shapes are supported (examples are under `testing/fixtures`):
- alias + named types
- nested structs (including across packages)
- pointers, arrays, slices, maps
- `time.Time` and `type MyTime time.Time`
- inline anonymous structs (including tags)

## Safety & Portability (read this)

Fastape is optimized for trusted bytes and in-memory usage.

Trusted bytes means: the data may represent user data, but no untrusted party can tamper with the encoded bytes (not hardened against malicious/corrupt input).

Portability notes:
- There is no defined endianness for numeric types: many primitives use `UnitTape`, which writes Go’s in-memory representation.
- Fixed-size structs may also be encoded via raw memory for speed, which can include padding bytes.
- `int`/`uint` sizes are architecture-dependent, so blobs may not be readable across 32-bit vs 64-bit builds.
- Because `UnitTape` relies on Go’s in-memory representation, compatibility across Go versions/toolchains is not guaranteed.
- Map encoding order is not stable (Go map iteration is randomized), so byte-for-byte determinism is not guaranteed.

Practical takeaway:
- For a single-machine database/cache where the same program reads what it wrote, this works well.
- If you need a stable, portable file/wire format, Fastape is not that today.

## Buffer and value pooling (optional)

If you want to reuse buffers and decoded values, `NewMarshalTap` composes a tape with pools:

```go
tap := fastape.NewMarshalTap(tape, valuePool)
buf, err := tap.Encode(v)
// ...
decoded, _, err := tap.Decode(buf.Bytes())
tap.Free(decoded)
buf.Free()
```

## CLI

The generator is a small CLI:

```sh
fastape [flags] <dir>
```

Supported flags:
- `-tags`: build tags passed to the parser (when loading the package)
- `-tapable`: also emit `Tape() fastape.Tape[T]` methods for generated types

## Benchmarks

See `BENCHMARK.md` for local `go test -bench` snapshots and a goserbench comparison.

## License

MIT
