package parse

// Operations on an arbitrary interface{}, typically read from JSON or XML.
// The basic pattern is to refer to keys by using "/" to denote parent/child.

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
)

type treeSetFunc func(key string, parent map[string]interface{}) error

const (
	TreeSeparator = "/"
)

func ReadJsonFile(filename string) (interface{}, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var tree interface{}
	err = json.Unmarshal(data, &tree)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

func TreeBool(path string, root interface{}, defaultValue bool) bool {
	_v, err := TreeValue(path, root)
	if err != nil {
		return defaultValue
	}
	switch v := _v.(type) {
	case bool:
		return v
	case int:
		return v > 0
	case string:
		return len(v) > 0 && v[0] == 'T' || v[0] == 't'
	}
	return defaultValue
}

func TreeInt(path string, root interface{}, defaultValue int) int {
	v, ok := FindTreeInt(path, root)
	if !ok {
		return defaultValue
	}
	return v
}

func TreeString(path string, root interface{}, defaultValue string) string {
	v, ok := FindTreeString(path, root)
	if !ok {
		return defaultValue
	}
	return v
}

func TreeStringSlice(path string, root interface{}) []string {
	_v, err := TreeValue(path, root)
	if err != nil {
		return nil
	}
	v, ok := _v.([]interface{})
	if !ok {
		return nil
	}
	var ans []string
	for _, _value := range v {
		if value, ok := _value.(string); ok {
			ans = append(ans, value)
		}
	}
	return ans
}

func TreeValue(path string, tree interface{}) (interface{}, error) {
	v, ok := FindTreeValue(path, tree)
	if !ok {
		return nil, noValueErr
	}
	return v, nil
}

// FindTreeInt() finds the value at the given path and answers it
// as an int, if it's convertible.
func FindTreeInt(path string, root interface{}) (int, bool) {
	_v, err := TreeValue(path, root)
	if err != nil {
		return 0, false
	}
	switch v := _v.(type) {
	case bool:
		if v {
			return 1, true
		}
		return 0, true
	case float64:
		return int(v), true
	case int:
		return v, true
	case string:
		i, err := strconv.Atoi(v)
		if err == nil {
			return i, true
		}
	}
	return 0, false
}

// FindTreeString() func follows a path down an arbitrary object to answer
// a string at the root. The path is "/" separated to enter maps.
// Note: This is identical to TreeString(), but I'm acknowledging, for
// all the convenience of that API in certain situations, I should have
// stuck to a more go-way of doing things.
func FindTreeString(path string, root interface{}) (string, bool) {
	_v, err := TreeValue(path, root)
	if err != nil {
		return "", false
	}
	v, ok := _v.(string)
	return v, ok
}

func FindTreeValue(_path string, _tree interface{}) (interface{}, bool) {
	tree, ok := _tree.(map[string]interface{})
	if !ok {
		return nil, false
	}
	path := strings.Split(_path, TreeSeparator)
	v, err := findTreeValue(path, tree)
	if err == nil {
		return v, true
	}
	return nil, false
}

func findTreeValue(path []string, root map[string]interface{}) (interface{}, error) {
	if len(path) < 1 {
		return "", badRequestErr
	}
	v := root[path[0]]
	if v == nil {
		return nil, noValueErr
	}
	if len(path) == 1 {
		return v, nil
	}
	if root, ok := v.(map[string]interface{}); ok {
		return findTreeValue(path[1:], root)
	}
	return nil, noMapErr
}

func SetTreeInt(path string, value int, tree interface{}) (interface{}, error) {
	fn := func(key string, parent map[string]interface{}) error {
		parent[key] = value
		return nil
	}
	return setTreeValue(path, tree, fn)
}

func SetTreeString(path, value string, tree interface{}) (interface{}, error) {
	fn := func(key string, parent map[string]interface{}) error {
		parent[key] = value
		return nil
	}
	return setTreeValue(path, tree, fn)
}

func setTreeValue(path string, tree interface{}, fn treeSetFunc) (interface{}, error) {
	if path == "" {
		return tree, badRequestErr
	}
	names := strings.Split(path, TreeSeparator)
	if tree == nil {
		tree = make(map[string]interface{})
	}
	return tree, setTreeValueOn(names, tree, fn)
}

func setTreeValueOn(path []string, tree interface{}, fn treeSetFunc) error {
	if len(path) < 1 || path[0] == "" {
		return badRequestErr
	}
	// This is an error condition that should never happen
	if tree == nil {
		return treeErr
	}
	switch t := tree.(type) {
	case map[string]interface{}:
		if len(path) > 1 {
			c, ok := t[path[0]]
			if c == nil || !ok {
				c = make(map[string]interface{})
				t[path[0]] = c
			}
			return setTreeValueOn(path[1:], c, fn)
		} else {
			return fn(path[0], t)
		}
	default:
		// Currently only support maps as containers
		return badRequestErr
	}
}
