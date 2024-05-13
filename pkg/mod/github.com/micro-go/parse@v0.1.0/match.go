package parse

import (
	"strings"
)

// ------------------------
// API
// ------------------------

// interface StringMatch describes an object that can match strings.
type StringMatch interface {
	Matches(cmp string) bool
}

// ------------------------
// MQTT Rules
// ------------------------

func NewMqttStringMatch(_pattern string) StringMatch {
	pattern := strings.Split(_pattern, mqttSeparator)
	return &mqtt_match{pattern}
}

const (
	mqttSeparator = "/"
)

type mqtt_match struct {
	pattern []string
}

func (m *mqtt_match) Matches(cmp string) bool {
	pn := NewStringToken(m.pattern...)
	cn := NewStringToken(strings.Split(cmp, mqttSeparator)...)

	for curr, err := pn.Next(); err == nil; curr, err = pn.Next() {
		switch curr {
		case "#":
			// Multi level wildcard. If I'm at the end,
			// it's a match, otherwise fast forward until
			// I find the next match.
			if pn.Empty() || cn.Empty() {
				return true
			}
			next, _ := pn.Next()
			if !Advance(next, cn) {
				return false
			}
		case "+":
			// Single level wildcard. Advance one.
			cn.Next()
		default:
			other, _ := cn.Next()
			if curr != other {
				return false
			}
		}
	}
	return cn.Empty()
}
