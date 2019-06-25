package series

import (
	"fmt"
	"log"
	"sort"
)

// [START Index modifications]

// Index contains index selection and conversion
type Index struct {
	s *Series
}

// Levels returns the number of levels in the index
func (idx Index) Levels() int {
	return idx.s.index.NumLevels()
}

// Len returns the number of items in each level of the index.
func (idx Index) Len() int {
	if len(idx.s.index.Levels) == 0 {
		return 0
	}
	return idx.s.index.Levels[0].Len()
}

// Swap swaps two labels at index level 0 and modifies the index in place.
func (idx Index) Swap(i, j int) {
	idx.s.Swap(i, j)
}

// Less compares two elements and returns true if the first is less than the second.
func (idx Index) Less(i, j int) bool {
	return idx.s.index.Levels[0].Labels.Less(i, j)
}

// Sort sorts the index by index level 0 and modifies the Series in place.
func (idx Index) Sort(asc bool) {
	if asc {
		sort.Stable(idx)
	} else {
		sort.Stable(sort.Reverse(idx))
	}
}

// At returns the index value at a specified index level and integer position.
func (idx Index) At(position int, level int) (interface{}, error) {
	if level >= idx.s.index.NumLevels() {
		return nil, fmt.Errorf("invalid index level: %d (len: %v)", level, idx.s.index.NumLevels())
	}
	if position >= idx.s.Len() {
		return nil, fmt.Errorf("invalid position: %d (len: %v)", position, idx.s.Len())
	}
	elem := idx.s.Element(position)
	return elem.Labels[level], nil
}

func (s *Series) rename(name string) {
	s = s.Copy()
	s.index.Levels[0].Name = name
}

// LevelToFloat64 converts the labels at a specified index level to float64 and returns a new Series.
func (idx Index) LevelToFloat64(level int) (*Series, error) {
	if level > idx.s.index.NumLevels() {
		return nil, fmt.Errorf("invalid index level: %d (len: %v)", level, idx.s.index.NumLevels())
	}
	s := idx.s.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToFloat64()
	return s, nil
}

// ToFloat64 converts the labels at index level 0 to float64 and returns a new Series.
func (idx Index) ToFloat64() *Series {
	if idx.s.index.NumLevels() == 0 {
		log.Println("Cannot convert empty index")
		return nil
	}
	s, _ := idx.LevelToFloat64(0)
	return s
}

// LevelToInt64 converts the labels at a specified index level to int64 and returns a new Series.
func (idx Index) LevelToInt64(level int) (*Series, error) {
	if level > idx.s.index.NumLevels() {
		return nil, fmt.Errorf("invalid index level: %d (len: %v)", level, idx.s.index.NumLevels())
	}
	s := idx.s.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToInt64()
	return s, nil
}

// ToInt64 converts the labels at index level 0 to int64 and returns a new Series.
func (idx Index) ToInt64() *Series {
	if idx.s.index.NumLevels() == 0 {
		log.Println("Cannot convert empty index")
		return nil
	}
	s, _ := idx.LevelToInt64(0)
	return s
}

// LevelToString converts the labels at a specified index level to string and returns a new Series.
func (idx Index) LevelToString(level int) (*Series, error) {
	if level > idx.s.index.NumLevels() {
		return nil, fmt.Errorf("invalid index level: %d (len: %v)", level, idx.s.index.NumLevels())
	}
	s := idx.s.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToString()
	return s, nil
}

// ToString converts the labels at index level 0 to string and returns a new Series.
func (idx Index) ToString() *Series {
	if idx.s.index.NumLevels() == 0 {
		log.Println("Cannot convert empty index")
		return nil
	}
	s, _ := idx.LevelToString(0)
	return s
}

// LevelToBool converts the labels at a specified index level to bool and returns a new Series.
func (idx Index) LevelToBool(level int) (*Series, error) {
	if level > idx.s.index.NumLevels() {
		return nil, fmt.Errorf("invalid index level: %d (len: %v)", level, idx.s.index.NumLevels())
	}
	s := idx.s.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToBool()
	return s, nil
}

// ToBool converts the labels at index level 0 to bool and returns a new Series.
func (idx Index) ToBool() *Series {
	if idx.s.index.NumLevels() == 0 {
		log.Println("Cannot convert empty index")
		return nil
	}
	s, _ := idx.LevelToBool(0)
	return s
}

// LevelToDateTime converts the labels at a specified index level to DateTime and returns a new Series.
func (idx Index) LevelToDateTime(level int) (*Series, error) {
	if level > idx.s.index.NumLevels() {
		return nil, fmt.Errorf("invalid index level: %d (len: %v)", level, idx.s.index.NumLevels())
	}
	s := idx.s.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToDateTime()
	return s, nil
}

// ToDateTime converts the labels at index level 0 to DateTime and returns a new Series.
func (idx Index) ToDateTime() *Series {
	if idx.s.index.NumLevels() == 0 {
		log.Println("Cannot convert empty index")
		return nil
	}
	s, _ := idx.LevelToDateTime(0)
	return s
}

// LevelToInterface converts the labels at a specified index level to interface and returns a new Series.
func (idx Index) LevelToInterface(level int) (*Series, error) {
	if level > idx.s.index.NumLevels() {
		return nil, fmt.Errorf("invalid index level: %d (len: %v)", level, idx.s.index.NumLevels())
	}
	s := idx.s.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToInterface()
	return s, nil
}

// ToInterface converts the labels at index level 0 to interface and returns a new Series.
func (idx Index) ToInterface() *Series {
	if idx.s.index.NumLevels() == 0 {
		log.Println("Cannot convert empty index")
		return nil
	}
	s, _ := idx.LevelToInterface(0)
	return s
}
