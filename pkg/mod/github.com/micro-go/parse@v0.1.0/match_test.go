package parse

import (
	"fmt"
	"testing"
)

func TestMatchAccept(t *testing.T) {
	cases := []struct {
		Pattern string
		Match   string
	}{
		{"a", "a"},
		{"a/b", "a/b"},
		{"a/b/c", "a/b/c"},

		{"a/+", "a/a"},
		{"a/+", "a/b"},
		{"a/+/c", "a/a/c"},

		{"a/#", "a/a/d"},
		{"a/#/d", "a/a/d"},
		{"a/#/d", "a/a/b/c/d"},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			m := NewMqttStringMatch(tc.Pattern)
			if !m.Matches(tc.Match) {
				fmt.Println("accept case", i, tc.Pattern, "should match to", tc.Match)
				t.Fatal()
			}
		})
	}
}

func TestMatchDecline(t *testing.T) {
	cases := []struct {
		Pattern string
		Match   string
	}{
		{"a", "b"},
		{"a/a", "b"},
		{"a/a", "a"},
		{"a", "a/b"},

		{"a/+", "a/b/c"},
		{"a/+/c", "a/a"},
		{"a/+/c", "a/a/d"},

		{"a/#/c", "a/a/d"},
		{"a/#/c", "a/a/b/d"},
		{"a/#/c", "a/a/c/d"},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			m := NewMqttStringMatch(tc.Pattern)
			if m.Matches(tc.Match) {
				fmt.Println("decline case", i, tc.Pattern, "should not match to", tc.Match)
				t.Fatal()
			}
		})
	}
}
