package dataframe

import (
	"fmt"
	"log"
	"sort"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/options"
)

// [START Index modifications]

// Values returns an []interface{} of the values at each level of the index
func (idx Index) Values() [][]interface{} {
	var ret [][]interface{}
	for j := 0; j < idx.df.IndexLevels(); j++ {
		var vals []interface{}
		for i := 0; i < idx.df.Len(); i++ {
			vals = append(vals, idx.df.index.Levels[j].Labels.Element(i).Value)
		}
		ret = append(ret, vals)
	}
	return ret
}

// return unique index labels and their positions, for use in stack()
func (idx Index) unique(levels ...int) (labels []string, startPositions []int) {
	g := idx.df.GroupByIndex(levels...)
	labels = g.Groups()
	for _, label := range labels {
		startPositions = append(startPositions, g.groups[label].Positions[0])
	}
	return
}

// Sort sorts the index by index level 0 and returns a new index.
func (idx Index) Sort(asc bool) {
	if asc {
		sort.Stable(idx)
	} else {
		sort.Stable(sort.Reverse(idx))
	}
	return
}

// Swap swaps two labels at index level 0 and modifies the index in place. Required by Sort interface.
func (idx Index) Swap(i, j int) {
	idx.df.InPlace.SwapRows(i, j)
}

// Less compares two elements and returns true if the first is less than the second. Required by Sort interface.
func (idx Index) Less(i, j int) bool {
	return idx.df.index.Levels[0].Labels.Less(i, j)
}

// Len returns the number of items in each level of the index.
func (idx Index) Len() int {
	if len(idx.df.index.Levels) == 0 {
		return 0
	}
	return idx.df.index.Levels[0].Len()
}

// At returns the index values at a specified row position and index level but returns nil if either integer is out of range.
func (idx Index) At(row int, level int) interface{} {
	if err := idx.df.ensureIndexLevelPositions([]int{level}); err != nil {
		if options.GetLogWarnings() {
			log.Printf("df.Index.At(): %v", err)
		}
		return nil
	}
	if err := idx.df.ensureRowPositions([]int{row}); err != nil {
		if options.GetLogWarnings() {
			log.Printf("df.Index.At(): %v", err)
		}
		return nil
	}
	return idx.df.index.Levels[level].Labels.Element(row).Value
}

// RenameLevel renames an index level in place but does not change anything if level is out of range.
func (idx Index) RenameLevel(level int, name string) error {
	if err := idx.df.ensureIndexLevelPositions([]int{level}); err != nil {
		return fmt.Errorf("df.Index.RenameLevel(): %v", err)
	}
	idx.df.index.Levels[level].Name = name
	idx.df.index.Refresh()
	return nil
}

// Reindex converts an index level to a default range []int{0, 1, 2,...n}
func (idx Index) Reindex(level int) error {
	if err := idx.df.ensureIndexLevelPositions([]int{level}); err != nil {
		return fmt.Errorf("df.Index.Reindex(): %v", err)
	}
	// ducks error because inputs are controlled
	idx.df.index.Levels[level] = index.NewDefaultLevel(idx.Len(), idx.df.index.Levels[level].Name)
	return nil
}

// null returns the integer position of all null labels in this index level.
func (idx Index) null(level int) []int {
	var ret []int
	for i := 0; i < idx.Len(); i++ {
		if idx.df.index.Levels[level].Labels.Element(i).Null {
			ret = append(ret, i)
		}
	}
	return ret
}

// DropNull drops null index values at the index level specified and modifies the DataFrame in place.
func (idx Index) DropNull(level int) error {
	if level >= idx.df.IndexLevels() {
		return fmt.Errorf("df.Index.DropNull(): invalid index level: %d (max: %v)", level, idx.df.IndexLevels()-1)
	}
	idx.df.InPlace.dropMany(idx.null(level))
	return nil
}

// SwapLevels swaps two levels in the index and modifies the DataFrame in place.
func (idx Index) SwapLevels(i, j int) error {
	if err := idx.df.index.SwapLevels(i, j); err != nil {
		return fmt.Errorf("df.Index.SwapLevels(): %v", err)
	}
	return nil
}

// InsertLevel inserts a level into the index and modifies the DataFrame in place.
func (idx Index) InsertLevel(pos int, values interface{}, name string) error {
	err := idx.df.index.InsertLevel(pos, values, name)
	if err != nil {
		return fmt.Errorf("df.Index.InsertLevel(): %v", err)
	}
	return nil
}

// AppendLevel adds a new index level to the end of the current index  and modifies the DataFrame in place.
func (idx Index) AppendLevel(values interface{}, name string) error {
	err := idx.InsertLevel(idx.Len(), values, name)
	if err != nil {
		return fmt.Errorf("df.Index.AppendLevel(): %v", err)
	}
	return nil
}

