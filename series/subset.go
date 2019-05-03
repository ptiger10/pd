package series

import "github.com/ptiger10/pd/internal/values"

// At subsets a Series by integer position
func (s Series) At(position int) Series {
	s = s.at(position)
	s.index.Refresh()
	return s
}

func (s Series) at(position int) Series {
	positions := []int{position}
	s.values = s.values.In(positions).(values.Values)
	for i, level := range s.index.Levels {
		s.index.Levels[i].Labels = level.Labels.In(positions).(values.Values)
	}
	return s
}

// an interface of valid (non-null) values; appropriate for type assertion
func (s Series) validVals() interface{} {
	valid := s.values.In(s.values.Valid())
	return valid.Vals()
}

// an interface slice of valid (non-null) values; appropriate for counting and indexing
func (s Series) validAll() []interface{} {
	valid := s.values.In(s.values.Valid())
	return valid.All()
}
