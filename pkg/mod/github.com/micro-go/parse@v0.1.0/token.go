package parse

// interface StringToken allows iterating over a series of strings.
type StringToken interface {
	Empty() bool
	Next() (string, error)
	// Answer the remaining items. This leaves the token empty.
	Remainder() ([]string, error)
}

// Answer a new iterating string token from a list of strings.
func NewStringToken(tokens ...string) StringToken {
	return &stringToken{tokens, 0}
}

type stringToken struct {
	tokens []string
	index  int
}

func (n *stringToken) Empty() bool {
	return n.index >= len(n.tokens)
}

func (n *stringToken) Next() (string, error) {
	if n.index >= len(n.tokens) {
		return "", emptyErr
	}
	idx := n.index
	n.index++
	return n.tokens[idx], nil
}

func (n *stringToken) Remainder() ([]string, error) {
	if n.index >= len(n.tokens) {
		return nil, emptyErr
	}
	idx := n.index
	n.index = len(n.tokens)
	return n.tokens[idx:], nil
}

// ----------------------------------
// StringToken functions
// ----------------------------------

// func Advance() advances the token until it finds the string, returning true if it's found.
// The token will be left on the next string.
func Advance(s string, t StringToken) bool {
	for curr, err := t.Next(); err == nil; curr, err = t.Next() {
		if curr == s {
			return true
		}
	}
	return false
}
