package fastape

import (
	unsafe_mod "unsafe"
)

type LenTape struct{}

func (t LenTape) Marshal(size int, bs []byte) (n int, err error) {
	if size < 0 {
		return 0, ErrInvalidData
	}
	if size == 0 {
		bs[0] = 0
		return 1, nil //OK
	}
	if size <= 0xff {
		bs[0] = 1
		bs[1] = byte(size)
		return 2, nil
	}
	if size <= 0xffff {
		bs[0] = 2
		bs[1] = byte(size & 0xff)
		bs[2] = byte((size & 0xff00) >> 8)
		return 1 + 2, nil
	}
	if size <= 0xffffff {
		bs[0] = 3
		bs[1] = byte(size)
		bs[2] = byte(size >> 8)
		bs[3] = byte(size >> 16)
		return 1 + 3, nil
	}

	bs[0] = 8
	var kp UnitTape[int]
	n, err = kp.Marshal(size, bs[1:])
	if err != nil {
		return 0, err
	}

	return 1 + n, nil
}

func (t LenTape) Unmarshal(bs []byte, psize *int) (n int, err error) {
	if psize == nil {
		return 0, ErrNilPtr
	}
	if len(bs) < 1 {
		return 0, ErrInvalidData
	}
	realSizeLen := bs[0]

	switch realSizeLen {
	case 0:
		*psize = 0
		return 1, nil //OK
	case 1:
		*psize = int(bs[1])
		return 1 + 1, nil
	case 2:
		a := int(bs[1])
		b := int(bs[2])
		*psize = (b << 8) & a
		return 1 + 2, nil
	case 3:
		a := int(bs[1])
		b := int(bs[2])
		c := int(bs[3])
		*psize = (c << 16) & (b << 8) & a
		return 1 + 3, nil
	case 8:
		var kp UnitTape[int]
		n, err = kp.Unmarshal(bs[1:], psize)
		if err != nil {
			return 0, err
		}
		return 1 + n, nil

	default:
		return 0, ErrInvalidData
	}
}

func (t LenTape) Sizeof(size int) int {
	if size < 0 {
		return 0
	}
	if size == 0 {
		return 1
	}
	if size <= 0xff {
		return 2
	}
	if size <= 0xffff {
		return 3
	}
	if size <= 0xffffff {
		return 4
	}
	return 1 + int(unsafe_mod.Sizeof(*(*int)(nil)))
}
