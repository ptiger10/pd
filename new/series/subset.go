package series

import (
	"reflect"

	"github.com/ptiger10/pd/new/internal/values"
)

type Element struct {
	Value interface{}
	Index []interface{}
	Kind  reflect.Kind
	Null  bool
}

func (s Series) Elem() Element {
	positions := []int{0}
	valElem := reflect.ValueOf(s.Values.In(positions)).Index(0)
	val := valElem.FieldByName("V").Interface()
	null := valElem.FieldByName("Null").Bool()

	var idx []interface{}
	for _, lvl := range s.Index.Levels {
		idxVal := reflect.ValueOf(lvl.Labels.In(positions)).Index(0).Interface()
		idx = append(idx, idxVal)
	}
	kind := s.Kind
	return Element{val, idx, kind, null}
}

func (s Series) At(position int) Series {
	return s.at(position)
}

func (s Series) at(position int) Series {
	positions := []int{position}
	s.Values = s.Values.In(positions).(values.Values)
	for i, level := range s.Index.Levels {
		s.Index.Levels[i].Labels = level.Labels.In(positions).(values.Values)
	}
	return s
}
