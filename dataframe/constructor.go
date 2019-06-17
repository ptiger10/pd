package dataframe

import (
	"fmt"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/series"
)

// New creates a new DataFrame with default column names.
func New(data []interface{}, config ...Config) (*DataFrame, error) {
	if data == nil {
		return &DataFrame{s: nil, index: index.New()}, nil
	}
	vals, _ := values.InterfaceFactory(data[0])
	valuesLen := vals.Values.Len()

	var s []*series.Series
	var idx index.Index
	var cols index.Columns
	var err error

	if config != nil {
		if len(config) > 1 {
			return nil, fmt.Errorf("dataframe.New(): can supply at most one Config (%d > 1)", len(config))
		}
		// Handling index
		idx, err = config[0].indexFactory()
		if err != nil {
			return nil, fmt.Errorf("dataframe.New(): %v", err)
		}
		//Handling columns
		cols, err = config[0].columnFactory()
		if err != nil {
			return nil, fmt.Errorf("dataframe.New(): %v", err)
		}
	} else {
		idx = index.NewDefault(valuesLen)
		cols = index.NewDefaultColumns(len(data))
	}

	// Handling Series
	for i := 0; i < len(data); i++ {
		sName := fmt.Sprint(cols.Levels[0].Labels[i])
		n, err := series.New(data[i], series.Config{ConfigInternalIndex: &idx, Name: sName})
		if err != nil {
			return nil, fmt.Errorf("dataframe.New(): %v", err)
		}
		s = append(s, n)
	}
	df := &DataFrame{
		s:     s,
		index: idx,
		cols:  cols,
	}

	if err := df.ensureAlignment(); err != nil {
		return nil, fmt.Errorf("dataframe.New(): %v", err)
	}

	return df, err
}

// create an Index from supplied config fields
func (config Config) indexFactory() (index.Index, error) {
	var idx index.Index
	if config.Index != nil && config.MultiIndex != nil {
		return index.Index{}, fmt.Errorf("indexFactory(): supplying both config.Index and config.MultiIndex is ambiguous; supply one or the other")
	}
	if config.Index != nil {
		newLevel, err := index.NewLevel(config.Index, config.IndexName)
		if err != nil {
			return index.Index{}, fmt.Errorf("indexFactory(): %v", err)
		}
		idx = index.New(newLevel)
	}
	if config.MultiIndex != nil {
		if config.MultiIndexNames != nil && len(config.MultiIndexNames) != len(config.MultiIndex) {
			return index.Index{}, fmt.Errorf(
				"indexFactory(): if MultiIndexNames is not nil, it must must have same length as MultiIndex: %d != %d",
				len(config.MultiIndexNames), len(config.MultiIndex))
		}
		var newLevels []index.Level
		for i := 0; i < len(config.MultiIndex); i++ {
			var levelName string
			if i < len(config.MultiIndexNames) {
				levelName = config.MultiIndexNames[i]
			} else {
				levelName = ""
			}
			newLevel, err := index.NewLevel(config.MultiIndex[i], levelName)
			if err != nil {
				return index.Index{}, fmt.Errorf("dataframe.New(): %v", err)
			}
			newLevels = append(newLevels, newLevel)
		}
		idx = index.New(newLevels...)
	}
	return idx, nil

}

// create Columns from supplied config fields
func (config Config) columnFactory() (index.Columns, error) {
	var columns index.Columns
	// Handling columns
	if config.Cols != nil && config.MultiCols != nil {
		return index.Columns{}, fmt.Errorf("columnFactory(): supplying both config.Index and config.MultiIndex is ambiguous; supply one or the other")
	}
	if config.Cols != nil {
		newLevel := index.NewColLevel(config.Cols, config.ColsName)
		columns = index.NewColumns(newLevel)
	}
	if config.MultiIndex != nil {
		if config.MultiColsNames != nil && len(config.MultiColsNames) != len(config.MultiCols) {
			return index.Columns{}, fmt.Errorf(
				"columnFactory(): if MultiColsNames is not nil, it must must have same length as MultiCols: %d != %d",
				len(config.MultiColsNames), len(config.MultiCols))
		}
		var newLevels []index.ColLevel
		for i := 0; i < len(config.MultiCols); i++ {
			var levelName string
			if i < len(config.MultiColsNames) {
				levelName = config.MultiColsNames[i]
			}
			newLevel := index.NewColLevel(config.MultiCols[i], levelName)
			newLevels = append(newLevels, newLevel)
		}
		columns = index.NewColumns(newLevels...)
	}
	return columns, nil
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

	if df.cols.Len() != df.Cols() {
		return fmt.Errorf("dataframe.New(): number of columns must match number of series: %d != %d",
			df.cols.Len(), df.Cols())
	}
	return nil
}

// Config customizes the DataFrame constructor.
type Config struct {
	Name            string
	Index           interface{}
	IndexName       string
	MultiIndex      []interface{}
	MultiIndexNames []string
	Cols            []interface{}
	ColsName        string
	MultiCols       [][]interface{}
	MultiColsNames  []string
}
