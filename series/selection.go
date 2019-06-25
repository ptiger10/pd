package series

import (
	"fmt"
	"log"

	"github.com/ptiger10/pd/options"
)

// [START Series methods]

// At returns the value at a single integer position, or an error if position is out of range.
//
// To return rows at one or more positions, use s.Select().Get())
func (s *Series) At(pos int) (interface{}, error) {
	sNew, err := s.selectByRows([]int{pos})
	if err != nil {
		return nil, fmt.Errorf("At(): %v", err)
	}
	return sNew.Element(0).Value, nil
}

// [END Series methods]

// [START Selection]

// XS returns a new Series with only the rows and index levels at the specified positions.
func (s *Series) XS(rowPositions []int, levelPositions []int) (*Series, error) {
	var err error
	s, err = s.Subset(rowPositions)
	if err != nil {
		return newEmptySeries(), fmt.Errorf("s.XS() rows: %v", err)
	}
	s, err = s.Index.Subset(levelPositions)
	if err != nil {
		return newEmptySeries(), fmt.Errorf("s.XS() index levels: %v", err)
	}
	return s, nil
}

// Set sets all the index labels values in the specified levels to val and returns a new Series.
func (idx Index) Set(levelPositions []int, val interface{}) (*Series, error) {
	s := idx.s.Copy()
	err := idx.s.ensureAlignment()
	if err != nil {
		return newEmptySeries(), fmt.Errorf("s.Index.Set(): %v", err)
	}
	for _, level := range levelPositions {
		for i := 0; i < s.Len(); i++ {
			err := s.index.Levels[level].Labels.Set(i, val)
			if err != nil {
				return newEmptySeries(), fmt.Errorf("s.Set() for val %v: %v", val, err)
			}
		}
	}
	return s, nil
}

// // SelectRows selects rows at the specified integer positions.
// func (s *Series) SelectRows(positions []int) Selection {
// 	swappable := len(positions) == 2
// 	return Selection{
// 		s:            s.Copy(),
// 		rowPositions: positions,
// 		swappable:    swappable,
// 		rowsOnly:     true,
// 	}
// }

// ByLabels returns the integer locations of all rows with the specified index label at level 0.
func (idx Index) ByLabels(labels []string, level int) []int {
	empty := make([]int, 0)
	include := make([]int, 0)
	for _, label := range labels {
		val, ok := idx.s.index.Levels[0].LabelMap[label]
		if !ok {
			if options.GetLogWarnings() {
				log.Printf("s.Index.ByLabels(): %v not in label map", label)
			}
			return empty
		}
		include = append(include, val...)
	}
	return include
}

// SwapLevel swaps the selected index levels and returns a new Series.
func (s *Series) SwapLevel(i, j int) *Series {
	s = s.Copy()
	s.index.Levels[i], s.index.Levels[j] = s.index.Levels[j], s.index.Levels[i]
	s.index.Refresh()
	return s
}

// ByLevelNames returns the integer locations of all index levels with the specified names.
func (idx Index) ByLevelNames(names []string) []int {
	empty := make([]int, 0)
	include := make([]int, 0)
	for _, name := range names {
		val, ok := idx.s.index.NameMap[name]
		if !ok {
			if options.GetLogWarnings() {
				log.Printf("s.Index.ByLevelNames(): %v not in name map", name)
			}
			return empty
		}
		include = append(include, val...)
	}
	return include
}

// Subset returns the Series with the specified index levels only.
func (idx Index) Subset(levelPositions []int) (*Series, error) {
	s := idx.s.Copy()
	index, err := idx.s.index.LevelsIn(levelPositions)
	if err != nil {
		return newEmptySeries(), fmt.Errorf("s.Index.Subset(): %v", err)
	}
	s.index = index
	return s, nil
}

