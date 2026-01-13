package inner

import "time"

type Key string

type Inner struct {
	ID   int
	When time.Time
}

type AliasInner = Inner
