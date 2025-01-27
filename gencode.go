package fastape

import (
	"reflect"

	"github.com/nazarifard/fastape/internal/gencode"
)

func GenCode(pkgName string, t reflect.Type, tapeName string) string {
	return gencode.GenCode(pkgName, t, tapeName)
}
