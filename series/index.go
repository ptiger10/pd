package series

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
	for j := 0; j < idx.s.NumLevels(); j++ {
		var vals []interface{}
		for i := 0; i < idx.s.Len(); i++ {
			vals = append(vals, idx.s.index.Levels[j].Labels.Element(i).Value)
		}
		ret = append(ret, vals)
	}
	return ret
}

// Sort sorts the index by index level 0 and modifies the index in place.
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
	idx.s.InPlace.Swap(i, j)
}

// Less compares two elements and returns true if the first is less than the second. Required by Sort interface.
func (idx Index) Less(i, j int) bool {
	return idx.s.index.Levels[0].Labels.Less(i, j)
}

// Len returns the number of items in each level of the index.
func (idx Index) Len() int {
	if len(idx.s.index.Levels) == 0 {
		return 0
	}
	return idx.s.index.Levels[0].Len()
}

// At returns the index value at a specified row position and index level but returns nil if either integer is out of range.
func (idx Index) At(row int, level int) interface{} {
	if err := idx.s.ensureLevelPositions([]int{level}); err != nil {
		if options.GetLogWarnings() {
			log.Printf("s.Index.At(): %v", err)
		}
		return nil
	}
	if err := idx.s.ensureRowPositions([]int{row}); err != nil {
		if options.GetLogWarnings() {
			log.Printf("s.Index.At(): %v", err)
		}
		return nil
	}
	return idx.s.index.Levels[level].Labels.Element(row).Value
}

// RenameLevel renames an index level in place but does not change anything if level is out of range.
func (idx Index) RenameLevel(level int, name string) error {
	if err := idx.s.ensureLevelPositions([]int{level}); err != nil {
		return fmt.Errorf("s.Index.RenameLevel(): %v", err)
	}
	idx.s.index.Levels[level].Name = name
	idx.s.index.Refresh()
	return nil
}

// Reindex converts an index level in place to an ordered default range []int{0, 1, 2,...n}
func (idx Index) Reindex(level int) error {
	if err := idx.s.ensureLevelPositions([]int{level}); err != nil {
		return fmt.Errorf("s.Index.Reindex(): %v", err)
	}
	// ducks error because inputs are controlled
	idx.s.index.Levels[level] = index.NewDefaultLevel(idx.Len(), idx.s.index.Levels[level].Name)
	return nil
}

// null returns the integer position of all null labels in this index level.
func (idx Index) null(level int) []int {
	var ret []int
	for i := 0; i < idx.Len(); i++ {
		if idx.s.index.Levels[level].Labels.Element(i).Null {
			ret = append(ret, i)
		}
	}
	return ret
}

// DropNull drops null index values at the index level specified and modifies the Series in place.
func (idx Index) DropNull(level int) error {
	if level >= idx.s.NumLevels() {
		return fmt.Errorf("s.Index.DropNull(): invalid index level: %d (max: %v)", level, idx.s.NumLevels()-1)
	}
	idx.s.InPlace.dropMany(idx.null(level))
	return nil
}

// SwapLevels swaps two levels in the index and modifies the Series in place.
func (idx Index) SwapLevels(i, j int) error {
	if err := idx.s.index.SwapLevels(i, j); err != nil {
		return fmt.Errorf("s.Index.SwapLevels(): invalid i: %v", err)
	}
	return nil
}

// InsertLevel inserts a level into the index and modifies the Series in place.
func (idx Index) InsertLevel(pos int, values interface{}, name string) error {
	if err := idx.s.index.InsertLevel(pos, values, name); err != nil {
		return fmt.Errorf("s.Index.InsertLevel(): invalid i: %v", err)
	}
	return nil
}

// AppendLevel adds a new index level to the end of the current index  and modifies the Series in place.
func (idx Index) AppendLevel(values interface{}, name string) error {
	err := idx.s.index.InsertLevel(idx.Len(), values, name)
	if err != nil {
		return fmt.Errorf("s.Index.AppendLevel(): %v", err)
	}
	return nil
}

// SubsetLevels modifies the Series in place with only the specified index levels.
func (idx Index) SubsetLevels(levelPositions []int) error {
	err := idx.s.ensureLevelPositions(levelPositions)
	if err != nil {
		return fmt.Errorf("s.Index.SubsetLevels(): %v", err)
	}
	if len(levelPositions) == 0 {
		return fmt.Errorf("s.Index.SubsetLevels(): no levels provided")
	}
	levels := make([]index.Level, 0)
	for _, position := range levelPositions {
		levels = append(levels, idx.s.index.Levels[position])
	}
	idx.s.index.Levels = levels
	idx.s.index.Refresh()
	return nil
}

// Set sets the label at the specified index row and level to val and modifies the Series in place.
// First converts val to be the same type as the index level.
func (idx Index) Set(row int, level int, val interface{}) error {
	err := idx.s.index.Set(row, level, val)
	if err != nil {
		return fmt.Errorf("s.Index.Set(): %v", err)
	}
	return nil
}

// DropLevel drops the specified index level and modifies the Series in place.
func (idx Index) DropLevel(level int) error {
	if err := idx.s.ensureLevelPositions([]int{level}); err != nil {
		return fmt.Errorf("s.Index.DropLevels(): %v", err)
	}
	idx.s.index.DropLevel(level)
	return nil
}

// SelectName returns the integer position of the index level at the first occurence of the supplied name, or -1 if not a valid index level name.
func (idx Index) SelectName(name string) int {
	v, ok := idx.s.index.NameMap[name]
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
		v, ok := idx.s.index.NameMap[name]
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

// Flip replaces the Series values with the labels at the supplied index level, and vice versa.
func (idx Index) Flip(level int) (*Series, error) {
	err := idx.s.ensureLevelPositions([]int{level})
	if err != nil {
		return newEmptySeries(), fmt.Errorf("s.Index.Flip(): %v", err)
	}
	s := idx.s.Copy()
	s.index.Levels[level].Labels, s.values = s.values, s.index.Levels[level].Labels
	s.index.Levels[level].Name, s.name = s.name, s.index.Levels[level].Name
	s.index.Levels[level].DataType, s.datatype = s.datatype, s.index.Levels[level].DataType
	s.index.Refresh()
	return s, nil
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

	if err := idx.s.ensureLevelPositions([]int{level}); err != nil {
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
	if err := idx.s.ensureLevelPositions([]int{level}); err != nil {
		return fmt.Errorf("df.Index.Convert(): %v", err)
	}
	err := idx.s.index.Levels[level].Convert(options.DT(dataType))
	if err != nil {
		return fmt.Errorf("df.Index.Convert(): %v", err)
	}
	return nil
}