// // SelectLabels selects rows at the specified index labels (for index level 0).
// // Appends out-of-range errors to Selection.err.
// func (s *Series) SelectLabels(labels []string) Selection {
// 	swappable := len(labels) == 2
// 	var positions []int
// 	var errList []string
// 	var err error
// 	for _, label := range labels {
// 		val, ok := s.index.Levels[0].LabelMap[label]
// 		if !ok {
// 			errList = append(errList, fmt.Sprintf("label %v not in index", label))
// 		} else {
// 			positions = append(positions, val...)
// 		}
// 	}
// 	if len(errList) < 1 {
// 		err = nil
// 	} else {
// 		err = fmt.Errorf(strings.Join(errList, " - "))
// 	}
// 	return Selection{
// 		s:            s.Copy(),
// 		rowPositions: positions,
// 		swappable:    swappable,
// 		rowsOnly:     true,
// 		err:          err,
// 	}
// }

// // SelectLevels selects index levels at the specified integer positions.
// func (s *Series) SelectLevels(positions []int) Selection {
// 	swappable := len(positions) == 2
// 	return Selection{
// 		s:              s.Copy(),
// 		levelPositions: positions,
// 		swappable:      swappable,
// 		levelsOnly:     true,
// 	}
// }

// // SelectLevelNames selects index levels at the specified index level names.
// func (s *Series) SelectLevelNames(names []string) Selection {
// 	swappable := len(names) == 2
// 	var positions []int
// 	var errList []string
// 	var err error
// 	for _, name := range names {
// 		val, ok := s.index.NameMap[name]
// 		if !ok {
// 			errList = append(errList, fmt.Sprintf("name %v not in index levels", name))
// 		} else {
// 			positions = append(positions, val...)
// 		}
// 	}
// 	if len(errList) < 1 {
// 		err = nil
// 	} else {
// 		err = fmt.Errorf(strings.Join(errList, " - "))
// 	}

// 	return Selection{
// 		s:              s.Copy(),
// 		levelPositions: positions,
// 		swappable:      swappable,
// 		levelsOnly:     true,
// 		err:            err,
// 	}
// }

// // SelectXS selects a cross-section of index rows and levels at the specified integer locations.
// func (s *Series) SelectXS(rows []int, levels []int) Selection {
// 	return Selection{
// 		s:              s.Copy(),
// 		rowPositions:   rows,
// 		levelPositions: levels,
// 	}
// }

// // GroupBy returns a Grouping for the selected levels.
// func (sel Selection) GroupBy() (Grouping, error) {
// 	if !sel.levelsOnly {
// 		return Grouping{}, fmt.Errorf("Selection.GroupBy() requires that only levels have been selected")
// 	}
// 	g := Grouping{s: sel.s, groups: make(map[string]*group)}
// 	for i := 0; i < sel.s.Len(); i++ {
// 		var levels []interface{}
// 		var labels []string
// 		for j := 0; j < len(sel.levelPositions); j++ {
// 			idx, err := sel.s.Index.At(i, sel.levelPositions[j])
// 			if err != nil {
// 				return Grouping{}, fmt.Errorf("series.GroupByIndex(): %v", err)
// 			}
// 			levels = append(levels, idx)
// 			labels = append(labels, fmt.Sprint(idx))
// 		}
// 		label := strings.Join(labels, " ")
// 		if g.groups[label] == nil {
// 			g.groups[label] = &group{}
// 		}
// 		g.groups[label].Positions = append(g.groups[label].Positions, i)
// 		g.groups[label].IndexLevels = levels
// 	}
// 	return g, nil
// }

// // ensure checks for errors on the Selection prior to calling other methods.
// func (sel Selection) ensure() error {
// 	if sel.err != nil {
// 		return sel.err
// 	}
// 	if err := sel.s.ensureAlignment(); err != nil {
// 		return fmt.Errorf("%v", err)
// 	}
// 	if err := sel.s.ensureRowPositions(sel.rowPositions); err != nil {
// 		return fmt.Errorf("%v", err)
// 	}
// 	if err := sel.s.ensureLevelPositions(sel.levelPositions); err != nil {
// 		return fmt.Errorf("%v", err)
// 	}
// 	return nil
// }

// // Series returns the Selection as a new Series.
// func (sel Selection) Series() (*Series, error) {
// 	return sel.series()
// }

