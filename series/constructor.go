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

	s.Filter = Filter{s: s}
	s.Index = Index{s: s, To: To{s: s, idx: true}}
	s.InPlace = InPlace{s: s}
	s.Apply = Apply{s: s}
	s.Math = Math{s: s}
	s.Select = Select{s: s}
	s.To = To{s: s}

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

func mustNew(data interface{}, index ...IndexLevel) *Series {
	s, err := New(data, index...)
	if err != nil {
		log.Fatalf("Internal error: mustNew(): %v", err)
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

// // Idx returns a options.ConstructorOption for use in the Series constructor New(),
// // and takes an optional Name.
// func Idx(data interface{}, options ...options.ConstructorOption) options.ConstructorOption {
// 	cfg := config.ConstructorConfig{}
// 	for _, option := range options {
// 		option(&cfg)
// 	}
// 	return func(c *config.ConstructorConfig) {
// 		idx := config.MiniIndex{
// 			Data:     data,
// 			DataType: cfg.DataType,
// 			Name:     cfg.Name,
// 		}
// 		c.Indices = append(c.Indices, idx)
// 	}
// }

// // New Series constructor
// //
// // Optional:
// //
// // - Name(string): If no name is supplied, no name will appear when Series is printed.
// // If multiple Name() options are supplied, only the final will be used.
// //
// // - Kind(options.DataType): Convert the Series values to the specified kind. Kind options: Float, Int, String, Bool, DateTime, Interface.
// // If multiple Kind() options are supplied, only the final will be used.
// //
// // - Index(interface{}, ...options.ConstructorOption): If no index is supplied, defaults to a single index of int64Values (0, 1, 2, ...n).
// // Index() also accepts an optional Name() and Kind().
// // If no name is supplied, index level will be unnamed.
// // If no kind is supplied, Index will default to the reflect.Kind() of its data.
// // If multiple Index() options are supplied, each will become its own Index level in a MultiIndex.
// func New(data interface{}, options ...options.ConstructorOption) (*Series, error) {
// 	// Setup
// 	cfg := config.ConstructorConfig{}

// 	for _, option := range options {
// 		option(&cfg)
// 	}
// 	suppliedKind := cfg.DataType
// 	var kind options.DataType
// 	name := cfg.Name

// 	var factory values.Factory
// 	var v values.Values
// 	var idx index.Index
// 	var err error

// 	// Values
// 	if data == nil {
// 		factory = values.Factory{Values: nil, DataType: options.None}
// 	} else {
// 		switch reflect.ValueOf(data).Kind() {
// 		case reflect.Float32, reflect.Float64,
// 			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
// 			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
// 			reflect.String,
// 			reflect.Bool,
// 			reflect.Struct:
// 			factory, err = values.ScalarFactory(data)

// 		case reflect.Slice:
// 			factory, err = values.SliceFactory(data)

// 		default:
// 			return nil, fmt.Errorf("unable to construct new Series: type not supported: %T", data)
// 		}
// 	}

// 	// Sets values and datatype based on the Values switch
// 	v = factory.Values
// 	kind = factory.DataType
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to construct new Series: unable to construct values: %v", err)
// 	}

// 	// options.ConstructorOptional kind conversion
// 	if suppliedKind != options.None {
// 		v, err = values.Convert(v, suppliedKind)
// 		if err != nil {
// 			return nil, fmt.Errorf("unable to construct new Series: %v", err)
// 		}
// 		kind = suppliedKind
// 	}
// 	// Index
// 	if data == nil {
// 		idx = index.New()
// 	} else {
// 		// Default case: no client-supplied Index
// 		requiredLen := v.Len()
// 		if cfg.Indices == nil {
// 			idx = index.Default(requiredLen)
// 		} else {
// 			// one or more client-supplied indices
// 			idx, err = indexFromMiniIndex(cfg.Indices, requiredLen)
// 			if err != nil {
// 				return nil, fmt.Errorf("unable to construct new Series: %v", err)
// 			}
// 		}
// 	}

// 	// Construct Series
// 	s := new(idx, v, kind, name)
// 	s.Filter = Filter{s: &s}
// 	s.Index = Index{s: &s, To: To{s: &s, idx: true}}
// 	s.InPlace = InPlace{s: &s}
// 	s.Apply = Apply{s: &s}
// 	s.Math = Math{s: &s}
// 	s.Select = Select{s: &s}
// 	s.To = To{s: &s}
// 	return s, err
// }

// func NewPointer(data interface{}, options ...options.ConstructorOption) (*Series, error) {
// 	// Setup
// 	cfg := config.ConstructorConfig{}

// 	for _, option := range options {
// 		option(&cfg)
// 	}
// 	suppliedKind := cfg.DataType
// 	var kind options.DataType
// 	name := cfg.Name

// 	var factory values.Factory
// 	var v values.Values
// 	var idx index.Index
// 	var err error

// 	// Values
// 	if data == nil {
// 		factory = values.Factory{Values: nil, DataType: options.None}
// 	} else {
// 		switch reflect.ValueOf(data).Kind() {
// 		case reflect.Float32, reflect.Float64,
// 			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
// 			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
// 			reflect.String,
// 			reflect.Bool,
// 			reflect.Struct:
// 			factory, err = values.ScalarFactory(data)

// 		case reflect.Slice:
// 			factory, err = values.SliceFactory(data)

// 		default:
// 			return &Series{}, fmt.Errorf("unable to construct new Series: type not supported: %T", data)
// 		}
// 	}

// 	// Sets values and kind based on the Values switch
// 	v = factory.Values
// 	kind = factory.DataType
// 	if err != nil {
// 		return &Series{}, fmt.Errorf("unable to construct new Series: unable to construct values: %v", err)
// 	}

// 	// options.ConstructorOptional kind conversion
// 	if suppliedKind != options.None {
// 		v, err = values.Convert(v, suppliedKind)
// 		if err != nil {
// 			return &Series{}, fmt.Errorf("unable to construct new Series: %v", err)
// 		}
// 		kind = suppliedKind
// 	}
// 	// Index
// 	if data == nil {
// 		idx = index.New()
// 	} else {
// 		// Default case: no client-supplied Index
// 		requiredLen := v.Len()
// 		if cfg.Indices == nil {
// 			idx = index.Default(requiredLen)
// 		} else {
// 			// one or more client-supplied indices
// 			idx, err = indexFromMiniIndex(cfg.Indices, requiredLen)
// 			if err != nil {
// 				return &Series{}, fmt.Errorf("unable to construct new Series: %v", err)
// 			}
// 		}
// 	}

// 	// Construct Series
// 	s := newPointer(idx, v, kind, name)
// 	s.Filter = Filter{s: s}
// 	s.Index = Index{s: s, To: To{s: s, idx: true}}
// 	s.InPlace = InPlace{s: s}
// 	s.Apply = Apply{s: s}
// 	s.Math = Math{s: s}
// 	s.Select = Select{s: s}
// 	s.To = To{s: s}
// 	return s, err
// }

// func new(idx index.Index, values values.Values, kind options.DataType, name string*Series {
// 	return Series{
// 		index:    idx,
// 		values:   values,
// 		datatype: kind,
// 		Name:     name,
// 	}
// }

// func newPointer(idx index.Index, values values.Values, kind options.DataType, name string) *Series {
// 	return &Series{
// 		index:    idx,
// 		values:   values,
// 		datatype: kind,
// 		Name:     name,
// 	}
// }

// // Calls New and panics if error. For use in testing
// func mustNew(data interface{}, options ...options.ConstructorOption*Series {
// 	s, err := New(data, options...)
// 	if err != nil {
// 		log.Printf("Internal error: %v\n", err)
// 		return nil
// 	}
// 	return s
// }

// // [START MiniIndex]

// // creates a full index from a mini client-supplied representation of an index level,
// // but returns an error if every index level is not the same length as requiredLen

// func indexFromMiniIndex(minis []config.MiniIndex, requiredLen int) (index.Index, error) {
// 	var levels []index.Level
// 	for _, miniIdx := range minis {
// 		level, err := index.NewLevel(miniIdx.Data, miniIdx.Name)
// 		if err != nil {
// 			return index.Index{}, fmt.Errorf("unable to construct index: %v", err)
// 		}
// 		labelLen := level.Len()
// 		if labelLen != requiredLen {
// 			return index.Index{}, fmt.Errorf("unable to construct index %v:"+
// 				"mismatch between supplied index length (%v) and expected length (%v)",
// 				miniIdx.Data, labelLen, requiredLen)
// 		}
// 		if miniIdx.DataType != options.None {
// 			level.Convert(miniIdx.DataType)
// 		}

// 		levels = append(levels, level)
// 	}
// 	idx := index.New(levels...)
// 	return idx, nil

// }

// // [END MiniIndex]
