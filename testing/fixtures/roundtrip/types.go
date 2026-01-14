package roundtrip

import "time"

type MyString string
type MyTime time.Time

type AliasSlice = []MyString
type NamedSlice []MyString

type AliasMap = map[MyString]int
type NamedMap map[MyString]int

type Inner struct {
	ID   int
	Note string
}

type NamedPtr *Inner

type Embedded struct {
	Enabled bool
}

// fastape:generate

type RoundTrip struct {
	Flag       bool
	Count      int
	Small      int8
	Medium     int16
	Large      int32
	Huge       int64
	UCount     uint
	U8         uint8
	U16        uint16
	U32        uint32
	U64        uint64
	F32        float32
	F64        float64
	Name       string
	Bytes      []byte
	AliasSlice AliasSlice
	NamedSlice NamedSlice
	AliasMap   AliasMap
	NamedMap   NamedMap
	Ptr        *Inner
	NamedPtr   NamedPtr
	When       time.Time
	WhenAlias  MyTime
	Inline     struct {
		Label string `json:"label"`
	}
	Embedded
}

func sampleRoundTrip() RoundTrip {
	ts := time.Unix(1700000000, 0)
	ptr := &Inner{ID: 42, Note: "note"}
	return RoundTrip{
		Flag:       true,
		Count:      7,
		Small:      -3,
		Medium:     1024,
		Large:      -200000,
		Huge:       9000000000,
		UCount:     13,
		U8:         8,
		U16:        16,
		U32:        32,
		U64:        64,
		F32:        3.25,
		F64:        -99.5,
		Name:       "roundtrip",
		Bytes:      []byte{0, 1, 2, 3},
		AliasSlice: AliasSlice{MyString("alpha"), MyString("beta")},
		NamedSlice: NamedSlice{MyString("gamma")},
		AliasMap:   AliasMap{MyString("k1"): 1},
		NamedMap:   NamedMap{MyString("k2"): 2},
		Ptr:        ptr,
		NamedPtr:   NamedPtr(ptr),
		When:       ts,
		WhenAlias:  MyTime(ts),
		Inline: struct {
			Label string `json:"label"`
		}{
			Label: "inline",
		},
		Embedded: Embedded{Enabled: true},
	}
}

func sampleRoundTripNil() RoundTrip {
	ts := time.Unix(1700000000, 0)
	return RoundTrip{
		Name:      "nil-case",
		When:      ts,
		WhenAlias: MyTime(ts),
	}
}
