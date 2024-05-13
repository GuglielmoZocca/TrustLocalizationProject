# parse
Simple parsing functions.

## Description

This currently supplies two simple objects to help with parsing:

StringToken is used to create a string iterator. It's a small convience over indexing a slice, but results in cleaner code.

Example:

```
	args := []string{"this", "is", "a", "sequence"}
	token := parse.NewStringToken(args...)
	for a, err := token.Next(); err == nil; a, err = token.Next() {
		switch a {
		}
	}
```
	
Matching is used to match two strings. It currently supports MQTT rules.

Example:

```
	pattern := "a/#"
	cmp := "a/a/d"
	m := NewMqttStringMatch(pattern)
	if m.Matches(cmp) {
		// success
	}
```