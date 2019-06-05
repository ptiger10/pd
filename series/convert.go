package series

import (
	"fmt"

	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/kinds"
)

// To converts the series values to the kind supplied
//
// Applies to All. If unsupported Kind is supplied, returns original Series.
func (s Series) To(kind kinds.Kind) Series {
	switch kind {
	case kinds.Float:
		s.values = s.values.ToFloat()
		s.kind = kinds.Float
	case kinds.Int:
		s.values = s.values.ToInt()
		s.kind = kinds.Int
	case kinds.String:
		s.values = s.values.ToString()
		s.kind = kinds.String
	case kinds.Bool:
		s.values = s.values.ToBool()
		s.kind = kinds.Bool
	case kinds.DateTime:
		s.values = s.values.ToDateTime()
		s.kind = kinds.DateTime
	case kinds.Interface:
		s.values = s.values.ToInterface()
		s.kind = kinds.Interface

	default:
		values.Warn(fmt.Errorf("unsupported conversion kind"), "original Series")
		return s
	}
	return s
}

// IndexTo converts the first level of the series index to the kind supplied.
//
// Applies to All. If unsupported Kind is supplied, returns error.
func (s Series) IndexTo(kind kinds.Kind) (Series, error) {
	return s.IndexLevelTo(0, kind)
}

// SetLevel sets the index level at position with the supplied level.
//
//
// func (s Series) SetLevel(position int, level index.Level) (Series, error) {

// }

// IndexLevelTo converts the specific integer level of the series index to the kind supplied
//
// Applies to All. If unsupported Kind or invalid level value is supplied, returns error.
func (s Series) IndexLevelTo(position int, kind kinds.Kind) (Series, error) {
	copyS := s.copy()
	if position >= len(s.index.Levels) {
		return Series{}, fmt.Errorf("unable to convert index at level %d: index out of range (Series has %d levels)", position, len(s.index.Levels))
	}
	lvl, err := copyS.index.Levels[position].Convert(kind)
	if err != nil {
		return Series{}, fmt.Errorf("unable to convert index to kind %v: unsupported kind", kind)
	}
	copyS.index.Levels[position] = lvl
	return copyS, nil
}
