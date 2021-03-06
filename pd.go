// Package pd (aka GoPandas) is a library for cleaning, aggregating, and transforming data.
// GoPandas combines a flexible API familiar to Python pandas users with the strengths of Go,
// including type safety, predictable error handling, and concurrent processing.
package pd

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ptiger10/pd/dataframe"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/series"
)

// Series constructs a new Series.
func Series(data interface{}, config ...Config) (*series.Series, error) {
	tmp := Config{}
	if config != nil {
		if len(config) > 1 {
			return series.MustNew(nil), fmt.Errorf("pd.Series(): can supply at most one Config (%d > 1)", len(config))
		}
		tmp = config[0]
	}
	sConfig := series.Config{
		Name: tmp.Name, DataType: tmp.DataType,
		Index: tmp.Index, IndexName: tmp.IndexName,
		MultiIndex: tmp.MultiIndex, MultiIndexNames: tmp.MultiIndexNames,
	}
	s, err := series.New(data, sConfig)
	if err != nil {
		return series.MustNew(nil), fmt.Errorf("pd.Series(): %v", err)
	}
	return s, nil

}

// DataFrame constructs a new DataFrame.
func DataFrame(data []interface{}, config ...Config) (*dataframe.DataFrame, error) {
	tmp := Config{}
	if config != nil {
		if len(config) > 1 {
			return dataframe.MustNew(nil), fmt.Errorf("pd.Series(): can supply at most one Config (%d > 1)", len(config))
		}
		tmp = config[0]
	}
	dfConfig := dataframe.Config{
		Name: tmp.Name, DataType: tmp.DataType,
		Index: tmp.Index, IndexName: tmp.IndexName,
		MultiIndex: tmp.MultiIndex, MultiIndexNames: tmp.MultiIndexNames,
		Col: tmp.Col, ColName: tmp.ColName,
		MultiCol: tmp.MultiCol, MultiColNames: tmp.MultiColNames,
	}
	df, err := dataframe.New(data, dfConfig)
	if err != nil {
		return dataframe.MustNew(nil), fmt.Errorf("pd.DataFrame(): %v", err)
	}
	return df, nil
}

// ReadInterface converts [][]interface{}{row1{col1, ...}...} into a DataFrame
func ReadInterface(input [][]interface{}, config ...ReadOptions) (*dataframe.DataFrame, error) {
	if len(input) == 0 {
		return dataframe.MustNew(nil), fmt.Errorf("ReadInterface(): Input must contain at least one row")
	}
	if len(input[0]) == 0 {
		return dataframe.MustNew(nil), fmt.Errorf("ReadInterface(): must contain at least one column")
	}

	data := make([][]interface{}, len(input))
	for i := 0; i < len(input); i++ {
		data[i] = make([]interface{}, len(input[0]))
		for m := 0; m < len(input[0]); m++ {
			data[i][m] = input[i][m]
		}
	}

	tmp := ReadOptions{}
	if config != nil {
		if len(config) > 1 {
			return dataframe.MustNew(nil), fmt.Errorf("ReadInterface(): can supply at most one ReadOptions (%d > 1)",
				len(config))
		}
		tmp = config[0]
	}

	var tmpMultiCol [][]interface{}
	if tmp.DropRows > len(data) {
		return dataframe.MustNew(nil), fmt.Errorf("ReadInterface(): DropRows cannot exceed the number of rows (%d > %d)",
			tmp.DropRows, len(data))
	}

	data = data[tmp.DropRows:]
	// header rows
	if tmp.HeaderRows > len(data) {
		return dataframe.MustNew(nil), fmt.Errorf("ReadInterface(): HeaderRows cannot exceed the number of rows (%d > %d)",
			tmp.HeaderRows, len(data))
	}

	tmpMultiCol = data[:tmp.HeaderRows]
	for m := 0; m < tmp.HeaderRows; m++ {
		tmpMultiCol[m] = tmpMultiCol[m][tmp.IndexCols:]
	}

	data = data[tmp.HeaderRows:]

	if tmp.IndexCols > len(data[0]) {
		return dataframe.MustNew(nil), fmt.Errorf("ReadInterface(): IndexCols cannot exceed the number of rows (%d > %d)",
			tmp.IndexCols, len(data))
	}

	tmpVals := make([][]interface{}, len(data[0])-tmp.IndexCols)
	tmpMultiIndex := make([][]interface{}, tmp.IndexCols)

	for m := 0; m < len(data[0])-tmp.IndexCols; m++ {
		tmpVals[m] = make([]interface{}, len(data))
	}
	for m := 0; m < tmp.IndexCols; m++ {
		tmpMultiIndex[m] = make([]interface{}, len(data))
	}

	// transpose index and values
	for i := 0; i < len(data); i++ {
		for m := 0; m < len(data[0]); m++ {
			if m < tmp.IndexCols {
				tmpMultiIndex[m][i] = data[i][m]
			} else {
				tmpVals[m-tmp.IndexCols][i] = data[i][m]
			}
		}
	}
	// convert [][]interface{} to []interface{} of []interface for compatibility with DataFrame constructor
	var (
		multiIndex []interface{}
		vals       []interface{}
	)
	for _, col := range tmpMultiIndex {
		multiIndex = append(multiIndex, col)
	}
	for _, col := range tmpVals {
		vals = append(vals, col)
	}
	multiCol := make([][]string, len(tmpMultiCol))

	if len(tmpMultiCol) > 0 {
		for j := 0; j < len(tmpMultiCol); j++ {
			multiCol[j] = make([]string, len(tmpMultiCol[0]))
			for m := 0; m < len(tmpMultiCol[0]); m++ {
				multiCol[j][m] = fmt.Sprint(tmpMultiCol[j][m])
			}
		}
	}

	// ducks error because all []interface{} values are supported and Config properties are controlled
	df, _ := DataFrame(vals, Config{Manual: tmp.Manual, MultiIndex: multiIndex, MultiCol: multiCol})

	for k, v := range tmp.DataTypes {
		colInt := df.SelectCol(k)
		if colInt != -1 {
			df.InPlace.SetCol(colInt, df.ColAt(colInt).Convert(v))
		}
	}
	for k, v := range tmp.IndexDataTypes {
		err := df.Index.Convert(v, k)
		if err != nil {
			if options.GetLogWarnings() {
				log.Printf("warning: ReadInterface() converting IndexDataTypes: %v", err)
			}
		}
	}
	df.RenameCols(tmp.Rename)

	return df, nil
}

