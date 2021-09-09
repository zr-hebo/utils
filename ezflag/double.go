package ezflag

import (
	"fmt"
)

type Float64Var struct {
	Var
	val float64
}

// NewFloat64Var build new int var
func NewFloat64Var(name string, defaultVal float64, usage string, required bool) *Float64Var {
	dv := &Float64Var{Var: createVar(name, usage, required), val: defaultVal}
	continueCommandLine.Float64Var(&dv.val, name, defaultVal, usage)

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
	val, err := dv.GetVal()
	if err != nil {
		panic(err.Error())
	}

	val = dv.val
	return
}
