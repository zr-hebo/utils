package ezflag

import (
	"flag"
	"fmt"
)

var (
	allVariables  = make([]Vari, 0, 8)
)

type Vari interface {
	Name() string
	Usage() string
	Required() bool
}

type Var struct {
	name     string
	usage    string
	required bool
}

func createVar(name, usage string, required bool) Var {
	v := Var{name: name, usage: usage, required: required}
	allVariables = append(allVariables, &v)
	return v
}

func (v *Var) Name() string {
	return v.name
}

func (v *Var) Usage() string {
	return v.usage
}

func (v *Var) Required() bool {
	return v.required
}

// PrintAllUsage print all flag usage
func PrintAllUsage()  {
	requiredVars := make([]Vari, 0, len(allVariables))
	optionVars := make([]Vari, 0, len(allVariables))
	for _, v := range allVariables {
		if v.Required() {
			requiredVars = append(requiredVars, v)
		} else {
			optionVars = append(optionVars, v)
		}
	}

	if len(requiredVars) > 0 {
		fmt.Println("Required:")
		for _, v := range requiredVars {
			fmt.Printf("--%s\t%s\n", v.Name(), v.Usage())
		}
	}

	if len(optionVars) > 0 {
		fmt.Println("Options:")
		for _, v := range optionVars {
			fmt.Printf("--%s\t%s\n", v.Name(), v.Usage())
		}
	}
}

func Parse()  {
	flag.Parse()
}
