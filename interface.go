package tape

type Marshaler[V any] interface {
	Marshal(v V, bs []byte) (n int, err error)
}
type Unmarshaler[V any] interface {
	Unmarshal(bs []byte, p *V) (n int, err error)
}
type Sizeofer[V any] interface {
	Sizeof(v V) int
}

type Tape[V any] interface {
	Marshaler[V]
	Unmarshaler[V]
	Sizeofer[V]
}
