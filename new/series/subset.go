package series

import (
	"reflect"

	"github.com/ptiger10/pd/new/internal/index"

	"github.com/ptiger10/pd/new/internal/values"
)

type Element struct {
	Value interface{}
	Index []interface{}
	Kind  reflect.Kind
	Null  bool
}

func (s Series) Elem() Element {
	valElem := reflect.ValueOf(s.Values.In([]int{0})).Elem()
	val := valElem.FieldByName("v")
	null := valElem.FieldByName("null").Bool()

	var idx []interface{}
	for _, lvl := range s.Index.Levels {
		idx = append(idx, lvl.Labels.In([]int{0}))
	}
	kind := s.Kind
	return Element{
		val,
		idx,
		kind,
		null,
	}
}

func (s Series) At(position int) Series {
	slicer := []int{position}
	copy := Series(s)
	copy.Values = s.Values.In(slicer).(values.Values)
	for i, level := range copy.Index.Levels {
		copy.Index.Levels[i].Labels = level.Labels.In(slicer).(index.Labels)
	}
	return copy
}
