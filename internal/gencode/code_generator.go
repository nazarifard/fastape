package gencode

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func GenCode(pkgName string, t reflect.Type, tapeName string) string {
	tapeTypeStr := tape(t)
	y := t.PkgPath()
	_ = y
	tapeTypeStr = strings.Replace(tapeTypeStr, pkgName+".", "", -1)
	generated_code = strings.Replace(generated_code, pkgName+".", "", -1)

	code := "package " + pkgName + "\n"
	if len(generated_array_map) > 0 {
		code += "import \"errors\"\n"
	}
	if hasTime {
		code += "import \"time\"\n"
	}
	for kt := range generatedTypes {
		pwd, _ := os.Getwd()
		//skip current local package
		if filepath.Base(pwd) == filepath.Base(kt.PkgPath()) {
			continue
		}

		switch kt.PkgPath() {
		case "", t.PkgPath(), pkgName, pwd, filepath.Base(pwd): //skip
		default:
			code += "import \"" + kt.PkgPath() + "\"\n"
		}
	}
	code += "import \"github.com/nazarifard/fastape\"" + "\n"

	code += "\n" + "type " + tapeName + " = " + tapeTypeStr

	compile_check_func_code := `
//check compileType error
var _ = func() bool{
   if false{
	 var tape %s
	 var v %s
	 _=tape.Sizeof(v)
	 _,_=tape.Roll(v, nil)
	 _,_=tape.Unroll(nil, &v)	
	 return true
    }
	return false
}()`
	code += fmt.Sprintf(compile_check_func_code, tapeName, fmt.Sprint(t))
	code = strings.Replace(code, pkgName+".", "", -1)

	code += "\n" + generated_code
	return code
}
