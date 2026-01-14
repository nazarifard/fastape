package fastape

import "testing"

// Regression tests (2026-01-14): LenTape varint encoding and SliceTape reuse safety.

func TestLenTape_RoundTripLargeValues(t *testing.T) {
	var tape LenTape

	// Values > 127 require multi-byte encoding.
	// Expect round-trip correctness.
	cases := []int{128, 129, 255, 1024, 16383, 16384}
	for _, want := range cases {
		t.Run("n="+itoa(want), func(t *testing.T) {
			buf := make([]byte, 32)
			n, err := tape.Roll(want, buf)
			if err != nil {
				t.Fatalf("roll: %v", err)
			}

			var got int
			m, err := tape.Unroll(buf[:n], &got)
			if err != nil {
				t.Fatalf("unroll: %v", err)
			}
			if n != m {
				t.Fatalf("size mismatch: roll=%d unroll=%d", n, m)
			}
			if got != want {
				t.Fatalf("roundtrip mismatch: want=%d got=%d", want, got)
			}
		})
	}
}

func TestSliceTape_UnrollIntoReusedSlice(t *testing.T) {
	var tape SliceTape[int, UnitTape[int]]

	encoded := []int{1}
	buf := make([]byte, tape.Sizeof(encoded))
	n, err := tape.Roll(encoded, buf)
	if err != nil {
		t.Fatalf("roll: %v", err)
	}

	// Reuse a larger slice. Unroll should decode the encoded length without over-reading.
	decoded := make([]int, 2)
	decoded[0], decoded[1] = -1, -2
	_, err = tape.Unroll(buf[:n], &decoded)
	if err != nil {
		t.Fatalf("unroll: %v", err)
	}
	if decoded[0] != 1 {
		t.Fatalf("decoded[0]=%d, want=1", decoded[0])
	}
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}

	var b [32]byte
	i := len(b)
	for v > 0 {
		i--
		b[i] = byte('0' + v%10)
		v /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
