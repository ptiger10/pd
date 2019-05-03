package series

import (
	"github.com/ptiger10/pd/kinds"
)

// As converts the series values to the kind supplied
func (s Series) As(kind kinds.Kind) Series {
	switch kind {
	case kinds.Float:
		s.values = s.values.ToFloat()
		s.Kind = kinds.Float
	case kinds.Int:
		s.values = s.values.ToInt()
		s.Kind = kinds.Int
	case kinds.String:
		s.values = s.values.ToString()
		s.Kind = kinds.String
	case kinds.Bool:
		s.values = s.values.ToBool()
		s.Kind = kinds.Bool
	case kinds.DateTime:
		s.values = s.values.ToDateTime()
		s.Kind = kinds.DateTime
	}
	return s
}

// IndexAs converts the first level of the series index to the kind supplied
func (s Series) IndexAs(kind kinds.Kind) Series {
	lvl, err := s.index.Levels[0].Convert(kind)
	if err != nil {
		return s
	}
	s.index.Levels[0] = lvl
	return s
}

// IndexLevelAs converts the specific integer level of the series index to the kind supplied
func (s Series) IndexLevelAs(level int, kind kinds.Kind) Series {
	lvl, err := s.index.Levels[level].Convert(kind)
	if err != nil {
		return s
	}
	s.index.Levels[level] = lvl
	return s
}