// SubsetLevels modifies the DataFrame in place with only the specified index levels.
func (idx Index) SubsetLevels(levelPositions []int) error {
	err := idx.df.index.SubsetLevels(levelPositions)
	if err != nil {
		return fmt.Errorf("df.Index.InsertLevel(): %v", err)
	}
	return nil
}

// Set sets the label at the specified index row and level to val and modifies the DataFrame in place.
// First converts val to be the same type as the index level.
func (idx Index) Set(row int, level int, val interface{}) error {
	err := idx.df.index.Set(row, level, val)
	if err != nil {
		return fmt.Errorf("df.Index.Set(): %v", err)
	}
	return nil
}

// DropLevel drops the specified index level and modifies the DataFrame in place.
func (idx Index) DropLevel(level int) error {
	if err := idx.df.ensureIndexLevelPositions([]int{level}); err != nil {
		return fmt.Errorf("s.Index.DropLevels(): %v", err)
	}
	idx.df.index.DropLevel(level)
	return nil
}

// SelectName returns the integer position of the index level at the first occurence of the supplied name, or -1 if not a valid index level name.
func (idx Index) SelectName(name string) int {
	v, ok := idx.df.index.NameMap[name]
	if !ok {
		if options.GetLogWarnings() {
			log.Printf("s.Index.SelectName(): name not in index level names: %v\n", name)
		}
		return -1
	}
	return v[0]
}

// SelectNames returns the integer positions of the index levels with the supplied names.
func (idx Index) SelectNames(names []string) []int {
	include := make([]int, 0)
	empty := make([]int, 0)
	for _, name := range names {
		v, ok := idx.df.index.NameMap[name]
		if !ok {
			if options.GetLogWarnings() {
				log.Printf("s.Index.SelectNames(): name not in index level names: %v\n", name)
			}
			return empty
		}
		include = append(include, v...)
	}
	return include
}

// Flip replaces the DataFrame values at the supplied column with the labels at the supplied index level, and vice versa.
func (idx Index) Flip(col int, level int) (*DataFrame, error) {
	err := idx.df.ensureIndexLevelPositions([]int{level})
	if err != nil {
		return newEmptyDataFrame(), fmt.Errorf("s.Index.Flip(): %v", err)
	}
	err = idx.df.ensureColumnPositions([]int{col})
	if err != nil {
		return newEmptyDataFrame(), fmt.Errorf("s.Index.Flip(): %v", err)
	}
	df := idx.df.Copy()
	df.index.Levels[level].Labels, df.vals[col].Values = df.vals[col].Values, df.index.Levels[level].Labels
	df.index.Levels[level].Name, df.name = df.name, df.index.Levels[level].Name
	df.index.Levels[level].DataType, df.vals[col].DataType = df.vals[col].DataType, df.index.Levels[level].DataType
	df.index.Refresh()
	return df, nil
}

// Filter an index level using a callback function test.
// The Filter function iterates over all index values in interface{} form and applies the callback test to each.
// The return value is a slice of integer positions of all the rows passing the test.
// The caller is responsible for handling the type assertion on the interface, though this step is not necessary if the datatype is known with certainty.
// For example, here are two ways to write a filter that returns all rows with the suffix "boo":
//
// #1 (safer) error check type assertion
//
//  s.Index.Filter(0, func(val interface{}) bool {
//		v, ok := val.(string)
//		if !ok {
// 			return false
//		}
//		if strings.HasSuffix(v, "boo") {
// 			return true
// 		}
// 		return false
// 	})
//
// Input:
// bamboo    0
// leaves    1
// taboo     2
//
// Output:
// []int{0,2}
//
// #2 (riskier) no error check
//
//  s.Filter(func(val interface{}) bool {
//		if strings.HasSuffix(val.(string), "boo") {
// 			return true
// 		}
// 		return false
// 	})
func (idx Index) Filter(level int, cmp func(interface{}) bool) []int {
	include := make([]int, 0)
	empty := make([]int, 0)

	if err := idx.df.ensureIndexLevelPositions([]int{level}); err != nil {
		if options.GetLogWarnings() {
			log.Printf("series.Index.Filter(): %v", err)
		}
		return empty
	}
	vals := idx.Values()[level]
	for i, val := range vals {
		if cmp(val) {
			include = append(include, i)
		}
	}
	return include
}

// Convert converts an index level to datatype in place.
func (idx Index) Convert(dataType string, level int) error {
	if err := idx.df.ensureIndexLevelPositions([]int{level}); err != nil {
		return fmt.Errorf("df.Index.Convert(): %v", err)
	}
	err := idx.df.index.Levels[level].Convert(options.DT(dataType))
	if err != nil {
		return fmt.Errorf("df.Index.Convert(): %v", err)
	}
	return nil
}
