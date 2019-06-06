package series

import (
	"fmt"
	"reflect"

	"github.com/ptiger10/pd/internal/config"
	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/kinds"
	"github.com/ptiger10/pd/opt"
)

// Idx returns a opt.ConstructorOption for use in the Series constructor New(),
// and takes an optional Name.
func Idx(data interface{}, options ...opt.ConstructorOption) opt.ConstructorOption {
	cfg := config.ConstructorConfig{}
	for _, option := range options {
		option(&cfg)
	}
	return func(c *config.ConstructorConfig) {
		idx := config.MiniIndex{
			Data: data,
			Kind: cfg.Kind,
			Name: cfg.Name,
		}
		c.Indices = append(c.Indices, idx)
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
// - Index(interface{}, ...opt.ConstructorOption): If no index is supplied, defaults to a single index of int64Values (0, 1, 2, ...n).
// Index() also accepts an optional Name() and Kind().
// If no name is supplied, index level will be unnamed.
// If no kind is supplied, Index will default to the reflect.Kind() of its data.
// If multiple Index() options are supplied, each will become its own Index level in a MultiIndex.
func New(data interface{}, options ...opt.ConstructorOption) (Series, error) {
	// Setup
	cfg := config.ConstructorConfig{}

	for _, option := range options {
		option(&cfg)
	}
	suppliedKind := cfg.Kind
	var kind kinds.Kind
	name := cfg.Name

	var factory values.Factory
	var v values.Values
	var idx index.Index
	var err error

	// Values
	switch reflect.ValueOf(data).Kind() {
	case reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.String,
		reflect.Bool,
		reflect.Struct:
		factory, err = values.ScalarFactory(data)

	case reflect.Slice:
		factory, err = values.SliceFactory(data)

	default:
		return Series{}, fmt.Errorf("unable to construct new Series: type not supported: %T", data)
	}

	// Sets values and kind based on the Values switch
	v = factory.Values
	kind = factory.Kind
	if err != nil {
		return Series{}, fmt.Errorf("unable to construct new Series: unable to construct values: %v", err)
	}

	// opt.ConstructorOptional kind conversion
	if suppliedKind != kinds.None {
		v, err = values.Convert(v, suppliedKind)
		if err != nil {
			return Series{}, fmt.Errorf("unable to construct new Series: %v", err)
		}
		kind = suppliedKind
	}
	// Index
	// Default case: no client-supplied Index
	requiredLen := v.Len()
	if cfg.Indices == nil {
		idx = index.Default(requiredLen)
	} else {
		// one or more client-supplied indices
		idx, err = indexFromMiniIndex(cfg.Indices, requiredLen)
		if err != nil {
			return Series{}, fmt.Errorf("unable to construct new Series: %v", err)
		}
	}

	// Construct Series
	s := new(idx, v, kind, name)
	s.Math = Math{s: &s}
	s.To = To{s: &s}
	s.Index = Index{s: &s, To: To{s: &s, idx: true}}
	s.Select = Select{s: &s}
	// s.IndexTo = IndexTo{s: &s}
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

// creates a full index from a mini client-supplied representation of an index level,
// but returns an error if every index level is not the same length as requiredLen

func indexFromMiniIndex(minis []config.MiniIndex, requiredLen int) (index.Index, error) {
	var levels []index.Level
	for _, miniIdx := range minis {
		level, err := index.NewLevel(miniIdx.Data, miniIdx.Name)
		if err != nil {
			return index.Index{}, fmt.Errorf("unable to construct index: %v", err)
		}
		labelLen := level.Labels.Len()
		if labelLen != requiredLen {
			return index.Index{}, fmt.Errorf("unable to construct index %v:"+
				"mismatch between supplied index length (%v) and expected length (%v)",
				miniIdx.Data, labelLen, requiredLen)
		}
		if miniIdx.Kind != kinds.None {
			level.Convert(miniIdx.Kind)
		}

		levels = append(levels, level)
	}
	idx := index.New(levels...)
	return idx, nil

}

// [END MiniIndex]
