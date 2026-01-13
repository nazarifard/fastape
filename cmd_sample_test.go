package fastape_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCmdSampleGoldens(t *testing.T) {
	t.Helper()

	repoRoot := repoRoot(t)
	baseDir := filepath.Join(repoRoot, "testing")
	basicDir := filepath.Join(baseDir, "fixtures", "basic")
	outerDir := filepath.Join(baseDir, "fixtures", "nested", "outer")
	fastapeCmd := filepath.Join(repoRoot, "cmd", "fastape")

	cases := []struct {
		name     string
		dir      string
		goFile   string
		typeName string
		golden   string
	}{
		{
			name:     "basic_edge",
			dir:      basicDir,
			goFile:   "types.go",
			typeName: "Edge",
			golden:   "basic_edge.golden",
		},
		{
			name:     "basic_fixed",
			dir:      basicDir,
			goFile:   "types.go",
			typeName: "Fixed",
			golden:   "basic_fixed.golden",
		},
		{
			name:     "basic_alias_slice",
			dir:      basicDir,
			goFile:   "types.go",
			typeName: "AliasSlice",
			golden:   "basic_alias_slice.golden",
		},
		{
			name:     "basic_named_slice",
			dir:      basicDir,
			goFile:   "types.go",
			typeName: "NamedSlice",
			golden:   "basic_named_slice.golden",
		},
		{
			name:     "basic_alias_array",
			dir:      basicDir,
			goFile:   "types.go",
			typeName: "AliasArray",
			golden:   "basic_alias_array.golden",
		},
		{
			name:     "nested_outer",
			dir:      outerDir,
			goFile:   "outer.go",
			typeName: "Outer",
			golden:   "nested_outer.golden",
		},
	}

	update := os.Getenv("UPDATE_GOLDEN") != ""
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := runFastape(t, fastapeCmd, tc.dir, tc.goFile, tc.typeName)
			goldenPath := filepath.Join(baseDir, "testdata", tc.golden)

			if update {
				if err := os.WriteFile(goldenPath, []byte(got), 0644); err != nil {
					t.Fatalf("write golden: %v", err)
				}
				return
			}

			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("read golden: %v", err)
			}
			if string(want) != got {
				t.Fatalf("golden mismatch for %s (set UPDATE_GOLDEN=1 to update)", tc.name)
			}
		})
	}
}

func TestCmdSampleRoundTripGenerated(t *testing.T) {
	repoRoot := repoRoot(t)
	fastapeCmd := filepath.Join(repoRoot, "cmd", "fastape")
	baseDir := filepath.Join(repoRoot, "testing", "fixtures")

	cases := []struct {
		name     string
		dir      string
		goFile   string
		typeName string
	}{
		{
			name:     "roundtrip",
			dir:      filepath.Join(baseDir, "roundtrip"),
			goFile:   "types.go",
			typeName: "RoundTrip",
		},
		{
			name:     "array",
			dir:      filepath.Join(baseDir, "array"),
			goFile:   "types.go",
			typeName: "ArraySample",
		},
		{
			name:     "nested",
			dir:      filepath.Join(baseDir, "nested_rt", "outer"),
			goFile:   "outer.go",
			typeName: "Outer",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			genPath := generateTapeFile(t, fastapeCmd, tc.dir, tc.goFile, tc.typeName)
			t.Cleanup(func() {
				_ = os.Remove(genPath)
			})

			cmd := exec.Command("go", "test", "-tags", "generated")
			cmd.Dir = tc.dir
			cmd.Env = os.Environ()
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("go test failed: %v\n%s", err, string(output))
			}
		})
	}
}

func runFastape(t *testing.T, fastapeCmd, dir, goFile, typeName string) string {
	t.Helper()

	genName := generatedFileName(typeName, goFile)
	genPath := filepath.Join(dir, genName)
	_ = os.Remove(genPath)
	defer func() {
		_ = os.Remove(genPath)
	}()

	cmd := exec.Command("go", "run", fastapeCmd, typeName)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOFILE="+goFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("fastape failed: %v\n%s", err, string(output))
	}

	data, err := os.ReadFile(genPath)
	if err != nil {
		t.Fatalf("read generated: %v", err)
	}

	return string(data)
}

func generateTapeFile(t *testing.T, fastapeCmd, dir, goFile, typeName string) string {
	t.Helper()

	genName := generatedFileName(typeName, goFile)
	genPath := filepath.Join(dir, genName)
	_ = os.Remove(genPath)

	cmd := exec.Command("go", "run", fastapeCmd, typeName)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOFILE="+goFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("fastape failed: %v\n%s", err, string(output))
	}

	if _, err := os.Stat(genPath); err != nil {
		t.Fatalf("generated file missing: %v", err)
	}
	return genPath
}

func generatedFileName(typeName, goFile string) string {
	name := strings.ToLower(typeName) + "_tape__gen"
	if strings.HasSuffix(goFile, "_test.go") {
		return name + "_test.go"
	}
	return name + ".go"
}

func repoRoot(t *testing.T) string {
	t.Helper()

	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve test file path")
	}
	return filepath.Dir(filePath)
}
