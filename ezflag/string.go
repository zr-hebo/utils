package ezflag

import (
	"flag"
	"fmt"
)

type StringVar struct {
	Var
	val string
}

// NewStringVar build new string var
func NewStringVar(name, defaultVal, usage string, required bool) *StringVar {
	sv := &StringVar{Var: createVar(name, usage, required), val: defaultVal}
	flag.StringVar(&sv.val, name, "", usage)

	return sv
}

func (sv *StringVar) GetVal() (val string , err error) {
	if sv.required && sv.val == "" {
		err = fmt.Errorf(sv.usage)
		return
	}

	val = sv.val
	return
}