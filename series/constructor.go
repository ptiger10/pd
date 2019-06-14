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
	Name     string
	DataType options.DataType
}

// New creates a new Series with the supplied values and n-level index.
func New(data interface{}, idx ...IndexLevel) (*Series, error) {
	// Handling values
	factory, err := values.InterfaceFactory(data)
	if err != nil {
		return nil, fmt.Errorf("series.New(): %v", err)
	}
	// Handling index
	var seriesIndex index.Index
	// Empty data: return empty index
	if data == nil {
		lvl, _ := index.NewLevel(nil, "")
		seriesIndex = index.New(lvl)
	} else if len(idx) != 0 {
		var levels []index.Level
		for i := 0; i < len(idx); i++ {
			// Any level with no values: create default index and supply name only
			if idx[i].Labels == nil {
				lvl := index.DefaultLevel(factory.Values.Len(), idx[i].Name)
				levels = append(levels, lvl)
			} else {
				// Create new level from label and name
				lvl, err := index.NewLevel(idx[i].Labels, idx[i].Name)
				// Optional type conversion
				if idx[i].DataType != options.None {
					lvl, err = lvl.Convert(idx[i].DataType)
					if err != nil {
						return nil, fmt.Errorf("series.New(): %v", err)
					}
				}
				levels = append(levels, lvl)
				if err != nil {
					return nil, fmt.Errorf("series.New(): %v", err)
				}
			}
		}
		seriesIndex = index.New(levels...)
		// No index supplied: return with default index
	} else {
		seriesIndex = index.Default(factory.Values.Len())
	}

	s := &Series{
		values:   factory.Values,
		index:    seriesIndex,
		datatype: factory.DataType,
	}

	if err := s.ensureAlignment(); err != nil {
		return nil, fmt.Errorf("series.New(): %v", err)
	}

	s.Filter = Filter{s: s}
	s.Index = Index{s: s}
	s.InPlace = InPlace{s: s}
	s.Apply = Apply{s: s}

	return s, nil
}

// NewWithConfig creates a new Series with the config struct, supplied values, and optional n-level index.
func NewWithConfig(config Config, data interface{}, idx ...IndexLevel) (*Series, error) {
	s, err := New(data, idx...)
	if err != nil {
		return nil, fmt.Errorf("series.NewWithConfig(): %v", err)
	}

	// Handling Config
	s.Name = config.Name
	if config.DataType != options.None {
		values.Convert(s.values, config.DataType)
	}
	return s, nil
}

// An IndexLevel is one level in a Series index.
type IndexLevel struct {
	Labels   interface{}
	Name     string
	DataType options.DataType
}

// Idx returns an anonymous IndexLevel with the supplied labels.
func Idx(labels interface{}) IndexLevel {
	return IndexLevel{
		Labels: labels,
	}
}

// MustNew returns a new Series or panics on error.
func MustNew(data interface{}, index ...IndexLevel) *Series {
	s, err := New(data, index...)
	if err != nil {
		log.Panicf("MustNew(): %v", err)
	}
	return s
}

func mustNewWithConfig(config Config, data interface{}, index ...IndexLevel) *Series {
	s, err := NewWithConfig(config, data, index...)
	if err != nil {
		log.Fatalf("Internal error: mustNewWithConfig(): %v", err)
	}
	return s
}
