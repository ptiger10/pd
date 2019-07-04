package series

import (
	"fmt"
	"log"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
)

// New creates a new Series with the supplied values and an optional config.
func New(data interface{}, config ...Config) (*Series, error) {
	var idx index.Index
	configuration := index.Config{} // Series config

	if data == nil {
		return newEmptySeries(), nil
	}

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
	container, err := values.InterfaceFactory(data)
	if err != nil {
		return newEmptySeries(), fmt.Errorf("series.New(): %v", err)
	}

	// Handling index
	// empty data: return empty index
	if lenValues := container.Values.Len(); lenValues == 0 {
		idx = index.New()
		// not empty data: use config
	} else {
		idx, err = index.NewFromConfig(configuration, lenValues)
		if err != nil {
			return newEmptySeries(), fmt.Errorf("series.New(): %v", err)
		}
	}

	s := &Series{
		values:   container.Values,
		index:    idx,
		datatype: container.DataType,
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

	s.Index = Index{s: s}
	s.InPlace = InPlace{s: s}

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
			log.Printf("series.MustNew(): %v", err)
		}
		return newEmptySeries()
	}
	return s
}

func newEmptySeries() *Series {
	// ducks error because InterfaceFactory supports nil data
	container, _ := values.InterfaceFactory(nil)
	s := &Series{index: index.New(), values: container.Values, datatype: container.DataType}
	s.Index = Index{s: s}
	s.InPlace = InPlace{s: s}
	return s
}

// Copy creates a new deep copy of a Series.
func (s *Series) Copy() *Series {
	idx := s.index.Copy()
	valsCopy := s.values.Copy()
	copyS := &Series{
		values:   valsCopy,
		index:    idx,
		datatype: s.datatype,
		name:     s.name,
	}
	copyS.Index = Index{s: copyS}
	copyS.InPlace = InPlace{s: copyS}
	return copyS
}

// [START semi-private methods]

// FromInternalComponents is a semi-private method for hydrating Series within the DataFrame module.
// The required inputs are not available to the caller.
func FromInternalComponents(container values.Container, index index.Index, name string) *Series {
	s := &Series{
		values:   container.Values,
		index:    index,
		datatype: container.DataType,
		name:     name,
	}
	s.Index = Index{s: s}
	s.InPlace = InPlace{s: s}
	return s
}

// ToInternalComponents is a semi-private method for using a Series within the DataFrame module.
// The required inputs are not available to the caller.
func (s *Series) ToInternalComponents() (values.Container, index.Index) {
	return values.Container{Values: s.values.Copy(), DataType: s.datatype}, s.index.Copy()
}

// [END semi-private methods]
