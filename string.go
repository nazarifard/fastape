package fastape

import "unsafe"

type StringTape struct {
	len LenTape
}

func (t StringTape) Sizeof(s string) int {
	if s == "" {
		return 1
	}
	strLen := len(s)
	return t.len.Sizeof(strLen) + strLen //8
}
func (t StringTape) Marshal(s string, bs []byte) (n int, err error) {
	if len(bs) == 0 {
		return 0, ErrNoSpaceLeft
	}
	if s == "" {
		bs[0] = 0
		return 1, nil //OK
	}
	strLen := len(s)
	n, err = t.len.Marshal(strLen, bs)
	if err != nil {
		return 0, err
	}
	n += copy(bs[n:], s)
	return n, err
}

func (t StringTape) Unmarshal(bs []byte, ps *string) (n int, err error) {
	if ps == nil {
		return 0, ErrNilPtr
	}
	if bs[0] == 0 {
		*ps = ""
		return 1, nil //OK
	}

	var strLen int
	n, err = t.len.Unmarshal(bs, &strLen)
	if err != nil {
		return 0, err
	}

	*ps = b2s(bs[n : n+strLen]) //todo more check
	return n + strLen, err
}

func b2s(b []byte) string {
	return unsafe.String(&b[0], len(b))
}

// https://github.com/golang/go/issues/53003#issuecomment-1145241692
// type StringHeader struct {
// 	Data *byte
// 	Len  int
// }
// func s2b(s string) []byte {
// 	header := (*StringHeader)(unsafe.Pointer(&s))
// 	bytes := *(*[]byte)(unsafe.Pointer(header))
// 	return bytes
// }
