package parse

import (
	"unicode"
)

// AsArguments() parses a single string into a slice based on typical
// command line argument rules (spaces separate tokens, except when quoted,
// and in that case quotes are removed).
// Note: seems very likely there's already something to handle this in go,
// but I'm not seeing it.
func AsArguments(s string) []string {
	b := &cla_builder{}
	for _, c := range s {
		b.push(c)
	}
	b.flush()
	return b.ans
}

type cla_builder struct {
	ans       []string
	cur       string
	in_quote  bool
	last_char string // Will be either escape or double quote
}

func (b *cla_builder) push(r rune) {
	c := string(r)
	if b.in_quote {
		if b.last_char == `\` {
			if c == `"` {
				b.cur += c
			} else {
				b.cur += b.last_char
				b.cur += c
			}
		} else if c == `"` {
			b.flush()
			b.in_quote = false
		} else if c == `\` {
			// wait to handle later
		} else {
			b.cur += c
		}
	} else if c == `"` {
		b.flush()
		b.in_quote = true
	} else if unicode.IsSpace(r) {
		b.flush()
	} else {
		b.cur += c
	}
	b.last_char = c
}

func (b *cla_builder) flush() {
	if len(b.cur) > 0 {
		b.cur = claEscape(b.cur)
		b.ans = append(b.ans, b.cur)
		b.cur = ""
	}
	b.in_quote = false
}

// claEscape returns the string with all characters escaped
func claEscape(s string) string {
	return s
}
