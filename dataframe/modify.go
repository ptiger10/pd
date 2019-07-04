package dataframe

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/ptiger10/pd/internal/index"
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
func (ip InPlace) InsertRow(row int, val []interface{}, idxLabels ...interface{}) error {
	// Handling empty DataFrame
	if Equal(ip.df, newEmptyDataFrame()) {
		df, err := New(val, Config{MultiIndex: idxLabels})
		if err != nil {
			return fmt.Errorf("DataFrame.InsertRow(): inserting into empty DataFrame requires creating a new DataFrame: %v", err)
		}
		ip.df.replace(df)
		return nil
	}

	// Handling errors
	if len(idxLabels) > ip.df.index.NumLevels() {
		return fmt.Errorf("DataFrame.InsertRow() len(idxLabels) must not exceed number of index levels: (%d != %d)",
			len(idxLabels), ip.df.index.NumLevels())
	}

	if row > ip.Len() {
		return fmt.Errorf("DataFrame.InsertRow(): invalid row: %d (max %v)", row, ip.Len())
	}

	if len(val) != ip.df.NumCols() {
		return fmt.Errorf("DataFrame.InsertRow(): len(val) must equal number of columns (%d != %d)", len(val), ip.df.NumCols())
	}

	for _, v := range idxLabels {
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
		if j < len(idxLabels) {
			ip.df.index.Levels[j].Labels.Insert(row, idxLabels[j])
		} else {
			ip.df.index.Levels[j].Labels.Insert(row, "")
		}
		// Reorder a default index
		if ip.df.index.Levels[j].IsDefault {
			// ducks error because index level is known to be in series.
			ip.df.Index.Reindex(j)
		} else {
			ip.df.index.Levels[j].Refresh()
		}
	}

	for m := 0; m < ip.df.NumCols(); m++ {
		ip.df.vals[m].Values.Insert(row, val[m])
	}

	return nil
}

// InsertColumn inserts a column with an indefinite number of column labels immediately before the specified column position and modifies the DataFrame in place.
func (ip InPlace) InsertColumn(col int, vals interface{}, colLabels ...string) error {
	// Handling empty DataFrame
	if Equal(ip.df, newEmptyDataFrame()) {
		df, err := New([]interface{}{vals}, Config{MultiCol: [][]string{colLabels}})
		if err != nil {
			return fmt.Errorf("DataFrame.InsertColumn(): inserting into empty DataFrame requires creating a new DataFrame: %v", err)
		}
		ip.df.replace(df)
		return nil
	}

	// Handling errors
	if len(colLabels) > ip.df.cols.NumLevels() {
		return fmt.Errorf("DataFrame.InsertColumn() len(colLabels) must not exceed number of column levels: (%d > %d)",
			len(colLabels), ip.df.cols.NumLevels())
	}

	if col > ip.df.NumCols() {
		return fmt.Errorf("DataFrame.InsertColumn(): invalid col: %d (max %v)", col, ip.df.NumCols())
	}
	container, err := values.InterfaceFactory(vals)
	if err != nil {
		return fmt.Errorf("DataFrame.InsertColumn(): %v", err)
	}
	if container.Values.Len() != ip.df.Len() {
		return fmt.Errorf("DataFrame.InsertColumn(): new column must contain same number of vals as all other columns (%d != %d)",
			container.Values.Len(), ip.df.Len())
	}
	// Insertion once errors have been handled

	for j := 0; j < ip.df.cols.NumLevels(); j++ {
		if j < len(colLabels) {
			ip.df.cols.Levels[j].Labels = append(ip.df.cols.Levels[j].Labels[:col], append([]string{colLabels[j]}, ip.df.cols.Levels[j].Labels[col:]...)...)
		} else {
			ip.df.cols.Levels[j].Labels = append(ip.df.cols.Levels[j].Labels[:col], append([]string{"NaN"}, ip.df.cols.Levels[j].Labels[col:]...)...)
		}
		// Reorder default columns
		if ip.df.cols.Levels[j].IsDefault {
			ip.df.cols.Levels[j].ResetDefault()
		} else {
			ip.df.cols.Levels[j].Refresh()
		}
	}
	ip.df.vals = append(ip.df.vals[:col], append([]values.Container{container}, ip.df.vals[col:]...)...)
	return nil
}

