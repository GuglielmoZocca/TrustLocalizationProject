package parse

import (
	"errors"
	"go/constant"
	"go/token"
	"go/types"
)

func SolveInt(s string) (int, error) {
	// Evaluate
	fs := token.NewFileSet()
	tv, err := types.Eval(fs, nil, token.NoPos, s)
	if err != nil {
		return 0, err
	}

	// Extract. Try both float and int
	vi := constant.ToInt(tv.Value)
	if vi.Kind() == constant.Int {
		i, _ := constant.Int64Val(vi)
		return int(i), nil
	}

	vf := constant.ToFloat(tv.Value)
	if vf.Kind() == constant.Float {
		f, _ := constant.Float64Val(vf)
		return int(f), nil
	}

	return 0, errors.New("Unsolvable: " + s)
}
