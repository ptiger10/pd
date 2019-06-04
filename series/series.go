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

func (s Series) copy() Series {
	idx := s.index.Copy()
	copyS := Series{
		values: s.values,
		index:  idx,
		kind:   s.kind,
		Name:   s.Name,
	}
	return copyS
}

// in subsets a Series to include only index items and values at the positions supplied,
// then returns as a new Series
func (s Series) in(positions []int) (Series, error) {
	if ok := s.ensureAlignment(); !ok {
		return s, fmt.Errorf("fatal error: Series values and index labels out of alignment: report issue and create new series")
	}
	newS := s.copy()
	values, err := newS.values.In(positions)
	if err != nil {
		return Series{}, fmt.Errorf("unable to get Series values at positions: %v", err)
	}
	newS.values = values
	for i, level := range newS.index.Levels {
		// Ducks error because positional alignment is ensured between values and all index levels
		newS.index.Levels[i].Labels, _ = level.Labels.In(positions)
	}
	newS.index.Refresh()
	return newS, nil
}
