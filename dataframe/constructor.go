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
			if cols.Len() != len(data) {
				return newEmptyDataFrame(), fmt.Errorf("dataframe.New(): dataframe out of alignment: the number of columns in each level must equal the number of Series: %d != %d",
					cols.Len(), len(data))
			}
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

// newFromComponents constructs a dataframe from its constituent parts but returns an empty dataframe if series is nil
func newFromComponents(s []*series.Series, idx index.Index, cols index.Columns, name string) *DataFrame {
	if s == nil {
		return newEmptyDataFrame()
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
		return nil, fmt.Errorf("internal error: newSingleIndexSeries(): %v", err)
	}
	if len(values) != len(idx) {
		return nil, fmt.Errorf("internal error: newSingleIndexSeries(): values must have same length as index: %d != %d", len(values), len(idx))
	}
	for i := 0; i < len(values); i++ {
		s, err := series.New(values[i], series.Config{Index: idx[i], Name: name})
		if err != nil {
			return nil, fmt.Errorf("internal error: newSingleIndexSeries(): %v", err)
		}
		ret.InPlace.Join(s)
	}
	return ret, nil
}

func (df *DataFrame) seriesAligned() error {
	if df.NumCols() == 0 {
		return nil
	}
	lvl0 := df.s[0].Len()
	for i := 1; i < df.NumCols(); i++ {
		if cmpLvl := df.s[i].Len(); lvl0 != cmpLvl {
			return fmt.Errorf("df.seriesAligned(): series %v must have same number of labels as series 0, %d != %d",
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

	if err := df.seriesAligned(); err != nil {
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
	var sCopy []*series.Series
	for i := 0; i < len(df.s); i++ {
		sCopy = append(sCopy, df.s[i].Copy())
	}
	idxCopy := df.index.Copy()
	colsCopy := df.cols.Copy()
	dfCopy := &DataFrame{
		s:     sCopy,
		index: idxCopy,
		cols:  colsCopy,
		name:  df.name,
	}
	// dfCopy.Apply = Apply{s: copyS}
	// dfCopy.Index = Index{s: copyS}
	// dfCopy.InPlace = InPlace{s: copyS}
	return dfCopy
}
