package gencode

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

var generatedTypes = make(map[reflect.Type]string)

func isFixedSize(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	case reflect.Array:
		return isFixedSize(t.Elem())
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			if !isFixedSize(t.Field(i).Type) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func isTimeType(t reflect.Type) bool {
	if t.Kind() != reflect.Struct {
		return false
	}
	var time time.Time
	timeType := reflect.TypeOf(time)
	for i := 0; i < t.NumField(); i++ {
		a := t.Field(i)
		b := timeType.Field(i)
		if a.Name != b.Name || a.Type != b.Type {
			return false
		}
	}
	return true
}

var hasTimeAlias bool
var hasStringAlias bool

func tape(t reflect.Type) (result string) {
	if tape, exists := generatedTypes[t]; exists {
		return tape
	}
	defer func() { generatedTypes[t] = result }()

	if isFixedSize(t) {
		return fmt.Sprintf("fastape.UnitTape[%s]", t)
	}
	if isTimeType(t) {
		if t.Name() == "time.Time" {
			return "fastape.TimeTape"
		} else {
			hasTimeAlias = true
			return fmt.Sprintf("TimeTape[%s]", t)
		}
	}
	switch t.Kind() {
	case reflect.String:
		if t.Name() == "string" {
			return "fastape.StringTape"
		} else {
			hasStringAlias = true
			return fmt.Sprintf("StringTape[%s]", t)
		}
	case reflect.Struct:
		var fields []string
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fields = append(fields, fmt.Sprintf("a%sTape %s", field.Name, tape(field.Type)))
		}
		return fmt.Sprintf("struct {\n%s\n}", strings.Join(fields, "\n"))

	case reflect.Array:
		return fmt.Sprintf("fastape.ArrayTape[%s, %s]", t.Elem(), tape(t.Elem()))
	case reflect.Slice:
		return fmt.Sprintf("fastape.SliceTape[%s, %s]", t.Elem(), tape(t.Elem()))
	case reflect.Map:
		return fmt.Sprintf("fastape.MapTape[%s, %s, %s, %s]", t.Key(), t.Elem(), tape(t.Key()), tape(t.Elem()))
	case reflect.Ptr:
		return fmt.Sprintf("fastape.PtrTape[%s, %s]", t.Elem(), tape(t.Elem()))
	default:
		fmt.Printf("fatal error: fastape doesnt supprt type: %s\n", t)
		os.Exit(-1)
		return
	}
}
