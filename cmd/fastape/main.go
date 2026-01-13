package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nazarifard/fastape/internal/gencode"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: fastape <typeName>\nexample: fastape int")
		os.Exit(-1)
	}

	typeName := os.Args[1]
	goFile := os.Getenv("GOFILE")
	if goFile == "" {
		fmt.Fprintln(os.Stderr, "GOFILE is not set; run via go:generate")
		os.Exit(-1)
	}

	filePath := goFile
	if !filepath.IsAbs(filePath) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(-1)
		}
		filePath = filepath.Join(cwd, goFile)
	}

	requestedTapeFileName := strings.ToLower(typeName) + "_tape__gen"
	if strings.HasSuffix(goFile, "_test.go") {
		requestedTapeFileName += "_test.go"
	} else {
		requestedTapeFileName += ".go"
	}

	code, err := gencode.GenCodeFromFile(filePath, typeName, typeName+"Tape")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}

	if err := os.WriteFile(requestedTapeFileName, []byte(code), 0666); err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(-1)
	}

	if err := exec.Command("gofmt", "-w", requestedTapeFileName).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
