package gencode

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func resetGenState() {
	generatedTypes = make(map[reflect.Type]string)
	generated_code = ""
	hasTime = false
	generated_array_map = make(map[int]bool)
	structCounter = 0
}

type genCodeTimeStruct struct {
	At time.Time
}

func TestGenCodeDedupesImports(t *testing.T) {
	resetGenState()
	code := GenCode("gencode", reflect.TypeOf(genCodeTimeStruct{}), "GenCodeTimeStructTape")
	if got := strings.Count(code, "import \"time\""); got != 1 {
		t.Fatalf("expected 1 time import, got %d\n%s", got, code)
	}
}
