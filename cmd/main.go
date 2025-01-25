package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run generate_tape.go <type_name>")
		return
	}

	pkgName, err := extractPackageName(".")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Package name: %s\n", pkgName)

	typeName := os.Args[1]
	tempFileName := "__temp_tape_test.go"

	// Generate temporary test file content
	content := generateTempTest(pkgName, typeName)

	// Write to temporary test file
	if err := os.WriteFile(tempFileName, []byte(content), 0644); err != nil {
		fmt.Printf("Error writing temp file: %v\n", err)
		return
	}

	// Run the temporary test file
	cmd := exec.Command("go", "test", "-tags=generate", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running temp test file: %v\n", err)
		return
	}

	// Remove the temporary test file
	if err := os.Remove(tempFileName); err != nil {
		fmt.Printf("Error removing temp test file: %v\n", err)
		return
	}
}

func generateTempTest(pkgName, typeName string) string {
	return fmt.Sprintf(`// +build generate

package %s

import (
    "testing"
    "github.com/nazarifard/fastape"
    "reflect"
)

func TestGenCode(t *testing.T) {
    var v %s
    fastape.GenCode(reflect.TypeOf(v))
}
`, pkgName, typeName)
}

func extractPackageName(dir string) (string, error) {
	var pkgName string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" && !strings.HasSuffix(path, "_test.go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			re := regexp.MustCompile(`(?m)^\s*package\s+(\w+)\s*$`)
			matches := re.FindStringSubmatch(string(content))
			if len(matches) > 1 {
				pkgName = matches[1]
				return filepath.SkipDir // Stop after finding the first package name
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if pkgName == "" {
		return "", fmt.Errorf("package name not found")
	}
	return pkgName, nil
}
