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
	for j := 0; j < idx.NumLevels(); j++ {
		var vals []interface{}
		for i := 0; i < idx.df.Len(); i++ {
			vals = append(vals, idx.df.index.Levels[j].Labels.Element(i).Value)
		}
		ret = append(ret, vals)
	}
	return ret
}

// Sort sorts the index by index level 0 and returns a new index.
func (idx Index) Sort(asc bool) *DataFrame {
	idx = idx.df.Copy().Index
	if asc {
		sort.Stable(idx)
	} else {
		sort.Stable(sort.Reverse(idx))
	}
	return idx.df
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

// NumLevels returns the number of levels in the index
func (idx Index) NumLevels() int {
	return idx.df.index.NumLevels()
}

// At returns the index values at a specified row position and index level but returns nil if either integer is out of range.
func (idx Index) At(row int, level int) interface{} {
	if level >= idx.NumLevels() {
		if options.GetLogWarnings() {
			log.Printf("df.Index.At(): invalid index level: %d (max: %v)", level, idx.NumLevels()-1)
		}
		return nil
	}
	if row >= idx.Len() {
		if options.GetLogWarnings() {
			log.Printf("df.Index.At(): invalid row: %d (max: %v)", row, idx.Len()-1)
		}
		return nil
	}
	elem := idx.df.Row(row)
	return elem.Labels[level]
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

// DropNull drops null index values at the index level specified and returns a new DataFrame.
func (idx Index) DropNull(level int) (*DataFrame, error) {
	if level >= idx.NumLevels() {
		return newEmptyDataFrame(), fmt.Errorf("df.Index.DropNull(): invalid index level: %d (max: %v)", level, idx.NumLevels()-1)
	}
	df := idx.df.Copy()
	df.InPlace.dropMany(idx.null(level))
	return df, nil
}

// SwapLevels swaps two levels in the index and returns a new DataFrame.
func (idx Index) SwapLevels(i, j int) (*DataFrame, error) {
	df := idx.df.Copy()
	if err := df.index.SwapLevels(i, j); err != nil {
		return newEmptyDataFrame(), fmt.Errorf("df.Index.SwapLevels(): invalid i: %v", err)
	}
	return df, nil
}

// InsertLevel inserts a level into the index and returns a new DataFrame.
func (idx Index) InsertLevel(pos int, values interface{}, name string) (*DataFrame, error) {
	df := idx.df.Copy()
	if err := df.index.InsertLevel(pos, values, name); err != nil {
		return newEmptyDataFrame(), fmt.Errorf("df.Index.InsertLevel(): invalid i: %v", err)
	}
	return df, nil
}

// AppendLevel adds a new index level to the end of the current index  and returns a new DataFrame.
func (idx Index) AppendLevel(values interface{}, name string) (*DataFrame, error) {
	df, err := idx.InsertLevel(idx.Len(), values, name)
	if err != nil {
		return newEmptyDataFrame(), fmt.Errorf("df.Index.AppendLevel(): %v", err)
	}
	return df, err
}

// SubsetLevels returns a new DataFrame with only the specified index levels.
func (idx Index) SubsetLevels(levelPositions []int) (*DataFrame, error) {
	err := idx.df.ensureIndexLevelPositions(levelPositions)
	if err != nil {
		return newEmptyDataFrame(), fmt.Errorf("s.Index.SubsetLevels(): %v", err)
	}
	df := idx.df.Copy()
	levels := make([]index.Level, 0)
	for _, position := range levelPositions {
		levels = append(levels, df.index.Levels[position])
	}
	df.index.Levels = levels
	df.index.Refresh()
	return df, nil
}

// Set sets the label at the specified index row and level to val and returns a new DataFrame.
// First converts val to be the same type as the index level.
func (idx Index) Set(row int, level int, val interface{}) (*DataFrame, error) {
	df := idx.df.Copy()
	err := df.index.Set(row, level, val)
	if err != nil {
		return newEmptyDataFrame(), fmt.Errorf("df.Index.Set(): %v", err)
	}
	return df, nil
}

// SetRows sets the label at the specified index rows and level to val and returns a new DataFrame.
// First converts val to be the same type as the index level.
func (idx Index) SetRows(rows []int, level int, val interface{}) (*DataFrame, error) {
	s := idx.df.Copy()
	err := s.index.SetRows(rows, level, val)
	if err != nil {
		return newEmptyDataFrame(), fmt.Errorf("s.Index.SetRows(): %v", err)
	}
	return s, nil
}

// DropLevel drops the specified index level and returns a new DataFrame.
func (idx Index) DropLevel(level int) (*DataFrame, error) {
	if err := idx.df.ensureIndexLevelPositions([]int{level}); err != nil {
		return newEmptyDataFrame(), fmt.Errorf("s.Index.DropLevels(): %v", err)
	}
	s := idx.df.Copy()
	s.index.DropLevel(level)
	return s, nil
}

// DropLevels drops the specified index levels and returns a new DataFrame.
func (idx Index) DropLevels(levelPositions []int) (*DataFrame, error) {
	if err := idx.df.ensureIndexLevelPositions(levelPositions); err != nil {
		return newEmptyDataFrame(), fmt.Errorf("s.Index.DropLevels(): %v", err)
	}
	s := idx.df.Copy()
	sort.IntSlice(levelPositions).Sort()
	for j, position := range levelPositions {
		s.index.DropLevel(position - j)
	}
	s.index.Refresh()
	return s, nil
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
// 0    bamboo
// 1    leaves
// 2    taboo
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

// LevelToFloat64 converts the labels at a specified index level to float64 and returns a new DataFrame.
func (idx Index) LevelToFloat64(level int) (*DataFrame, error) {
	if level >= idx.NumLevels() {
		return newEmptyDataFrame(), fmt.Errorf("invalid index level: %d (len: %v)", level, idx.NumLevels())
	}
	s := idx.df.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToFloat64()
	return s, nil
}

// LevelToInt64 converts the labels at a specified index level to int64 and returns a new DataFrame.
func (idx Index) LevelToInt64(level int) (*DataFrame, error) {
	if level > idx.NumLevels() {
		return newEmptyDataFrame(), fmt.Errorf("invalid index level: %d (len: %v)", level, idx.NumLevels())
	}
	s := idx.df.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToInt64()
	return s, nil
}

// LevelToString converts the labels at a specified index level to string and returns a new DataFrame.
func (idx Index) LevelToString(level int) (*DataFrame, error) {
	if level > idx.NumLevels() {
		return newEmptyDataFrame(), fmt.Errorf("invalid index level: %d (len: %v)", level, idx.NumLevels())
	}
	s := idx.df.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToString()
	return s, nil
}

// LevelToBool converts the labels at a specified index level to bool and returns a new DataFrame.
func (idx Index) LevelToBool(level int) (*DataFrame, error) {
	if level > idx.NumLevels() {
		return newEmptyDataFrame(), fmt.Errorf("invalid index level: %d (len: %v)", level, idx.NumLevels())
	}
	s := idx.df.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToBool()
	return s, nil
}

// LevelToDateTime converts the labels at a specified index level to DateTime and returns a new DataFrame.
func (idx Index) LevelToDateTime(level int) (*DataFrame, error) {
	if level > idx.NumLevels() {
		return newEmptyDataFrame(), fmt.Errorf("invalid index level: %d (len: %v)", level, idx.NumLevels())
	}
	s := idx.df.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToDateTime()
	return s, nil
}

// LevelToInterface converts the labels at a specified index level to interface and returns a new DataFrame.
func (idx Index) LevelToInterface(level int) (*DataFrame, error) {
	if level > idx.NumLevels() {
		return newEmptyDataFrame(), fmt.Errorf("invalid index level: %d (len: %v)", level, idx.NumLevels())
	}
	s := idx.df.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToInterface()
	return s, nil
}
