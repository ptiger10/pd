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

// Rows selects rows at the specified integer positions.
func (sel Select) Rows(positions []int) Selection {
	swappable := len(positions) == 2
	return Selection{
		s:              sel.s.copy(),
		rowPositions:   positions,
		levelPositions: values.MakeIntRange(0, sel.s.index.Len()),
		swappableRows:  swappable,
	}
}

// Labels selects rows at the specified index labels (for index level 0).
// Appends out-of-range errors to Selection.err.
func (sel Select) Labels(labels []string) Selection {
	swappable := len(labels) == 2
	var positions []int
	var errList []string
	for _, label := range labels {
		val, ok := sel.s.index.Levels[0].LabelMap[label]
		if !ok {
			errList = append(errList, fmt.Sprintf("label %v not in index", label))
		} else {
			positions = append(positions, val...)
		}
	}
	return Selection{
		s:              sel.s.copy(),
		rowPositions:   positions,
		levelPositions: values.MakeIntRange(0, sel.s.index.Len()),
		swappableRows:  swappable,
		err:            fmt.Errorf(strings.Join(errList, " - ")),
	}
}

// Levels selects index levels at the specified integer positions.
func (sel Select) Levels(positions []int) Selection {
	swappable := len(positions) == 2
	return Selection{
		s:               sel.s.copy(),
		rowPositions:    values.MakeIntRange(0, sel.s.Len()),
		levelPositions:  positions,
		swappableLevels: swappable,
	}
}

