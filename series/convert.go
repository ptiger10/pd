package series

import (
	"fmt"

	"github.com/ptiger10/pd/kinds"
)

// As converts the series values to the kind supplied
//
// Applies to All. If unsupported Kind is supplied, returns original Series.
func (s Series) As(kind kinds.Kind) Series {
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
	default:
		return s
	}
	return s
}

// IndexAs converts the first level of the series index to the kind supplied
//
// Applies to All. If unsupported Kind is supplied, returns original Series.
func (s Series) IndexAs(kind kinds.Kind) Series {
	lvl, err := s.index.Levels[0].Convert(kind)
	if err != nil {
		return s
	}
	s.index.Levels[0] = lvl
	return s
}

// IndexLevelAs converts the specific integer level of the series index to the kind supplied
//
// Applies to All. If unsupported Kind, returns original Series. If invalid invalid level is supplied, returns error.
func (s Series) IndexLevelAs(level int, kind kinds.Kind) (Series, error) {
	if level >= len(s.index.Levels) {
		return Series{}, fmt.Errorf("Unable to convert index at level %d: index out of range", level)
	}
	lvl, err := s.index.Levels[level].Convert(kind)
	if err != nil {
		return s, nil
	}
	s.index.Levels[level] = lvl
	return s, nil
}
