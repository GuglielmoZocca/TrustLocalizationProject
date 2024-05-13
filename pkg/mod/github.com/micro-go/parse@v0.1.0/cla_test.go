package parse

import (
	"fmt"
	"testing"
)

func TestAsArguments(t *testing.T) {
	cases := []struct {
		From string
		Want []string
	}{
		{"parse it all", []string{"parse", "it", "all"} },
		{`parse "it" all`, []string{"parse", "it", "all"} },
		{`parse "yo momma" all`, []string{"parse", "yo momma", "all"} },
		{`parse "yo \"momma\"" all`, []string{"parse", `yo "momma"`, "all"} },
		{`parse "path\to\file"`, []string{"parse", `path\to\file`} },
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			have := AsArguments(tc.From)
			if arraysMatch(have, tc.Want) == false {
				fmt.Println("No match have", have, "b", tc.Want)
				t.Fail()
			}
		})
	}
}
