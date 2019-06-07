package series

import (
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
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
		s, _ := New([]string{"foo", "bar"}, Idx([]string{"A", "B"}), Idx([]int{1, 2}))
		newS, err := s.Insert(test.pos, test.val, test.idx)
		if err != nil {
			t.Errorf("s.Insert(): %v", err)
		}
		if !seriesEquals(newS, test.want) {
			t.Errorf("s.insert() returned %v, want %v", s, test.want)
		}
		if seriesEquals(newS, s) {
			t.Error("s.insert() maintained reference to original Series, want fresh copy")
		}
	}
}

func TestAppend(t *testing.T) {
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
		newS := s.Append(test.val, test.idx)
		if !seriesEquals(newS, test.want) {
			t.Errorf("s.Append() returned %v, want %v", s, test.want)
		}
		if seriesEquals(newS, s) {
			t.Error("s.Append() maintained reference to original Series, want fresh copy")
		}
	}
}

func TestDrop(t *testing.T) {
	var tests = []struct {
		pos  int
		want Series
	}{
		{0, mustNew([]string{"bar"}, Idx([]string{"B"}), Idx([]int{2}))},
		{1, mustNew([]string{"foo"}, Idx([]string{"A"}), Idx([]int{1}))},
	}
	for _, test := range tests {
		s, _ := New([]string{"foo", "bar"}, Idx([]string{"A", "B"}), Idx([]int{1, 2}))
		newS, err := s.Drop(test.pos)
		if err != nil {
			t.Errorf("s.Drop(): %v", err)
		}
		if !seriesEquals(newS, test.want) {
			t.Errorf("s.Drop() returned %v, want %v", s, test.want)
		}
		if seriesEquals(newS, s) {
			t.Error("s.Drop() maintained reference to original Series, want fresh copy")
		}
	}
}

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

func Test_InPlace_Sort(t *testing.T) {
	var tests = []struct {
		desc string
		orig Series
		asc  bool
		want Series
	}{
		{"float", mustNew([]float64{3, 5, 1}), true, mustNew([]float64{1, 3, 5}, Idx([]int{2, 0, 1}))},
		{"float reverse", mustNew([]float64{3, 5, 1}), false, mustNew([]float64{5, 3, 1}, Idx([]int{1, 0, 2}))},

		{"int", mustNew([]int{3, 5, 1}), true, mustNew([]int{1, 3, 5}, Idx([]int{2, 0, 1}))},
		{"int reverse", mustNew([]int{3, 5, 1}), false, mustNew([]int{5, 3, 1}, Idx([]int{1, 0, 2}))},

		{"string", mustNew([]string{"3", "5", "1"}), true, mustNew([]string{"1", "3", "5"}, Idx([]int{2, 0, 1}))},
		{"string reverse", mustNew([]string{"3", "5", "1"}), false, mustNew([]string{"5", "3", "1"}, Idx([]int{1, 0, 2}))},

		{"bool", mustNew([]bool{false, true, false}), true, mustNew([]bool{false, false, true}, Idx([]int{0, 2, 1}))},
		{"bool reverse", mustNew([]bool{false, true, false}), false, mustNew([]bool{true, false, false}, Idx([]int{1, 0, 2}))},

		{
			"datetime",
			mustNew([]time.Time{time.Date(2019, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2019, 3, 1, 0, 0, 0, 0, time.UTC), time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)}),
			true,
			mustNew([]time.Time{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2019, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2019, 3, 1, 0, 0, 0, 0, time.UTC)}, Idx([]int{2, 0, 1})),
		},
		{
			"datetime reverse",
			mustNew([]time.Time{time.Date(2019, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2019, 3, 1, 0, 0, 0, 0, time.UTC), time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)}),
			false,
			mustNew([]time.Time{time.Date(2019, 3, 1, 0, 0, 0, 0, time.UTC), time.Date(2019, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)}, Idx([]int{1, 0, 2})),
		},
	}
	for _, test := range tests {
		s := test.orig
		s.InPlace.Sort(test.asc)
		if !seriesEquals(s, test.want) {
			t.Errorf("series.Sort() test %v returned %v, want %v", test.desc, s, test.want)
		}
	}
}
