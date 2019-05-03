package pd

import (
	"log"

	"github.com/ptiger10/pd/series"
)

// Series is the default Series constructor.
// For more configuration options (e.g., custom index), use pd/series.New()
func Series(data interface{}) series.Series {
	s, err := series.New(data)
	if err != nil {
		log.Printf("Unable to create series: %v", err)
	}
	return s
}
