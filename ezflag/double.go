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
func NewFloat64Var(name string, defaultVal float64, usage string, required bool) *Float64Var {
	dv := &Float64Var{Var: createVar(name, usage, required), val: defaultVal}
	flag.Float64Var(&dv.val, name, defaultVal, usage)

	return dv
}

func (dv *Float64Var) GetVal() (val float64, err error) {
	if dv.required && dv.val == 0 {
		err = fmt.Errorf("%s cannot be zero value, %s", dv.name, dv.usage)
		return
	}

	val = dv.val
	return
}

func (dv *Float64Var) MustGetVal() (val float64) {
	if dv.required && dv.val == 0 {
		panic(fmt.Sprintf("%s cannot be zero value, %s", dv.name, dv.usage))
	}

	val = dv.val
	return
}
