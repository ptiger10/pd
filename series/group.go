package series

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/ptiger10/pd/options"

	"github.com/ptiger10/pd/internal/index"
)

// A Grouping relates group labels to integer positions.
type Grouping struct {
	s      *Series
	groups map[string]*group
}

type group struct {
	Index     index.Index
	Positions []int
}

func (grp *group) copy() *group {
	pos := make([]int, 0)
	for _, p := range grp.Positions {
		pos = append(pos, p)
	}
	return &group{Positions: pos, Index: grp.Index.Copy()}
}

// copy a grouping
func (g Grouping) copy() Grouping {
	grps := make(map[string]*group)
	for k, v := range g.groups {
		grps[k] = v.copy()
	}
	return Grouping{
		s:      g.s.Copy(),
		groups: grps,
	}
}

var wg sync.WaitGroup

func (g Grouping) asyncMath(fn func(*Series) float64) *Series {
	g = g.copy()
	if g.Len() == 0 {
		return newEmptySeries()
	}

	// synchronous option
	if !options.GetAsync() {
		ret := newEmptySeries()
		for _, group := range g.Groups() {
			s := g.math(group, fn)
			ret.InPlace.Join(s)
		}
		return ret
	}

	// asynchronous option
	ch := make(chan calcReturn, g.Len())
	for i, group := range g.Groups() {
		wg.Add(1)
		go g.awaitMath(ch, i, group, fn)
	}
	wg.Wait()
	close(ch)
	container := make([]calcReturn, g.Len())
	// iterating over channel range returns nil Series if pointer is provided instead of value
	for result := range ch {
		container = append(container, result)
	}
	sort.Slice(container, func(i, j int) bool {
		return container[i].n < container[j].n
	})

	s := newEmptySeries()
	for _, result := range container {
		s.InPlace.Join(&result.s)
	}
	s.index.Refresh()
	return s
}

type calcReturn struct {
	s Series
	n int
}

func (g Grouping) awaitMath(ch chan<- calcReturn, n int, group string, fn func(*Series) float64) {
	s := g.math(group, fn)
	ret := calcReturn{s: *s, n: n}
	ch <- ret
	wg.Done()
}

func (g Grouping) math(group string, fn func(*Series) float64) *Series {
	positions := g.groups[group].Positions
	rows := g.s.subsetRows(positions)
	calc := fn(rows)
	s := MustNew(calc)
	s.index = g.groups[group].Index
	return s
}

// Groups returns all valid group labels in the Grouping.
func (g Grouping) Groups() []string {
	var keys []string
	for k := range g.groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Len returns the number of groups in the Grouping.
func (g Grouping) Len() int {
	return len(g.groups)
}

// Group returns the Series with the given group label, or an error if that label does not exist.
func (g Grouping) Group(label string) *Series {
	group, ok := g.groups[label]
	if !ok {
		return newEmptySeries()
	}
	s := g.s.subsetRows(group.Positions)
	return s
}

func newEmptyGrouping() Grouping {
	groups := make(map[string]*group)
	s := newEmptySeries()
	return Grouping{s: s, groups: groups}
}

// GroupByIndex groups a Series by one or more of its index levels. If no int is provided, all index levels are used.
func (s *Series) GroupByIndex(levelPositions ...int) Grouping {
	groups := make(map[string]*group)
	if len(levelPositions) != 0 {
		var err error
		s, err = s.Index.SubsetLevels(levelPositions)
		if err != nil {
			if options.GetLogWarnings() {
				log.Printf("s.GroupByIndex() %v\n", err)
			}
			return newEmptyGrouping()
		}
	}

	for i := 0; i < s.Len(); i++ {
		row := s.subsetRows([]int{i})
		labels := row.Element(0).Labels
		var strLabels []string
		for _, label := range labels {
			strLabels = append(strLabels, fmt.Sprint(label))
		}
		groupLabel := strings.Join(strLabels, " ")

		if _, ok := groups[groupLabel]; !ok {
			groups[groupLabel] = &group{Index: row.index}
		}
		groups[groupLabel].Positions = append(groups[groupLabel].Positions, i)
	}
	return Grouping{s: s, groups: groups}
}

// First returns the first occurence of each grouping in the Series.
func (g Grouping) First() *Series {
	first := func(group string) *Series {
		position := g.groups[group].Positions[0]
		s := g.s.subsetRows([]int{position})
		return s
	}
	ret := newEmptySeries()
	for _, group := range g.Groups() {
		s := first(group)
		ret.InPlace.Join(s)
	}
	return ret
}

// Last returns the last occurence of each grouping in the Series.
func (g Grouping) Last() *Series {
	last := func(group string) *Series {
		lastIdx := len(g.groups[group].Positions) - 1
		position := g.groups[group].Positions[lastIdx]
		s := g.s.subsetRows([]int{position})
		return s
	}
	ret := newEmptySeries()
	for _, group := range g.Groups() {
		s := last(group)
		ret.InPlace.Join(s)
	}
	return ret
}

// Sum for each group in the Grouping.
func (g Grouping) Sum() *Series {
	return g.asyncMath((*Series).Sum)
}

// Mean for each group in the Grouping.
func (g Grouping) Mean() *Series {
	return g.asyncMath((*Series).Mean)
}

// Min for each group in the Grouping.
func (g Grouping) Min() *Series {
	return g.asyncMath((*Series).Min)
}

// Max for each group in the Grouping.
func (g Grouping) Max() *Series {
	return g.asyncMath((*Series).Max)
}

// Median for each group in the Grouping.
func (g Grouping) Median() *Series {
	return g.asyncMath((*Series).Median)
}

// Std for each group in the Grouping.
func (g Grouping) Std() *Series {
	return g.asyncMath((*Series).Std)
}
