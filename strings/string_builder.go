package strings

import (
	"fmt"
	"strings"
)

func BuildString(args ...string) string {
	if len(args) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, arg := range args {
		fmt.Fprint(&sb, arg)
	}
	return sb.String()
}

func ConcatValueInSlice(vals []interface{}) string {
	if vals == nil {
		return "NULL"
	}

	var sb strings.Builder
	var sep = ""
	for _, val := range vals {
		fmt.Fprint(&sb, sep)
		fmt.Fprint(&sb, fmt.Sprint(val))
		sep = "#"
	}

	return sb.String()
}
