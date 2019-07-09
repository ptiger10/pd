package pd

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ptiger10/pd/dataframe"
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
}

// ReadOptions are options for reading in files from other formats
type ReadOptions struct {
	DropRows        int
	NumHeaderRows   int
	NumIndexColumns int
	DataTypes       map[string]string
	IndexDataTypes  map[int]string
	ColumnDataTypes map[int]string
	Rename          map[string]string
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
	records, error := reader.ReadAll()
	if error != nil {
		return dataframe.MustNew(nil), fmt.Errorf("ReadCSV(): %s", err)
	}

	tmpVals := make([][]string, 0)
	tmpMultiIndex := make([][]string, 0)

	var multiCol [][]string
	if tmp.DropRows > len(records) {
		return dataframe.MustNew(nil), fmt.Errorf("ReadCSV(): %s", err)
	}
	records = records[tmp.DropRows:]
	// header rows
	multiCol = records[:tmp.NumHeaderRows]
	for m := 0; m < tmp.NumHeaderRows; m++ {
		multiCol[m] = multiCol[m][tmp.NumIndexColumns:]
	}
	if tmp.NumHeaderRows > len(records) {
		return dataframe.MustNew(nil), fmt.Errorf("ReadCSV(): %s", err)
	}
	records = records[tmp.NumHeaderRows:]

	// transpose index and values
	for i := 0; i < len(records); i++ {
		for m := 0; m < len(records[0]); m++ {
			if m < tmp.NumIndexColumns {
				if m >= len(tmpMultiIndex) {
					tmpMultiIndex = append(tmpMultiIndex, make([]string, len(records)))
				}
				tmpMultiIndex[m][i] = records[i][m]
			} else {
				if m-tmp.NumIndexColumns >= len(tmpVals) {
					tmpVals = append(tmpVals, make([]string, len(records)))
				}
				tmpVals[m-tmp.NumIndexColumns][i] = records[i][m]
			}
		}
	}
	// convert [][]string to []interface{} of []string for compatability with DataFrame constructor
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

	df, err := DataFrame(vals, Config{MultiIndex: multiIndex, MultiCol: multiCol})
	if err != nil {
		return dataframe.MustNew(nil), fmt.Errorf("ReadCSV(): %s", err)
	}
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
				log.Printf("pd.ReadCSV() converting IndexDataTypes: %v", err)
			}
		}
	}
	df.RenameCols(tmp.Rename)

	return df, nil
}
