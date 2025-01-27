package main

import (
	"errors"
	"time"

	"github.com/nazarifard/fastape"
)

type InfoTape = fastape.SliceTape[map[MyString][3]struct {
	*MyTime
	int
}, fastape.MapTape[MyString, [3]struct {
	*MyTime
	int
}, fastape.NamedTape[MyString, string, fastape.StringTape],
	Array_3[[3]struct {
		*MyTime
		int
	}, struct {
		*MyTime
		int
	}, struct_0001Tape]]]

// check compileType error
var _ = func() bool {
	if false {
		var tape InfoTape
		var v []map[MyString][3]struct {
			*MyTime
			int
		}
		_ = tape.Sizeof(v)
		_, _ = tape.Roll(v, nil)
		_, _ = tape.Unroll(nil, &v)
		return true
	}
	return false
}()

type Array_3[A interface{ ~[3]V }, V any, VT fastape.Tape[V]] struct{ vt VT }

func (array Array_3[A, V, VT]) Sizeof(a A) int {
	size := 0
	for i := range a {
		size += array.vt.Sizeof(a[i])
	}
	return size
}

func (array Array_3[A, V, VT]) Roll(a A, bs []byte) (n int, err error) {
	for i := range len(a) {
		m, err := array.vt.Roll(a[i], bs[n:])
		if err != nil {
			return 0, err
		}
		n += m
	}
	return n, nil
}

func (array Array_3[A, V, VT]) Unroll(bs []byte, p *A) (n int, err error) {
	if p == nil {
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

type struct_0001 = struct {
	*MyTime
	int
}

func (t struct_0001Tape) Sizeof(v struct_0001) int {
	size := 0
	size += t.aMyTime.Sizeof(v.MyTime)
	size += t.aint.Sizeof(v.int)
	return size
}

func (t struct_0001Tape) Roll(v struct_0001, bs []byte) (n int, err error) {
	var m int
	m, err = t.aMyTime.Roll(v.MyTime, bs[n:])
	if err != nil {
		return
	}
	n += m
	m, err = t.aint.Roll(v.int, bs[n:])
	if err != nil {
		return
	}
	n += m
	return
}

func (t struct_0001Tape) Unroll(bs []byte, v *struct_0001) (n int, err error) {
	var m int
	m, err = t.aMyTime.Unroll(bs[n:], &v.MyTime)
	if err != nil {
		return
	}
	n += m
	m, err = t.aint.Unroll(bs[n:], &v.int)
	if err != nil {
		return
	}
	n += m
	return
}

type struct_0001Tape struct {
	aMyTime fastape.PtrTape[MyTime, fastape.NamedTape[MyTime, time.Time, fastape.TimeTape]]
	aint    fastape.UnitTape[int]
}
