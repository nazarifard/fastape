package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: fastape <typeName>\nexample: fastape int")
		return
	}

	pkgName := os.Getenv("GOPACKAGE")
	if pkgName == "" {
		pkgName = "main"
	}
	typeName := os.Args[1]
	tempFileName := "temp_tape_test.go"

	// Generate temporary test file content
	requestedTapeFileName := strings.ToLower(typeName) + "_tape__gen.go"
	content := generateTempTest(pkgName, typeName, requestedTapeFileName)

	// Write to temporary test file
	if err := os.WriteFile(tempFileName, []byte(content), 0644); err != nil {
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
	cmd := exec.Command("go", "test", "-tags=generate", "-run=TapeGenCode", ".")
	cmd.Stdout = devNull
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return
	}
	defer os.Remove(tempFileName)

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
	return fmt.Sprintf(`// +build generate

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
