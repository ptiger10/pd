package series

import (
	"fmt"
	"log"
	"sort"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
)

// [START Index modifications]

// Index contains index selection and conversion
type Index struct {
	s *Series
}

// Sort sorts the index by index level 0 and returns a new index.
func (idx Index) Sort(asc bool) *Series {
	idx = idx.s.Copy().Index
	if asc {
		sort.Stable(idx)
	} else {
		sort.Stable(sort.Reverse(idx))
	}
	return idx.s
}

// Swap swaps two labels at index level 0 and modifies the index in place.
func (idx Index) Swap(i, j int) {
	idx.s.InPlace.Swap(i, j)
}

// Less compares two elements and returns true if the first is less than the second.
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

// NumLevels returns the number of levels in the index
func (idx Index) NumLevels() int {
	return idx.s.index.NumLevels()
}

// At returns the index value at a specified row position and index level but returns nil if either integer is out of range.
func (idx Index) At(row int, level int) interface{} {
	if level >= idx.NumLevels() {
		if options.GetLogWarnings() {
			log.Printf("s.Index.At(): invalid index level: %d (max: %v)", level, idx.NumLevels()-1)
		}
		return nil
	}
	if row >= idx.Len() {
		if options.GetLogWarnings() {
			log.Printf("s.Index.At(): invalid row: %d (max: %v)", row, idx.Len()-1)
		}
		return nil
	}
	elem := idx.s.Element(row)
	return elem.Labels[level]
}

// RenameLevel renames an index level in place but does not change anything if level is out of range.
func (idx Index) RenameLevel(level int, name string) error {
	if err := idx.ensureLevelPositions([]int{level}); err != nil {
		return fmt.Errorf("s.Index.RenameLevel(): %v", err)
	}
	idx.s.index.Levels[level].Name = name
	idx.s.index.Refresh()
	return nil
}

// Reindex converts an index level to a default range []int{0, 1, 2,...n}
func (idx Index) Reindex(level int) error {
	if err := idx.ensureLevelPositions([]int{level}); err != nil {
		return fmt.Errorf("s.Index.Reindex(): %v", err)
	}
	// ducks error because inputs are controlled
	idxVals := values.MakeIntRange(0, idx.Len())
	newLvl := index.MustNewLevel(idxVals, idx.s.index.Levels[level].Name)
	idx.s.index.Levels[level] = newLvl
	idx.s.index.Refresh()
	return nil
}

// Values returns an interface{}, ready for type assertion, of all values at the specified index level, but returns nil if level is out of range.
func (idx Index) Values(level int) interface{} {
	if err := idx.ensureLevelPositions([]int{level}); err != nil {
		if options.GetLogWarnings() {
			log.Printf("s.Index.Values(): %v\n", err)
		}
		return nil
	}
	return idx.s.index.Levels[level].Labels.Vals()
}

