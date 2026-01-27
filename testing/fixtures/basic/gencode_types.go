package basic

import "time"

// fastape:generate

type simpleStruct struct {
	Name string
	Age  int
}

// fastape:generate

type namedStruct struct {
	Label MyString
	When  MyTime
}

// fastape:generate

type complexStruct struct {
	Ptr    *inner
	Tags   []MyString
	Lookup map[MyString]*inner
	Codes  [2]string
}

// fastape:generate

type timeStruct struct {
	At time.Time
}
