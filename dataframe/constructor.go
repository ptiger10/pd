package dataframe

import (
	"fmt"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/series"
)

// New creates a new DataFrame with default column names.
func New(data []interface{}, index ...series.IndexLevel) (*DataFrame, error) {
	var s []*series.Series
	for _, val := range data {
		n, err := series.New(val, index...)
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
		s:     s,
		index: idx,
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
		df.cols = config.Columns
	}
	df.Name = config.Name
	for i := 0; i < len(df.cols); i++ {
		df.s[i].Name = df.cols[i]
	}
	return df, nil
}

// Config customizes the new DataFrame constructor.
type Config struct {
	Name     string
	Columns  []string
	DataType options.DataType
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
