package dataframe

import (
	"fmt"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/series"
)

// New creates a new DataFrame with default column names.
func New(data []interface{}, index ...series.IndexLevel) (*DataFrame, error) {
	if data == nil {
		return nil, fmt.Errorf("dataframe.New(): data cannot be nil")
	}
	var s []*series.Series
	columnSlice := values.NewDefaultColumns(len(data))
	cols := make(map[string][]int, len(columnSlice))
	for i, val := range columnSlice {
		cols[val] = append(cols[val], i)
	}
	for i := 0; i < len(data); i++ {
		n, err := series.NewWithConfig(series.Config{Name: columnSlice[i]}, data[i], index...)
		if err != nil {
			return nil, fmt.Errorf("dataframe.New(): %v", err)
		}
		s = append(s, n)
	}
	var length int
	if len(data) > 0 && data[0] != nil {
		length = s[0].Len()
	}
	idx, err := indexFactory(index, length, data == nil)
	if err != nil {
		return nil, fmt.Errorf("dataframe.New(): %v", err)
	}
	df := &DataFrame{
		s:      s,
		index:  idx,
		colMap: cols,
	}
	return df, nil
}

// NewWithConfig creates a new DataFrame with the config struct, supplied values, and optional n-level index.
func NewWithConfig(config Config, data []interface{}, index ...series.IndexLevel) (*DataFrame, error) {
	df, err := New(data, index...)
	if err != nil {
		return nil, fmt.Errorf("dataframe.NewWithConfig(): %v", err)
	}
	if config.Columns != nil {
		if len(config.Columns) != len(df.s) {
			return nil, fmt.Errorf("dataframe.NewWithConfig(): number of columns must match number of series: %d != %d",
				len(config.Columns), len(df.s))
		}
		for i, val := range config.Columns {
			df.colMap[val] = append(df.colMap[val], i)
		}
	}
	df.Name = config.Name
	for i := 0; i < len(config.Columns); i++ {
		df.s[i].Name = config.Columns[i]
	}
	return df, nil
}

// Config customizes the new DataFrame constructor.
type Config struct {
	Name           string
	Cols           []interface{}
	ColsName       string
	MultiCols      [][]interface{}
	MultiColsNames []string
	DataType       options.DataType
}

// indexFactory creates an index from supplied IndexLevels.
// Duplicated from series to maintain index.Index encapsulation.
func indexFactory(idx []series.IndexLevel, length int, nullData bool) (index.Index, error) {
	// Handling index
	var ret index.Index
	// Empty data: return empty index
	if nullData {
		lvl, _ := index.NewLevel(nil, "")
		ret = index.New(lvl)
	} else if len(idx) != 0 {
		var levels []index.Level
		for i := 0; i < len(idx); i++ {
			// Any level with no values: create default index and supply name only
			if idx[i].Labels == nil {
				lvl := index.DefaultLevel(length, idx[i].Name)
				levels = append(levels, lvl)
			} else {
				// Create new level from label and name
				lvl, err := index.NewLevel(idx[i].Labels, idx[i].Name)
				// Optional type conversion
				if idx[i].DataType != options.None {
					lvl, err = lvl.Convert(idx[i].DataType)
					if err != nil {
						return index.Index{}, fmt.Errorf("internal.IndexFactory(): %v", err)
					}
				}
				levels = append(levels, lvl)
				if err != nil {
					return index.Index{}, fmt.Errorf("internal.IndexFactory(): %v", err)
				}
			}
		}
		ret = index.New(levels...)
		// No index supplied: return with default index
	} else {
		ret = index.Default(length)
	}
	return ret, nil
}
