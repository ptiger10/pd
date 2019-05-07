package values

import (
	"fmt"

	"github.com/ptiger10/pd/kinds"
)

// Convert a collection of values from one type to another, and coerce to null if a value cannot be converted sensibly
func Convert(currentVals Values, kind kinds.Kind) (Values, error) {
	var vals Values
	switch kind {
	case kinds.None:
		return nil, fmt.Errorf("Unable to convert values: must supply a valid Kind")
	case kinds.Float:
		vals = currentVals.ToFloat()
	case kinds.Int:
		vals = currentVals.ToInt()
	case kinds.String:
		vals = currentVals.ToString()
	case kinds.Bool:
		vals = currentVals.ToBool()
	case kinds.DateTime:
		vals = currentVals.ToDateTime()
	case kinds.Interface:
		vals = currentVals.ToInterface()
	default:
		return nil, fmt.Errorf("Unable to convert values: kind not supported: %v", kind)
	}
	return vals, nil
}
