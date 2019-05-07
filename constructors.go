package pd

import (
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/series"
)

// Series is the default Series constructor.
// For more configuration options (e.g., custom index), use pd/series.New()
func Series(data interface{}) series.Series {
	s, err := series.New(data)
	if err != nil {
		values.Warn(err, "nil Series")
	}
	return s
}
