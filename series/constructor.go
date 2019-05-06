package series

import (
	"fmt"
	"reflect"

	"github.com/ptiger10/pd/kinds"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
)

// An Option is an optional parameter in the Series constructor
type Option func(*seriesConfig)
type seriesConfig struct {
	indices []miniIndex
	kind    kinds.Kind
	name    string
}

// Kind will convert either values or an index level to the specified kind
func Kind(kind kinds.Kind) Option {
	return func(c *seriesConfig) {
		c.kind = kind
	}
}

// Name will name either values or an index level
func Name(n string) Option {
	return func(c *seriesConfig) {
		c.name = n
	}

}

// Index returns a Option for use in the Series constructor New(),
// and takes an optional Name.
func Index(data interface{}, options ...Option) Option {
	config := seriesConfig{}
	for _, option := range options {
		option(&config)
	}
	return func(c *seriesConfig) {
		idx := miniIndex{
			data: data,
			kind: config.kind,
			name: config.name,
		}
		c.indices = append(c.indices, idx)
	}
}

// New Series constructor
// Optional
// - Index(): If no index is supplied, defaults to a single index of int64Values (0, 1, 2, ...n)
// - Name(): If no name is supplied, no name will appear when Series is printed
// - Kind(): Convert the Series values to the specified kind
// If passing []interface{}, must supply a type expectation for the Series.
// Options: Float, Int, String, Bool, DateTime
func New(data interface{}, options ...Option) (Series, error) {
	// Setup
	config := seriesConfig{}

	for _, option := range options {
		option(&config)
	}
	suppliedKind := config.kind
	var kind kinds.Kind
	name := config.name

	var factory values.Factory
	var v values.Values
	var idx index.Index
	var err error

	// Values
	switch reflect.ValueOf(data).Kind() {
	case reflect.Slice:
		factory, err = values.SliceFactory(data)

	default:
		return Series{}, fmt.Errorf("Unable to construct new Series: type not supported: %T", data)
	}

	// Sets values and kind based on the Values switch
	v = factory.V
	kind = factory.Kind
	if err != nil {
		return Series{}, fmt.Errorf("Unable to construct new Series: unable to construct values: %v", err)
	}

	// Optional kind conversion
	if suppliedKind != kinds.Invalid {
		v, err = values.Convert(v, suppliedKind)
		if err != nil {
			return Series{}, fmt.Errorf("Unable to construct new Series: %v", err)
		}
	}
	// Index
	// Default case: no client-supplied Index
	requiredLen := v.Len()
	if config.indices == nil {
		idx = index.Default(requiredLen)
	} else {
		idx, err = indexFromMiniIndex(config.indices, requiredLen)
		if err != nil {
			return Series{}, fmt.Errorf("Unable to construct new Series: %v", err)
		}
	}

	// Construct Series
	s := Series{
		index:  idx,
		values: v,
		kind:   kind,
		Name:   name,
	}

	return s, err
}

// [START MiniIndex]

// an untyped representation of one index level.
// It is used for unpacking client-supplied index data and optional metadata
type miniIndex struct {
	data interface{}
	kind kinds.Kind
	name string
}

// creates a full index from a mini client-supplied representation of an index level,
// but returns an error if every index level is not the same length as requiredLen

func indexFromMiniIndex(minis []miniIndex, requiredLen int) (index.Index, error) {
	var levels []index.Level
	for _, miniIdx := range minis {
		if reflect.ValueOf(miniIdx.data).Kind() != reflect.Slice {
			return index.Index{}, fmt.Errorf("Unable to construct index: custom index must be a Slice: unsupported index type: %T", miniIdx.data)
		}
		level, err := index.NewLevelFromSlice(miniIdx.data, miniIdx.name)
		if err != nil {
			return index.Index{}, fmt.Errorf("Unable to construct index: %v", err)
		}
		labelLen := level.Labels.Len()
		if labelLen != requiredLen {
			return index.Index{}, fmt.Errorf("Unable to construct index %v:"+
				"mismatch between supplied index length (%v) and expected length (%v)",
				miniIdx.data, labelLen, requiredLen)
		}
		if miniIdx.kind != kinds.Invalid {
			level.Convert(miniIdx.kind)
		}

		levels = append(levels, level)
	}
	idx := index.New(levels...)
	return idx, nil

}

// [END MiniIndex]
