package series

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

// A Grouping relates group labels to integer positions.
type Grouping struct {
	s      *Series
	groups map[string]*group
}

type group struct {
	IndexLevels []interface{}
	Positions   []int
}

// Sum for each group in the Grouping.
func (g Grouping) Sum() *Series {
	s, _ := New(nil)
	for _, group := range g.Groups() {
		positions := g.groups[group].Positions
		sum := g.s.mustIn(positions).Sum()
		newS := MustNew(sum, Config{MultiIndex: g.groups[group].IndexLevels})
		s.InPlace.Join(newS)
	}
	s.index.Refresh()
	return s
}

// func (g group) buildIndex() []interface{} {
// 	var idxLevels []interface{}
// 	for _, lvl := range g.IndexLevels {
// 		idxLevels = append(idxLevels, Idx(lvl))
// 	}
// 	return idxLevels
// }

// Groups returns all valid group labels in the Grouping.
func (g Grouping) Groups() []string {
	var keys []string
	for k := range g.groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Group returns the Series with the given group label, or an error if that label does not exist.
func (g Grouping) Group(label string) (*Series, error) {
	group, ok := g.groups[label]
	if !ok {
		return nil, fmt.Errorf("Grouping.Group(): %v not a valid group label", label)
	}
	s, err := g.s.in(group.Positions)
	if err != nil {
		return nil, fmt.Errorf("Grouping.Group(): %v", err)
	}
	return s, nil
}

// GroupByIndex groups a Series by all of its index levels.
func (s *Series) GroupByIndex() Grouping {

	g := Grouping{s: s, groups: make(map[string]*group)}
	for i := 0; i < s.Len(); i++ {
		var levels []interface{}
		var labels []string
		for j := 0; j < s.index.Len(); j++ {
			idx, err := s.Index.At(i, j)
			if err != nil {
				log.Printf("series.GroupByIndex(): %v", err)
				return Grouping{}
			}
			levels = append(levels, idx)
			labels = append(labels, fmt.Sprint(idx))
		}
		label := strings.Join(labels, " ")
		if g.groups[label] == nil {
			g.groups[label] = &group{}
		}
		g.groups[label].Positions = append(g.groups[label].Positions, i)
		g.groups[label].IndexLevels = levels
	}
	return g
}
