package inner

import "time"

type Key string

type Payload = struct {
	ID   int
	When time.Time
	Note string
}
