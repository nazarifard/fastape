package main

import "time"

// firstly install fastape code generator
// go install github.com/nazarifard/fastape/cmd/fastape
type MyTime time.Time
type MyString string
type Info = []map[MyString][3]struct {
	*MyTime
	int
}

//go:generate fastape Info
