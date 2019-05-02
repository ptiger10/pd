package index

import (
	"fmt"
	"reflect"

	"github.com/ptiger10/pd/new/kinds"
)

// Convert an index level from one kind to another
func (lvl Level) Convert(kind reflect.Kind) (Level, error) {
	var convertedLvl Level
	switch kind {
	case kinds.None: // this checks for the pseduo-nil type
		return Level{}, fmt.Errorf("Unable to convert index level: must supply a valid Kind")
	case kinds.Float:
		convertedLvl = lvl.toFloat()
	case kinds.Int:
		convertedLvl = lvl.toInt()
	case kinds.String:
		convertedLvl = lvl.toString()
	case kinds.Bool:
		convertedLvl = lvl.toBool()
	case kinds.DateTime:
		convertedLvl = lvl.toDateTime()
	default:
		return Level{}, fmt.Errorf("Unable to convert level: kind not supported: %v", kind)
	}
	convertedLvl.Refresh()
	return convertedLvl, nil
}

func (lvl Level) toFloat() Level {
	lvl.Labels = lvl.Labels.ToFloat()
	lvl.Kind = kinds.Float
	return lvl
}

func (lvl Level) toInt() Level {
	lvl.Labels = lvl.Labels.ToInt()
	lvl.Kind = kinds.Int
	return lvl
}

func (lvl Level) toString() Level {
	lvl.Labels = lvl.Labels.ToString()
	lvl.Kind = kinds.String
	return lvl
}

func (lvl Level) toBool() Level {
	lvl.Labels = lvl.Labels.ToBool()
	lvl.Kind = kinds.Bool
	return lvl
}

func (lvl Level) toDateTime() Level {
	lvl.Labels = lvl.Labels.ToDateTime()
	lvl.Kind = kinds.DateTime
	return lvl
}
