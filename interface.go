package fastape

type Roller[V any] interface {
	Roll(v V, bs []byte) (n int, err error)
}
type Unroller[V any] interface {
	Unroll(bs []byte, p *V) (n int, err error)
}
type Sizeofer[V any] interface {
	Sizeof(v V) int
}

type Tape[V any] interface {
	Roller[V]
	Unroller[V]
	Sizeofer[V]
}
