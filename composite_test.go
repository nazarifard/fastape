package fastape

import (
	"testing"
	"time"

	"github.com/sanity-io/litter"
)

type Composite struct {
	Name     string
	BirthDay time.Time
	Phone    string
	Siblings int
	Spouse   bool
	Money    float64
	// Ptr      *Block
	// Slice    []Block
	// Map      map[string]Block
}

type compositeTape struct {
	NameTape     StringTape
	BirthDayTape TimeTape
	PhoneTape    StringTape
	SiblingsTape UnitTape[int]
	SpouseTape   UnitTape[bool]
	MoneyTape    UnitTape[float64]

	// PtrTape  PtrT[Block, BlockT]
	// SliceTape SliceT[Block, BlockT]
	// MapTape  MapT[string, Block, StringT, BlockT]
}

func (t compositeTape) Sizeof(p Composite) int {
	return t.NameTape.Sizeof(p.Name) +
		t.BirthDayTape.Sizeof(p.BirthDay) +
		t.PhoneTape.Sizeof(p.Phone) +
		t.SiblingsTape.Sizeof(p.Siblings) +
		t.SpouseTape.Sizeof(p.Spouse) +
		t.MoneyTape.Sizeof(p.Money)
}

func (t compositeTape) Roll(a Composite, bs []byte) (n int) {
	k := 0
	k, _ = t.NameTape.Roll(a.Name, bs[n:])
	n += k
	k, _ = t.BirthDayTape.Roll(a.BirthDay, bs[n:])
	n += k
	k, _ = t.PhoneTape.Roll(a.Phone, bs[n:])
	n += k
	k, _ = t.SiblingsTape.Roll(a.Siblings, bs[n:])
	n += k
	k, _ = t.SpouseTape.Roll(a.Spouse, bs[n:])
	n += k
	k, _ = t.MoneyTape.Roll(a.Money, bs[n:])
	n += k

	// k, _ = t.PtrTape.Roll(a.Ptr, bs[n:])
	// n += k
	// k, _ = t.SliceTape.Roll(a.Slice, bs[n:])
	// n += k
	// k, _ = t.MapTape.Roll(a.Map, bs[n:])
	// n += k
	return
}

func (t compositeTape) Unroll(bs []byte, a *Composite) (n int, err error) {
	var l int
	l, err = t.NameTape.Unroll(bs[n:], &a.Name)
	n += l
	if err != nil {
		return 0, err
	}

	l, err = t.BirthDayTape.Unroll(bs[n:], &a.BirthDay)
	n += l
	if err != nil {
		return 0, err
	}

	l, err = t.PhoneTape.Unroll(bs[n:], &a.Phone)
	n += l
	if err != nil {
		return 0, err
	}

	l, err = t.SiblingsTape.Unroll(bs[n:], &a.Siblings)
	n += l
	if err != nil {
		return 0, err
	}

	l, err = t.SpouseTape.Unroll(bs[n:], &a.Spouse)
	n += l
	if err != nil {
		return 0, err
	}

	l, err = t.MoneyTape.Unroll(bs[n:], &a.Money)
	n += l
	if err != nil {
		return 0, err
	}

	// l, err = t.PtrTape.Unroll(bs[n:], &a.Ptr)
	// n += l
	// if err != nil {
	// 	return 0, err
	// }

	// l, err = t.SliceTape.Unroll(bs[n:], &a.Slice)
	// n += l
	// if err != nil {
	// 	return 0, err
	// }

	// l, err = t.MapTape.Unroll(bs[n:], &a.Map)
	// n += l
	// if err != nil {
	// 	return 0, err
	// }

	return
}

var compo = Composite{
	Name:     "Bahador Nazari Fard",
	BirthDay: time.Now(),
	Phone:    "05123486142",
	Siblings: 7,
	Spouse:   false,
	Money:    987654321.98765432,

	// Map:   map[string]Block{"aaa": block, "bbb": block, "ccc": block, "ddd": block},
	// Slice: []Block{block, block, block},
	// Ptr:   &block,
}

func TestComposite(t *testing.T) {
	var compositeTape compositeTape
	var newCompo Composite
	buff := make([]byte, 1000)
	_ = compositeTape.Roll(compo, buff)
	_, _ = compositeTape.Unroll(buff, &newCompo)
	if newCompo.Name != compo.Name ||
		newCompo.Money != compo.Money ||
		newCompo.Phone != compo.Phone ||
		!newCompo.BirthDay.Equal(compo.BirthDay) ||
		newCompo.Siblings != compo.Siblings ||
		newCompo.Spouse != compo.Spouse {

		litter.Dump(compo)
		litter.Dump(newCompo)
		t.Errorf("mixTape failed")
	}

	// litter.Dump(newCompo)

}
