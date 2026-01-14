//go:build generated

package roundtrip

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/nazarifard/fastape"
)

func TestRoundTrip(t *testing.T) {
	runRoundTrip(t, RoundTrip{}, sampleRoundTrip())
	runRoundTrip(t, RoundTrip{}, sampleRoundTripNil())
}

func runRoundTrip(t *testing.T, tapable fastape.Tapable[RoundTrip], want RoundTrip) {
	t.Helper()

	tape := tapable.Tape()
	buf1 := make([]byte, tape.Sizeof(want))
	n1, err := tape.Roll(want, buf1)
	if err != nil {
		t.Fatalf("roll1: %v", err)
	}

	var got RoundTrip
	n2, err := tape.Unroll(buf1, &got)
	if err != nil {
		t.Fatalf("unroll1: %v", err)
	}

	buf2 := make([]byte, tape.Sizeof(got))
	n3, err := tape.Roll(got, buf2)
	if err != nil {
		t.Fatalf("roll2: %v", err)
	}

	var got2 RoundTrip
	n4, err := tape.Unroll(buf2, &got2)
	if err != nil {
		t.Fatalf("unroll2: %v", err)
	}

	if n1 != n2 || n1 != n3 || n3 != n4 {
		t.Fatalf("size mismatch: n1=%d n2=%d n3=%d n4=%d", n1, n2, n3, n4)
	}
	if !bytes.Equal(buf1[:n1], buf2[:n3]) {
		t.Fatalf("encoded bytes mismatch")
	}
	if !reflect.DeepEqual(want, got2) {
		t.Fatalf("roundtrip mismatch:\nwant=%#v\ngot=%#v", want, got2)
	}
}
