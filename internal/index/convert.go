package index

import (
	"fmt"

	"github.com/ptiger10/pd/options"
)

// Convert an index level from one kind to another, then refresh the LabelMap
func (lvl Level) Convert(kind options.DataType) (Level, error) {
	var convertedLvl Level
	switch kind {
	case options.None:
		return Level{}, fmt.Errorf("unable to convert index level: must supply a valid Kind")
	case options.Float64:
		convertedLvl = lvl.ToFloat()
	case options.Int64:
		convertedLvl = lvl.ToInt()
	case options.String:
		convertedLvl = lvl.ToString()
	case options.Bool:
		convertedLvl = lvl.ToBool()
	case options.DateTime:
		convertedLvl = lvl.ToDateTime()
	case options.Interface:
		convertedLvl = lvl.ToInterface()
	default:
		return Level{}, fmt.Errorf("unable to convert level: kind not supported: %v", kind)
	}
	return convertedLvl, nil
}

// ToFloat converts an index level to Float
func (lvl Level) ToFloat() Level {
	lvl.Labels = lvl.Labels.ToFloat()
	lvl.DataType = options.Float64
	lvl.Refresh()
	return lvl
}

// ToInt converts an index level to Int
func (lvl Level) ToInt() Level {
	lvl.Labels = lvl.Labels.ToInt()
	lvl.DataType = options.Int64
	lvl.Refresh()
	return lvl
}

// ToString converts an index level to String
func (lvl Level) ToString() Level {
	lvl.Labels = lvl.Labels.ToString()
	lvl.DataType = options.String
	lvl.Refresh()
	return lvl
}

// ToBool converts an index level to Bool
func (lvl Level) ToBool() Level {
	lvl.Labels = lvl.Labels.ToBool()
	lvl.DataType = options.Bool
	lvl.Refresh()
	return lvl
}

// ToDateTime converts an index level to DateTime
func (lvl Level) ToDateTime() Level {
	lvl.Labels = lvl.Labels.ToDateTime()
	lvl.DataType = options.DateTime
	lvl.Refresh()
	return lvl
}

// ToInterface converts an index level to Interface
func (lvl Level) ToInterface() Level {
	lvl.Labels = lvl.Labels.ToInterface()
	lvl.DataType = options.Interface
	lvl.Refresh()
	return lvl
}
