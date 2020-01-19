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
	iv := &StringVar{Var: createVar(name, usage, required), val: defaultVal}
	flag.StringVar(&iv.val, name, "", usage)

	return iv
}

func (iv *IntVar) GetVal() (val int, err error) {
	if iv.required && iv.val == 0 {
		err = fmt.Errorf("%s cannot be zero value, %s", iv.name, iv.usage)
		return
	}

	val = iv.val
	return
}