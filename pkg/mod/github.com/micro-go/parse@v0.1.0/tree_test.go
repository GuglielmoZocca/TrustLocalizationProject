package parse

import (
	"reflect"
	"testing"
)

// No wildcards

func TestTree1(t *testing.T) {
	a, _ := SetTreeString("k", "v", nil)
	testTreeEquals(t, a, testTree1)
}

func TestTree2(t *testing.T) {
	a, _ := SetTreeString("m/k", "v", nil)
	testTreeEquals(t, a, testTree2)
}

func TestTree3(t *testing.T) {
	a, _ := SetTreeInt("k", 404, nil)
	testTreeEquals(t, a, testTree3)
}

func TestTree4(t *testing.T) {
	a, _ := SetTreeInt("m/k", 404, nil)
	testTreeEquals(t, a, testTree4)
}

func TestTree5(t *testing.T) {
	v, ok := FindTreeString("m/k", testTree2)
	if !ok || v != "v" {
		t.Fail()
	}
}

func TestTree6(t *testing.T) {
	v, ok := FindTreeInt("m/k", testTree4)
	if !ok || v != 404 {
		t.Fail()
	}
}

// Helpers

func testTreeEquals(t *testing.T, a, b interface{}) {
	if !cmpTrees(a, b) {
		t.Fail()
	}
}

func cmpTrees(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil {
		return false
	} else if b == nil {
		return false
	}
	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)
	if ta != tb {
		return false
	}
	switch av := a.(type) {
	case map[string]interface{}:
		return cmpTreeMaps(av, b.(map[string]interface{}))
	case []interface{}:
		return cmpTreeSlices(av, b.([]interface{}))
	default:
		return av == b
	}
	return false
}

func cmpTreeMaps(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for ka, va := range a {
		vb, ok := b[ka]
		if !ok {
			return false
		}
		cmp := cmpTrees(va, vb)
		if !cmp {
			return false
		}
	}
	return true
}

func cmpTreeSlices(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i, va := range a {
		cmp := cmpTrees(va, b[i])
		if !cmp {
			return false
		}
	}
	return true
}

// Data

var (
	// Note that we have to be careful in how we construct the test data,
	// or else it won't type-match.
	testTree1 = map[string]interface{}{"k": "v"}
	testTree2 = map[string]interface{}{"m": testTree1}

	testTree3 = map[string]interface{}{"k": 404}
	testTree4 = map[string]interface{}{"m": testTree3}
)
