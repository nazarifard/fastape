package fastape

type MapTape[K comparable, V any, KT Tape[K], VT Tape[V]] struct {
	lenTape LenTape
	kt      KT
	vt      VT
}

func (t MapTape[K, V, KT, VT]) Sizeof(pm map[K]V) int {
	if pm == nil {
		return 1
	}

	mLen := len(pm)
	n := t.lenTape.Sizeof(mLen)
	for k, v := range pm {
		n += t.kt.Sizeof(k) + t.vt.Sizeof(v)
	}
	return n
}

func (mt MapTape[K, V, KT, VT]) Marshal(m map[K]V, bs []byte) (n int, err error) {
	if len(bs) == 0 {
		return 0, ErrNoSpaceLeft
	}
	if m == nil {
		bs[0] = 0
		return 1, nil //OK
	}

	n, err = mt.lenTape.Marshal(len(m), bs)
	if err != nil {
		return 0, err
	}

	t := 0
	for k, v := range m {
		t, err = mt.kt.Marshal(k, bs[n:])
		if err != nil {
			return 0, err
		}
		n += t
		t, err = mt.vt.Marshal(v, bs[n:])
		if err != nil {
			return 0, err
		}
		n += t
	}
	return n, err
}
func (mt MapTape[K, V, KT, VT]) Unmarshal(bs []byte, pm *map[K]V) (n int, err error) {
	if pm == nil {
		return 0, ErrNilPtr
	}
	if len(bs) == 0 {
		return 0, ErrInvalidData
	}
	if bs[0] == 0 {
		*pm = nil
		return 1, nil //OK
	}

	var k K
	var size int
	n, err = mt.lenTape.Unmarshal(bs, &size)
	if err != nil {
		return 0, err
	}
	if size == 0 {
		*pm = nil
		return n, nil //OK
	}
	if *pm == nil {
		*pm = make(map[K]V)
	}
	for i := 0; i < size; i++ {
		t, err := mt.kt.Unmarshal(bs[n:], &k)
		if err != nil {
			return 0, err
		}
		n += t

		var v V
		t, err = mt.vt.Unmarshal(bs[n:], &v)
		if err != nil {
			return 0, err
		}
		(*pm)[k] = v
		n += t
	}
	return n, err
}
