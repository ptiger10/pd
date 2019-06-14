package dataframe

import (
	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/series"
)

// A DataFrame is a 2D collection of one or more Series.
type DataFrame struct {
	s    []*series.Series
	cols []index.Index
}
