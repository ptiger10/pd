package dataframe

import (
	"fmt"
	"sort"

	"github.com/ptiger10/pd/internal/values"

	"github.com/ptiger10/pd/options"
)

// Rename the DataFrame.
func (df *DataFrame) Rename(name string) {
	df.name = name
}

// replace one DataFrame with another in place.
func (df *DataFrame) replace(df2 *DataFrame) {
	df.name = df2.name
	df.vals = df2.vals
	df.index = df2.index
	df.cols = df2.cols
}

// [START InPlace]

// // Sort sorts the series by its values and modifies the DataFrame in place.
// func (ip InPlace) Sort(asc bool) {
// 	if asc {
// 		sort.Stable(ip)
// 	} else {
// 		sort.Stable(sort.Reverse(ip))
// 	}
// }

// Len returns the length of the underlying DataFrame (required by Sort interface)
func (ip InPlace) Len() int {
	return ip.df.Len()
}

// SwapRows swaps the selected rows in place.
func (ip InPlace) SwapRows(i, j int) {
	for m := 0; m < ip.df.NumCols(); m++ {
		ip.df.vals[m].Values.Swap(i, j)
	}
	for l := 0; l < ip.df.IndexLevels(); l++ {
		ip.df.index.Levels[l].Labels.Swap(i, j)
		ip.df.index.Levels[l].Refresh()
	}
}

// SwapColumns swaps the selected columns in place.
func (ip InPlace) SwapColumns(i, j int) {
	ip.df.vals[i], ip.df.vals[j] = ip.df.vals[j], ip.df.vals[i]
	for l := 0; l < ip.df.ColLevels(); l++ {
		ip.df.cols.Levels[l].Labels[i], ip.df.cols.Levels[l].Labels[j] = ip.df.cols.Levels[l].Labels[j], ip.df.cols.Levels[l].Labels[i]
		ip.df.cols.Levels[l].Refresh()
	}
}

// Less returns true if the value at i > j in col.
func (ip InPlace) Less(col int, i, j int) bool {
	return ip.df.vals[col].Values.Less(i, j)
}

// InsertRow inserts a new row into the DataFrame immediately before the specified integer position and modifies the DataFrame in place.
// If the original DataFrame is empty, replaces it with a new DataFrame.
func (ip InPlace) InsertRow(pos int, val []interface{}, idx []interface{}) error {
	// Handling empty DataFrame
	if Equal(ip.df, newEmptyDataFrame()) {
		newDf, err := New(val, Config{MultiIndex: idx})
		if err != nil {
			return fmt.Errorf("DataFrame.InsertRow(): inserting into empty DataFrame requires creating a new DataFrame: %v", err)
		}
		ip.df.replace(newDf)
		return nil
	}

	// Handling errors
	if err := ip.df.ensureAlignment(); err != nil {
		return fmt.Errorf("DataFrame.InsertRow(): %v", err)
	}
	if len(idx) != ip.df.index.NumLevels() {
		return fmt.Errorf("DataFrame.InsertRow() len(idx) must equal number of index levels: supplied %v want %v",
			len(idx), ip.df.index.NumLevels())
	}

	if pos > ip.Len() {
		return fmt.Errorf("DataFrame.InsertRow(): invalid position: %d (max %v)", pos, ip.Len())
	}

	if len(val) != ip.df.NumCols() {
		return fmt.Errorf("DataFrame.InsertRow(): len(val) must equal number of columns (%d != %d)", len(val), ip.df.NumCols())
	}

	for _, v := range idx {
		if _, err := values.InterfaceFactory(v); err != nil {
			return fmt.Errorf("DataFrame.InsertRow(): %v", err)
		}
	}

	for _, v := range val {
		if _, err := values.InterfaceFactory(v); err != nil {
			return fmt.Errorf("DataFrame.InsertRow(): %v", err)
		}
	}

	// Insertion once errors have been handled
	for j := 0; j < ip.df.index.NumLevels(); j++ {
		ip.df.index.Levels[j].Labels.Insert(pos, idx[j])
		ip.df.index.Levels[j].Refresh()
	}
	for m := 0; m < ip.df.NumCols(); m++ {
		ip.df.vals[m].Values.Insert(pos, val[m])
	}

	return nil
}

// AppendRow adds a row at a specified integer position and modifies the DataFrame in place.
func (ip InPlace) AppendRow(val []interface{}, idx []interface{}) error {
	err := ip.df.InPlace.InsertRow(ip.Len(), val, idx)
	if err != nil {
		return fmt.Errorf("DataFrame.AppendRow(): %v", err)
	}
	return nil
}

// SetRow sets the values in the specified row to val and modifies the DataFrame in place. First converts val to be the same type as the index level.
func (ip InPlace) SetRow(row int, val interface{}) error {
	if err := ip.df.ensureRowPositions([]int{row}); err != nil {
		return fmt.Errorf("DataFrame.SetRow(): %v", err)
	}

	if _, err := values.InterfaceFactory(val); err != nil {
		return fmt.Errorf("DataFrame.SetRow(): %v", err)
	}
	for m := 0; m < ip.df.NumCols(); m++ {
		ip.df.vals[m].Values.Set(row, val)
	}
	return nil
}

