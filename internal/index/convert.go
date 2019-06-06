package index

import (
	"fmt"

	"github.com/ptiger10/pd/kinds"
)

// Convert an index level from one kind to another, then refresh the LabelMap
func (lvl Level) Convert(kind kinds.Kind) (Level, error) {
	var convertedLvl Level
	switch kind {
	case kinds.None:
		return Level{}, fmt.Errorf("unable to convert index level: must supply a valid Kind")
	case kinds.Float:
		convertedLvl = lvl.ToFloat()
	case kinds.Int:
		convertedLvl = lvl.ToInt()
	case kinds.String:
		convertedLvl = lvl.ToString()
	case kinds.Bool:
		convertedLvl = lvl.ToBool()
	case kinds.DateTime:
		convertedLvl = lvl.ToDateTime()
	case kinds.Interface:
		convertedLvl = lvl.ToInterface()
	default:
		return Level{}, fmt.Errorf("unable to convert level: kind not supported: %v", kind)
	}
	return convertedLvl, nil
}

// ToFloat converts an index level to Float
func (lvl Level) ToFloat() Level {
	lvl.Labels = lvl.Labels.ToFloat()
	lvl.Kind = kinds.Float
	lvl.Refresh()
	return lvl
}

// ToInt converts an index level to Int
func (lvl Level) ToInt() Level {
	lvl.Labels = lvl.Labels.ToInt()
	lvl.Kind = kinds.Int
	lvl.Refresh()
	return lvl
}

// ToString converts an index level to String
func (lvl Level) ToString() Level {
	lvl.Labels = lvl.Labels.ToString()
	lvl.Kind = kinds.String
	lvl.Refresh()
	return lvl
}

// ToBool converts an index level to Bool
func (lvl Level) ToBool() Level {
	lvl.Labels = lvl.Labels.ToBool()
	lvl.Kind = kinds.Bool
	lvl.Refresh()
	return lvl
}

// ToDateTime converts an index level to DateTime
func (lvl Level) ToDateTime() Level {
	lvl.Labels = lvl.Labels.ToDateTime()
	lvl.Kind = kinds.DateTime
	lvl.Refresh()
	return lvl
}

// ToInterface converts an index level to Interface
func (lvl Level) ToInterface() Level {
	lvl.Labels = lvl.Labels.ToInterface()
	lvl.Kind = kinds.Interface
	lvl.Refresh()
	return lvl
}
