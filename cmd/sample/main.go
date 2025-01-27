package main

import (
	"fmt"
	"time"
)

// firstly install fastape code generator
// go install github.com/nazarifard/fastape/cmd/fastape
type MyTime time.Time
type MyString string
type Info = []map[MyString][3]struct {
	*MyTime
	int
}

//go:generate fastape Info
func main() {
	now := MyTime(time.Now())
	var bInfo Info
	aInfo := Info{{"user1": {{&now, 123}}}}

	tape := InfoTape{} //use auto generated tape
	bs := make([]byte, tape.Sizeof(aInfo))
	tape.Roll(aInfo, bs)
	tape.Unroll(bs, &bInfo)

	if now == *aInfo[0]["user1"][0].MyTime {
		fmt.Println("time is OK")
	}
}
