package index

import (
	"fmt"
	"reflect"

	"github.com/ptiger10/pd/new/kinds"
)

func (lvl Level) Convert(kind reflect.Kind) (Level, error) {
	switch kind {
	case kinds.None: // this checks for the pseduo-nil type
		return Level{}, fmt.Errorf("Unable to convert index level: must supply a valid Kind")
	case kinds.Float:
		return lvl.ToFloat(), nil
	case kinds.Int:
		return lvl.ToInt(), nil
	case kinds.String:
		return lvl.ToString(), nil
	case kinds.Bool:
		return lvl.ToBool(), nil
	case kinds.DateTime:
		return lvl.ToDateTime(), nil
	default:
		return Level{}, fmt.Errorf("Unable to convert level: kind not supported: %v", kind)

	}
}

func (lvl Level) ToFloat() Level {
	lvl.Labels = lvl.Labels.ToFloat()
	lvl.Kind = kinds.Float
	return lvl
}

func (lvl Level) ToInt() Level {
	lvl.Labels = lvl.Labels.ToInt()
	lvl.Kind = kinds.Int
	return lvl
}

func (lvl Level) ToString() Level {
	lvl.Labels = lvl.Labels.ToString()
	lvl.Kind = kinds.String
	return lvl
}

func (lvl Level) ToBool() Level {
	lvl.Labels = lvl.Labels.ToBool()
	lvl.Kind = kinds.Bool
	return lvl
}

func (lvl Level) ToDateTime() Level {
	lvl.Labels = lvl.Labels.ToDateTime()
	lvl.Kind = kinds.DateTime
	return lvl
}
