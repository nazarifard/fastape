package fastape

type LenTape struct{}

func (t *LenTape) Roll(size int, bs []byte) (n int, err error) {
	if size < 0 {
		return 0, ErrInvalidData
	}
	if len(bs) == 0 {
		return 0, ErrNoSpaceLeft
	}

	// Unsigned varint (LEB128-style): 7 bits per byte, continuation bit 0x80.
	u := uint(size)
	i := 0
	for u >= 0x80 {
		if i >= len(bs) {
			return 0, ErrNoSpaceLeft
		}
		bs[i] = byte(u&0x7F) | 0x80
		u >>= 7
		i++
	}
	if i >= len(bs) {
		return 0, ErrNoSpaceLeft
	}
	bs[i] = byte(u)
	return i + 1, nil
}

func (t *LenTape) Unroll(bs []byte, size *int) (n int, err error) {
	if size == nil {
		return 0, ErrNilPtr
	}
	if len(bs) < 1 {
		return 0, ErrInvalidData
	}

	var out uint64
	var shift uint
	for n = 0; n < len(bs) && n < 10; n++ {
		b := bs[n]
		out |= uint64(b&0x7F) << shift
		if b&0x80 == 0 {
			maxInt := int(^uint(0) >> 1)
			if out > uint64(maxInt) {
				return 0, ErrInvalidData
			}
			*size = int(out)
			return n + 1, nil
		}
		shift += 7
	}
	return 0, ErrInvalidData
}

func (t LenTape) Sizeof(size int) int {
	if size < 0 {
		return 0
	}
	i := 1
	for size = size >> 7; size > 0; size >>= 7 {
		i++
	}
	return i
}