func (idx Index) String() string {
	return fmt.Sprintf("{Index | Len: %d, NumLevels: %d}", idx.Len(), idx.NumLevels())
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

// DropNull drops null index values at the index level specified.
func (idx Index) DropNull(level int) (*Series, error) {
	if level >= idx.NumLevels() {
		return newEmptySeries(), fmt.Errorf("s.Index.DropNull(): invalid index level: %d (max: %v)", level, idx.NumLevels()-1)
	}
	s := idx.s.Copy()
	s.InPlace.dropMany(idx.null(level))
	return s, nil
}

// SwapLevels swaps two levels in the index.
func (idx Index) SwapLevels(i, j int) (*Series, error) {
	if err := idx.ensureLevelPositions([]int{i}); err != nil {
		return newEmptySeries(), fmt.Errorf("s.Index.SwapLevels(): invalid i: %v", err)
	}
	if err := idx.ensureLevelPositions([]int{j}); err != nil {
		return newEmptySeries(), fmt.Errorf("s.Index.SwapLevels(): invalid j: %v", err)
	}
	s := idx.s.Copy()
	s.index.Levels[i], s.index.Levels[j] = s.index.Levels[j], s.index.Levels[i]
	s.index.Refresh()
	return s, nil
}

// InsertLevel inserts a level into the index.
func (idx Index) InsertLevel(pos int, values interface{}, name string) (*Series, error) {
	if pos > idx.NumLevels() {
		return newEmptySeries(), fmt.Errorf("s.Index.InsertLevel(): invalid index level: %d (max: %v)", pos, idx.NumLevels()-1)
	}
	lvl, err := index.NewLevel(values, name)
	if err != nil {
		return newEmptySeries(), fmt.Errorf("s.Index.InsertLevel(): %v", err)
	}
	if lvl.Len() != idx.Len() {
		return newEmptySeries(), fmt.Errorf("s.Index.InsertLevel(): values to insert must have same length as existing index: %d != %d", lvl.Len(), idx.Len())
	}
	s := idx.s.Copy()
	s.index.Levels = append(s.index.Levels[:pos], append([]index.Level{lvl}, s.index.Levels[pos:]...)...)
	s.index.Refresh()
	return s, nil
}

func (idx Index) ensureLevelPositions(levelPositions []int) error {
	if len(levelPositions) == 0 {
		return fmt.Errorf("no levels provided")
	}

	for _, pos := range levelPositions {
		if pos >= idx.NumLevels() {
			return fmt.Errorf("invalid index level: %d (max: %v)", pos, idx.NumLevels()-1)
		}
	}
	return nil
}

// SubsetLevels returns a Series with only the specified index levels.
func (idx Index) SubsetLevels(levelPositions []int) (*Series, error) {
	err := idx.ensureLevelPositions(levelPositions)
	if err != nil {
		return newEmptySeries(), fmt.Errorf("s.Index.SubsetLevels(): %v", err)
	}
	s := idx.s.Copy()
	levels := make([]index.Level, 0)
	for _, position := range levelPositions {
		levels = append(levels, s.index.Levels[position])
	}
	s.index.Levels = levels
	s.index.Refresh()
	return s, nil
}

// Set sets the label at the specified index row and level to val. First converts val to be the same type as the index level.
func (idx Index) Set(row int, level int, val interface{}) (*Series, error) {
	s := idx.s.Copy()
	err := s.index.Set(row, level, val)
	if err != nil {
		return newEmptySeries(), fmt.Errorf("s.Index.Set(): %v", err)
	}
	return s, nil
}

// SetRows sets the label at the specified index rows and level to val. First converts val to be the same type as the index level.
func (idx Index) SetRows(rows []int, level int, val interface{}) (*Series, error) {
	s := idx.s.Copy()
	err := s.index.SetRows(rows, level, val)
	if err != nil {
		return newEmptySeries(), fmt.Errorf("s.Index.SetRows(): %v", err)
	}
	return s, nil
}

// DropLevel drops the specified index level.
func (idx Index) DropLevel(level int) (*Series, error) {
	if err := idx.ensureLevelPositions([]int{level}); err != nil {
		return newEmptySeries(), fmt.Errorf("s.Index.DropLevels(): %v", err)
	}
	s := idx.s.Copy()
	s.index.DropLevel(level)
	return s, nil
}

// DropLevels drops the specified index levels.
func (idx Index) DropLevels(levelPositions []int) (*Series, error) {
	if err := idx.ensureLevelPositions(levelPositions); err != nil {
		return newEmptySeries(), fmt.Errorf("s.Index.DropLevels(): %v", err)
	}
	s := idx.s.Copy()
	sort.IntSlice(levelPositions).Sort()
	for j, position := range levelPositions {
		s.index.DropLevel(position - j)
	}
	s.index.Refresh()
	return s, nil
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
	err := idx.ensureLevelPositions([]int{level})
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

// LevelToFloat64 converts the labels at a specified index level to float64 and returns a new Series.
func (idx Index) LevelToFloat64(level int) (*Series, error) {
	if level >= idx.NumLevels() {
		return newEmptySeries(), fmt.Errorf("invalid index level: %d (len: %v)", level, idx.NumLevels())
	}
	s := idx.s.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToFloat64()
	return s, nil
}

// LevelToInt64 converts the labels at a specified index level to int64 and returns a new Series.
func (idx Index) LevelToInt64(level int) (*Series, error) {
	if level > idx.NumLevels() {
		return newEmptySeries(), fmt.Errorf("invalid index level: %d (len: %v)", level, idx.NumLevels())
	}
	s := idx.s.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToInt64()
	return s, nil
}

// LevelToString converts the labels at a specified index level to string and returns a new Series.
func (idx Index) LevelToString(level int) (*Series, error) {
	if level > idx.NumLevels() {
		return newEmptySeries(), fmt.Errorf("invalid index level: %d (len: %v)", level, idx.NumLevels())
	}
	s := idx.s.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToString()
	return s, nil
}

// LevelToBool converts the labels at a specified index level to bool and returns a new Series.
func (idx Index) LevelToBool(level int) (*Series, error) {
	if level > idx.NumLevels() {
		return newEmptySeries(), fmt.Errorf("invalid index level: %d (len: %v)", level, idx.NumLevels())
	}
	s := idx.s.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToBool()
	return s, nil
}

// LevelToDateTime converts the labels at a specified index level to DateTime and returns a new Series.
func (idx Index) LevelToDateTime(level int) (*Series, error) {
	if level > idx.NumLevels() {
		return newEmptySeries(), fmt.Errorf("invalid index level: %d (len: %v)", level, idx.NumLevels())
	}
	s := idx.s.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToDateTime()
	return s, nil
}

// LevelToInterface converts the labels at a specified index level to interface and returns a new Series.
func (idx Index) LevelToInterface(level int) (*Series, error) {
	if level > idx.NumLevels() {
		return newEmptySeries(), fmt.Errorf("invalid index level: %d (len: %v)", level, idx.NumLevels())
	}
	s := idx.s.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToInterface()
	return s, nil
}
