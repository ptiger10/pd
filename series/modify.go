package series

import (
	"fmt"
	"sort"

	"github.com/ptiger10/pd/options"
)

// Rename the Series.
func (s *Series) Rename(name string) {
	s.name = name
}

// [START InPlace]

// InPlace contains methods for modifying a Series in place.
type InPlace struct {
	s *Series
}

// Sort sorts the series by its values and modifies the Series in place.
func (ip InPlace) Sort(asc bool) {
	if asc {
		sort.Stable(ip)
	} else {
		sort.Stable(sort.Reverse(ip))
	}
}

// Len returns the length of the underlying Series (required by Sort interface)
func (ip InPlace) Len() int {
	return ip.s.Len()
}

// Swap swaps the selected rows in place.
func (ip InPlace) Swap(i, j int) {
	ip.s.values.Swap(i, j)
	for lvl := 0; lvl < ip.s.index.NumLevels(); lvl++ {
		ip.s.index.Levels[lvl].Labels.Swap(i, j)
		ip.s.index.Levels[lvl].Refresh()
	}
}

func (ip InPlace) Less(i, j int) bool {
	return ip.s.values.Less(i, j)
}

// Insert inserts a new row into the Series immediately before the specified integer position and modifies the Series in place.
// If the original Series is empty, replaces it with a new Series.
func (ip InPlace) Insert(pos int, val interface{}, idx []interface{}) error {
	// Handling empty Series
	if Equal(ip.s, newEmptySeries()) {
		newS, err := New(val, Config{MultiIndex: idx})
		if err != nil {
			return fmt.Errorf("Series.Insert(): inserting into empty Series requires creating a new Series: %v", err)
		}
		ip.s.replace(newS)
		return nil
	}

	if err := ip.s.ensureAlignment(); err != nil {
		return fmt.Errorf("Series.Insert(): %v", err)
	}

	if len(idx) != ip.s.index.NumLevels() {
		return fmt.Errorf("Series.Insert() len(idx) must equal number of index levels: supplied %v want %v",
			len(idx), ip.s.index.NumLevels())
	}
	for j := 0; j < ip.s.index.NumLevels(); j++ {
		err := ip.s.index.Levels[j].Labels.Insert(pos, idx[j])
		if err != nil {
			return fmt.Errorf("Series.Insert(): %v", err)
		}
		ip.s.index.Levels[j].Refresh()
	}
	// ducks error safely due to index alignment check
	ip.s.values.Insert(pos, val)
	return nil
}

// Append adds a row at a specified integer position and modifies the Series in place.
func (ip InPlace) Append(val interface{}, idx []interface{}) error {
	err := ip.s.InPlace.Insert(ip.s.Len(), val, idx)
	if err != nil {
		return fmt.Errorf("Series.Append(): %v", err)
	}
	return nil
}

// Set sets all the values in the specified rows to val and modifies the Series in place.
// If an error would be encountered in any row position, the entire operation is cancelled before it starts.
func (ip InPlace) Set(rowPositions []int, val interface{}) error {
	if err := ip.s.ensureRowPositions(rowPositions); err != nil {
		return fmt.Errorf("Series.Set(): %v", err)
	}
	// ducks error safely due to index alignment check
	for _, row := range rowPositions {
		ip.s.values.Set(row, val)
	}
	return nil
}

// Drop drops the row at the specified integer position and modifies the Series in place.
// If an error would be encountered in any row position, the entire operation is cancelled before it starts.
func (ip InPlace) Drop(rowPositions []int) error {
	if err := ip.dropMany(rowPositions); err != nil {
		return fmt.Errorf("Series.Drop(): %v", err)
	}
	return nil
}

// DropNull drops all null values and modifies the Series in place.
func (ip InPlace) DropNull() {
	ip.dropMany(ip.s.null())
	return
}

// dropMany drops multiple rows and modifies the Series in place.
func (ip InPlace) dropMany(positions []int) error {
	if err := ip.s.ensureRowPositions(positions); err != nil {
		return err
	}
	// ducks error safely due to index alignment check
	sort.IntSlice(positions).Sort()
	for i, position := range positions {
		ip.s.InPlace.dropOne(position - i)
	}
	if ip.Len() == 0 {
		ip.s.replace(newEmptySeries())
	}
	return nil
}

// dropOne drops a row at a specified integer position and modifies the Series in place.
func (ip InPlace) dropOne(pos int) error {
	for i := 0; i < ip.s.index.NumLevels(); i++ {
		// ducks errors safely due to index alignment check in dropMany
		ip.s.index.Levels[i].Labels.Drop(pos)
		ip.s.index.Levels[i].Refresh()
	}
	ip.s.values.Drop(pos)
	return nil
}

