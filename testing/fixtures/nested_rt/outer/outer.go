package outer

import (
	"time"

	"github.com/nazarifard/fastape/testing/fixtures/nested_rt/inner"
)

type Outer struct {
	Primary inner.Payload
	Alias   inner.Payload
	Keys    []inner.Key
	Lookup  map[inner.Key]*inner.Payload
	Ref     *inner.Payload
}

func sampleOuter() Outer {
	ref := &inner.Payload{ID: 99, When: sampleTime(), Note: "ref"}
	return Outer{
		Primary: inner.Payload{ID: 1, When: sampleTime(), Note: "primary"},
		Alias:   inner.Payload{ID: 2, When: sampleTime(), Note: "alias"},
		Keys:    []inner.Key{"k1", "k2"},
		Lookup:  map[inner.Key]*inner.Payload{"k1": ref},
		Ref:     ref,
	}
}

func sampleOuterNil() Outer {
	return Outer{
		Primary: inner.Payload{ID: 3, When: sampleTime(), Note: "zero"},
		Alias:   inner.Payload{ID: 0, When: sampleTime(), Note: ""},
	}
}

func sampleTime() time.Time {
	return time.Unix(1700001000, 0)
}
