package basic

import "time"

type MyString string
type MyTime time.Time

// fastape:generate

type AliasSlice = []MyString

// fastape:generate

type NamedSlice []MyString

type AliasMap = map[MyString]int
type NamedMap map[MyString]int

type inner struct {
	ID   int
	Note string
}

type AliasPtr = *inner
type NamedPtr *inner

// fastape:generate

type AliasArray = [2]string

// fastape:generate

type NamedArray [2]string

type Embedded struct {
	Flag bool
}

// fastape:generate

type Fixed struct {
	A int32
	B uint64
	C [4]byte
}

// fastape:generate

type Edge struct {
	Title      string
	AliasSlice AliasSlice
	NamedSlice NamedSlice
	AliasMap   AliasMap
	NamedMap   NamedMap
	AliasPtr   AliasPtr
	NamedPtr   NamedPtr
	AliasArray AliasArray
	NamedArray NamedArray
	When       time.Time
	WhenAlias  MyTime
	Inline     struct {
		Label string `json:"label"`
	}
	Embedded
}
