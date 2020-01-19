package ezflag

import (
	"flag"
	"fmt"
)

type Float64Var struct {
	Var
	val float64
}

// NewFloat64Var build new int var
func NewFloat64Var(defaultVal float64, name, usage string, required bool) *Float64Var {
	dv := &Float64Var{Var: createVar(name, usage, required), val: defaultVal}
	flag.Float64Var(&dv.val, name, defaultVal, usage)

	return dv
}

func (dv *Float64Var) GetVal() (val float64, err error) {
	if dv.required && dv.val == 0 {
		err = fmt.Errorf(dv.usage)
		return
	}

	val = dv.val
	return
}