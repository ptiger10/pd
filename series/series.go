package series

import (
	"fmt"
	"reflect"

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
	Math   Math
}

// Kind is the data kind of the Series' values. Mimics reflect.Kind with the addition of time.Time
func (s Series) Kind() string {
	return fmt.Sprint(s.kind)
}

func (s Series) copy() Series {
	idx := s.index.Copy()
	valsCopy := s.values.Copy()
	copyS := &Series{
		values: valsCopy,
		index:  idx,
		kind:   s.kind,
		Name:   s.Name,
	}
	copyS.Math = Math{s: copyS}
	return *copyS
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

func seriesEquals(s1, s2 Series) bool {
	sameIndex := reflect.DeepEqual(s1.index, s2.index)
	sameValues := reflect.DeepEqual(s1.values, s2.values)
	sameName := s1.Name == s2.Name
	sameKind := s1.kind == s2.kind
	if sameIndex && sameValues && sameName && sameKind {
		return true
	}
	return false
}
