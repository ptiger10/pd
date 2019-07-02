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
		return &DataFrame{vals: nil, index: index.New(), cols: index.NewColumns()}, nil
	}
	// Handling config
	if config != nil {
		if len(config) > 1 {
			return newEmptyDataFrame(), fmt.Errorf("dataframe.New(): can supply at most one Config (%d > 1)", len(config))
		}
		tmp = config[0]
		configuration = index.Config{
			Name:  tmp.Name,
			Index: tmp.Index, IndexName: tmp.IndexName,
			MultiIndex: tmp.MultiIndex, MultiIndexNames: tmp.MultiIndexNames,
			Col: tmp.Col, ColName: tmp.ColName,
			MultiCol: tmp.MultiCol, MultiColNames: tmp.MultiColNames,
		}
	}

	// Handling values
	for i := 0; i < len(data); i++ {
		v, err := values.InterfaceFactory(data[i])
		if err != nil {
			return newEmptyDataFrame(), fmt.Errorf("dataframe.New(): %v", err)
		}
		vals = append(vals, v)
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

// newFromComponents constructs a dataframe from its constituent parts but returns an empty dataframe if series is nil
func newFromComponents(vals []values.Container, idx index.Index, cols index.Columns, name string) *DataFrame {
	if vals == nil {
		return newEmptyDataFrame()
	}
	return &DataFrame{
		vals:  vals,
		index: idx,
		cols:  cols,
		name:  name,
	}
}

// newSingleIndexSeries constructs a Series with a single-level index from raw values and index slices. Used to convert DataFrames to Series.
func newSingleIndexSeries(data []interface{}, idx index.Index, name string) (*series.Series, error) {
	ret, err := series.New(nil)
	if err != nil {
		return nil, fmt.Errorf("internal error: newSingleIndexSeries(): %v", err)
	}
	if len(data) != idx.Len() {
		return nil, fmt.Errorf("internal error: newSingleIndexSeries(): values must have same length as index: %d != %d", len(data), idx.Len())
	}
	for i := 0; i < len(data); i++ {
		vals := values.MustCreateValuesFromInterface(data[i])
		s := series.FromInternalComponents(vals.Values, idx, vals.DataType, name)
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

// returns an error if any index levels have different lengths
// or if there is a mismatch between the number of values and index items
func (df *DataFrame) ensureAlignment() error {
	if err := df.index.Aligned(); err != nil {
		return fmt.Errorf("dataframe out of alignment: %v", err)
	}
	if labels := df.index.Levels[0].Len(); df.Len() != labels {
		return fmt.Errorf("dataframe out of alignment: dataframe must have same number of values as index labels (%d != %d)", df.Len(), labels)
	}

	if err := df.valsAligned(); err != nil {
		return fmt.Errorf("dataframe out of alignment: %v", err)
	}

	if df.cols.Len() != df.NumCols() {
		return fmt.Errorf("dataframe.New(): number of columns must match number of series: %d != %d",
			df.cols.Len(), df.NumCols())
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
	// dfCopy.Index = Index{s: copyS}
	// dfCopy.InPlace = InPlace{s: copyS}
	return dfCopy
}

// hydrateSeries converts a column of values.Values into a Series with the same index as df.
func (df *DataFrame) hydrateSeries(col int) *series.Series {
	return series.FromInternalComponents(
		df.vals[col].Values, df.index, df.vals[col].DataType, df.cols.Names()[col])
}
