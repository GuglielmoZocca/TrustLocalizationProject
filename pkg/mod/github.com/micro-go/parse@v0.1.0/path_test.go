package parse

import (
	"fmt"
	"testing"
)

func TestPathOne(t *testing.T) {
	cases := []struct {
		Path  string
		Rules []interface{}
		Want  string
	}{
		{"a", []interface{}{KeepRule}, "a"},
		{"a", []interface{}{SkipRule}, ""},
		{"a/b", []interface{}{KeepRule}, "a"},
		{"a/b", []interface{}{KeepRuleT{}}, "a"},
		{"a/b", []interface{}{KeepRule, KeepRule}, "a/b"},
		{"a/b", []interface{}{SkipRule, KeepRule}, "b"},
		{"a/b/c", []interface{}{KeepRule, SkipRule}, "a"},
		{"a", []interface{}{KeepAllRule}, "a"},
		{"a/b", []interface{}{KeepAllRule}, "a/b"},
		{"a/b/c", []interface{}{KeepAllRule}, "a/b/c"},
		{"a/b/c", []interface{}{KeepRule, ReplaceRuleT{"E"}}, "a/E"},
		{"a/b/c", []interface{}{KeepRule, ReplaceRuleT{"E"}, KeepRule}, "a/E/c"},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			opts := PathOpts{Separator: "/"}
			have := PathOne(tc.Path, opts, tc.Rules...)
			if have != tc.Want {
				fmt.Println("Mismtach have", have, "want", tc.Want)
				t.Fail()
			}
		})
	}
}

func TestPathTwo(t *testing.T) {
	cases := []struct {
		Path  string
		Rules []interface{}
		Want  []string
	}{
		{"a", []interface{}{KeepRule}, []string{"a"}},
		{"a/b", []interface{}{KeepRule}, []string{"a"}},
		{"a/b", []interface{}{KeepAllRule}, []string{"a", "b"}},
		{"a/b/c", []interface{}{KeepAllRule}, []string{"a", "b"}},
		{"a/b/c", []interface{}{KeepRule, SkipRule}, []string{"a"}},
		{"a/b/c", []interface{}{KeepRule, SkipRule, KeepRule}, []string{"a", "c"}},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			opts := PathOpts{Separator: "/"}
			havea, haveb := PathTwo(tc.Path, opts, tc.Rules...)
			var have []string
			if havea != "" {
				have = append(have, havea)
			}
			if haveb != "" {
				have = append(have, haveb)
			}
			if arraysMatch(tc.Want, have) == false {
				fmt.Println("Mismtach have", have, "want", tc.Want)
				t.Fail()
			}
		})
	}
}

func TestPathExtractSepAccept(t *testing.T) {
	cases := []struct {
		Path  string
		Rules []interface{}
		Want  []string
	}{
		{"a/b", []interface{}{KeepRule, KeepRule}, []string{"a", "b"}},
		{"a/b/c", []interface{}{KeepRule, KeepRule}, []string{"a", "b"}},
		{"a/b/c", []interface{}{KeepRule, KeepRule, SkipRule}, []string{"a", "b"}},
		{"a/b/c", []interface{}{KeepRule, SkipRule, KeepRule}, []string{"a", "c"}},
		{"bub/1/a/b/c", []interface{}{KeepRule, "1", "a", "b", KeepRule}, []string{"bub", "c"}},
		{"/bub/1/a/b/c", []interface{}{KeepRule, "1", "a", "b", KeepRule}, []string{"bub", "c"}},
		{"bub/1/a/b/c", []interface{}{SkipRule, "1", "a", "b", KeepRule}, []string{"c"}},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			have, err := ExtractSep(tc.Path, "/", tc.Rules...)
			if err != nil {
				fmt.Println(err)
				t.Fail()
			} else if arraysMatch(tc.Want, have) == false {
				t.Fail()
			}
		})
	}
}

func TestPathExtractSepDecline(t *testing.T) {
	cases := []struct {
		Path  string
		Rules []interface{}
		Want  []string
	}{
		{"bub/1/a/b/c", []interface{}{"plus", "1", "a", "b", KeepRule}, []string{"c"}},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			have, err := ExtractSep(tc.Path, "/", tc.Rules...)
			if err == nil && arraysMatch(tc.Want, have) == true {
				t.Fail()
			}
		})
	}
}

// Helpers

func arraysMatch(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
