package index

import (
	"fmt"

	"github.com/ptiger10/pd/datatypes"
)

// Convert an index level from one kind to another, then refresh the LabelMap
func (lvl Level) Convert(kind datatypes.DataType) (Level, error) {
	var convertedLvl Level
	switch kind {
	case datatypes.None:
		return Level{}, fmt.Errorf("unable to convert index level: must supply a valid Kind")
	case datatypes.Float64:
		convertedLvl = lvl.ToFloat()
	case datatypes.Int64:
		convertedLvl = lvl.ToInt()
	case datatypes.String:
		convertedLvl = lvl.ToString()
	case datatypes.Bool:
		convertedLvl = lvl.ToBool()
	case datatypes.DateTime:
		convertedLvl = lvl.ToDateTime()
	case datatypes.Interface:
		convertedLvl = lvl.ToInterface()
	default:
		return Level{}, fmt.Errorf("unable to convert level: kind not supported: %v", kind)
	}
	return convertedLvl, nil
}

// ToFloat converts an index level to Float
func (lvl Level) ToFloat() Level {
	lvl.Labels = lvl.Labels.ToFloat()
	lvl.DataType = datatypes.Float64
	lvl.Refresh()
	return lvl
}

// ToInt converts an index level to Int
func (lvl Level) ToInt() Level {
	lvl.Labels = lvl.Labels.ToInt()
	lvl.DataType = datatypes.Int64
	lvl.Refresh()
	return lvl
}

// ToString converts an index level to String
func (lvl Level) ToString() Level {
	lvl.Labels = lvl.Labels.ToString()
	lvl.DataType = datatypes.String
	lvl.Refresh()
	return lvl
}

// ToBool converts an index level to Bool
func (lvl Level) ToBool() Level {
	lvl.Labels = lvl.Labels.ToBool()
	lvl.DataType = datatypes.Bool
	lvl.Refresh()
	return lvl
}

// ToDateTime converts an index level to DateTime
func (lvl Level) ToDateTime() Level {
	lvl.Labels = lvl.Labels.ToDateTime()
	lvl.DataType = datatypes.DateTime
	lvl.Refresh()
	return lvl
}

// ToInterface converts an index level to Interface
func (lvl Level) ToInterface() Level {
	lvl.Labels = lvl.Labels.ToInterface()
	lvl.DataType = datatypes.Interface
	lvl.Refresh()
	return lvl
}