// SetRows sets all the values in the specified rows to val and modifies the DataFrame in place. First converts val to be the same type as the index level.
// If an error would be encountered in any row position, the entire operation is cancelled before it starts.
func (ip InPlace) SetRows(rowPositions []int, val interface{}) error {
	if err := ip.df.ensureRowPositions(rowPositions); err != nil {
		return fmt.Errorf("DataFrame.SetRows(): %v", err)
	}
	if _, err := values.InterfaceFactory(val); err != nil {
		return fmt.Errorf("DataFrame.SetRows(): %v", err)
	}
	for m := 0; m < ip.df.NumCols(); m++ {
		for _, row := range rowPositions {
			ip.df.vals[m].Values.Set(row, val)
		}
	}
	return nil
}

// DropRow drops the row at the specified integer position and modifies the DataFrame in place.
func (ip InPlace) DropRow(row int) error {
	if err := ip.dropMany([]int{row}); err != nil {
		return fmt.Errorf("DataFrame.DropRow(): %v", err)
	}
	return nil
}

// DropRows drops the rows at the specified integer position and modifies the DataFrame in place.
// If an error would be encountered in any row position, the entire operation is cancelled before it starts.
func (ip InPlace) DropRows(rowPositions []int) error {
	if err := ip.dropMany(rowPositions); err != nil {
		return fmt.Errorf("DataFrame.DropRows(): %v", err)
	}
	return nil
}

// // DropDuplicates drops any rows containing duplicate index + DataFrame values and modifies the DataFrame in place.
// func (ip InPlace) DropDuplicates() {
// 	g := ip.df.GroupByIndex()
// 	var toDrop []int
// 	for _, group := range g.Groups() {
// 		// only inspect groups with at least one position
// 		if positions := g.groups[group].Positions; len(positions) > 0 {
// 			exists := make(map[interface{}]bool)
// 			for _, pos := range positions {
// 				if exists[ip.df.At(pos)] {
// 					toDrop = append(toDrop, pos)
// 				} else {
// 					exists[ip.df.At(pos)] = true
// 				}
// 			}
// 		}
// 	}
// 	// ducks error because position is assumed to be in DataFrame
// 	if len(toDrop) != 0 {
// 		ip.DropRows(toDrop)
// 	}
// }

// // DropDuplicates drops any rows containing duplicate index + DataFrame values and returns a new DataFrame.
// func (df *DataFrame) DropDuplicates() *DataFrame {
// 	df = df.Copy()
// df.InPlace.DropDuplicates()
// 	return df
// }

// // DropNull drops all null values and modifies the DataFrame in place.
// func (ip InPlace) DropNull() {
// 	ip.dropMany(ip.df.null())
// 	return
// }

// dropMany drops multiple rows and modifies the DataFrame in place.
func (ip InPlace) dropMany(positions []int) error {
	if err := ip.df.ensureRowPositions(positions); err != nil {
		return err
	}
	sort.IntSlice(positions).Sort()
	for i, position := range positions {
		ip.df.InPlace.dropOne(position - i)
	}
	if ip.Len() == 0 {
		ip.df.replace(newEmptyDataFrame())
	}
	return nil
}

// dropOne drops a row at a specified integer position and modifies the DataFrame in place.
// Should be called via dropMany to catch errors.
func (ip InPlace) dropOne(pos int) {
	for i := 0; i < ip.df.index.NumLevels(); i++ {
		ip.df.index.Levels[i].Labels.Drop(pos)
		ip.df.index.Levels[i].Refresh()
	}
	for m := 0; m < ip.df.NumCols(); m++ {
		ip.df.vals[m].Values.Drop(pos)
	}
	return
}

// ToFloat64 converts DataFrame values to float64 in place.
func (ip InPlace) ToFloat64() {
	for m := 0; m < ip.df.NumCols(); m++ {
		ip.df.vals[m].Values = ip.df.vals[m].Values.ToFloat64()
		ip.df.vals[m].DataType = options.Float64
	}
}

// ToInt64 converts DataFrame values to int64 in place.
func (ip InPlace) ToInt64() {
	for m := 0; m < ip.df.NumCols(); m++ {
		ip.df.vals[m].Values = ip.df.vals[m].Values.ToInt64()
		ip.df.vals[m].DataType = options.Int64
	}
}

// ToString converts DataFrame values to string in place.
func (ip InPlace) ToString() {
	for m := 0; m < ip.df.NumCols(); m++ {
		ip.df.vals[m].Values = ip.df.vals[m].Values.ToString()
		ip.df.vals[m].DataType = options.String
	}
}

// ToBool converts DataFrame values to bool in place.
func (ip InPlace) ToBool() {
	for m := 0; m < ip.df.NumCols(); m++ {
		ip.df.vals[m].Values = ip.df.vals[m].Values.ToBool()
		ip.df.vals[m].DataType = options.Bool
	}
}

