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

// New creates a new Series with the supplied values and default index.
func New(data interface{}, config ...Config) (*Series, error) {
	// Handling values
	factory, err := values.InterfaceFactory(data)
	if err != nil {
		return nil, fmt.Errorf("series.New(): %v", err)
	}
	// Handling index
	var seriesIndex index.Index
	// Empty data: return empty index
	if data == nil {
		seriesIndex = index.New()
		// Not empty data: return with default index
	} else {
		seriesIndex = index.NewDefault(factory.Values.Len())
	}

	s := &Series{
		values:   factory.Values,
		index:    seriesIndex,
		datatype: factory.DataType,
	}

	s.Filter = Filter{s: s}
	s.Index = Index{s: s}
	s.InPlace = InPlace{s: s}
	s.Apply = Apply{s: s}

	if config != nil {
		if len(config) > 1 {
			return nil, fmt.Errorf("series.New(): can supply at most one Config (%d > 1)", len(config))
		}
		s, err = s.configure(config[0])
	}
	return s, err
}

// configure configures an existing Series with supplied config struct
func (s *Series) configure(config Config) (*Series, error) {
	var err error
	// Handling name
	s.Name = config.Name
	// If Series has no values, ignore other Config fields
	if s.Len() == 0 {
		return s, nil
	}

	// Optional datatype conversion
	if config.DataType != options.None {
		s.values, err = values.Convert(s.values, config.DataType)
		if err != nil {
			return nil, fmt.Errorf("series.NewWithConfig(): %v", err)
		}
		s.datatype = config.DataType
	}

	// Handling index
	if config.Index != nil && config.MultiIndex != nil {
		return nil, fmt.Errorf("series.NewWithConfig(): supplying both config.Index and config.MultiIndex is ambiguous; supply one or the other")
	}
	if config.Index != nil {
		newLevel, err := index.NewLevel(config.Index, config.IndexName)
		if err != nil {
			return nil, fmt.Errorf("series.NewWithConfig(): %v", err)
		}
		s.index = index.New(newLevel)
	}
	if config.MultiIndex != nil {
		if config.MultiIndexNames != nil && len(config.MultiIndexNames) != len(config.MultiIndex) {
			return nil, fmt.Errorf(
				"series.NewWithConfig(): if MultiIndexNames is not nil, it must must have same length as MultiIndex: %d != %d",
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
				return nil, fmt.Errorf("series.NewWithConfig(): %v", err)
			}
			newLevels = append(newLevels, newLevel)
		}
		s.index = index.New(newLevels...)
	}

	if err := s.ensureAlignment(); err != nil {
		return nil, fmt.Errorf("series.NewWithConfig(): %v", err)
	}
	return s, nil
}

// MustNew returns a new Series or logs an error and returns a blank Series.
func MustNew(data interface{}, config ...Config) *Series {
	s, err := New(data, config...)
	if err != nil {
		if options.GetLogWarnings() {
			log.Printf("MustNew(): %v", err)
		}
	}
	return s
}
