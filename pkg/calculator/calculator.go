package calculator

import (
	"go/token"
	"go/types"
	"strconv"
)

func Calculate(expression string) (float64, error) {
	fs := token.NewFileSet()
	tv, err := types.Eval(fs, nil, token.NoPos, expression)
	if err != nil {
		return .0, err
	}
	res, err := strconv.ParseFloat(tv.Value.String(), 64)
	if err != nil {
		return .0, err
	}
	return res, nil
}
