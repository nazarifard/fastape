package fastape

type Tapable[T any] interface {
	Tape() Tape[T]
}
