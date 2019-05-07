package series

import (
	"fmt"
	"reflect"

	"github.com/ptiger10/pd/kinds"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
)

// A ConstructorOption is an optional parameter in the Series constructor.
type ConstructorOption func(*constructorConfig)
type constructorConfig struct {
	indices []miniIndex
	kind    kinds.Kind
	name    string
}

// Kind will convert either values or an index level to the specified kind
func Kind(kind kinds.Kind) ConstructorOption {
	return func(c *constructorConfig) {
		c.kind = kind
	}
}

// Name will name either values or an index level
func Name(n string) ConstructorOption {
	return func(c *constructorConfig) {
		c.name = n
	}

}

// Index returns a ConstructorOption for use in the Series constructor New(),
// and takes an optioanl Name.
func Index(data interface{}, options ...ConstructorOption) ConstructorOption {
	config := constructorConfig{}
	for _, option := range options {
		option(&config)
	}
	return func(c *constructorConfig) {
		idx := miniIndex{
			data: data,
			kind: config.kind,
			name: config.name,
		}
		c.indices = append(c.indices, idx)
	}
}

// New Series constructor
//
// Optional:
//
// - Name(string): If no name is supplied, no name will appear when Series is printed.
// If multiple Name() options are supplied, only the final will be used.
//
// - Kind(kinds.Kind): Convert the Series values to the specified kind. Kind options: Float, Int, String, Bool, DateTime, Interface.
// If multiple Kind() options are supplied, only the final will be used.
//
// - Index(interface{}, ...ConstructorOption): If no index is supplied, defaults to a single index of int64Values (0, 1, 2, ...n).
// Index() also accepts an optional Name() and Kind().
// If no name is supplied, index level will be unnamed.
// If no kind is supplied, Index will default to the reflect.Kind() of its data.
// If multiple Index() options are supplied, each will become its own Index level in a MultiIndex.
func New(data interface{}, options ...ConstructorOption) (Series, error) {
	// Setup
	config := constructorConfig{}

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

	// ConstructorOptional kind conversion
	if suppliedKind != kinds.None {
		v, err = values.Convert(v, suppliedKind)
		if err != nil {
			return Series{}, fmt.Errorf("Unable to construct new Series: %v", err)
		}
		kind = suppliedKind
	}
	// Index
	// Default case: no client-supplied Index
	requiredLen := v.Len()
	if config.indices == nil {
		idx = index.Default(requiredLen)
	} else {
		// one or more client-supplied indices
		idx, err = indexFromMiniIndex(config.indices, requiredLen)
		if err != nil {
			return Series{}, fmt.Errorf("Unable to construct new Series: %v", err)
		}
	}

	// Construct Series
	s := new(idx, v, kind, name)
	return s, err
}

func new(idx index.Index, values values.Values, kind kinds.Kind, name string) Series {
	return Series{
		index:  idx,
		values: values,
		kind:   kind,
		Name:   name,
	}
}

// [START MiniIndex]

// An untyped representation of one index level.
// It is used for unpacking client-supplied index data and Constructoroptional metadata.
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
		if miniIdx.kind != kinds.None {
			level.Convert(miniIdx.kind)
		}

		levels = append(levels, level)
	}
	idx := index.New(levels...)
	return idx, nil

}

// [END MiniIndex]
