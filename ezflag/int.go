package ezflag

import (
	"flag"
	"fmt"
)

type IntVar struct {
	Var
	val int
}

// NewIntVar build new int var
func NewIntVar(name, defaultVal, usage string, required bool) *StringVar {
	sv := &StringVar{Var: createVar(name, usage, required), val: defaultVal}
	flag.StringVar(&sv.val, name, "", usage)

	return sv
}

func (sv *IntVar) GetVal() (val int, err error) {
	if sv.required && sv.val == 0 {
		err = fmt.Errorf(sv.usage)
		return
	}

	val = sv.val
	return
}