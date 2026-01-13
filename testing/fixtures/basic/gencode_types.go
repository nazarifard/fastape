package basic

import "time"

type simpleStruct struct {
	Name string
	Age  int
}

type namedStruct struct {
	Label MyString
	When  MyTime
}

type complexStruct struct {
	Ptr    *inner
	Tags   []MyString
	Lookup map[MyString]*inner
	Codes  [2]string
}

type timeStruct struct {
	At time.Time
}
