package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nazarifard/fastape/internal/gencode"
)

func main() {
	var tags string
	flag.StringVar(&tags, "tags", "", "build tags passed to the parser")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: fastape [flags] <dir>\nexample: fastape .")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	dir := flag.Arg(0)
	absDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}

	outputs, err := gencode.GenCodeFromDir(absDir, tags)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}

	for name, code := range outputs {
		target := filepath.Join(absDir, name)
		if err := os.WriteFile(target, []byte(code), 0666); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			os.Exit(-1)
		}
		if err := exec.Command("gofmt", "-w", target).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}
