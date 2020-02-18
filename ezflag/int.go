package ezflag

import (
	"flag"
	"fmt"
)

type IntVar struct {
	Var
	val int64
}

// NewIntVar build new int var
func NewIntVar(name string, defaultVal int, usage string, required bool) *IntVar {
	iv := &IntVar{Var: createVar(name, usage, required), val: int64(defaultVal)}
	flag.Int64Var(&iv.val, name, int64(defaultVal), usage)

	return iv
}

func (iv *IntVar) GetVal() (val int, err error) {
	if iv.required && iv.val == 0 {
		err = fmt.Errorf("%s cannot be zero value, %s", iv.name, iv.usage)
		return
	}

	val = int(iv.val)
	return
}