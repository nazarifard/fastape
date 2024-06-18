package fastape

type Marshaler[V any] interface {
	Roll(v V, bs []byte) (n int, err error)
}
type Unmarshaler[V any] interface {
	Unroll(bs []byte, p *V) (n int, err error)
}
type Sizeofer[V any] interface {
	Sizeof(v V) int
}

type Tape[V any] interface {
	Marshaler[V]
	Unmarshaler[V]
	Sizeofer[V]
}
