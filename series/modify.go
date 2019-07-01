package series

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/ptiger10/pd/internal/values"

	"github.com/ptiger10/pd/options"
)

// Rename the Series.
func (s *Series) Rename(name string) {
	s.name = name
}

// replace one Series with another in place.
func (s *Series) replace(s2 *Series) {
	s.name = s2.name
	s.datatype = s2.datatype
	s.values = s2.values
	s.index = s2.index
}

// [START InPlace]

// InPlace contains methods for modifying a Series in place.
type InPlace struct {
	s *Series
}

func (ip InPlace) String() string {
	printer := "{InPlace Handler}\n"
	printer += "Methods:\n"
	t := reflect.TypeOf(InPlace{})
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		printer += fmt.Sprintln(method.Name)
	}
	return printer
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

	// Handling errors
	if err := ip.s.ensureAlignment(); err != nil {
		return fmt.Errorf("Series.Insert(): %v", err)
	}
	if len(idx) != ip.s.index.NumLevels() {
		return fmt.Errorf("Series.Insert() len(idx) must equal number of index levels: supplied %v want %v",
			len(idx), ip.s.index.NumLevels())
	}

	if pos > ip.Len() {
		return fmt.Errorf("Series.Insert(): invalid position: %d (max %v)", pos, ip.Len())
	}

	for _, v := range idx {
		if _, err := values.InterfaceFactory(v); err != nil {
			return fmt.Errorf("Series.Insert(): %v", err)
		}
	}
	if _, err := values.InterfaceFactory(val); err != nil {
		return fmt.Errorf("Series.Insert(): %v", err)
	}

	// Insertion once errors have been handled
	for j := 0; j < ip.s.index.NumLevels(); j++ {
		ip.s.index.Levels[j].Labels.Insert(pos, idx[j])
		ip.s.index.Levels[j].Refresh()
	}
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

// Set sets the values in the specified row to val and modifies the Series in place. First converts val to be the same type as the index level.
func (ip InPlace) Set(row int, val interface{}) error {
	if err := ip.s.ensureRowPositions([]int{row}); err != nil {
		return fmt.Errorf("Series.Set(): %v", err)
	}

	if _, err := values.InterfaceFactory(val); err != nil {
		return fmt.Errorf("Series.Set(): %v", err)
	}
	ip.s.values.Set(row, val)
	return nil
}

// SetRows sets all the values in the specified rows to val and modifies the Series in place. First converts val to be the same type as the index level.
// If an error would be encountered in any row position, the entire operation is cancelled before it starts.
func (ip InPlace) SetRows(rowPositions []int, val interface{}) error {
	if err := ip.s.ensureRowPositions(rowPositions); err != nil {
		return fmt.Errorf("Series.SetRows(): %v", err)
	}
	if _, err := values.InterfaceFactory(val); err != nil {
		return fmt.Errorf("Series.SetRows(): %v", err)
	}

	for _, row := range rowPositions {
		ip.s.values.Set(row, val)
	}
	return nil
}

// Drop drops the row at the specified integer position and modifies the Series in place.
func (ip InPlace) Drop(row int) error {
	if err := ip.dropMany([]int{row}); err != nil {
		return fmt.Errorf("Series.Drop(): %v", err)
	}
	return nil
}

// DropRows drops the rows at the specified integer position and modifies the Series in place.
// If an error would be encountered in any row position, the entire operation is cancelled before it starts.
func (ip InPlace) DropRows(rowPositions []int) error {
	if err := ip.dropMany(rowPositions); err != nil {
		return fmt.Errorf("Series.DropRows(): %v", err)
	}
	return nil
}

// DropDuplicates drops any rows containing duplicate index + Series values and modifies the Series in place.
func (ip InPlace) DropDuplicates() {
	g := ip.s.GroupByIndex()
	var toDrop []int
	for _, group := range g.Groups() {
		// only inspect groups with at least one position
		if positions := g.groups[group].Positions; len(positions) > 0 {
			exists := make(map[interface{}]bool)
			for _, pos := range positions {
				if exists[ip.s.At(pos)] {
					toDrop = append(toDrop, pos)
				} else {
					exists[ip.s.At(pos)] = true
				}
			}
		}
	}
	// ducks error because position is assumed to be in Series
	if len(toDrop) != 0 {
		ip.DropRows(toDrop)
	}
}

// DropDuplicates drops any rows containing duplicate index + Series values and returns a new Series.
func (s *Series) DropDuplicates() *Series {
	s = s.Copy()
	s.InPlace.DropDuplicates()
	return s
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
// Should be called via dropMany to catch errors.
func (ip InPlace) dropOne(pos int) {
	for i := 0; i < ip.s.index.NumLevels(); i++ {
		ip.s.index.Levels[i].Labels.Drop(pos)
		ip.s.index.Levels[i].Refresh()
	}
	ip.s.values.Drop(pos)
	return
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

// Set sets the value in the specified rows to val and returns a new Series.
func (s *Series) Set(row int, val interface{}) (*Series, error) {
	s = s.Copy()
	err := s.InPlace.Set(row, val)
	return s, err
}

// SetRows sets all the values in the specified rows to val and returns a new Series.
func (s *Series) SetRows(rowPositions []int, val interface{}) (*Series, error) {
	s = s.Copy()
	err := s.InPlace.SetRows(rowPositions, val)
	return s, err
}

// Drop drops the row at the specified integer position and returns a new Series.
func (s *Series) Drop(row int) (*Series, error) {
	s = s.Copy()
	err := s.InPlace.Drop(row)
	return s, err
}

// DropRows drops the rows at the specified integer position and returns a new Series.
func (s *Series) DropRows(positions []int) (*Series, error) {
	s = s.Copy()
	err := s.InPlace.DropRows(positions)
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
