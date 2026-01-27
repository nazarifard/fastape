# Golden Provenance

The `.golden` files under `testing/testdata` capture AST generator output for
annotated fixture packages. Legacy outputs remain in git history for comparison.

Generator
- Tool: `cmd/fastape` (AST-based, uses `internal/gencode.GenCodeFromDir`)
- Command: `fastape <dir>` (example: `fastape .`)

Annotations
- Mark types with `// fastape:generate` on its own comment group.
- Keep a blank line between the marker and any doc comment/type declaration to
  avoid it being picked up by docs.

Outputs
- `fastape_gen.go` for non-test files.
- `fastape_gen_test.go` for annotated types declared in `_test.go` files.

Goldens
- `basic.golden` from `testing/fixtures/basic`
- `nested_outer.golden` from `testing/fixtures/nested/outer`

Changelog
- AST refresh: regenerated goldens using the AST generator; removed legacy quirks
  (time import side effects, quoted struct tags, array newline formatting, named
  array byte-size bugs, and package name stripping).
- Package generator: switched to per-package output with annotation scanning;
  goldens updated to reflect `fastape_gen.go` output.
