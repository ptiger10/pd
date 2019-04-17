package pd

import "github.com/ptiger10/pd/series"

// Series is the default Series constructor.
// For more configuration options (e.g., custom index), use pd/series.New()
func Series(data interface{}) *series.Series {
	return series.New(data)
}
