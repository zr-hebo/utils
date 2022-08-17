package ezflag

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	allVariables = make([]Vari, 0, 8)
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
func PrintAllUsage() {
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

func getKeyFromFlag(flag string) string {
	fields := strings.SplitAfterN(flag, "=", 2)
	return strings.TrimLeft(fields[0], "-")
}

func Parse(reservedKeys []string) error {
	continueCommandLine.Usage = NoOperationUsage

	reservedArgs := make([]string, 0, len(reservedKeys))
	userArgs := make([]string, 0, len(os.Args[1:]))
	reservedKeyInMap := make(map[string]bool)
	for _, key := range reservedKeys {
		reservedKeyInMap[key] = true
	}
	for _, arg := range os.Args[1:] {
		argKey := getKeyFromFlag(arg)
		if reservedKeyInMap[argKey] {
			reservedArgs = append(reservedArgs, arg)
		} else {
			userArgs = append(userArgs, arg)
		}
	}
	os.Args = append([]string{os.Args[0]}, reservedArgs...)
	flag.Parse()

	errMsgs := make([]string, 0, 4)
	for _, userArg := range userArgs {
		oneErr := continueCommandLine.Parse([]string{userArg})
		if oneErr != nil {
			errMsgs = append(errMsgs, oneErr.Error())
		}
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "; "))
	}
	return nil
}

func NoOperationUsage() {
}
