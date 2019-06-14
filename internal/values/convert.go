package values

import (
	"fmt"

	"github.com/ptiger10/pd/options"
)

// Convert a collection of values from one type to another, and coerce to null if a value cannot be converted sensibly
func Convert(currentVals Values, dataType options.DataType) (Values, error) {
	var vals Values
	switch dataType {
	case options.None:
		return nil, fmt.Errorf("unable to convert values: must supply a valid Kind")
	case options.Float64:
		vals = currentVals.ToFloat64()
	case options.Int64:
		vals = currentVals.ToInt64()
	case options.String:
		vals = currentVals.ToString()
	case options.Bool:
		vals = currentVals.ToBool()
	case options.DateTime:
		vals = currentVals.ToDateTime()
	case options.Interface:
		vals = currentVals.ToInterface()
	default:
		return nil, fmt.Errorf("unable to convert values: kind not supported: %v", dataType)
	}
	return vals, nil
}
