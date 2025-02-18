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
		fmt.Println("Usage: fastape <typeName>\nexample: fastape int")
		os.Exit(-1)
	}
	var err error

	pkgName := os.Getenv("GOPACKAGE")
	if pkgName == "" {
		pkgName, err = grepPackageName()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	}
	typeName := os.Args[1]
	tempMainFile := "temp_main__gen_test.go"

	requestedTapeFileName := strings.ToLower(typeName) + "_tape__gen"
	revokedFileName := os.Getenv("GOFILE")
	Len := len(revokedFileName)
	if Len > 8 && string(revokedFileName[Len-8:]) == "_test.go" {
		requestedTapeFileName += "_test.go"
	} else {
		requestedTapeFileName += ".go"
	}

	// Generate temporary test file content
	content := generateTempTest(pkgName, typeName, requestedTapeFileName)

	// Write to temporary test file
	if err := os.WriteFile(tempMainFile, []byte(content), 0644); err != nil {
		fmt.Printf("Error writing temp file: %v\n", err)
		return
	}

	// Run the temporary test file
	devNull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}
	defer devNull.Close()

	//mandatory command
	//that will generate required output
	cmd := exec.Command("go", "get", "github.com/nazarifard/fastape")
	cmd.Stdout = devNull
	cmd.Stderr = devNull
	cmd.Run()

	//mandatory command
	//that will generate required output
	cmd = exec.Command("go", "test", "-run=TapeGenCode", ".")
	cmd.Stdout = devNull
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return
	}
	defer os.Remove(tempMainFile)

	// goimports
	cmd = exec.Command("goimports", "-w", requestedTapeFileName)
	cmd.Stdout = devNull
	cmd.Stderr = devNull
	cmd.Run()

	//gofmt
	cmd = exec.Command("gofmt", "-w", requestedTapeFileName)
	cmd.Stdout = devNull
	cmd.Stderr = devNull
	cmd.Run()
}

func generateTempTest(pkgName, typeName, fileName string) string {
	return fmt.Sprintf(`
package %s
import (
    "os"
    "testing"
    "github.com/nazarifard/fastape"
    "reflect"
)

func TestTapeGenCode(t *testing.T) {
    var v %s
    code := fastape.GenCode("%s", reflect.TypeOf(v), "%s")
    err := os.WriteFile("%s", []byte(code), 0666)
    _ = err
}
`, pkgName, typeName, pkgName, typeName+"Tape", fileName)
}

// grepPackageName searches through *.go files in the current directory
// and returns the first valid package name it finds.
func grepPackageName() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	// Compile regexes for package declaration and comments/whitespace
	packageRegex := regexp.MustCompile(`^\s*package\s+(\w+)`)
	commentOrWhitespaceRegex := regexp.MustCompile(`^\s*(//.*)?$`)

	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return "", fmt.Errorf("failed to list .go files: %v", err)
	}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return "", fmt.Errorf("failed to read file %s: %v", file, err)
		}
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if commentOrWhitespaceRegex.MatchString(line) {
				continue
			}
			matches := packageRegex.FindStringSubmatch(line)
			if len(matches) > 1 {
				return matches[1], nil
			}
			break
		}
	}
	return "", fmt.Errorf("no valid package declaration found in current folder")
}
