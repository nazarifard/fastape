package array

import "time"

type Inner struct {
	ID   int
	When time.Time
	Note string
}

// fastape:generate

type ArraySample = [2]Inner

func sampleArray() ArraySample {
	return ArraySample{
		{ID: 1, When: time.Unix(1700002000, 0), Note: "one"},
		{ID: 2, When: time.Unix(1700003000, 0), Note: "two"},
	}
}
