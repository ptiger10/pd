package series

import (
	"fmt"
	"sort"

	"github.com/ptiger10/pd/internal/values"
)

// At subsets a Series by stringified index label at index level 0.
func (s Series) At(label string) Series {
	sNew, err := s.atLevelLabel(0, label)
	if err != nil {
		values.Warn(err, "Original Series")
		return s
	}
	sNew.index.Refresh()
	return sNew
}

func (s Series) atLevelLabel(level int, label string) (Series, error) {
	var err error
	positions, ok := s.index.Levels[level].LabelMap[label]
	if !ok {
		return Series{}, fmt.Errorf("Label not in index level %v: %v", level, label)
	}
	s.values, err = s.values.In(positions)
	if err != nil {
		return Series{}, fmt.Errorf("Unable to get Series values at positions %v: %v", positions, err)
	}
	for i, level := range s.index.Levels {
		s.index.Levels[i].Labels, err = level.Labels.In(positions)
		if err != nil {
			return Series{}, fmt.Errorf("Unable to get Series index labels at level %v and positions %v: %v", level, positions, err)
		}
	}
	return s, nil
}

// A Selection is a portion of a Series, and is typically used as an intermediate step in manipulating or analyzing data,
// such as getting, setting, or dropping.
type Selection struct {
	s              Series
	levelPositions []int
	rowPositions   []int
	category       string
	swappable      bool
}

// A SelectionOption is an optional parameter in a Series selector.
type SelectionOption func(*selectionConfig)
type selectionConfig struct {
	levelPositions []int
	levelNames     []string
	rowPositions   []int
	rowLabels      []string
}

// ByIndexLevels selects one or more index levels by their integer positions
func ByIndexLevels(positions []int) SelectionOption {
	return func(c *selectionConfig) {
		c.levelPositions = positions
	}
}

// ByIndexNames selects one or more index levels by their names
func ByIndexNames(names []string) SelectionOption {
	return func(c *selectionConfig) {
		c.levelNames = names
	}
}

// ByRows selects one or more rows by their integer positions
func ByRows(positions []int) SelectionOption {
	return func(c *selectionConfig) {
		c.rowPositions = positions
	}
}

// ByLabels selects one or more rows by their stringified index labels
func ByLabels(labels []string) SelectionOption {
	return func(c *selectionConfig) {
		c.rowLabels = labels
	}
}

