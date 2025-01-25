package fastape

type LenTape struct{}

func (t *LenTape) Roll(size int, bs []byte) (n int, err error) {
	if size < 0 {
		return 0, ErrInvalidData
	}

	bs[0] = byte(size)
	i := 0
	for size = size >> 7; size > 0; size >>= 7 {
		bs[i] &= 0x80
		i++
		bs[i] = byte(size) & 0x7F
	}
	return i + 1, nil
}

func (t *LenTape) Unroll(bs []byte, size *int) (n int, err error) {
	if size == nil {
		return 0, ErrNilPtr
	}
	if len(bs) < 1 {
		return 0, ErrInvalidData
	}

	for n = 0; n < len(bs) && bs[n]&0x80 != 0 && n < 9; n++ {
		*size |= int(bs[n]&0x7F) << (n * 7)
	}
	if n >= 9 {
		return 0, ErrInvalidData //error
	}
	*size |= int(bs[n]&0x7F) << (n * 7)

	return n + 1, nil
}

func (t LenTape) Sizeof(size int) int {
	i := 1
	for size = size >> 7; size > 0; size >>= 7 {
		i++
	}
	return i
}
