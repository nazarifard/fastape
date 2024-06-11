package tape

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

func (t compositeTape) Marshal(a Composite, bs []byte) (n int) {
	k := 0
	k, _ = t.NameTape.Marshal(a.Name, bs[n:])
	n += k
	k, _ = t.BirthDayTape.Marshal(a.BirthDay, bs[n:])
	n += k
	k, _ = t.PhoneTape.Marshal(a.Phone, bs[n:])
	n += k
	k, _ = t.SiblingsTape.Marshal(a.Siblings, bs[n:])
	n += k
	k, _ = t.SpouseTape.Marshal(a.Spouse, bs[n:])
	n += k
	k, _ = t.MoneyTape.Marshal(a.Money, bs[n:])
	n += k

	// k, _ = t.PtrTape.Marshal(a.Ptr, bs[n:])
	// n += k
	// k, _ = t.SliceTape.Marshal(a.Slice, bs[n:])
	// n += k
	// k, _ = t.MapTape.Marshal(a.Map, bs[n:])
	// n += k
	return
}

func (t compositeTape) Unmarshal(bs []byte, a *Composite) (n int, err error) {
	var l int
	l, err = t.NameTape.Unmarshal(bs[n:], &a.Name)
	n += l
	if err != nil {
		return 0, err
	}

	l, err = t.BirthDayTape.Unmarshal(bs[n:], &a.BirthDay)
	n += l
	if err != nil {
		return 0, err
	}

	l, err = t.PhoneTape.Unmarshal(bs[n:], &a.Phone)
	n += l
	if err != nil {
		return 0, err
	}

	l, err = t.SiblingsTape.Unmarshal(bs[n:], &a.Siblings)
	n += l
	if err != nil {
		return 0, err
	}

	l, err = t.SpouseTape.Unmarshal(bs[n:], &a.Spouse)
	n += l
	if err != nil {
		return 0, err
	}

	l, err = t.MoneyTape.Unmarshal(bs[n:], &a.Money)
	n += l
	if err != nil {
		return 0, err
	}

	// l, err = t.PtrTape.Unmarshal(bs[n:], &a.Ptr)
	// n += l
	// if err != nil {
	// 	return 0, err
	// }

	// l, err = t.SliceTape.Unmarshal(bs[n:], &a.Slice)
	// n += l
	// if err != nil {
	// 	return 0, err
	// }

	// l, err = t.MapTape.Unmarshal(bs[n:], &a.Map)
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
	_ = compositeTape.Marshal(compo, buff)
	_, _ = compositeTape.Unmarshal(buff, &newCompo)
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
