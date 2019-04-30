package index

import (
	"reflect"
)

type Index struct {
	Levels []Level
}

type Level struct {
	Kind     reflect.Kind
	Labels   Labels
	LabelMap map[string][]int
}

type Labels interface {
	In([]int) interface{}
}

type LabelMap map[string][]int
