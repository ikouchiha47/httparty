package errorrs

import "errors"

var (
	InvalidFile        = errors.New("InvalidFile")
	InvalidFormat      = errors.New("InvalidFormat")
	UnknownContentType = errors.New("UnknownContentType")
	IndexOutOfRange    = errors.New("IndexOutOfRange")
	DuplicateEntry     = errors.New("DuplicateEntry")
	CyclicDependency   = errors.New("CyclicDependency")
)