// Unpack the supplied options and try to categorize the caller's intention.
func (config *selectionConfig) unpack(s Series) (Selection, error) {
	var sel = Selection{s: s}
	noSelection := (config.levelPositions == nil && config.levelNames == nil && config.rowPositions == nil && config.rowLabels == nil)
	multipleLevelIdentifiers := (config.levelPositions != nil && config.levelNames != nil)
	multipleRowIdentifiers := (config.rowPositions != nil && config.rowLabels != nil)
	levelsOnly := (!noSelection && config.rowPositions == nil && config.rowLabels == nil)
	rowsOnly := (!noSelection && config.levelPositions == nil && config.levelNames == nil)
	levelsAndLabels := (!rowsOnly && config.rowLabels != nil)
	crossSection := (!noSelection && !levelsOnly && !rowsOnly)

	if noSelection {
		// return all row positions
		sel.rowPositions = values.MakeIntRange(0, s.Len())
		sel.category = "all"
		return sel, nil
	}

	if multipleLevelIdentifiers {
		return Selection{}, fmt.Errorf("Cannot process Selection. The combination of integer positions and names is ambiguous. Provide at most one form of selecting index levels")
	}

	if multipleRowIdentifiers {
		return Selection{}, fmt.Errorf("Cannot process Selection. The combination of integer positions and labels is ambiguous. Provide at most one form of selecting rows")
	}

	if config.levelPositions != nil {
		err := s.ensureLevelPositions(config.levelPositions)
		if err != nil {
			return Selection{}, fmt.Errorf("Cannot process level Selection: %v", err)
		}
		sel.levelPositions = config.levelPositions
	} else {
		for _, name := range config.levelNames {
			val, ok := s.index.NameMap[name]
			if !ok {
				return Selection{}, fmt.Errorf("Cannot process Selection. Level name %v not in index", name)
			}
			sel.levelPositions = append(sel.levelPositions, val...)
		}
	}

	if config.rowPositions != nil {
		sel.rowPositions = config.rowPositions
		err := s.ensureRowPositions(config.rowPositions)
		if err != nil {
			return Selection{}, fmt.Errorf("Cannot process row Selection: %v", err)
		}
	} else {
		var lvl int
		// no index level provided; defaults to first level
		if rowsOnly {
			lvl = 0
		} else {
			// multiple levels and row labels
			if levelsAndLabels && len(config.levelPositions) > 1 {
				return Selection{}, fmt.Errorf("Cannot process Selection. The combination of multiple levels with row labels is ambiguous. Provide row integer values instead")
			}
			// a single index level provided
			lvl = sel.levelPositions[0]
		}
		for _, label := range config.rowLabels {
			val, ok := s.index.Levels[lvl].LabelMap[label]
			if !ok {
				return Selection{}, fmt.Errorf("Cannot process Selection. Label value %v not in index level %v", label, lvl)
			}
			sel.rowPositions = append(sel.rowPositions, val...)
		}
	}

	if levelsOnly {
		sel.category = "levelsOnly"
		if len(sel.levelPositions) == 2 {
			sel.swappable = true
		}
	}
	if rowsOnly {
		sel.category = "rowsOnly"
		if len(sel.rowPositions) == 2 {
			sel.swappable = true
		}
	}
	if crossSection {
		sel.category = "xs"
	}

	if !noSelection && !levelsOnly && !rowsOnly && !crossSection {
		return Selection{}, fmt.Errorf("Cannot process Selection. Unable to categorize intention of the caller")
	}

	sel.rowPositions = sort.IntSlice(sel.rowPositions)
	sel.levelPositions = sort.IntSlice(sel.levelPositions)
	return sel, nil
}

// Select a portion of a Series (index levels and/or rows), based on either integer or string-based inputs. Options:
//
// Select index level(s): ByIndexLevels([]int), ByIndexNames([]string)
//
// Select row(s): ByRows([]int), ByLabels([]string)
//
// If no options are passed, selects the entire Series. If multiple of the same type of option are passed, only the last one is used.
func (s Series) Select(options ...SelectionOption) (Selection, error) {
	// Setup
	config := selectionConfig{}
	for _, option := range options {
		option(&config)
	}
	sel, err := config.unpack(s)
	if err != nil {
		return Selection{}, err
	}
	return sel, nil
}

// Get returns the Series underpinning this Selection
func (sel Selection) Get() Series {
	return sel.get()
}

func (s Series) ensureRowPositions(positions []int) error {
	_, err := s.values.In(positions)
	if err != nil {
		return fmt.Errorf("Bad row position data supplied: %v", err)
	}
	return nil
}

func (s Series) ensureLevelPositions(positions []int) error {
	_, err := s.index.In(positions)
	if err != nil {
		return fmt.Errorf("Bad index level position data supplied: %v", err)
	}
	return nil
}

// Create a new Series based on the Selection. Ducks .In() errors because those are checked by the unpacker on calls to s.Select().
func (sel Selection) get() Series {
	s := sel.s.copy()
	var err error
	switch sel.category {
	case "all":
		return s
	case "levelsOnly":
		s.index, _ = s.index.In(sel.levelPositions)
	case "rowsOnly":

		if err != nil {

		}
		for i, level := range s.index.Levels {
			s.index.Levels[i].Labels, _ = level.Labels.In(sel.rowPositions)
		}
	case "xs":
		s.index, _ = s.index.In(sel.levelPositions)
		s.values, _ = s.values.In(sel.rowPositions)
		for i := 0; i < len(sel.levelPositions); i++ {
			s.index.Levels[i].Labels, _ = s.index.Levels[i].Labels.In(sel.rowPositions)
		}
	default:
		values.Warn(fmt.Errorf("Unable to categorize intention of caller"), "Original Series")
		return s
	}
	return s
}
