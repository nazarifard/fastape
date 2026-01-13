package fastape_test

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nazarifard/fastape/internal/gencode"
)

func TestGenCodeGoldens(t *testing.T) {
	t.Helper()

	repoRoot := repoRoot(t)
	fixtureFile := filepath.Join(repoRoot, "testing", "fixtures", "basic", "types.go")
	goldenDir := filepath.Join(repoRoot, "testing", "testdata")

	cases := []struct {
		name     string
		typeName string
		tapeName string
		golden   string
		legacy   bool
	}{
		{
			name:     "fixed_struct",
			typeName: "Fixed",
			tapeName: "FixedTape",
			golden:   "basic_fixed.golden",
			legacy:   true,
		},
		{
			name:     "simple_struct",
			typeName: "simpleStruct",
			tapeName: "SimpleStructTape",
			golden:   "simple_struct.golden",
		},
		{
			name:     "named_struct",
			typeName: "namedStruct",
			tapeName: "NamedStructTape",
			golden:   "named_struct.golden",
		},
		{
			name:     "complex_struct",
			typeName: "complexStruct",
			tapeName: "ComplexStructTape",
			golden:   "complex_struct.golden",
		},
		{
			name:     "slice_of_named",
			typeName: "[]MyString",
			tapeName: "MyStringSliceTape",
			golden:   "slice_of_named.golden",
		},
		{
			name:     "anon_struct",
			typeName: "struct { Title string; Scores []int }",
			tapeName: "AnonStructTape",
			golden:   "anon_struct.golden",
		},
		{
			name:     "anon_struct_tagged",
			typeName: "struct { Embedded; Name string `json:\"name\"` }",
			tapeName: "AnonStructTaggedTape",
			golden:   "anon_struct_tagged.golden",
		},
		{
			name:     "named_slice",
			typeName: "NamedSlice",
			tapeName: "NamedSliceTape",
			golden:   "basic_named_slice.golden",
			legacy:   true,
		},
		{
			name:     "alias_array",
			typeName: "AliasArray",
			tapeName: "AliasArrayTape",
			golden:   "basic_alias_array.golden",
			legacy:   true,
		},
		{
			name:     "edge_struct",
			typeName: "Edge",
			tapeName: "EdgeTape",
			golden:   "basic_edge.golden",
			legacy:   true,
		},
	}

	update := os.Getenv("UPDATE_GOLDEN") != ""
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := gencode.GenCodeFromFile(fixtureFile, tc.typeName, tc.tapeName)
			if err != nil {
				t.Fatalf("generate code: %v", err)
			}
			formatted := formatCode(t, []byte(got))

			goldenPath := filepath.Join(goldenDir, tc.golden)
			if update {
				if tc.legacy {
					return
				}
				if err := os.MkdirAll(goldenDir, 0755); err != nil {
					t.Fatalf("create golden dir: %v", err)
				}
				if err := os.WriteFile(goldenPath, formatted, 0644); err != nil {
					t.Fatalf("write golden: %v", err)
				}
				return
			}

			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("read golden: %v", err)
			}
			if string(want) != string(formatted) {
				t.Fatalf("golden mismatch for %s (set UPDATE_GOLDEN=1 to update)", tc.name)
			}
		})
	}
}

func TestGenCodeDedupesImports(t *testing.T) {
	t.Helper()

	repoRoot := repoRoot(t)
	fixtureFile := filepath.Join(repoRoot, "testing", "fixtures", "basic", "types.go")

	code, err := gencode.GenCodeFromFile(fixtureFile, "timeStruct", "TimeStructTape")
	if err != nil {
		t.Fatalf("generate code: %v", err)
	}
	if got := strings.Count(code, "import \"time\""); got != 1 {
		t.Fatalf("expected 1 time import, got %d\n%s", got, code)
	}
}

func formatCode(t *testing.T, src []byte) []byte {
	t.Helper()

	formatted, err := format.Source(src)
	if err != nil {
		return src
	}
	return formatted
}
