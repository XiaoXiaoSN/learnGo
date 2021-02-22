package main

import (
	"fmt"
	"go/constant"
	"go/token"
	"go/types"
)

func main() {
	// fmt.Println(calu("$1 * ($2 + 3) / 7", 10, 30))
	fmt.Println(calu("8 * (a_1 + 3) - 7", 10))
	fmt.Println(calu("10"))
}

func calu(expr string, params ...float64) interface{} {

	for i, v := range params {
		algebra := fmt.Sprintf("a_%d", i+1)

		c := types.NewConst(token.NoPos, nil, algebra, types.Typ[types.Float64], constant.MakeFloat64(v))
		types.Universe.Insert(c)
	}

	fs := token.NewFileSet()
	tv, err := types.Eval(fs, nil, token.NoPos, expr)
	if err != nil {
		panic(err)
	}

	// result, succ := constant.Float64Val(tv.Value)
	// if !succ {
	// 	panic("assert float64 fail")
	// }

	// return result
	return tv.Value
}
