package fastape

type SliceTape[V any, VT Tape[V]] struct {
	len LenTape
	vt  VT
}

func (t SliceTape[V, VT]) Sizeof(s []V) int {
	if s == nil {
		return 1
	}
	sLen := len(s)

	n := t.len.Sizeof(sLen)
	for i := 0; i < sLen; i++ {
		n += t.vt.Sizeof(s[i])
	}
	return n
}

func (t SliceTape[V, VT]) Marshal(s []V, bs []byte) (n int, err error) {
	if len(bs) == 0 {
		return 0, ErrNoSpaceLeft
	}
	if s == nil {
		bs[0] = 0
		return 1, nil //OK
	}
	sLen := len(s)
	n, err = t.len.Marshal(sLen, bs)
	if err != nil {
		return 0, err
	}
	var vt VT
	var k int
	for i := 0; i < sLen; i++ {
		k, err = vt.Marshal(s[i], bs[n:])
		if err != nil {
			return 0, err
		}
		n += k
	}
	return
}

func (t SliceTape[V, VT]) Unmarshal(bs []byte, ps *[]V) (n int, err error) {
	if ps == nil {
		return 0, ErrNilPtr
	}
	if len(bs) == 0 {
		return 0, ErrInvalidData
	}
	if bs[0] == 0 {
		*ps = nil
		return 1, nil //OK
	}

	var sLen int
	n, err = t.len.Unmarshal(bs, &sLen)
	if err != nil {
		return 0, err
	}

	if *ps == nil || len(*ps) < sLen {
		*ps = make([]V, sLen)
	}

	var mv VT
	var k int
	for i := 0; i < len(*ps); i++ {
		k, err = mv.Unmarshal(bs[n:], &(*ps)[i])
		if err != nil {
			return 0, err
		}
		n += k
	}

	return
}
