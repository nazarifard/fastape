# Golden Provenance

The `.golden` files under `testing/testdata` capture legacy generator output.
They are intentionally preserved as a baseline for validating the new generator.

Generator
- Path: `cmd/fastape/fastape`
- SHA256: `ee5d30c4f2c5466f9d6aef35a6f84ff8475d43e0e0537cd62f6955146d3b92c4`
- BuildID: `5c42d0c07aca7a9bb9800a89bcd39fc07af14e22`
- Git: `d62b3636c98e75fecf1f711fb822cce624422d88` (`d62b363-dirty`)

Environment
- Date (UTC): `2026-01-13T21:49:31Z`
- Go: `go1.25.5`
- GOOS/GOARCH: `linux/amd64`

Command Pattern
- Run from the fixture directory with `GOFILE` set:
  `GOFILE=<file.go> cmd/fastape/fastape <TypeName>`
- Output file: `<lower(type)>_tape__gen.go`

Known Legacy Quirks Captured
- Time import is added for any non-time struct (legacy `hasTime` side effect).
- Struct tags render as quoted string literals (e.g. `"json:\"label\""`).
- Array types are emitted with a leading newline; array fields split across lines.
- Named arrays use byte-size in the underlying type (e.g. `[32]string` for `[2]string`).
- External package types may be unqualified after package name stripping.

Additional Goldens
- `simple_struct.golden`, `named_struct.golden`, `complex_struct.golden`,
  `slice_of_named.golden`, `anon_struct.golden`, and `anon_struct_tagged.golden`
  are produced by `internal/gencode.GenCodeFromFile` using fixture types in
  `testing/fixtures/basic`. These mirror the legacy reflect generator output
  (same quirks) but are generated without the CLI.

Changelog
- AST refresh: regenerated goldens using the AST generator; removed legacy quirks
  (time import side effects, quoted struct tags, array newline formatting, named
  array byte-size bugs, and package name stripping). Legacy outputs remain in git
  history for comparison.
