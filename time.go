package fastape

import (
	"time"
)

type TimeTape struct {
	time UnitTape[int64]
}

func (tt TimeTape) Sizeof(t time.Time) int {
	//u := int64(0)
	return 8 //t.time.Sizeof(u)
}
func (tt TimeTape) Marshal(t time.Time, bs []byte) (n int, err error) {
	nano := t.UnixNano()
	return tt.time.Marshal(nano, bs)
}
func (t TimeTape) Unmarshal(bs []byte, p *time.Time) (n int, err error) {
	if p == nil {
		return 0, ErrNilPtr
	}

	var nano int64
	n, err = t.time.Unmarshal(bs, &nano)
	if err != nil {
		return 0, err
	}
	*p = time.Unix(0, nano)
	return n, err
}