// ReadCSV converts a CSV file into a DataFrame.
func ReadCSV(path string, config ...ReadOptions) (*dataframe.DataFrame, error) {
	tmp := ReadOptions{}
	if config != nil {
		if len(config) > 1 {
			return dataframe.MustNew(nil), fmt.Errorf("ReadCSV(): can supply at most one ReadOptions (%d > 1)",
				len(config))
		}
		tmp = config[0]
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return dataframe.MustNew(nil), fmt.Errorf("ReadCSV(): %s", err)
	}
	reader := csv.NewReader(bytes.NewReader(data))
	if tmp.Delimiter != 0 {
		reader.Comma = tmp.Delimiter
	}

	records, err := reader.ReadAll()
	if err != nil {
		return dataframe.MustNew(nil), fmt.Errorf("ReadCSV(): %v", err)
	}
	if len(records) == 0 {
		return dataframe.MustNew(nil), fmt.Errorf("ReadCSV(): input must contain at least one row")
	}

	// convert to [][]interface
	var interfaceRecords [][]interface{}
	for j := 0; j < len(records); j++ {
		interfaceRecords = append(interfaceRecords, make([]interface{}, len(records[0])))
		for m := 0; m < len(records[0]); m++ {
			// optional interpolation if not in Manual mode
			if !tmp.Manual {
				interfaceRecords[j][m] = values.InterpolateString(records[j][m])
			} else {
				interfaceRecords[j][m] = records[j][m]
			}
		}
	}

	df, err := ReadInterface(interfaceRecords, tmp)
	if err != nil {
		return dataframe.MustNew(nil), fmt.Errorf("ReadCSV(): %v", err)
	}
	return df, nil
}

// Config customizes the construction of either a DataFrame or Series.
type Config struct {
	Name            string
	DataType        options.DataType
	Index           interface{}
	IndexName       string
	MultiIndex      []interface{}
	MultiIndexNames []string
	Col             []string
	ColName         string
	MultiCol        [][]string
	MultiColNames   []string
	Manual          bool
}

// ReadOptions are options for reading in files from other formats
type ReadOptions struct {
	DropRows        int
	HeaderRows      int
	IndexCols       int
	Manual          bool
	DataTypes       map[string]string
	IndexDataTypes  map[int]string
	ColumnDataTypes map[int]string
	Rename          map[string]string
	Delimiter       rune
}
