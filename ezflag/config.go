package ezflag

import (
	"flag"
	"os"
)

var (
	continueCommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
)
