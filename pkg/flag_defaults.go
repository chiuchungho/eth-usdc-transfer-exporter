package pkg

import (
	"fmt"

	"github.com/namsral/flag"
)

func GetFlagDefaults(f *flag.FlagSet) string {
	var defaults string
	format := "  -%s=%s: %s\n"
	f.VisitAll(func(flag *flag.Flag) {
		defaults += fmt.Sprintf(format, flag.Name, flag.DefValue, flag.Usage)
	})
	return defaults
}
