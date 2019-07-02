package dataframe

import (
	"github.com/ptiger10/pd/internal/index"
)

// func (g Grouping) String() string {
// 	printer := fmt.Sprintf("{Grouping | NumGroups: %v, Groups: [%v]}\n", len(g.groups), strings.Join(g.Groups(), ", "))
// 	return printer
// }

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
		df:     g.df.Copy(),
		groups: grps,
	}
}

// var wg sync.WaitGroup

// func (g Grouping) asyncMath(fn func(*DataFrame) float64) *DataFrame {
// 	g = g.copy()
// 	if g.Len() == 0 {
// 		return newEmptyDataFrame()
// 	}

// 	// synchronous option
// 	if !options.GetAsync() {
// 		ret := newEmptyDataFrame()
// 		for _, group := range g.Groups() {
// 			df := g.math(group, fn)
// 			ret.InPlace.Join(df)
// 		}
// 		return ret
// 	}

// 	// asynchronous option
// 	ch := make(chan calcReturn, g.Len())
// 	for i, group := range g.Groups() {
// 		wg.Add(1)
// 		go g.awaitMath(ch, i, group, fn)
// 	}
// 	wg.Wait()
// 	close(ch)
// 	container := make([]calcReturn, g.Len())
// 	// iterating over channel range returns nil DataFrame if pointer is provided instead of value
// 	for result := range ch {
// 		container = append(container, result)
// 	}
// 	sort.Slice(container, func(i, j int) bool {
// 		return container[i].n < container[j].n
// 	})

// 	s := newEmptyDataFrame()
// 	for _, result := range container {
// 		s.InPlace.Join(&result.s)
// 	}
// 	s.index.Refresh()
// 	return s
// }

// type calcReturn struct {
// 	s DataFrame
// 	n int
// }

// func (g Grouping) awaitMath(ch chan<- calcReturn, n int, group string, fn func(*DataFrame) float64) {
// 	s := g.math(group, fn)
// 	ret := calcReturn{s: *s, n: n}
// 	ch <- ret
// 	wg.Done()
// }

// func (g Grouping) math(group string, fn func(*DataFrame) float64) *DataFrame {
// 	positions := g.groups[group].Positions
// 	rows := g.df.subsetRows(positions)
// 	calc := fn(rows)
// 	s := MustNew(calc)
// 	s.index = g.groups[group].Index
// 	return s
// }

// // Groups returns all valid group labels in the Grouping.
// func (g Grouping) Groups() []string {
// 	var keys []string
// 	for k := range g.groups {
// 		keys = append(keys, k)
// 	}
// 	sort.Strings(keys)
// 	return keys
// }

// // Len returns the number of groups in the Grouping.
// func (g Grouping) Len() int {
// 	return len(g.groups)
// }

// // Group returns the DataFrame with the given group label, or an error if that label does not exist.
// func (g Grouping) Group(label string) *DataFrame {
// 	group, ok := g.groups[label]
// 	if !ok {
// 		if options.GetLogWarnings() {
// 			log.Printf("s.Grouping.Group(): label %v not in g.Groups()", label)
// 		}
// 		return newEmptyDataFrame()
// 	}
// 	s := g.df.subsetRows(group.Positions)
// 	return s
// }

func newEmptyGrouping() Grouping {
	groups := make(map[string]*group)
	df := newEmptyDataFrame()
	return Grouping{df: df, groups: groups}
}

// // GroupByIndex groups a DataFrame by one or more of its index levels. If no int is provided, all index levels are used.
// func (df *DataFrame) GroupByIndex(levelPositions ...int) Grouping {
// 	groups := make(map[string]*group)
// 	if len(levelPositions) != 0 {
// 		var err error
// 		df, err = df.Index.SubsetLevels(levelPositions)
// 		if err != nil {
// 			if options.GetLogWarnings() {
// 				log.Printf("df.GroupByIndex() %v\n", err)
// 			}
// 			return newEmptyGrouping()
// 		}
// 	}

// 	for i := 0; i < df.Len(); i++ {
// 		row := df.subsetRows([]int{i})
// 		labels := row.Element(0).Labels
// 		var strLabels []string
// 		for _, label := range labels {
// 			strLabels = append(strLabels, fmt.Sprint(label))
// 		}
// 		groupLabel := strings.Join(strLabels, " ")

// 		if _, ok := groups[groupLabel]; !ok {
// 			groups[groupLabel] = &group{Index: row.index}
// 		}
// 		groups[groupLabel].Positions = append(groups[groupLabel].Positions, i)
// 	}
// 	return Grouping{df: df, groups: groups}
// }

// // First returns the first occurence of each grouping in the DataFrame.
// func (g Grouping) First() *DataFrame {
// 	first := func(group string) *DataFrame {
// 		position := g.groups[group].Positions[0]
// 		s := g.df.subsetRows([]int{position})
// 		return s
// 	}
// 	ret := newEmptyDataFrame()
// 	for _, group := range g.Groups() {
// 		s := first(group)
// 		ret.InPlace.Join(s)
// 	}
// 	return ret
// }

// // Last returns the last occurence of each grouping in the DataFrame.
// func (g Grouping) Last() *DataFrame {
// 	last := func(group string) *DataFrame {
// 		lastIdx := len(g.groups[group].Positions) - 1
// 		position := g.groups[group].Positions[lastIdx]
// 		s := g.df.subsetRows([]int{position})
// 		return s
// 	}
// 	ret := newEmptyDataFrame()
// 	for _, group := range g.Groups() {
// 		s := last(group)
// 		ret.InPlace.Join(s)
// 	}
// 	return ret
// }

// // Sum for each group in the Grouping.
// func (g Grouping) Sum() *DataFrame {
// 	return g.asyncMath((*DataFrame).Sum)
// }

// // Mean for each group in the Grouping.
// func (g Grouping) Mean() *DataFrame {
// 	return g.asyncMath((*DataFrame).Mean)
// }

// // Min for each group in the Grouping.
// func (g Grouping) Min() *DataFrame {
// 	return g.asyncMath((*DataFrame).Min)
// }

// // Max for each group in the Grouping.
// func (g Grouping) Max() *DataFrame {
// 	return g.asyncMath((*DataFrame).Max)
// }

// // Median for each group in the Grouping.
// func (g Grouping) Median() *DataFrame {
// 	return g.asyncMath((*DataFrame).Median)
// }

// // Std for each group in the Grouping.
// func (g Grouping) Std() *DataFrame {
// 	return g.asyncMath((*DataFrame).Std)
// }