// LevelNames selects the index levels at the specified names.
func (sel Select) LevelNames(names []string) Selection {
	swappable := len(names) == 2
	var positions []int
	var errList []string
	for _, name := range names {
		val, ok := sel.s.index.NameMap[name]
		if !ok {
			errList = append(errList, fmt.Sprintf("name %v not in index levels", name))
		} else {
			positions = append(positions, val...)
		}
	}

	return Selection{
		s:               sel.s.copy(),
		rowPositions:    values.MakeIntRange(0, sel.s.Len()),
		levelPositions:  positions,
		swappableLevels: swappable,
		err:             fmt.Errorf(strings.Join(errList, " - ")),
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

// Series returns the selected Series as a copy.
func (sel Selection) Series() (Series, error) {
	return sel.series()
}

// series returns the rows or index levels specified in the Selection.
func (sel Selection) series() (Series, error) {
	if err := sel.ensure(); err != nil {
		return Series{}, fmt.Errorf("Selection.series(): %v", err)
	}
	sel.s, _ = sel.s.in(sel.rowPositions)
	sel.s.index, _ = sel.s.index.In(sel.levelPositions)
	return sel.s, nil
}

// SwapLevels swaps the selected index levels. If Selection is not swappable, returns error.
func (sel Selection) SwapLevels() (Series, error) {
	if err := sel.ensure(); err != nil {
		return Series{}, fmt.Errorf("Selection.SwapLevels(): %v", err)
	}
	if !sel.swappableLevels {
		return sel.s, fmt.Errorf("selection is not swappable - must select exactly two levels")
	}
	lvl := sel.s.index.Levels
	lvl[sel.levelPositions[0]], lvl[sel.levelPositions[1]] = lvl[sel.levelPositions[1]], lvl[sel.levelPositions[0]]
	sel.s.index.Refresh()
	return sel.s, nil
}

// SwapRows swaps the selected rows. If Selection is not swappable, returns error.
func (sel Selection) SwapRows() (Series, error) {
	if err := sel.ensure(); err != nil {
		return Series{}, fmt.Errorf("Selection.SwapLevels(): %v", err)
	}
	if !sel.swappableRows {
		return sel.s, fmt.Errorf("selection is not swappable - must select exactly two rows")
	}
	s := sel.s
	r1 := sel.rowPositions[0]
	r2 := sel.rowPositions[1]

	for i := 0; i < len(s.index.Levels); i++ {
		r1v := s.index.Levels[i].Labels.Element(r1).Value
		r2v := s.index.Levels[i].Labels.Element(r2).Value
		s.index.Levels[i].Labels.Set(r1, r2v)
		s.index.Levels[i].Labels.Set(r2, r1v)
		s.index.Levels[i].Refresh()
	}
	r1v := s.values.Element(r1).Value
	r2v := s.values.Element(r2).Value
	s.values.Set(r1, r2v)
	s.values.Set(r2, r1v)
	return s, nil
}

// A Selection is a portion of a Series for use as an intermediate step in transforming data,
// such as getting, setting, swapping, or dropping.
type Selection struct {
	s               Series
	levelPositions  []int
	rowPositions    []int
	groups          []string
	category        derivedIntent
	swappableRows   bool
	swappableLevels bool
	err             error
}

type derivedIntent string

const (
	catAll        derivedIntent = "Select All"
	catLevelsOnly derivedIntent = "Select Levels"
	catRowsOnly   derivedIntent = "Select Rows"
	catXS         derivedIntent = "Select Cross-Section"
)

func (sel Selection) String() string {
	return fmt.Sprintf(
		"Selection Object. To print underlying Series, call .Get()\nDerivedIntent: %v\nRows: %v\nLevels: %v\nError: %v",
		sel.category, sel.rowPositions, sel.levelPositions, sel.err != nil)
}

// // Unpack the supplied options and try to categorize the caller's intention.
// func (s Series) unpack(cfg config.SelectionConfig) Selection {
// 	var sel = Selection{s: s.copy()}
// 	// [START check input for errors]
// 	noSelection := (cfg.LevelPositions == nil && cfg.LevelNames == nil && cfg.RowPositions == nil && cfg.RowLabels == nil)
// 	multipleLevelIdentifiers := (cfg.LevelPositions != nil && cfg.LevelNames != nil)
// 	multipleRowIdentifiers := (cfg.RowPositions != nil && cfg.RowLabels != nil)
// 	levelsOnly := (!noSelection && cfg.RowPositions == nil && cfg.RowLabels == nil)
// 	rowsOnly := (!noSelection && cfg.LevelPositions == nil && cfg.LevelNames == nil)
// 	levelsAndLabels := (!rowsOnly && cfg.RowLabels != nil)
// 	crossSection := (!noSelection && !levelsOnly && !rowsOnly)

// 	if noSelection {
// 		// return all row positions
// 		sel.rowPositions = values.MakeIntRange(0, s.Len())
// 		sel.category = catAll
// 		return sel
// 	}

// 	if multipleLevelIdentifiers {
// 		err := errors.New("the combination of integer positions and names is ambiguous. Provide at most one form of selecting index levels")
// 		values.Warn(
// 			fmt.Errorf("Cannot process level Selection: %v", err),
// 			"invalid Selection (will return error if called)")
// 		sel.err = err
// 		return sel
// 	}

// 	if multipleRowIdentifiers {
// 		err := errors.New("the combination of integer positions and labels is ambiguous. Provide at most one form of selecting rows")
// 		values.Warn(
// 			fmt.Errorf("Cannot process row Selection: %v", err),
// 			"invalid Selection (will return error if called)")
// 		sel.err = err
// 		return sel
// 	}
// 	// [END check input for errors]

// 	// Handling ByLevels
// 	if cfg.LevelPositions != nil {
// 		err := s.ensureLevelPositions(cfg.LevelPositions)
// 		if err != nil {
// 			values.Warn(
// 				fmt.Errorf("Cannot process level Selection: %v", err),
// 				"invalid Selection (will return error if called)")
// 			sel.err = err
// 			return sel
// 		}
// 		sel.levelPositions = cfg.LevelPositions
// 	} else {
// 		// Handling ByLevelNames
// 		for _, name := range cfg.LevelNames {
// 			val, ok := s.index.NameMap[name]
// 			if !ok {
// 				err := fmt.Errorf("level name %v not in index", name)
// 				values.Warn(
// 					fmt.Errorf("Cannot process level Selection: %v", err),
// 					"invalid Selection (will return error if called)")
// 				sel.err = err
// 				return sel
// 			}
// 			sel.levelPositions = append(sel.levelPositions, val...)
// 		}
// 	}
// 	// Handling ByRows
// 	if cfg.RowPositions != nil {
// 		err := s.ensureRowPositions(cfg.RowPositions)
// 		if err != nil {
// 			values.Warn(
// 				fmt.Errorf("Cannot process level Selection: %v", err),
// 				"invalid Selection (will return error if called)")
// 			sel.err = err
// 			return sel
// 		}
// 		sel.rowPositions = cfg.RowPositions
// 	} else {
// 		// Handling ByLabels
// 		var lvl int
// 		// no index level provided; defaults to first level
// 		if rowsOnly {
// 			lvl = 0
// 		} else {
// 			// multiple levels and row labels
// 			if levelsAndLabels && len(cfg.LevelPositions) > 1 {
// 				err := errors.New("the combination of multiple levels with row labels is ambiguous. To index on multiple levels, provide row integer values instead with opt.ByRows()")
// 				values.Warn(
// 					fmt.Errorf("Cannot process level Selection: %v", err),
// 					"invalid Selection (will return error if called)")
// 				sel.err = err
// 				return sel

// 			}
// 			// a single index level provided
// 			lvl = sel.levelPositions[0]
// 		}
// 		for _, label := range cfg.RowLabels {
// 			val, ok := s.index.Levels[lvl].LabelMap[label]
// 			if !ok {
// 				err := fmt.Errorf("label value %v not in index level %v", label, lvl)
// 				values.Warn(
// 					fmt.Errorf("Cannot process level Selection: %v", err),
// 					"invalid Selection (will return error if called)")
// 				sel.err = err
// 				return sel
// 			}
// 			sel.rowPositions = append(sel.rowPositions, val...)
// 		}
// 	}

// 	// Infer category and swappable
// 	if levelsOnly {
// 		sel.category = catLevelsOnly
// 		if len(sel.levelPositions) == 2 {
// 			sel.swappable = true
// 		}
// 	}
// 	if rowsOnly {
// 		sel.category = catRowsOnly
// 		if len(sel.rowPositions) == 2 {
// 			sel.swappable = true
// 		}
// 	}
// 	if crossSection {
// 		sel.category = catXS
// 	}

// 	if !noSelection && !levelsOnly && !rowsOnly && !crossSection {
// 		sel.category = "Unknown"
// 		return sel
// 	}

// 	sel.rowPositions = sort.IntSlice(sel.rowPositions)
// 	sel.levelPositions = sort.IntSlice(sel.levelPositions)
// 	return sel
// }

// // Select a portion of a Series (index levels and/or rows), based on either integer or string-based inputs. Options:
// //
// // - Select index level(s): opt.ByIndexLevels([]int), opt.ByIndexNames([]string)
// //
// // - Select row(s): opt.ByRows([]int), opt.ByLabels([]string)
// //
// // If no options are passed, selects the entire Series. If multiple of the same type of option are passed, only the last one is used.
// //
// // The following option combinations are ambiguous:
// //
// // - Both ByIndexLevels() and ByIndexNames(): to select index level(s), use one or the other.
// //
// // - Both ByRows() and ByLabels(): to want to select row(s), use one or the other.
// //
// // - An index level selector with more than 1 item and ByLabels(): to select multiple index levels and multiple index rows, use ByRows().
// //
// // If the caller passes invalid options, a warning will be logged, and attempts to call Selection methods will return an error.
// // func (s Series) Select(options ...opt.SelectionOption) Selection {
// // 	// Setup
// // 	cfg := config.SelectionConfig{}
// // 	for _, option := range options {
// // 		option(&cfg)
// // 	}
// // 	sel := s.unpack(cfg)
// // 	return sel
// // }

// // Get returns the Series specified in this Selection.
// //
// // Always returns a new Series.
// // func (sel Selection) Get() (Series, error) {
// // 	if sel.err != nil {
// // 		return sel.s, sel.err
// // 	}
// // 	return sel.get()
// // }

// // // returns the rows or index levels specified in the Selection and
// // // ducks .In() errors because those are checked by the unpacker on calls to s.Select().
// // func (sel Selection) get() (Series, error) {
// // 	switch sel.category {
// // 	case "Select All":
// // 		return sel.s, nil
// // 	case "Select Levels":
// // 		sel.s.index, _ = sel.s.index.In(sel.levelPositions)
// // 		return sel.s, nil
// // 	case "Select Rows":
// // 		sel.s, _ = sel.s.in(sel.rowPositions)
// // 		return sel.s, nil
// // 	case "Select Cross-Section":
// // 		sel.s, _ = sel.s.in(sel.rowPositions)
// // 		sel.s.index, _ = sel.s.index.In(sel.levelPositions)
// // 		return sel.s, nil
// // 	default:
// // 		return sel.s, fmt.Errorf("unable to categorize intention of caller")
// // 	}
// // }

// // Swap swaps the selected rows or index levels. If Selection is not swappable, returns error.
// //
// // Always returns a new Series.
// // func (sel Selection) Swap() (Series, error) {
// // 	if !sel.swappable {
// // 		return sel.s, fmt.Errorf("selection is not swappable - must select exactly two of either rows or levels")
// // 	}
// // 	s := sel.s
// // 	switch sel.category {
// // 	case "Select Levels":
// // 		lvl := s.index.Levels
// // 		lvl[sel.levelPositions[0]], lvl[sel.levelPositions[1]] = lvl[sel.levelPositions[1]], lvl[sel.levelPositions[0]]
// // 		s.index.Refresh()
// // 		return s, nil
// // 	case "Select Rows":
// // 		r1 := sel.rowPositions[0]
// // 		r2 := sel.rowPositions[1]

// // 		for i := 0; i < len(s.index.Levels); i++ {
// // 			r1v := s.index.Levels[i].Labels.Element(r1).Value
// // 			r2v := s.index.Levels[i].Labels.Element(r2).Value
// // 			s.index.Levels[i].Labels.Set(r1, r2v)
// // 			s.index.Levels[i].Labels.Set(r2, r1v)
// // 			s.index.Levels[i].Refresh()
// // 		}
// // 		r1v := s.values.Element(r1).Value
// // 		r2v := s.values.Element(r2).Value
// // 		s.values.Set(r1, r2v)
// // 		s.values.Set(r2, r1v)
// // 		return s, nil
// // 	default:
// // 		return sel.s, fmt.Errorf("")
// // 	}
// // }

// // // Set sets all the value in a Selection to val
// // func (sel Selection) Set(val interface{}) (Series, error) {
// // 	switch sel.category {
// // 	case "Select Rows":
// // 		for _, row := range sel.rowPositions {
// // 			err := sel.s.values.Set(row, val)
// // 			if err != nil {
// // 				return Series{}, fmt.Errorf("Selection.Set(): %v", err)
// // 			}
// // 		}
// // 	}
// // 	return sel.s, nil
// // }

// // // [END Selection]
