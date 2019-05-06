package constructors

import (
	"strings"

	"github.com/ptiger10/pd/kinds"
	"github.com/ptiger10/pd/options"

	"github.com/ptiger10/pd/internal/values"
)

// [START Convenience Functions]

func isNullString(s string) bool {
	nullStrings := options.StringNullValues
	for _, ns := range nullStrings {
		if strings.TrimSpace(s) == ns {
			return true
		}
	}
	return false
}

// [END Convenience Functions]

// [START Constructor Functions]

// SliceString converts []string -> ValuesFactory with values.StringValues
func SliceString(vals []string) ValuesFactory {
	var v values.StringValues
	for _, val := range vals {
		if isNullString(val) {
			v = append(v, values.String(options.StringNullFiller, true))
		} else {
			v = append(v, values.String(val, false))
		}

	}
	return ValuesFactory{v, kinds.String}
}

// [END Constructor Functions]
