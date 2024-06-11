package fastape

import "errors"

var ErrNoSpaceLeft = errors.New("no space left")
var ErrInvalidData = errors.New("invalid data")
var ErrNilPtr = errors.New("target pointer is nil")
