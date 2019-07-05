package dataframe

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/series"
)

type group struct {
	Index     []interface{}
	Positions []int
}

func (grp *group) copy() *group {
	pos := make([]int, len(grp.Positions))
	idx := make([]interface{}, len(grp.Index))
	for i, p := range grp.Positions {
		pos[i] = p
	}
	for i, ind := range grp.Index {
		idx[i] = ind
	}
	return &group{Positions: pos, Index: idx}
}

// copy a grouping
func (g Grouping) copy() Grouping {
	grps := make(map[string]*group)
	for k, v := range g.groups {
		grps[k] = v.copy()
	}
	return Grouping{
		df:     g.df.Copy(),
		groups: grps,
	}
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

// Group returns the DataFrame with the given group label, or an error if that label does not exist.
func (g Grouping) Group(label string) *DataFrame {
	group, ok := g.groups[label]
	if !ok {
		if options.GetLogWarnings() {
			log.Printf("s.Grouping.Group(): label %v not in g.Groups()", label)
		}
		return newEmptyDataFrame()
	}
	s := g.df.subsetRows(group.Positions)
	return s
}

func newEmptyGrouping() Grouping {
	groups := make(map[string]*group)
	df := newEmptyDataFrame()
	return Grouping{df: df, groups: groups, err: true}
}

// GroupByIndex groups a DataFrame by one or more of its index levels. If no level is provided, all index levels are used.
func (df *DataFrame) GroupByIndex(levelPositions ...int) Grouping {
	if len(levelPositions) != 0 {
		var err error
		df, err = df.Index.SubsetLevels(levelPositions)
		if err != nil {
			if options.GetLogWarnings() {
				log.Printf("df.GroupByIndex() %v\n", err)
			}
			return newEmptyGrouping()
		}
	}

	// Default: use all label level positions
	return df.groupby()
}

// GroupBy groups a DataFrame by one or more columns.
// If no column is supplied or an invalid column is supplied, an empty grouping is returned.
func (df *DataFrame) GroupBy(cols ...int) Grouping {
	if len(cols) == 0 {
		if options.GetLogWarnings() {
			log.Print("df.GroupBy(): empty cols, returning empty Grouping\n")
		}
		return newEmptyGrouping()
	}
	if len(cols) == df.NumCols() {
		if options.GetLogWarnings() {
			log.Print("df.GroupBy(): at least one column must be excluded from the Grouping\n")
		}
		return newEmptyGrouping()
	}
	if err := df.ensureColumnPositions(cols); err != nil {
		if options.GetLogWarnings() {
			log.Printf("df.GroupBy(): %v\n", err)
		}
		return newEmptyGrouping()
	}
	df = df.Copy()
	df.InPlace.replaceIndex(cols)

	return df.groupby()
}

func (ip InPlace) replaceIndex(cols []int) {

	lengthArchive := ip.df.IndexLevels()
	// set new levels
	ip.setIndexes(cols)

	// Drop old levels
	for j := len(cols); j < len(cols)+lengthArchive; j++ {
		// use lower-level inplace method
		// duck error because level is certain to be in index
		ip.df.index.DropLevel(j)
	}
	ip.df.index.Refresh()
}

func (df *DataFrame) groupby() Grouping {
	groups := make(map[string]*group)
	for i := 0; i < df.Len(); i++ {
		labels := df.Row(i).Labels
		var strLabels []string
		for _, label := range labels {
			strLabels = append(strLabels, fmt.Sprint(label))
		}
		groupLabel := strings.Join(strLabels, " | ")

		// create group with groupLabel and index labels if it is not already within groups map
		if _, ok := groups[groupLabel]; !ok {
			groups[groupLabel] = &group{Index: labels}
		}
		groups[groupLabel].Positions = append(groups[groupLabel].Positions, i)
	}
	return Grouping{df: df, groups: groups}
}

// First returns the first occurence of each grouping in the DataFrame.
func (g Grouping) First() *DataFrame {
	first := func(group string) *DataFrame {
		position := g.groups[group].Positions[0]
		df := g.df.subsetRows([]int{position})
		return df
	}
	ret := newEmptyDataFrame()
	for _, group := range g.Groups() {
		df := first(group)
		ret.InPlace.appendDataFrameRow(df)
	}
	return ret
}

// Last returns the last occurence of each grouping in the DataFrame.
func (g Grouping) Last() *DataFrame {
	last := func(group string) *DataFrame {
		lastIdx := len(g.groups[group].Positions) - 1
		position := g.groups[group].Positions[lastIdx]
		df := g.df.subsetRows([]int{position})
		return df
	}
	ret := newEmptyDataFrame()
	for _, group := range g.Groups() {
		df := last(group)
		ret.InPlace.appendDataFrameRow(df)
	}
	return ret
}

var wg sync.WaitGroup

func (g Grouping) asyncMath(fn func(*DataFrame) *series.Series) *DataFrame {
	g = g.copy()
	if g.Len() == 0 {
		return newEmptyDataFrame()
	}

	// synchronous option
	if !options.GetAsync() {
		ret := newEmptyDataFrame()
		for _, group := range g.Groups() {
			df := g.math(group, fn)
			ret.InPlace.appendDataFrameRow(df)
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

	df := newEmptyDataFrame()
	for _, result := range container {
		df.InPlace.appendDataFrameRow((&result.df))
	}
	df.index.Refresh()
	return df
}

type calcReturn struct {
	df DataFrame
	n  int
}

func (g Grouping) awaitMath(ch chan<- calcReturn, n int, group string, fn func(*DataFrame) *series.Series) {
	df := g.math(group, fn)
	ret := calcReturn{df: *df, n: n}
	ch <- ret
	wg.Done()
}

func (g Grouping) math(group string, fn func(*DataFrame) *series.Series) *DataFrame {
	positions := g.groups[group].Positions
	rows := g.df.subsetRows(positions)
	// df := newEmptyDataFrame()
	for m := 0; m < g.df.NumCols(); m++ {

	}
	calc := fn(rows)
	calc.Rename(group)
	fmt.Println(calc)
	// construct new DataFrame

	// df := MustNew([]interface{}{calc})
	// df.index = g.groups[group].Index
	return MustNew(nil)
}

// Sum for each group in the Grouping.
func (g Grouping) Sum() *DataFrame {
	return g.asyncMath((*DataFrame).Sum)
}

// Mean for each group in the Grouping.
func (g Grouping) Mean() *DataFrame {
	return g.asyncMath((*DataFrame).Mean)
}

// Min for each group in the Grouping.
func (g Grouping) Min() *DataFrame {
	return g.asyncMath((*DataFrame).Min)
}

// Max for each group in the Grouping.
func (g Grouping) Max() *DataFrame {
	return g.asyncMath((*DataFrame).Max)
}

// Median for each group in the Grouping.
func (g Grouping) Median() *DataFrame {
	return g.asyncMath((*DataFrame).Median)
}

// Std for each group in the Grouping.
func (g Grouping) Std() *DataFrame {
	return g.asyncMath((*DataFrame).Std)
}