// // series returns the rows or index levels specified in the Selection.
// func (sel Selection) series() (*Series, error) {
// 	if err := sel.ensure(); err != nil {
// 		return nil, fmt.Errorf("Selection.series(): %v", err)
// 	}
// 	if sel.rowPositions == nil {
// 		sel.rowPositions = values.MakeIntRange(0, sel.s.Len())
// 	}
// 	if sel.levelPositions == nil {
// 		sel.levelPositions = values.MakeIntRange(0, sel.s.index.Len())
// 	}
// 	s, _ := sel.s.selectByRows(sel.rowPositions)
// 	s.index, _ = s.index.LevelsIn(sel.levelPositions)
// 	return s, nil
// }

// // if swappable is true, then either sel.rowsOnly or sel.levelsOnly must be true
// func (sel Selection) swap() (*Series, error) {
// 	if err := sel.ensure(); err != nil {
// 		return nil, fmt.Errorf("Selection.Swap(): %v", err)
// 	}
// 	if !sel.swappable {
// 		return sel.s, fmt.Errorf("selection is not swappable: must select exactly two of either rows or levels")
// 	}
// 	if sel.rowsOnly {
// 		// swap Rows
// 		r1 := sel.rowPositions[0]
// 		r2 := sel.rowPositions[1]
// 		sel.s.values.Swap(r1, r2)

// 		for i := 0; i < len(sel.s.index.Levels); i++ {
// 			sel.s.index.Levels[i].Labels.Swap(r1, r2)
// 			sel.s.index.Levels[i].Refresh()
// 		}

// 	} else {
// 		// swap Levels
// 		lvl := sel.s.index.Levels
// 		lvl[sel.levelPositions[0]], lvl[sel.levelPositions[1]] = lvl[sel.levelPositions[1]], lvl[sel.levelPositions[0]]
// 		sel.s.index.Refresh()
// 	}
// 	return sel.s, nil
// }

// // Set sets all the values in a Selection to val and returns a new Series.
// func (sel Selection) Set(val interface{}) (*Series, error) {
// 	if err := sel.ensure(); err != nil {
// 		return nil, fmt.Errorf("Selection.Set(): %v", err)
// 	}
// 	if sel.rowsOnly {
// 		for _, row := range sel.rowPositions {
// 			sel.s.values.Set(row, val)
// 		}
// 	} else if sel.levelsOnly {
// 		for _, level := range sel.levelPositions {
// 			for i := 0; i < sel.s.index.Levels[0].Len(); i++ {
// 				sel.s.index.Levels[level].Labels.Set(i, val)
// 			}
// 		}
// 	} else {
// 		for _, row := range sel.rowPositions {
// 			sel.s.values.Set(row, val)
// 		}
// 		for _, level := range sel.levelPositions {
// 			for i := 0; i < sel.s.index.Levels[0].Len(); i++ {
// 				sel.s.index.Levels[level].Labels.Set(i, val)
// 			}
// 		}
// 	}
// 	return sel.s, nil
// }

// // Drop drops all the values or index levels in a Selection and returns a new Series.
// // Will not drop an index level if it is the last remaining level.
// func (sel Selection) Drop() (*Series, error) {
// 	if err := sel.ensure(); err != nil {
// 		return nil, fmt.Errorf("Selection.Drop(): %v", err)
// 	}
// 	if sel.rowsOnly {
// 		sel.s.InPlace.dropRows(sel.rowPositions)
// 	} else if sel.levelsOnly {
// 		for _, level := range sel.levelPositions {
// 			sel.s.index.Drop(level)
// 		}
// 	} else {
// 		sel.s.InPlace.dropRows(sel.rowPositions)
// 		for _, level := range sel.levelPositions {
// 			sel.s.index.Drop(level)
// 		}
// 	}
// 	return sel.s, nil
// }

// // A Selection is a portion of a Series for use as an intermediate step in transforming data,
// // such as getting, setting, swapping, or dropping it.
// type Selection struct {
// 	s              *Series
// 	levelPositions []int
// 	rowPositions   []int
// 	swappable      bool
// 	levelsOnly     bool
// 	rowsOnly       bool
// 	err            error
// }

// func (sel Selection) String() string {
// 	return fmt.Sprintf(
// 		"Selection{rows: %v, levels: %v, swappable: %v, hasError: %v}",
// 		sel.rowPositions, sel.levelPositions, sel.swappable, sel.err != nil)
// }
