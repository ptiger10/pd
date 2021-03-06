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

	if len(data) == 0 {
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

	// Handling map
	isSplit, extractedData, extractedColumns := values.MapSplitter(data)
	if isSplit {
		data = extractedData
		configuration.Col = extractedColumns
	}

	// Handling values
	vals, err = values.InterfaceSliceFactory(data, tmp.Manual, configuration.DataType)
	if err != nil {
		return newEmptyDataFrame(), fmt.Errorf("dataframe.New(): %v", err)
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

	df.Columns = Columns{df: df}
	df.Index = Index{df: df}
	df.InPlace = InPlace{df: df}

	if err := df.ensureAlignment(); err != nil {
		return newEmptyDataFrame(), fmt.Errorf("dataframe.New(): %v", err)
	}

	return df, err
}

func newEmptyDataFrame() *DataFrame {
	df := &DataFrame{vals: nil, index: index.New(), cols: index.NewColumns()}
	df.Columns = Columns{df: df}
	df.Index = Index{df: df}
	df.InPlace = InPlace{df: df}
	return df
}

// MustNew constructs a new DataFrame or logs an error and returns an empty DataFrame.
func MustNew(data []interface{}, config ...Config) *DataFrame {
	df, err := New(data, config...)
	if err != nil {
		if options.GetLogWarnings() {
			log.Printf("dataframe.MustNew(): %v", err)
		}
		return newEmptyDataFrame()
	}
	return df
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
	df.Columns = Columns{df: df}
	df.Index = Index{df: df}
	df.InPlace = InPlace{df: df}

	return df
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
	dfCopy.Columns = Columns{df: dfCopy}
	dfCopy.Index = Index{df: dfCopy}
	dfCopy.InPlace = InPlace{df: dfCopy}
	return dfCopy
}

// hydrateSeries converts a column of values.Values into a Series with the same index as df.
func (df *DataFrame) hydrateSeries(col int) *series.Series {
	return series.FromInternalComponents(
		df.vals[col], df.index, df.cols.Name(col))
}
