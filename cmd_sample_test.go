package fastape_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCmdSampleGoldens(t *testing.T) {
	t.Helper()

	repoRoot := repoRoot(t)
	baseDir := filepath.Join(repoRoot, "testing")
	basicDir := filepath.Join(baseDir, "fixtures", "basic")
	outerDir := filepath.Join(baseDir, "fixtures", "nested", "outer")
	fastapeRoot := repoRoot
	cases := []struct {
		name   string
		dir    string
		golden string
	}{
		{
			name:   "basic",
			dir:    basicDir,
			golden: "basic.golden",
		},
		{
			name:   "nested_outer",
			dir:    outerDir,
			golden: "nested_outer.golden",
		},
	}

	update := os.Getenv("UPDATE_GOLDEN") != ""
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			generatePackageFiles(t, fastapeRoot, tc.dir, "")
			got := readGeneratedFile(t, tc.dir, "fastape_gen.go")
			goldenPath := filepath.Join(baseDir, "testdata", tc.golden)

			if update {
				if err := os.WriteFile(goldenPath, got, 0644); err != nil {
					t.Fatalf("write golden: %v", err)
				}
				return
			}

			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("read golden: %v", err)
			}
			if string(want) != string(got) {
				t.Fatalf("golden mismatch for %s (set UPDATE_GOLDEN=1 to update)", tc.name)
			}
		})
	}
}

func TestCmdSampleRoundTripGenerated(t *testing.T) {
	repoRoot := repoRoot(t)
	baseDir := filepath.Join(repoRoot, "testing", "fixtures")
	fastapeRoot := repoRoot

	cases := []struct {
		name string
		dir  string
	}{
		{
			name: "roundtrip",
			dir:  filepath.Join(baseDir, "roundtrip"),
		},
		{
			name: "array",
			dir:  filepath.Join(baseDir, "array"),
		},
		{
			name: "nested",
			dir:  filepath.Join(baseDir, "nested_rt", "outer"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			generatePackageFiles(t, fastapeRoot, tc.dir, "")
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

func generatePackageFiles(t *testing.T, fastapeRoot, dir, tags string) []string {
	t.Helper()

	removeFiles([]string{
		filepath.Join(dir, "fastape_gen.go"),
		filepath.Join(dir, "fastape_gen_test.go"),
	})

	runFastape(t, fastapeRoot, dir, tags)

	paths := []string{
		filepath.Join(dir, "fastape_gen.go"),
		filepath.Join(dir, "fastape_gen_test.go"),
	}
	var existing []string
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			existing = append(existing, path)
		}
	}
	if len(existing) == 0 {
		t.Fatalf("generated files missing in %s", dir)
	}
	return existing
}

func runFastape(t *testing.T, fastapeRoot, dir, tags string) {
	t.Helper()

	args := []string{"run", "./cmd/fastape"}
	if tags != "" {
		args = append(args, "-tags", tags)
	}
	args = append(args, dir)

	cmd := exec.Command("go", args...)
	cmd.Dir = fastapeRoot
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("fastape failed: %v\n%s", err, string(output))
	}
}

func readGeneratedFile(t *testing.T, dir, name string) []byte {
	t.Helper()

	path := filepath.Join(dir, name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read generated: %v", err)
	}
	return data
}

func removeFiles(paths []string) {
	for _, path := range paths {
		_ = os.Remove(path)
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()

	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve test file path")
	}
	return filepath.Dir(filePath)
}
