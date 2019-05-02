package series

import (
	"github.com/ptiger10/pd/new/internal/values"
)

// At subsets a Series by integer position
func (s Series) At(position int) Series {
	s = s.at(position)
	s.Index.Refresh()
	return s
}

func (s Series) at(position int) Series {
	positions := []int{position}
	s.Values = s.Values.In(positions).(values.Values)
	for i, level := range s.Index.Levels {
		s.Index.Levels[i].Labels = level.Labels.In(positions).(values.Values)
	}
	return s
}
