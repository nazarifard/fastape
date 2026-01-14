package outer

import "github.com/nazarifard/fastape/testing/fixtures/nested/inner"

// fastape:generate

type Outer struct {
	Primary inner.Inner
	Alias   inner.AliasInner
	Keys    []inner.Key
	Lookup  map[inner.Key]*inner.Inner
	Items   [2]inner.Inner
}
