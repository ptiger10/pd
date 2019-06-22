package dataframe

import (
	"fmt"
	"log"
	"strings"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/series"
)

// New creates a new DataFrame with default column names.
func New(data []interface{}, config ...Config) (*DataFrame, error) {
	var s []*series.Series
	var idx index.Index
	var cols index.Columns
	configuration := index.Config{}
	tmp := Config{}
	var err error

	if data == nil {
		return &DataFrame{s: nil, index: index.New(), cols: index.NewColumns()}, nil
	}
	// Handling config
	if config != nil {
		if len(config) > 1 {
			return nil, fmt.Errorf("dataframe.New(): can supply at most one Config (%d > 1)", len(config))
		}
		tmp = config[0]
		configuration = index.Config{
			Name:  tmp.Name,
			Index: tmp.Index, IndexName: tmp.IndexName,
			MultiIndex: tmp.MultiIndex, MultiIndexNames: tmp.MultiIndexNames,
			Cols: tmp.Cols, ColsName: tmp.ColsName,
			MultiCols: tmp.MultiCols, MultiColsNames: tmp.MultiColsNames,
		}
	}

	// Values length
	vals, err := values.InterfaceFactory(data[0])
	if err != nil {
		return newEmptyDataFrame(), fmt.Errorf("dataframe.New(): %v", err)
	}
	valuesLen := vals.Values.Len()

	// Handling index
	idx, err = index.NewFromConfig(configuration, valuesLen)
	if err != nil {
		return newEmptyDataFrame(), fmt.Errorf("dataframe.New(): %v", err)
	}
	//Handling columns
	cols, err = index.NewColumnsFromConfig(configuration, len(data))
	if err != nil {
		return newEmptyDataFrame(), fmt.Errorf("dataframe.New(): %v", err)
	}

	// Handling Series
	for i := 0; i < len(data); i++ {
		var sNameSlice []string
		for _, col := range cols.Levels {
			sNameSlice = append(sNameSlice, fmt.Sprint(col.Labels[i]))
		}
		sName := strings.Join(sNameSlice, " | ")
		sConfig := series.Config{
			Name: sName, DataType: tmp.DataType,
			Index: tmp.Index, IndexName: tmp.IndexName,
			MultiIndex: tmp.MultiIndex, MultiIndexNames: tmp.MultiIndexNames,
		}
		n, err := series.New(data[i], sConfig)
		if err != nil {
			return newEmptyDataFrame(), fmt.Errorf("dataframe.New(): %v", err)
		}
		s = append(s, n)
	}
	df := &DataFrame{
		s:     s,
		index: idx,
		cols:  cols,
		name:  configuration.Name,
	}

	if err := df.ensureAlignment(); err != nil {
		return newEmptyDataFrame(), fmt.Errorf("dataframe.New(): %v", err)
	}

	return df, err
}

func newEmptyDataFrame() *DataFrame {
	return MustNew(nil)
}

// MustNew constructs a new DataFrame or logs an error and returns an empty DataFrame.
func MustNew(data []interface{}, config ...Config) *DataFrame {
	s, err := New(data, config...)
	if err != nil {
		if options.GetLogWarnings() {
			log.Printf("dataframe.MustNew(): %v", err)
		}
		return newEmptyDataFrame()
	}
	return s
}

func newFromComponents(s []*series.Series, idx index.Index, cols index.Columns, name string) *DataFrame {
	if s == nil {
		df, _ := New(nil)
		return df
	}
	return &DataFrame{
		s:     s,
		index: idx,
		cols:  cols,
		name:  name,
	}
}

// newSingleIndexSeries constructs a Series with a single-level index from raw values and index slices. Used to convert DataFrames to Series.
func newSingleIndexSeries(values []interface{}, idx []interface{}, name string) (*series.Series, error) {
	ret, err := series.New(nil)
	if err != nil {
		return nil, fmt.Errorf("internal error: newFromSeries(): %v", err)
	}
	if len(values) != len(idx) {
		return nil, fmt.Errorf("internal error: newFromSeries(): values must have same length as index: %d != %d", len(values), len(idx))
	}
	for i := 0; i < len(values); i++ {
		n, err := series.New(values[i], series.Config{Index: idx[i], Name: name})
		if err != nil {
			return nil, fmt.Errorf("internal error: newFromSeries(): %v", err)
		}
		ret.InPlace.Join(n)
	}
	return ret, nil
}

// returns an error if any index levels have different lengths
// or if there is a mismatch between the number of values and index items
func (df *DataFrame) ensureAlignment() error {
	if err := df.index.Aligned(); err != nil {
		return fmt.Errorf("dataframe out of alignment: %v", err)
	}
	if labels := df.index.Levels[0].Len(); df.Len() != labels {
		return fmt.Errorf("dataframe out of alignment: dataframe must have same number of values as index labels (%d != %d)", df.Len(), labels)
	}

	if df.cols.Len() != df.NumCols() {
		return fmt.Errorf("dataframe.New(): number of columns must match number of series: %d != %d",
			df.cols.Len(), df.NumCols())
	}
	return nil
}

// Config customizes the DataFrame constructor.
type Config struct {
	Name            string
	DataType        options.DataType
	Index           interface{}
	IndexName       string
	MultiIndex      []interface{}
	MultiIndexNames []string
	Cols            []interface{}
	ColsName        string
	MultiCols       [][]interface{}
	MultiColsNames  []string
}
