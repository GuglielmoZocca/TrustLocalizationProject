package parse

import (
	"strings"
)

// PathOne() answers a single item from the path input, reconstructed from
// the results. There's no error reporting, if you need more complete information
// use a more complex form.
func PathOne(path string, opts PathOpts, rules ...interface{}) string {
	ext, err := Path(path, opts, rules...)
	if err != nil || len(ext) < 1 {
		return ""
	}
	if len(ext) == 1 {
		return ext[0]
	}
	ans := ext[0]
	for i := 1; i < len(ext); i++ {
		ans += opts.Separator + ext[i]
	}
	return ans
}

// PathTwo() answers two items from the path input. There's no error reporting,
// if you need more complete information use a more complex form.
func PathTwo(path string, opts PathOpts, rules ...interface{}) (string, string) {
	ext, _ := Path(path, opts, rules...)
	a := ""
	b := ""
	if len(ext) > 0 {
		a = ext[0]
	}
	if len(ext) > 1 {
		b = ext[1]
	}
	return a, b
}

// func Path() takes a list of strings and returns portions of it, according to the rules.
// Rules can contain either strings or a path rule.
// Strings:
//		The node must match the string exactly. It is not included in the response.
// Path rules:
//		KeepRule. Node can be anything. It is returned in the response.
//		KeepAllRule. All remaining nodes are added to the response and returned.
//		ReplaceRuleT. Replace this item with the Replace.With.
//		SkipRule. Node can be anything. It is not included in the response.
func Path(_path string, opts PathOpts, rules ...interface{}) ([]string, error) {
	if opts.Separator == "" {
		return nil, badRequestErr
	}
	src := strings.Split(_path, opts.Separator)
	var resp []string
	srci := 0
	for _, k := range rules {
		// Skip empties in the source
		for srci < len(src) && src[srci] == "" {
			srci++
		}
		if srci >= len(src) {
			return nil, ruleMismatchErr
		}

		switch rule := k.(type) {
		case string:
			if src[srci] != rule {
				return nil, stringMismatchErr
			}
		case KeepRuleT:
			resp = append(resp, src[srci])
		case KeepAllRuleT:
			// This can only be the final element.
			resp = append(resp, src[srci:len(src)]...)
			return resp, nil
		case ReplaceRuleT:
			resp = append(resp, rule.With)
		case SkipRuleT:
		default:
			return nil, invalidRuleErr
		}
		srci++
	}
	return resp, nil
}

type PathOpts struct {
	Separator string
}

// ---------- DEPRECATED API

// func ExtractSep() is a wrapper on Extract(), first splitting the string based on the separator.
func ExtractSep(path, separator string, rules ...interface{}) ([]string, error) {
	if separator == "" {
		return nil, badRequestErr
	}
	all := strings.Split(path, separator)
	return Extract(all, rules...)
}

// func Extract() takes a list of strings and returns portions of it, according to the rules.
// Rules can contain either strings or an ExtractRule.
// Strings:
//		The node must match the string exactly. It is not included in the response.
// ExtractRules:
//		KeepRule. Node can be anything. It is returned in the response.
//		SkipRule. Node can be anything. It is not included in the response.
func Extract(src []string, rules ...interface{}) ([]string, error) {
	var resp []string
	srci := 0
	for _, k := range rules {
		// Skip empties in the source
		for srci < len(src) && src[srci] == "" {
			srci++
		}
		if srci >= len(src) {
			return nil, ruleMismatchErr
		}

		switch rule := k.(type) {
		case string:
			if src[srci] != rule {
				return nil, stringMismatchErr
			}
		case KeepRuleT:
			resp = append(resp, src[srci])
		case KeepAllRuleT:
			// This can only be the final element.
			resp = append(resp, src[srci:len(src)]...)
			return resp, nil
		case ReplaceRuleT:
			resp = append(resp, rule.With)
		case SkipRuleT:
		default:
			return nil, invalidRuleErr
		}
		srci++
	}
	return resp, nil
}

// ------------------------------------------------------------
// RULES

type KeepRuleT struct{}
type KeepAllRuleT struct{}
type ReplaceRuleT struct {
	With string
}
type SkipRuleT struct{}

var (
	KeepRule    = KeepRuleT{}
	KeepAllRule = KeepAllRuleT{}
	SkipRule    = SkipRuleT{}
)