// ToDateTime converts DataFrame values to datetime in place.
func (ip InPlace) ToDateTime() {
	for m := 0; m < ip.df.NumCols(); m++ {
		ip.df.vals[m].Values = ip.df.vals[m].Values.ToDateTime()
		ip.df.vals[m].DataType = options.DateTime
	}
}

// ToInterface converts DataFrame values to interface in place.
func (ip InPlace) ToInterface() {
	for m := 0; m < ip.df.NumCols(); m++ {
		ip.df.vals[m].Values = ip.df.vals[m].Values.ToInterface()
		ip.df.vals[m].DataType = options.Interface
	}
}

// [END InPlace]

// [START Copy]

// // Sort sorts the series by its values and returns a new DataFrame.
// func (df *DataFrame) Sort(asc bool) *DataFrame {
// 	df = df.Copy()
// df.InPlace.Sort(asc)
// 	return df
// }

// SwapRows swaps the selected rows and returns a new DataFrame.
func (df *DataFrame) SwapRows(i, j int) (*DataFrame, error) {
	df = df.Copy()
	if i >= df.Len() {
		return newEmptyDataFrame(), fmt.Errorf("invalid position: %d (max %v)", i, df.Len()-1)
	}
	if j >= df.Len() {
		return newEmptyDataFrame(), fmt.Errorf("invalid position: %d (max %v)", j, df.Len()-1)
	}
	df.InPlace.SwapRows(i, j)
	return df, nil
}

// SwapColumns swaps the selected rows and returns a new DataFrame.
func (df *DataFrame) SwapColumns(i, j int) (*DataFrame, error) {
	df = df.Copy()
	if i >= df.NumCols() {
		return newEmptyDataFrame(), fmt.Errorf("invalid position: %d (max %v)", i, df.Len()-1)
	}
	if j >= df.NumCols() {
		return newEmptyDataFrame(), fmt.Errorf("invalid position: %d (max %v)", j, df.Len()-1)
	}
	df.InPlace.SwapColumns(i, j)
	return df, nil
}

// InsertRow inserts a new row into the DataFrame immediately before the specified integer position and returns a new DataFrame.
func (df *DataFrame) InsertRow(pos int, val []interface{}, idx []interface{}) (*DataFrame, error) {
	df = df.Copy()
	err := df.InPlace.InsertRow(pos, val, idx)
	return df, err
}

// AppendRow adds a row at the end and returns a new DataFrame.
func (df *DataFrame) AppendRow(val []interface{}, idx []interface{}) (*DataFrame, error) {
	df, err := df.InsertRow(df.Len(), val, idx)
	if err != nil {
		return newEmptyDataFrame(), fmt.Errorf("DataFrame.AppendRow(): %v", err)
	}
	return df, nil
}

// SetRow sets the value in the specified rows to val and returns a new DataFrame.
func (df *DataFrame) SetRow(row int, val interface{}) (*DataFrame, error) {
	df = df.Copy()
	err := df.InPlace.SetRow(row, val)
	return df, err
}

// SetRows sets all the values in the specified rows to val and returns a new DataFrame.
func (df *DataFrame) SetRows(rowPositions []int, val interface{}) (*DataFrame, error) {
	df = df.Copy()
	err := df.InPlace.SetRows(rowPositions, val)
	return df, err
}

// DropRow drops the row at the specified integer position and returns a new DataFrame.
func (df *DataFrame) DropRow(row int) (*DataFrame, error) {
	df = df.Copy()
	err := df.InPlace.DropRow(row)
	return df, err
}

// DropRows drops the rows at the specified integer position and returns a new DataFrame.
func (df *DataFrame) DropRows(positions []int) (*DataFrame, error) {
	df = df.Copy()
	err := df.InPlace.DropRows(positions)
	return df, err
}

// // DropNull drops all null values and modifies the DataFrame in place.
// func (df *DataFrame) DropNull() *DataFrame {
// 	df = df.Copy()
// 	df.InPlace.DropNull()
// 	return df
// }

// ToFloat64 converts DataFrame values to float64 and returns a new DataFrame.
func (df *DataFrame) ToFloat64() *DataFrame {
	df = df.Copy()
	df.InPlace.ToFloat64()
	return df
}

// ToInt64 converts DataFrame values to int64 and returns a new DataFrame.
func (df *DataFrame) ToInt64() *DataFrame {
	df = df.Copy()
	df.InPlace.ToInt64()
	return df
}

// ToString converts DataFrame values to string and returns a new DataFrame.
func (df *DataFrame) ToString() *DataFrame {
	df = df.Copy()
	df.InPlace.ToString()
	return df
}

// ToBool converts DataFrame values to bool and returns a new DataFrame.
func (df *DataFrame) ToBool() *DataFrame {
	df = df.Copy()
	df.InPlace.ToBool()
	return df
}

// ToDateTime converts DataFrame values to time.Time and returns a new DataFrame.
func (df *DataFrame) ToDateTime() *DataFrame {
	df = df.Copy()
	df.InPlace.ToDateTime()
	return df
}

// ToInterface converts DataFrame values to interface and returns a new DataFrame.
func (df *DataFrame) ToInterface() *DataFrame {
	df = df.Copy()
	df.InPlace.ToInterface()
	return df
}

// [END Copy]
