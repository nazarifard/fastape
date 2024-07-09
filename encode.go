package fastape

import "github.com/nazarifard/syncpool"

// type Tap[V any] interface {
// 	Encode(v V) (buf syncpool.Buffer, err error)
// 	Decode(bs []byte) (v *V, n int, err error)
// 	Free(*V)
// }

func NewMarshalTap[V any, T Tape[V]](tape T, pool *syncpool.Pool[V]) *tap[V, T] {
	return &tap[V, T]{
		tape:       tape,
		vPool:      pool,
		bufferPool: syncpool.NewBufferPool(),
	}
}

type tap[V any, T Tape[V]] struct {
	tape       T
	vPool      *syncpool.Pool[V]
	bufferPool syncpool.BufferPool
}

func (t *tap[V, T]) Free(v *V) {
	t.vPool.Put(v)
}

func (t *tap[V, T]) Encode(v V) (buf syncpool.Buffer, err error) {
	size := t.tape.Sizeof(v)
	buf = t.bufferPool.Get(size)
	_, err = t.tape.Roll(v, buf.Bytes())
	if err != nil {
		buf.Free()
	}
	return
}

func (t *tap[V, T]) Decode(bs []byte) (v *V, n int, err error) {
	v = t.vPool.Get()
	n, err = t.tape.Unroll(bs, v)
	return v, n, err
}

// func Encode[V any, T Tape[V]](v V, t T) (buf syncpool.Buffer, err error) {
// 	buf = bufferPool.Get(t.Sizeof(v))
// 	_, err = t.Roll(v, buf.Bytes())
// 	if err != nil {
// 		buf.Free()
// 		return buf, err
// 	}
// 	return
// }
// func Decode[V any, T Tape[V]](bs []byte, t T, v *V) (n int, err error) {
// 	//n, err = t.Unroll(bs, v)
// 	return
// }
