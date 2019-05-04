package series

import (
	"fmt"
	"log"

	"github.com/ptiger10/pd/internal/index"

	"github.com/jinzhu/copier"
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
		log.Print("Unsupported kind - returning original Series")
		return s
	}
	return s
}

func (s Series) copySeries() Series {
	idx := index.Index{}
	copier.Copy(&idx, &s.index)
	idx.Levels = make([]index.Level, len(s.index.Levels))
	for i := 0; i < len(s.index.Levels); i++ {
		copier.Copy(&idx.Levels[i], &s.index.Levels[i])
	}

	copyS := Series{
		values: s.values,
		index:  idx,
		kind:   s.kind,
		Name:   s.Name,
	}
	return copyS
}

// IndexAs converts the first level of the series index to the kind supplied.
//
// Applies to All. If unsupported Kind is supplied, returns error.
func (s Series) IndexAs(kind kinds.Kind) (Series, error) {
	return s.IndexLevelAs(0, kind)
}

// SetLevel sets the index level at position with the supplied level.
//
//
// func (s Series) SetLevel(position int, level index.Level) (Series, error) {

// }

// IndexLevelAs converts the specific integer level of the series index to the kind supplied
//
// Applies to All. If unsupported Kind or invalid level value is supplied, returns error.
func (s Series) IndexLevelAs(position int, kind kinds.Kind) (Series, error) {
	copyS := s.copySeries()
	if position >= len(s.index.Levels) {
		return Series{}, fmt.Errorf("Unable to convert index at level %d: index out of range (Series has %d levels)", position, len(s.index.Levels))
	}
	lvl, err := copyS.index.Levels[position].Convert(kind)
	if err != nil {
		return Series{}, fmt.Errorf("Unable to convert index to kind %v: unsupported kind", kind)
	}
	copyS.index.Levels[position] = lvl
	return copyS, nil
}
