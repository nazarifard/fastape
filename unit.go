package tape

import (
	"unsafe"
	unsafe_mod "unsafe"
)

type UnitTape[V any] struct {
}

func (_ UnitTape[V]) Marshal(v V, bs []byte) (n int, err error) {
	n = int(unsafe_mod.Sizeof(*(*V)(nil)))
	if len(bs) < n {
		return 0, ErrNoSpaceLeft
	}

	p := (*byte)(unsafe_mod.Pointer(&v))

	if n == 1 {
		bs[0] = *p
		return 1, nil //t 1 byte
	}

	n = copy(bs, unsafe_mod.Slice(p, n))
	return
}

func (_ UnitTape[V]) Unmarshal(bs []byte, pv *V) (n int, err error) {
	if pv == nil {
		return 0, ErrNilPtr
	}
	n = int(unsafe_mod.Sizeof(*(*V)(nil)))
	if len(bs) < n {
		return 0, ErrInvalidData
	}

	p := (*byte)(unsafe_mod.Pointer(pv))

	if n == 1 {
		*p = bs[0]
		return 1, nil //t 1 byte
	}

	n = copy(unsafe_mod.Slice(p, n), bs)
	return
}

func (_ UnitTape[V]) Sizeof(v V) int {
	return int(unsafe.Sizeof(*(*V)(nil)))
}
