package parse

import (
	"errors"
)

var (
	badRequestErr         = errors.New("Bad request")
	emptyErr              = errors.New("Token empty")
	invalidExtractRuleErr = errors.New("Invalid ExtractRule")
	invalidRuleErr        = errors.New("Invalid rule")
	noValueErr            = errors.New("No value")
	noMapErr              = errors.New("No map")
	ruleMismatchErr       = errors.New("Rule mismatch")
	stringMismatchErr     = errors.New("String mismatch")
	treeErr               = errors.New("Unknown tree error")
)
