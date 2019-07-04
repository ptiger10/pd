package dataframe

import (
	"fmt"
	"log"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/series"
)

// New creates a new DataFrame with default column names.
func New(data []interface{}, config ...Config) (*DataFrame, error) {
	var vals []values.Container
	var idx index.Index
	var cols index.Columns
	configuration := index.Config{}
	tmp := Config{}
	var err error

	if data == nil {
		return newEmptyDataFrame(), nil
	}
	// Handling config
	if config != nil {
		if len(config) > 1 {
			return newEmptyDataFrame(), fmt.Errorf("dataframe.New(): can supply at most one Config (%d > 1)", len(config))
		}
		tmp = config[0]
		configuration = index.Config{
			Name:     tmp.Name,
			DataType: tmp.DataType,
			Index:    tmp.Index, IndexName: tmp.IndexName,
			MultiIndex: tmp.MultiIndex, MultiIndexNames: tmp.MultiIndexNames,
			Col: tmp.Col, ColName: tmp.ColName,
			MultiCol: tmp.MultiCol, MultiColNames: tmp.MultiColNames,
		}
	}

	// Handling values
	for i := 0; i < len(data); i++ {
		container, err := values.InterfaceFactory(data[i])
		if err != nil {
			return newEmptyDataFrame(), fmt.Errorf("dataframe.New(): %v", err)
		}
		// optional DataType conversion
		if configuration.DataType != options.None {
			container.Values, err = values.Convert(container.Values, configuration.DataType)
			if err != nil {
				return newEmptyDataFrame(), fmt.Errorf("dataframe.New(): %v", err)
			}
			container.DataType = configuration.DataType
		}
		vals = append(vals, container)
	}

	// Handling index
	idx, err = index.NewFromConfig(configuration, vals[0].Values.Len())
	if err != nil {
		return newEmptyDataFrame(), fmt.Errorf("dataframe.New(): %v", err)
	}
	//Handling columns
	cols, err = index.NewColumnsFromConfig(configuration, len(data))
	if err != nil {
		return newEmptyDataFrame(), fmt.Errorf("dataframe.New(): %v", err)
	}

	df := &DataFrame{
		vals:  vals,
		index: idx,
		cols:  cols,
		name:  configuration.Name,
	}

	// df.Columns = Columns{df: df}
	df.Index = Index{df: df}
	df.InPlace = InPlace{df: df}

	if err := df.ensureAlignment(); err != nil {
		return newEmptyDataFrame(), fmt.Errorf("dataframe.New(): %v", err)
	}

	return df, err
}

func newEmptyDataFrame() *DataFrame {
	df := &DataFrame{vals: nil, index: index.New(), cols: index.NewColumns()}
	// df.Columns = Columns{df: df}
	df.Index = Index{df: df}
	df.InPlace = InPlace{df: df}
	return df
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

// newFromComponents constructs a dataframe from its constituent parts but returns an empty dataframe if series is nil
func newFromComponents(vals []values.Container, idx index.Index, cols index.Columns, name string) *DataFrame {
	if vals == nil {
		return newEmptyDataFrame()
	}
	df := &DataFrame{
		vals:  vals,
		index: idx,
		cols:  cols,
		name:  name,
	}
	// df.Columns = Columns{df: df}
	df.Index = Index{df: df}
	df.InPlace = InPlace{df: df}

	return df
}

// deriveSeries constructs a Series with a single-level index from raw values and index slices. Used to convert DataFrames to Series.
func deriveSeries(values []interface{}, idx []interface{}, name string) (*series.Series, error) {
	ret, err := series.New(nil)
	if err != nil {
		return nil, fmt.Errorf("internal error: deriveSeries(): %v", err)
	}
	if len(values) != len(idx) {
		return nil, fmt.Errorf("internal error: deriveSeries(): values must have same length as index: %d != %d", len(values), len(idx))
	}
	for i := 0; i < len(values); i++ {
		s, err := series.New(values[i], series.Config{Index: idx[i], Name: name})
		if err != nil {
			return nil, fmt.Errorf("internal error: deriveSeries(): %v", err)
		}
		ret.InPlace.Join(s)
	}
	return ret, nil
}

func (df *DataFrame) valsAligned() error {
	if df.NumCols() == 0 {
		return nil
	}
	lvl0 := df.vals[0].Values.Len()
	for i := 1; i < df.NumCols(); i++ {
		if cmpLvl := df.vals[i].Values.Len(); lvl0 != cmpLvl {
			return fmt.Errorf("df.valsAligned(): values container at %v must have same number of labels as container 0, %d != %d",
				i, cmpLvl, lvl0)
		}
	}
	return nil
}

// Copy creates a new deep copy of a Series.
func (df *DataFrame) Copy() *DataFrame {
	var valsCopy []values.Container
	for i := 0; i < len(df.vals); i++ {
		valsCopy = append(valsCopy, df.vals[i].Copy())
	}
	idxCopy := df.index.Copy()
	colsCopy := df.cols.Copy()
	dfCopy := &DataFrame{
		vals:  valsCopy,
		index: idxCopy,
		cols:  colsCopy,
		name:  df.name,
	}
	dfCopy.InPlace = InPlace{df: dfCopy}
	dfCopy.Index = Index{df: dfCopy}
	// dfCopy.Columns = Columns{df: dfCopy}
	return dfCopy
}

// hydrateSeries converts a column of values.Values into a Series with the same index as df.
func (df *DataFrame) hydrateSeries(col int) *series.Series {
	return series.FromInternalComponents(
		df.vals[col], df.index, df.cols.Name(col))
}
