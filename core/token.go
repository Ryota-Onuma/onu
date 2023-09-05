package core

type Token struct {
	Type    int
	Literal string
}

const (
	EOF     = 0
	UNKNOWN = -1
)
