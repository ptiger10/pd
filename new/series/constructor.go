package series

import (
	"fmt"
	"log"
	"reflect"

	"github.com/ptiger10/pd/new/kinds"

	"github.com/ptiger10/pd/new/internal/index"
	constructIdx "github.com/ptiger10/pd/new/internal/index/constructors"
	"github.com/ptiger10/pd/new/internal/values"
	constructVal "github.com/ptiger10/pd/new/internal/values/constructors"
)

type Option func(*seriesConfig)
type seriesConfig struct {
	indices []index.MiniIndex
	kind    reflect.Kind
	name    string
}

func Kind(t reflect.Kind) Option {
	return func(c *seriesConfig) {
		c.kind = t
	}
}

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
		idx := index.MiniIndex{
			Data: data,
			Kind: config.kind,
			Name: config.name,
		}
		c.indices = append(c.indices, idx)
	}
}

// Calls New and panics if error. For use in testing
func mustNew(data interface{}, options ...Option) Series {
	s, err := New(data, options...)
	if err != nil {
		log.Panic(err)
	}
	return s
}

// New Series constructor
// Optional
// - Index(): If no index is supplied, defaults to a single index of IntValues (0, 1, 2, ...n)
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
	kind := config.kind
	name := config.name

	var v values.Values
	var idx index.Index
	var err error

	// Values
	switch reflect.ValueOf(data).Kind() {
	case reflect.Slice:
		v, kind, err = constructVal.ValuesFromSlice(data)

	default:
		return Series{}, fmt.Errorf("Unable to construct new Series: type not supported: %T", data)
	}

	// Optional client-provided conversion
	if kind != kinds.None {
		v, err = constructVal.Convert(v, kind)
		if err != nil {
			return Series{}, fmt.Errorf("Unable to construct new Series: %v", err)
		}
	}
	// Index
	// Default case: no client-supplied Index
	if config.indices == nil {
		idx = constructIdx.Default(v.Len())
	} else {
		idx, err = constructIdx.IndexFromMiniIndex(config.indices)
		if err != nil {
			return Series{}, fmt.Errorf("Unable to construct new Series: %v", err)
		}
	}

	// Construct Series
	s := Series{
		Index:  idx,
		Values: v,
		Kind:   kind,
		Name:   name,
	}

	return s, err
}

func (s Series) Set(subset Series, data interface{}) {

}
