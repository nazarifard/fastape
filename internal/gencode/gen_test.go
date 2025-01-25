package gencode

import (
	"reflect"
	"testing"
)

func TestGen(t *testing.T) {
	type S struct {
		B struct {
			x int
			y string
		}
		z *int
	}
	GenCode("main", reflect.TypeOf(S{}))
}
