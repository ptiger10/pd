package series

import (
	"github.com/ptiger10/pd/datatypes"
	"github.com/ptiger10/pd/internal/index"
)

// A Group is a combination of group labels and the integer positions associated with each label.
type Group struct {
	s        *Series
	groups   index.LabelMap
	idxKinds []datatypes.DataType
}

// Sum all the groups
func (g Group) Sum() Series {
	s, _ := NewPointer(nil)
	for k, v := range g.groups {
		s2 := g.s.mustIn(v).copy()
		d := s2.Math.Sum()
		// var idxVals []
		// for _, word := range strings.Split(k, "//") {

		// }
		newS := mustNew(d, Idx(k))
		s.InPlace.Join(newS)
	}
	return *s
}

// GroupByIndex groups a Series by index level 0.
func (s Series) GroupByIndex() Group {
	return Group{
		s:        &s,
		groups:   s.index.Levels[0].LabelMap,
		idxKinds: s.index.Kinds(),
	}
}
