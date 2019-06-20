package series

import (
	"fmt"
	"log"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
)

// The Config struct can be used in the custom Series constructor to name the Series or specify its data type.
// Basic usage: New("foo", series.Config{Name: "bar"})
type Config struct {
	Name            string
	DataType        options.DataType
	Index           interface{}
	IndexName       string
	MultiIndex      []interface{}
	MultiIndexNames []string
}

// New creates a new Series with the supplied values and an optional config.
func New(data interface{}, config ...Config) (*Series, error) {
	var idx index.Index
	configuration := index.Config{} // Series config

	// Handling config
	if config != nil {
		if len(config) > 1 {
			return newEmptySeries(), fmt.Errorf("series.New(): can supply at most one Config (%d > 1)", len(config))
		}
		tmp := config[0]
		configuration = index.Config{
			Name: tmp.Name, DataType: tmp.DataType,
			Index: tmp.Index, IndexName: tmp.IndexName,
			MultiIndex: tmp.MultiIndex, MultiIndexNames: tmp.MultiIndexNames,
		}
	}

	// Handling values
	factory, err := values.InterfaceFactory(data)
	if err != nil {
		return newEmptySeries(), fmt.Errorf("series.New(): %v", err)
	}

	// Handling index
	// empty data: return empty index
	if lenValues := factory.Values.Len(); lenValues == 0 {
		idx = index.New()
		// not empty data: use config
	} else {
		idx, err = index.NewFromConfig(configuration, lenValues)
		if err != nil {
			return newEmptySeries(), fmt.Errorf("series.New(): %v", err)
		}
	}

	s := &Series{
		values:   factory.Values,
		index:    idx,
		datatype: factory.DataType,
		name:     configuration.Name,
	}

	// Optional datatype conversion
	if configuration.DataType != options.None {
		s.values, err = values.Convert(s.values, configuration.DataType)
		if err != nil {
			return newEmptySeries(), fmt.Errorf("series.New(): %v", err)
		}
		s.datatype = configuration.DataType
	}

	s.Filter = Filter{s: s}
	s.Index = Index{s: s}
	s.InPlace = InPlace{s: s}
	s.Apply = Apply{s: s}

	// Alignment check
	if err := s.ensureAlignment(); err != nil {
		return newEmptySeries(), fmt.Errorf("series.New(): %v", err)
	}

	return s, err
}

// MustNew returns a new Series or logs an error and returns an empty Series.
func MustNew(data interface{}, config ...Config) *Series {
	s, err := New(data, config...)
	if err != nil {
		if options.GetLogWarnings() {
			log.Printf("dataframe.MustNew(): %v", err)
		}
		return newEmptySeries()
	}
	return s
}

func newEmptySeries() *Series {
	s, _ := New(nil)
	return s
}
