package fastape

import (
	"reflect"

	"github.com/nazarifard/fastape/internal/gencode"
)

func GenCode(pkgName string, t reflect.Type) {
	gencode.GenCode(pkgName, t)
}
