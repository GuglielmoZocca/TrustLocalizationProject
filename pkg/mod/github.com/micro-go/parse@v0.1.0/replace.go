package parse

import (
	"fmt"
	"strconv"
	"strings"
)

// ReplacePairs() takes a string and applies a series of replace pairs.
// For each pair in the arguments, the first value is found and replaced with the second.
// The second argument can be anything that can reduce down to a string, including
// a function with no arguments that returns a string.
func ReplacePairs(s string, pairs ...interface{}) string {
	even := true
	var last interface{}
	for _, p := range pairs {
		if even {
			last = p
		} else {
			first, err1 := tostring(last)
			second, err2 := tostring(p)
			if err1 == nil && err2 == nil {
				s = strings.Replace(s, first, second, -1)
			}
		}
		even = !even
	}
	return s
}

func tostring(i interface{}) (string, error) {
	if i == nil {
		return "", nil
	}
	// Optimized
	switch t := i.(type) {
	case string:
		return t, nil
	case int:
		return strconv.Itoa(t), nil
	case func() string:
		return t(), nil
	}
	// Fallback
	return fmt.Sprintf("%v", i), nil
}
