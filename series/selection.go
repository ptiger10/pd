package series

import (
	"fmt"
	"strings"

	"github.com/ptiger10/pd/internal/values"
)

// [START Series methods]

// At returns the value at a single integer position, or an error if position is out of range.
//
// To return rows at one or more positions, use s.Select().Get())
func (s Series) At(pos int) (interface{}, error) {
	sNew, err := s.in([]int{pos})
	if err != nil {
		return nil, fmt.Errorf("At(): %v", err)
	}
	return sNew.Element(0).Value, nil
}

// [END Series methods]

// [START Selection]

// Select contains methods that return a Selection.
type Select struct {
	s *Series
}

// ByRows selects rows at the specified integer positions.
func (sel Select) ByRows(positions []int) Selection {
	swappable := len(positions) == 2
	return Selection{
		s:            sel.s.copy(),
		rowPositions: positions,
		swappable:    swappable,
		rowsOnly:     true,
	}
}

// ByLabels selects rows at the specified index labels (for index level 0).
// Appends out-of-range errors to Selection.err.
func (sel Select) ByLabels(labels []string) Selection {
	swappable := len(labels) == 2
	var positions []int
	var errList []string
	var err error
	for _, label := range labels {
		val, ok := sel.s.index.Levels[0].LabelMap[label]
		if !ok {
			errList = append(errList, fmt.Sprintf("label %v not in index", label))
		} else {
			positions = append(positions, val...)
		}
	}
	if len(errList) < 1 {
		err = nil
	} else {
		err = fmt.Errorf(strings.Join(errList, " - "))
	}
	return Selection{
		s:            sel.s.copy(),
		rowPositions: positions,
		swappable:    swappable,
		rowsOnly:     true,
		err:          err,
	}
}

// ByLevels selects index levels at the specified integer positions.
func (sel Select) ByLevels(positions []int) Selection {
	swappable := len(positions) == 2
	return Selection{
		s:              sel.s.copy(),
		levelPositions: positions,
		swappable:      swappable,
		levelsOnly:     true,
	}
}

// ByLevelNames selects the index levels at the specified names.
func (sel Select) ByLevelNames(names []string) Selection {
	swappable := len(names) == 2
	var positions []int
	var errList []string
	var err error
	for _, name := range names {
		val, ok := sel.s.index.NameMap[name]
		if !ok {
			errList = append(errList, fmt.Sprintf("name %v not in index levels", name))
		} else {
			positions = append(positions, val...)
		}
	}
	if len(errList) < 1 {
		err = nil
	} else {
		err = fmt.Errorf(strings.Join(errList, " - "))
	}

	return Selection{
		s:              sel.s.copy(),
		levelPositions: positions,
		swappable:      swappable,
		levelsOnly:     true,
		err:            err,
	}
}

// XS selects a cross-section of index rows and levels at the specified integer locations.
func (sel Select) XS(rows []int, levels []int) Selection {
	return Selection{
		s:              sel.s.copy(),
		rowPositions:   rows,
		levelPositions: levels,
	}
}