// ToFloat64 converts Series values to float64 in place.
func (ip InPlace) ToFloat64() {
	ip.s.values = ip.s.values.ToFloat64()
	ip.s.datatype = options.Float64
}

// ToInt64 converts Series values to int64 in place.
func (ip InPlace) ToInt64() {
	ip.s.values = ip.s.values.ToInt64()
	ip.s.datatype = options.Int64
}

// ToString converts Series values to string in place.
func (ip InPlace) ToString() {
	ip.s.values = ip.s.values.ToString()
	ip.s.datatype = options.String
}

// ToBool converts Series values to bool in place.
func (ip InPlace) ToBool() {
	ip.s.values = ip.s.values.ToBool()
	ip.s.datatype = options.Bool
}

// ToDateTime converts Series values to datetime in place.
func (ip InPlace) ToDateTime() {
	ip.s.values = ip.s.values.ToDateTime()
	ip.s.datatype = options.DateTime
}

// ToInterface converts Series values to interface in place.
func (ip InPlace) ToInterface() {
	ip.s.values = ip.s.values.ToInterface()
	ip.s.datatype = options.Interface
}

// [END InPlace]

// [START Copy]

// Sort sorts the series by its values and returns a new Series.
func (s *Series) Sort(asc bool) *Series {
	s = s.Copy()
	s.InPlace.Sort(asc)
	return s
}

// Swap swaps the selected rows and returns a new Series.
func (s *Series) Swap(i, j int) (*Series, error) {
	s = s.Copy()
	if i >= s.Len() {
		return newEmptySeries(), fmt.Errorf("invalid position: %d (max %v)", i, s.Len()-1)
	}
	if j >= s.Len() {
		return newEmptySeries(), fmt.Errorf("invalid position: %d (max %v)", j, s.Len()-1)
	}
	s.InPlace.Swap(i, j)
	return s, nil
}

// Insert inserts a new row into the Series immediately before the specified integer position and returns a new Series.
func (s *Series) Insert(pos int, val interface{}, idx []interface{}) (*Series, error) {
	s = s.Copy()
	err := s.InPlace.Insert(pos, val, idx)
	return s, err
}

// Append adds a row at the end and returns a new Series.
func (s *Series) Append(val interface{}, idx []interface{}) (*Series, error) {
	s, err := s.Insert(s.Len(), val, idx)
	if err != nil {
		return newEmptySeries(), fmt.Errorf("Series.Append(): %v", err)
	}
	return s, nil
}

// Set sets all the values in the specified rows to val and returns a new Series.
func (s *Series) Set(rowPositions []int, val interface{}) (*Series, error) {
	s = s.Copy()
	for _, row := range rowPositions {
		err := s.values.Set(row, val)
		if err != nil {
			return newEmptySeries(), fmt.Errorf("s.Set() for val %v: %v", val, err)
		}
	}
	return s, nil
}

// Drop drops the row at the specified integer position and returns a new Series.
func (s *Series) Drop(positions []int) (*Series, error) {
	s = s.Copy()
	err := s.InPlace.Drop(positions)
	return s, err
}

// DropNull drops all null values and modifies the Series in place.
func (s *Series) DropNull() *Series {
	s = s.Copy()
	s.InPlace.DropNull()
	return s
}

// ToFloat64 converts Series values to float64 and returns a new Series.
func (s *Series) ToFloat64() *Series {
	s = s.Copy()
	s.InPlace.ToFloat64()
	return s
}

// ToInt64 converts Series values to int64 and returns a new Series.
func (s *Series) ToInt64() *Series {
	s = s.Copy()
	s.InPlace.ToInt64()
	return s
}

// ToString converts Series values to string and returns a new Series.
func (s *Series) ToString() *Series {
	s = s.Copy()
	s.InPlace.ToString()
	return s
}

// ToBool converts Series values to bool and returns a new Series.
func (s *Series) ToBool() *Series {
	s = s.Copy()
	s.InPlace.ToBool()
	return s
}

// ToDateTime converts Series values to time.Time and returns a new Series.
func (s *Series) ToDateTime() *Series {
	s = s.Copy()
	s.InPlace.ToDateTime()
	return s
}

// ToInterface converts Series values to interface and returns a new Series.
func (s *Series) ToInterface() *Series {
	s = s.Copy()
	s.InPlace.ToInterface()
	return s
}

// [END Copy]
