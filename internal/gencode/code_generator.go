package gencode

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

func GenCode(pkgName string, t reflect.Type) {

	tapeName := t.Name() + "Tape"
	tapeTypeStr := tape(t)

	code := "package " + pkgName + "\n"
	code += "import \"github.com/nazarifard/fastape\"" + "\n"
	for kt := range generatedTypes {
		if kt.PkgPath() != "" && kt.PkgPath() != t.PkgPath() {
			code += "import \"" + kt.PkgPath() + "\"\n"
		}
	}
	if hasStringAlias || hasTimeAlias {
		if !strings.Contains(code, "\"unsafe\"") {
			code += "import \"unsafe\"\n"
		}
	}
	if hasStringAlias {
		if !strings.Contains(code, "\"time\"") {
			code += "import \"time\"\n"
		}
	}

	if t.PkgPath() == pkgName {
		//remove current pkgName
		tapeTypeStr = strings.Replace(tapeTypeStr, pkgName+".", "", -1)
	}

	code += fmt.Sprintf("type %s %s\n\n", tapeName, tapeTypeStr)
	if strings.HasPrefix(tapeTypeStr, "struct") {
		code += generateStructMethods(tapeTypeStr, tapeName)
	} else {
		code += generateNonStructMethods(tapeName, tapeTypeStr, reflect.TypeOf(t).Name())
	}

	if hasTimeAlias {
		code += "\n" + timeTapeCode + "\n"
	}
	if hasStringAlias {
		code += "\n" + stringTapeCode + "\n"
	}

	filename := "__" + tapeName + "_gen.go"
	os.WriteFile(filename, []byte(code), 0644)
}

var timeTapeCode = `
type TimeTape[T any] struct {
	tt fastape.TimeTape
}
func (tt TimeTape[T]) Sizeof(t T) int {
	t0 := (*time.Time)(unsafe.Pointer(&t))
	return tt.tt.Sizeof(*t0)
}
func (tt TimeTape[T]) Roll(t T, bs []byte) (n int, err error) {
	t0 := (*time.Time)(unsafe.Pointer(&t))
	return tt.tt.Roll(*t0, bs)
}
func (tt TimeTape[T]) Unroll(bs []byte, p *T) (n int, err error) {
	p0 := (*time.Time)(unsafe.Pointer(p))
	return tt.tt.Unroll(bs, p0)
}
`

var stringTapeCode = `
type StringTape[T ~string] struct {
	st fastape.StringTape
}
func (tt StringTape[T]) Sizeof(t T) int {
	return tt.st.Sizeof(string(t))
}
func (tt StringTape[T]) Roll(t T, bs []byte) (n int, err error) {
	return tt.st.Roll(string(t), bs)
}
func (tt StringTape[T]) Unroll(bs []byte, p *T) (n int, err error) {
	return tt.st.Unroll(bs, (*string)(unsafe.Pointer(p)))
}
`