// ensure checks for errors on the Selection prior to calling other methods.
func (sel Selection) ensure() error {
	if sel.err != nil {
		return sel.err
	}
	if err := sel.s.ensureAlignment(); err != nil {
		return fmt.Errorf("%v", err)
	}
	if err := sel.s.ensureRowPositions(sel.rowPositions); err != nil {
		return fmt.Errorf("%v", err)
	}
	if err := sel.s.ensureLevelPositions(sel.levelPositions); err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

// Series returns the Selection as a new Series.
func (sel Selection) Series() (Series, error) {
	return sel.series()
}

// series returns the rows or index levels specified in the Selection.
func (sel Selection) series() (Series, error) {
	if err := sel.ensure(); err != nil {
		return Series{}, fmt.Errorf("Selection.series(): %v", err)
	}
	if sel.rowPositions == nil {
		sel.rowPositions = values.MakeIntRange(0, sel.s.Len())
	}
	if sel.levelPositions == nil {
		sel.levelPositions = values.MakeIntRange(0, sel.s.index.Len())
	}
	sel.s, _ = sel.s.in(sel.rowPositions)
	sel.s.index, _ = sel.s.index.In(sel.levelPositions)
	return sel.s, nil
}

// if sel.swappable is true, then either sel.rowsOnly or sel.levelsOnly must be true
func (sel Selection) swap() (Series, error) {
	if err := sel.ensure(); err != nil {
		return Series{}, fmt.Errorf("Selection.Swap(): %v", err)
	}
	if !sel.swappable {
		return sel.s, fmt.Errorf("selection is not swappable: must select exactly two of either rows or levels")
	}
	if sel.rowsOnly {
		// swap Rows
		r1 := sel.rowPositions[0]
		r2 := sel.rowPositions[1]
		sel.s.values.Swap(r1, r2)

		for i := 0; i < len(sel.s.index.Levels); i++ {
			sel.s.index.Levels[i].Labels.Swap(r1, r2)
			sel.s.index.Levels[i].Refresh()
		}

	} else {
		// swap Levels
		lvl := sel.s.index.Levels
		lvl[sel.levelPositions[0]], lvl[sel.levelPositions[1]] = lvl[sel.levelPositions[1]], lvl[sel.levelPositions[0]]
		sel.s.index.Refresh()
	}
	return sel.s, nil
}

// Swap swaps the selected rows or index levels and returns a new Series. If Selection is not swappable, returns error.
func (sel Selection) Swap() (Series, error) {
	return sel.swap()
}

// Set sets all the values in a Selection to val and returns a new Series.
func (sel Selection) Set(val interface{}) (Series, error) {
	if err := sel.ensure(); err != nil {
		return Series{}, fmt.Errorf("Selection.Set(): %v", err)
	}
	if sel.rowsOnly {
		for _, row := range sel.rowPositions {
			sel.s.values.Set(row, val)
		}
	} else if sel.levelsOnly {
		for _, level := range sel.levelPositions {
			for i := 0; i < sel.s.index.Levels[0].Len(); i++ {
				sel.s.index.Levels[level].Labels.Set(i, val)
			}
		}
	} else {
		for _, row := range sel.rowPositions {
			sel.s.values.Set(row, val)
		}
		for _, level := range sel.levelPositions {
			for i := 0; i < sel.s.index.Levels[0].Len(); i++ {
				sel.s.index.Levels[level].Labels.Set(i, val)
			}
		}
	}
	return sel.s, nil
}

// Drop drops all the values or index levels in a Selection and returns a new Series.
// Will not drop an index level if it is the last remaining level.
func (sel Selection) Drop() (Series, error) {
	if err := sel.ensure(); err != nil {
		return Series{}, fmt.Errorf("Selection.Drop(): %v", err)
	}
	if sel.rowsOnly {
		sel.s.InPlace.dropRows(sel.rowPositions)
	} else if sel.levelsOnly {
		for _, level := range sel.levelPositions {
			sel.s.index.Drop(level)
		}
	} else {
		sel.s.InPlace.dropRows(sel.rowPositions)
		for _, level := range sel.levelPositions {
			sel.s.index.Drop(level)
		}
	}
	return sel.s, nil
}

// A Selection is a portion of a Series for use as an intermediate step in transforming data,
// such as getting, setting, swapping, or dropping it.
type Selection struct {
	s              Series
	levelPositions []int
	rowPositions   []int
	swappable      bool
	levelsOnly     bool
	rowsOnly       bool
	err            error
}

func (sel Selection) String() string {
	return fmt.Sprintf(
		"Selection {rows: %v, levels: %v, swappable: %v, hasError: %v}. \nTo print underlying Series, call .Series() ",
		sel.rowPositions, sel.levelPositions, sel.swappable, sel.err != nil)
}
