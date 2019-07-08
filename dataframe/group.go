package dataframe

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/araddon/dateparse"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/series"
)

type group struct {
	Positions     []int
	FirstPosition int
}

func (grp *group) copy() *group {
	pos := make([]int, len(grp.Positions))
	for i, p := range grp.Positions {
		pos[i] = p
	}
	return &group{Positions: pos, FirstPosition: grp.FirstPosition}
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

// SortedGroups returns all valid group labels in the Grouping, sorted in alphabetical order.
func (g Grouping) SortedGroups() []string {
	var keys []string
	for k := range g.groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Groups returns all valid group labels in the Grouping, in their original group position.
func (g Grouping) Groups() []string {
	type groupContainer struct {
		grp   *group
		label string
	}
	var orderedGroups []groupContainer
	for k, v := range g.groups {
		orderedGroups = append(orderedGroups, groupContainer{grp: v, label: k})
	}
	sort.Slice(orderedGroups, func(i, j int) bool {
		if orderedGroups[i].grp.FirstPosition < orderedGroups[j].grp.FirstPosition {
			return true
		}
		return false
	})
	var labels []string
	for _, grp := range orderedGroups {
		labels = append(labels, grp.label)
	}
	return labels
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
		df = df.Copy()
		err := df.Index.SubsetLevels(levelPositions)
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
		// use lower-level method to change index in place and duck error because level is certain to be in index
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
		groupLabel := strings.Join(strLabels, values.GetMultiColNameSeparator())

		// create group with groupLabel and index labels if it is not already within groups map
		if _, ok := groups[groupLabel]; !ok {
			groups[groupLabel] = &group{FirstPosition: i}
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

type calcReturn struct {
	df *DataFrame
	n  int
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
			df := g.mathSingleGroup(group, fn)
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
	var returnedData []calcReturn
	for result := range ch {
		returnedData = append(returnedData, result)
	}
	sort.Slice(returnedData, func(i, j int) bool {
		return returnedData[i].n < returnedData[j].n
	})

	df := newEmptyDataFrame()
	for _, result := range returnedData {
		df.InPlace.appendDataFrameRow((result.df))
	}
	df.index.Refresh()
	return df
}

func (g Grouping) awaitMath(ch chan<- calcReturn, n int, group string, fn func(*DataFrame) *series.Series) {
	df := g.mathSingleGroup(group, fn)
	ret := calcReturn{df: df, n: n}
	ch <- ret
	wg.Done()
}

func parseStringIntoValuesContainer(s string) values.Container {
	var container values.Container
	if intVal, err := strconv.Atoi(s); err == nil {
		container = values.MustCreateValuesFromInterface(intVal)
	} else if floatVal, err := strconv.ParseFloat(s, 64); err == nil {
		container = values.MustCreateValuesFromInterface(floatVal)
	} else if boolVal, err := strconv.ParseBool(s); err == nil {
		container = values.MustCreateValuesFromInterface(boolVal)
	} else if dateTimeVal, err := dateparse.ParseAny(s); err == nil {
		container = values.MustCreateValuesFromInterface(dateTimeVal)
	} else {
		container = values.MustCreateValuesFromInterface(s)
	}
	return container
}

func (g Grouping) mathSingleGroup(group string, fn func(*DataFrame) *series.Series) *DataFrame {
	positions := g.groups[group].Positions
	rows := g.df.subsetRows(positions)
	calc := fn(rows)
	calc.Rename(group)
	df := transposeSeries(calc)
	return df
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
