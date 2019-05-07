package series

import (
	"fmt"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/kinds"
)

// A Series is a 1-D data container with a labeled index, static type, and the ability to handle null values
type Series struct {
	index  index.Index
	values values.Values
	kind   kinds.Kind
	Name   string
}

// Kind is the data kind of the Series' values. Mimics reflect.Kind with the addition of time.Time
func (s Series) Kind() string {
	return fmt.Sprint(s.kind)
}

func (s Series) in(positions []int) (Series, error) {
	if ok := s.ensureAlignment(); !ok {
		return s, fmt.Errorf("fatal error: Series values and index labels out of alignment: report issue and create new series")
	}
	values, err := s.values.In(positions)
	if err != nil {
		return Series{}, fmt.Errorf("unable to get Series values at positions: %v", err)
	}
	s.values = values
	for i, level := range s.index.Levels {
		// Ducks error because positional alignment is ensured between values and all index levels
		s.index.Levels[i].Labels, _ = level.Labels.In(positions)
	}
	s.index.Refresh()
	return s, nil
}