// AppendRow adds a row at a specified integer position and modifies the DataFrame in place.
func (ip InPlace) AppendRow(val []interface{}, idxLabels ...interface{}) error {
	err := ip.df.InPlace.InsertRow(ip.Len(), val, idxLabels...)
	if err != nil {
		return fmt.Errorf("DataFrame.AppendRow(): %v", err)
	}
	return nil
}

// AppendColumn adds a row at a specified integer position and modifies the DataFrame in place.
func (ip InPlace) AppendColumn(val interface{}, colLabels ...string) error {
	err := ip.df.InPlace.InsertColumn(ip.Len(), val, colLabels...)
	if err != nil {
		return fmt.Errorf("DataFrame.AppendColumn(): %v", err)
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

// SetColumn sets the values in the specified column to val and modifies the DataFrame in place.
func (ip InPlace) SetColumn(col int, val interface{}) error {
	if err := ip.df.ensureColumnPositions([]int{col}); err != nil {
		return fmt.Errorf("DataFrame.SetColumn(): %v", err)
	}

	container, err := values.InterfaceFactory(val)
	if err != nil {
		return fmt.Errorf("DataFrame.SetColumn(): %v", err)
	}
	if container.Values.Len() != ip.df.Len() {
		return fmt.Errorf("DataFrame.SetColumn(): val must have same number of values as columns in the existing dataframe (%d != %d)",
			container.Values.Len(), ip.df.Len())
	}
	ip.df.vals[col] = container
	return nil
}

// SetColumns sets the values in the specified columns to val and modifies the DataFrame in place.
func (ip InPlace) SetColumns(columnPositions []int, val interface{}) error {
	if err := ip.df.ensureColumnPositions(columnPositions); err != nil {
		return fmt.Errorf("DataFrame.SetColumn(): %v", err)
	}

	container, err := values.InterfaceFactory(val)
	if err != nil {
		return fmt.Errorf("DataFrame.SetColumn(): %v", err)
	}
	if container.Values.Len() != ip.df.Len() {
		return fmt.Errorf("DataFrame.SetColumn(): val must have same number of values as columns in the existing dataframe (%d != %d)",
			container.Values.Len(), ip.df.Len())
	}
	for _, col := range columnPositions {
		ip.df.vals[col] = container
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

// Hash computes a unique identifer for each Row.
func (r Row) hash() string {
	jsonBytes, _ := json.Marshal(r)
	h := sha1.New()
	h.Write(jsonBytes)
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// DropDuplicates drops any rows containing duplicate index + DataFrame values and modifies the DataFrame in place.
func (ip InPlace) DropDuplicates() {
	var toDrop []int
	g := ip.df.GroupByIndex()
	for _, group := range g.Groups() {
		// only inspect groups with at least one position
		if positions := g.groups[group].Positions; len(positions) > 0 {
			exists := make(map[interface{}]bool)
			for _, pos := range positions {
				if exists[ip.df.Row(pos).hash()] {
					toDrop = append(toDrop, pos)
				} else {
					exists[ip.df.Row(pos).hash()] = true
				}
			}
		}
	}
	// ducks error because position is assumed to be in DataFrame
	ip.DropRows(toDrop)
}

// DropDuplicates drops any rows containing duplicate index + DataFrame values and returns a new DataFrame.
func (df *DataFrame) DropDuplicates() *DataFrame {
	df = df.Copy()
	df.InPlace.DropDuplicates()
	return df
}

// DropNull drops all null values and modifies the DataFrame in place. If an invalid column is provided, returns original DataFrame.
func (ip InPlace) DropNull(cols ...int) {
	if err := ip.df.ensureColumnPositions(cols); err != nil {
		if options.GetLogWarnings() {
			log.Printf("df.DropNull(): %v", err)
		}
		return
	}
	ip.dropMany(ip.df.null(cols...))
	return
}

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

// DropColumn drops a column at a specified integer position and modifies the DataFrame in place.
func (ip InPlace) DropColumn(col int) error {

	// Handling errors
	if err := ip.df.ensureColumnPositions([]int{col}); err != nil {
		return fmt.Errorf("DataFrame.DropColumn(): %v", err)
	}

	for j := 0; j < ip.df.cols.NumLevels(); j++ {
		ip.df.cols.Levels[j].Labels = append(ip.df.cols.Levels[j].Labels[:col], ip.df.cols.Levels[j].Labels[col+1:]...)
		ip.df.cols.Levels[j].Refresh()
	}

	ip.df.vals = append(ip.df.vals[:col], ip.df.vals[col+1:]...)
	if ip.df.NumCols() == 0 {
		ip.df.replace(newEmptyDataFrame())
	}
	return nil
}

// DropColumns drops the columns at the specified integer positions and modifies the DataFrame in place.
func (ip InPlace) DropColumns(columnPositions []int) error {
	if err := ip.df.ensureColumnPositions(columnPositions); err != nil {
		return fmt.Errorf("DataFrame.DropColumns(): %v", err)
	}
	sort.IntSlice(columnPositions).Sort()
	for i, position := range columnPositions {
		// ducks error because all columns are assumed to be safe after aggregate error check above
		ip.df.InPlace.DropColumn(position - i)
	}
	return nil
}

func (ip InPlace) setIndex(col int) {
	container := ip.df.vals[col]
	newLevel := index.Level{Name: ip.df.cols.Name(col), Labels: container.Values, DataType: container.DataType}
	ip.df.index.Levels = append(ip.df.index.Levels, newLevel)
	// ducks error because levels are certain to be in index
	ip.df.index.SwapLevels(0, ip.df.IndexLevels()-1)

	ip.DropColumn(col)
	return
}

// SetIndex sets a column as an index level, drops the column, and modifies the DataFrame in place. If col is the only column, nothing happens.
func (ip InPlace) SetIndex(col int) error {
	if err := ip.df.ensureColumnPositions([]int{col}); err != nil {
		return fmt.Errorf("DataFrame.SetIndex(): %v", err)
	}
	if ip.df.NumCols() == 1 {
		return nil
	}
	ip.setIndex(col)
	return nil
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
func (df *DataFrame) InsertRow(row int, val []interface{}, idxLabels ...interface{}) (*DataFrame, error) {
	df = df.Copy()
	err := df.InPlace.InsertRow(row, val, idxLabels...)
	return df, err
}

// InsertColumn inserts a new column into the DataFrame immediately before the specified integer position and returns a new DataFrame.
func (df *DataFrame) InsertColumn(row int, val interface{}, colLabels ...string) (*DataFrame, error) {
	df = df.Copy()
	err := df.InPlace.InsertColumn(row, val, colLabels...)
	return df, err
}

// AppendRow adds a row at the end and returns a new DataFrame.
func (df *DataFrame) AppendRow(val []interface{}, idxLabels ...interface{}) (*DataFrame, error) {
	df, err := df.InsertRow(df.Len(), val, idxLabels...)
	if err != nil {
		return newEmptyDataFrame(), fmt.Errorf("DataFrame.AppendRow(): %v", err)
	}
	return df, nil
}

// AppendColumn adds a column at the end and returns a new DataFrame.
func (df *DataFrame) AppendColumn(val interface{}, colLabels ...string) (*DataFrame, error) {
	df, err := df.InsertColumn(df.Len(), val, colLabels...)
	if err != nil {
		return newEmptyDataFrame(), fmt.Errorf("DataFrame.AppendColumn(): %v", err)
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

// SetColumn sets all the values in the specified columns to val and returns a new DataFrame.
func (df *DataFrame) SetColumn(col int, val interface{}) (*DataFrame, error) {
	df = df.Copy()
	err := df.InPlace.SetColumn(col, val)
	return df, err
}

// SetColumns sets all the values in the specified columns to val and returns a new DataFrame.
func (df *DataFrame) SetColumns(columnPositions []int, val interface{}) (*DataFrame, error) {
	df = df.Copy()
	err := df.InPlace.SetColumns(columnPositions, val)
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

// DropColumn drops the column at the specified integer position and returns a new DataFrame.
func (df *DataFrame) DropColumn(col int) (*DataFrame, error) {
	df = df.Copy()
	err := df.InPlace.DropColumn(col)
	return df, err
}

// DropColumns drops the column at the specified integer position and returns a new DataFrame.
func (df *DataFrame) DropColumns(columnPositions []int) (*DataFrame, error) {
	df = df.Copy()
	err := df.InPlace.DropColumns(columnPositions)
	return df, err
}

// DropNull drops all null values and returns a new DataFrame. If an invalid column is provided, returns a copy of the original DataFrame.
func (df *DataFrame) DropNull(cols ...int) *DataFrame {
	df = df.Copy()
	df.InPlace.DropNull(cols...)
	return df
}

// SetIndex sets a column as an index level, drops the column, and returns a new DataFrame. If col is the only column, nothing happens.
func (df *DataFrame) SetIndex(col int) (*DataFrame, error) {
	df = df.Copy()
	err := df.InPlace.SetIndex(col)
	return df, err
}

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
