package tape

type PtrTape[V any, VT Tape[V]] struct {
	vt VT
}

func (t PtrTape[V, VT]) Marshal(pv *V, bs []byte) (n int, err error) {
	if len(bs) == 0 {
		return 0, ErrNoSpaceLeft
	}
	if pv == nil {
		bs[0] = 0
		return 1, nil //OK
	}
	bs[0] = 1
	n, err = t.vt.Marshal(*pv, bs[1:])
	if err != nil {
		return 0, err
	}
	return 1 + n, nil
}

func (t PtrTape[V, VT]) Unmarshal(bs []byte, ppv **V) (n int, err error) {
	if ppv == nil {
		return 0, ErrNilPtr
	}

	if len(bs) == 0 {
		return 0, ErrInvalidData
	}
	if bs[0] == 0 {
		*ppv = nil
		return 1, nil //OK
	}
	var v V
	*ppv = &v
	n, err = t.vt.Unmarshal(bs[1:], *ppv)
	if err != nil {
		return 0, err
	}
	return 1 + n, nil
}

func (t PtrTape[V, VT]) Sizeof(pv *V) int {
	if pv == nil {
		return 1
	}
	return 1 + t.vt.Sizeof(*pv)
}
