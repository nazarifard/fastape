package fastape

import (
	"unsafe"
	unsafe_mod "unsafe"
)

type UnitTape[V any] struct {
}

func (_ UnitTape[V]) Roll(v V, bs []byte) (n int, err error) {
	n = int(unsafe_mod.Sizeof(v))
	if len(bs) < n {
		return 0, ErrNoSpaceLeft
	}

	p := (*V)(unsafe_mod.Pointer(&bs[0]))
	*p = v
	return
}

func (_ UnitTape[V]) Unroll(bs []byte, pv *V) (n int, err error) {
	if pv == nil {
		return 0, ErrNilPtr
	}
	n = int(unsafe_mod.Sizeof(*pv))
	if len(bs) < n {
		return 0, ErrInvalidData
	}

	p := (*V)(unsafe_mod.Pointer(&bs[0]))
	*pv = *p
	return
}

func (_ UnitTape[V]) Sizeof(v V) int {
	return int(unsafe.Sizeof(v))
}
