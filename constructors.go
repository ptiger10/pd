package pd

import (
	"fmt"

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
