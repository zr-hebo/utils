package strings

import (
	"bytes"
	"fmt"
)

func BuildString(args ...string) string {
	if len(args) == 0 {
		return ""
	}

	strLen := 4
	for _, arg := range args {
		strLen += len(arg)
	}

	sb := bytes.NewBuffer(make([]byte, 0, strLen))
	for _, arg := range args {
		sb.WriteString(arg)
	}
	return HackString(sb.Bytes())
}

func ConcatValueInSlice(vals []interface{}) string {
	if vals == nil {
		return "NULL"
	}

	sb := bytes.NewBuffer(make([]byte, 0, 32))
	var sep = ""
	for _, val := range vals {
		sb.WriteString(sep)
		sb.WriteString(fmt.Sprint(val))
		sep = "#"
	}

	return HackString(sb.Bytes())
}
