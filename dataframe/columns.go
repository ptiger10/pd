package dataframe

import (
	"fmt"
	"log"

	"github.com/ptiger10/pd/internal/index"

	"github.com/ptiger10/pd/options"
)

// Values returns an []string of the values at each level of the cols.
func (col Columns) Values() [][]string {
	var ret [][]string
	for j := 0; j < col.df.ColLevels(); j++ {
		ret = append(ret, col.df.cols.Levels[j].Labels)
	}
	return ret
}

// SwapLevels swaps two column levels and modifies the cols in place.
func (col Columns) SwapLevels(i, j int) error {
	if err := col.df.ensureColumnLevelPositions([]int{i, j}); err != nil {
		return fmt.Errorf("Columns.SwapLevels(): %v", err)
	}
	col.df.cols.Levels[i], col.df.cols.Levels[j] = col.df.cols.Levels[j], col.df.cols.Levels[i]
	col.df.cols.Refresh()
	return nil
}

// At returns the cols values at a specified col level and column position but returns nil if either integer is out of range.
func (col Columns) At(level int, column int) string {
	if err := col.df.ensureColumnLevelPositions([]int{level}); err != nil {
		if options.GetLogWarnings() {
			log.Printf("Columns.At(): %v", err)
		}
		return ""
	}
	if err := col.df.ensureColumnPositions([]int{column}); err != nil {
		if options.GetLogWarnings() {
			log.Printf("Columns.At(): %v", err)
		}
		return ""
	}
	return col.df.cols.Levels[level].Labels[column]
}

// RenameLevel renames an cols level in place but does not change anything if level is out of range.
func (col Columns) RenameLevel(level int, name string) error {
	if err := col.df.ensureColumnLevelPositions([]int{level}); err != nil {
		return fmt.Errorf("df.cols.RenameLevel(): %v", err)
	}
	col.df.cols.Levels[level].Name = name
	col.df.cols.Refresh()
	return nil
}

// InsertLevel inserts a level into the cols and modifies the DataFrame in place.
func (col Columns) InsertLevel(pos int, labels []string, name string) error {
	if pos > col.df.ColLevels() {
		return fmt.Errorf("invalid column level: %d (max: %v)", pos, col.df.ColLevels()-1)
	}
	lvl := index.NewColLevel(labels, name)
	if len(labels) != col.df.NumCols() {
		return fmt.Errorf("df.cols.InsertLevel(): len(labels) must equal number of columns (%d != %d)",
			len(labels), col.df.NumCols())
	}
	col.df.cols.Levels = append(col.df.cols.Levels[:pos], append([]index.ColLevel{lvl}, col.df.cols.Levels[pos:]...)...)
	col.df.cols.Refresh()
	return nil
}

// AppendLevel adds a new cols level to the end of the current cols  and modifies the DataFrame in place.
func (col Columns) AppendLevel(labels []string, name string) error {
	err := col.InsertLevel(col.df.ColLevels(), labels, name)
	if err != nil {
		return fmt.Errorf("df.cols.AppendLevel(): %v", err)
	}
	return nil
}

// SubsetLevels modifies the DataFrame in place with only the specified cols levels.
func (col Columns) SubsetLevels(levelPositions []int) error {

	err := col.df.ensureColumnLevelPositions(levelPositions)
	if err != nil {
		return fmt.Errorf("df.cols.SubsetLevels(): %v", err)
	}
	if len(levelPositions) == 0 {
		return fmt.Errorf("df.cols.SubsetLevels(): no levels provided")
	}

	levels := make([]index.ColLevel, 0)
	for _, position := range levelPositions {
		levels = append(levels, col.df.cols.Levels[position])
	}
	col.df.cols.Levels = levels
	col.df.cols.Refresh()
	return nil
}

// DropLevel drops the specified cols level and modifies the DataFrame in place.
func (col Columns) DropLevel(level int) error {
	if err := col.df.ensureColumnLevelPositions([]int{level}); err != nil {
		return fmt.Errorf("Columns.DropLevel(): %v", err)
	}
	col.df.cols.Levels = append(col.df.cols.Levels[:level], col.df.cols.Levels[level+1:]...)
	return nil
}

// SelectName returns the integer position of the cols level at the first occurence of the supplied name, or -1 if not a valid cols level name.
func (col Columns) SelectName(name string) int {
	v, ok := col.df.cols.NameMap[name]
	if !ok {
		if options.GetLogWarnings() {
			log.Printf("Columns.SelectName(): name not in cols level names: %v\n", name)
		}
		return -1
	}
	return v[0]
}

// SelectNames returns the integer positions of the cols levels with the supplied names.
func (col Columns) SelectNames(names []string) []int {
	include := make([]int, 0)
	empty := make([]int, 0)
	for _, name := range names {
		v, ok := col.df.cols.NameMap[name]
		if !ok {
			if options.GetLogWarnings() {
				log.Printf("Columns.SelectNames(): name not in cols level names: %v\n", name)
			}
			return empty
		}
		include = append(include, v...)
	}
	return include
}
