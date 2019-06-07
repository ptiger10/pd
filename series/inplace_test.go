package series

import "testing"

func TestInsertInPlace(t *testing.T) {
	var tests = []struct {
		pos  int
		val  interface{}
		idx  []interface{}
		want Series
	}{
		{0, "baz", []interface{}{"C", 3},
			mustNew([]string{"baz", "foo", "bar"}, Idx([]string{"C", "A", "B"}), Idx([]int{3, 1, 2}))},
		{1, "baz", []interface{}{"C", 3},
			mustNew([]string{"foo", "baz", "bar"}, Idx([]string{"A", "C", "B"}), Idx([]int{1, 3, 2}))},
		{2, "baz", []interface{}{"C", 3},
			mustNew([]string{"foo", "bar", "baz"}, Idx([]string{"A", "B", "C"}), Idx([]int{1, 2, 3}))},
	}
	for _, test := range tests {
		s, err := New([]string{"foo", "bar"}, Idx([]string{"A", "B"}), Idx([]int{1, 2}))
		if err != nil {
			t.Error(err)
		}
		s.InPlace.Insert(test.pos, test.val, test.idx)
		if !seriesEquals(s, test.want) {
			t.Errorf("s.insert() returned %v, want %v", s, test.want)
		}
	}
}

func TestAppendInPlace(t *testing.T) {
	var tests = []struct {
		val  interface{}
		idx  []interface{}
		want Series
	}{
		{"baz", []interface{}{"C", 3},
			mustNew([]string{"foo", "bar", "baz"}, Idx([]string{"A", "B", "C"}), Idx([]int{1, 2, 3}))},
	}
	for _, test := range tests {
		s, _ := New([]string{"foo", "bar"}, Idx([]string{"A", "B"}), Idx([]int{1, 2}))
		s.InPlace.Append(test.val, test.idx)
		if !seriesEquals(s, test.want) {
			t.Errorf("s.Append() returned %v, want %v", s, test.want)
		}
	}
}

func TestDropInPlace(t *testing.T) {
	var tests = []struct {
		pos  int
		want Series
	}{
		{0, mustNew([]string{"bar"}, Idx([]string{"B"}), Idx([]int{2}))},
		{1, mustNew([]string{"foo"}, Idx([]string{"A"}), Idx([]int{1}))},
	}
	for _, test := range tests {
		s, _ := New([]string{"foo", "bar"}, Idx([]string{"A", "B"}), Idx([]int{1, 2}))
		s.InPlace.Drop(test.pos)
		if !seriesEquals(s, test.want) {
			t.Errorf("s.insert() returned %v, want %v", s, test.want)
		}
	}
}
