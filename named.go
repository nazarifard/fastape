package fastape

import (
	"unsafe"
)

//example
//type P *int

type NamedTape[N any, T any, TTape Tape[T]] struct{ ttape TTape }

func (t NamedTape[N, T, TTape]) Sizeof(v N) int {
	return t.ttape.Sizeof(*(*T)(unsafe.Pointer(&v)))
}

func (t NamedTape[N, T, TTape]) Roll(v N, bs []byte) (n int, err error) {
	return t.ttape.Roll(*(*T)(unsafe.Pointer(&v)), bs)
}

func (t NamedTape[N, T, TTape]) Unroll(bs []byte, p *N) (n int, err error) {
	return t.ttape.Unroll(bs, (*T)(unsafe.Pointer(p)))
}
