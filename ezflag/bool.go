package ezflag

type BoolVar struct {
	Var
	val bool
}

// NewBoolVar build new bool var
func NewBoolVar(name string, defaultVal bool, usage string, required bool) *BoolVar {
	bv := &BoolVar{Var: createVar(name, usage, required), val: defaultVal}
	continueCommandLine.BoolVar(&bv.val, name, defaultVal, usage)
	return bv
}

func (bv *BoolVar) GetVal() (val bool , err error) {
	val = bv.val
	return
}

func (bv *BoolVar) MustGetVal() (val bool) {
	val = bv.val
	return
}