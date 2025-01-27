package gencode

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

var generatedTypes = make(map[reflect.Type]string)
var hasTime = false

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
			hasTime = true
			return false
		}
	}
	return true
}

var generated_code string

func UnderlyingType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Ptr:
		return fmt.Sprintf("*%s", t.Elem())
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", t.Key(), t.Elem())
	case reflect.Slice:
		return fmt.Sprintf("[]%s", t.Elem())
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", t.Size(), t.Elem())
	case reflect.Struct:
		s := ""
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			s += fmt.Sprintf("%s %s;", field.Name, field.Type)
		}
		return "struct {" + s + "}"
	case reflect.String:
		return "string"
	default:
		return fmt.Sprint(t)
	}
}

var structCounter = 0
var generated_array_map = make(map[int]bool)

func tape(t reflect.Type) (result string) {
	if tape, exists := generatedTypes[t]; exists {
		return tape
	}
	underType := ""
	defer func() {
		if underType != "" {
			result = fmt.Sprintf("fastape.NamedTape[%s, %s, %s]", t, underType, result)
		}
		generatedTypes[t] = result
	}()

	if isFixedSize(t) {
		return fmt.Sprintf("fastape.UnitTape[%s]", t)
	}
	if isTimeType(t) {
		if t.Name() != "time.Time" {
			underType = "time.Time"
		}
		return "fastape.TimeTape"
	}
	//if t.Kind() != reflect.Struct && t.Name() != "" && fmt.Sprint(t.Name()) != fmt.Sprint(t) {
	if t.Name() != "" && fmt.Sprint(t.Name()) != fmt.Sprint(t) {
		underType = UnderlyingType(t)
	}

	switch t.Kind() {
	case reflect.String:
		return "fastape.StringTape"
	case reflect.Struct:
		structName := t.Name()
		if t.Name() == "" {
			structCounter++
			structName = "struct_" + fmt.Sprintf("%04d", structCounter)
			generated_code += "\ntype " + structName + " = " + t.String() + "\n"
		}
		tapeName := structName + "Tape"

		sizeof_code := fmt.Sprintf("\nfunc (t %s) Sizeof(v %s) int {\n size := 0", tapeName, structName)
		roll_code := fmt.Sprintf("\nfunc (t %s) Roll(v %s, bs []byte) (n int, err error) {	var m int", tapeName, structName)
		unroll_code := fmt.Sprintf("\nfunc (t %s) Unroll(bs []byte, v *%s) (n int, err error) { var m int", tapeName, structName)

		var tapeFields []string
		//var realFields []string
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			// fieldName := field.Name
			// if field.Anonymous {
			// 	fieldName = ""
			// }
			tapeFields = append(tapeFields, fmt.Sprintf("a"+"%s %s", field.Name, tape(field.Type)))
			//realFields = append(realFields, fmt.Sprintf("%s %s", fieldName, field.Type))

			sizeof_code += fmt.Sprintf("\nsize += t.a%s.Sizeof(v.%s)", field.Name, field.Name)
			roll_code += fmt.Sprintf("\nm,err=t.a%s.Roll(v.%s,bs[n:]); if err!=nil{return};n+=m", field.Name, field.Name)
			unroll_code += fmt.Sprintf("\nm, err = t.a%s.Unroll(bs[n:], &v.%s); if err!=nil{return};n+=m", field.Name, field.Name)
		}
		sizeof_code += "\nreturn size}\n"
		roll_code += "\nreturn}\n"
		unroll_code += "\nreturn}\n"

		generated_code += sizeof_code
		generated_code += roll_code
		generated_code += unroll_code

		newTapeStruct := fmt.Sprintf("struct {\n%s\n}", strings.Join(tapeFields, "\n"))
		generated_code += "\ntype " + tapeName + " " + newTapeStruct + "\n"
		underType = ""
		return tapeName
	case reflect.Array:
		if !generated_array_map[t.Len()] {
			generated_code += GenArrayCode(t.Len())
			generated_array_map[t.Len()] = true
		}
		return fmt.Sprintf("\nArray_%d[%s, %s, %s]", t.Len(), t, t.Elem(), tape(t.Elem()))
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

func GenArrayCode(n int) string {
	const code_Array = `

type Array_%d[A interface{ ~[%d]V }, V any, VT fastape.Tape[V]] struct{ vt VT }

func (array Array_%d[A, V, VT]) Sizeof(a A) int {
	size := 0
	for i := range a {
		size += array.vt.Sizeof(a[i])
	}
	return size
}

func (array Array_%d[A, V, VT]) Roll(a A, bs []byte) (n int, err error) {
	for i := range len(a) {
		m, err := array.vt.Roll(a[i], bs[n:])
		if err != nil {
			return 0, err
		}
		n += m
	}
	return n, nil
}

func (array Array_%d[A, V, VT]) Unroll(bs []byte, p *A) (n int, err error) {
    if p==nil {
	return 0, errors.New("target pointer is nil")
	}
	for i := range len(*p) {
		m, err := array.vt.Unroll(bs[n:], &(*p)[i])
		if err != nil {
			return 0, err
		}
		n += m
	}
	return n, nil
}

`
	return fmt.Sprintf(code_Array, n, n, n, n, n)
}
