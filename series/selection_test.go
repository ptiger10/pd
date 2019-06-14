package series

import (
	"testing"
)

// func TestMain(m *testing.M) {
// 	options.SetLogWarnings(false)
// 	exitCode := m.Run()
// 	options.RestoreDefaults()
// 	os.Exit(exitCode)
// }

func TestAt_singleIdx(t *testing.T) {
	var tests = []struct {
		input int
		want  interface{}
	}{
		{0, "hot"},
		{1, "dog"},
		{2, "log"},
	}
	s, _ := New([]string{"hot", "dog", "log"}, Idx([]int{10, 100, 1000}))
	for _, test := range tests {
		got, err := s.At(test.input)
		if err != nil {
			t.Errorf("Returned %v want nil", err)
		}
		if got != test.want {
			t.Errorf("Returned Series %v, want %v", got, test.want)
		}
	}
}

func TestAt_singleIndex_multipleCalls(t *testing.T) {
	s, _ := New([]string{"hot", "dog", "log"}, Idx([]int{10, 100, 1000}))
	origS, _ := New([]string{"hot", "dog", "log"}, Idx([]int{10, 100, 1000}))
	_, err := s.At(0)
	if err != nil {
		t.Errorf("Returned %v want nil", err)
	}
	_, err = s.At(1)
	if err != nil {
		t.Errorf("Returned %v want nil", err)
	}
	if !seriesEquals(s, origS) {
		t.Errorf("s.At() modified s, want fresh copy")
	}
}

func TestSelect_Drop(t *testing.T) {
	s, _ := New(
		[]string{"foo", "bar", "baz"},
		IndexLevel{Labels: []string{"qux", "quux", "corge"}, Name: "1"},
		IndexLevel{Labels: []string{"A", "B", "C"}, Name: "2"},
	)
	origS, _ := New(
		[]string{"foo", "bar", "baz"},
		IndexLevel{Labels: []string{"qux", "quux", "corge"}, Name: "1"},
		IndexLevel{Labels: []string{"A", "B", "C"}, Name: "2"},
	)
	var tests = []struct {
		desc      string
		selection Selection
		want      *Series
	}{
		{"ByRows", s.SelectRows([]int{0, 2}), mustNew([]string{"bar"}, IndexLevel{Labels: []string{"quux"}, Name: "1"}, IndexLevel{Labels: []string{"B"}, Name: "2"})},
		{"ByLabels", s.SelectLabels([]string{"qux", "corge"}), mustNew([]string{"bar"}, IndexLevel{Labels: []string{"quux"}, Name: "1"}, IndexLevel{Labels: []string{"B"}, Name: "2"})},
		{"ByLevels", s.SelectLevels([]int{0}), mustNew([]string{"foo", "bar", "baz"}, IndexLevel{Labels: []string{"A", "B", "C"}, Name: "2"})},
		{"ByLevelNames", s.SelectLevelNames([]string{"1"}), mustNew([]string{"foo", "bar", "baz"}, IndexLevel{Labels: []string{"A", "B", "C"}, Name: "2"})},
	}

	for _, test := range tests {
		newS, err := test.selection.Drop()
		if err != nil {
			t.Errorf("Select.%v.Drop(): %v", test.desc, err)
		}
		if !seriesEquals(newS, test.want) {
			t.Errorf("Select.%v.Drop() returned %v, want %v", test.desc, newS, test.want)
		}
		if !seriesEquals(s, origS) {
			t.Errorf("Select.%v.Drop() modified s, want fresh copy", test.desc)
		}
	}
}

func TestSelect_Rows_Swap(t *testing.T) {
	s, _ := New([]string{"hot", "dog", "log"}, Idx([]int{10, 100, 1000}))
	newS, err := s.SelectRows([]int{0, 2}).Swap()
	if err != nil {
		t.Errorf("SelectRows.Swap(): %v", err)
	}
	wantS, _ := New([]string{"log", "dog", "hot"}, Idx([]int{1000, 100, 10}))
	if !seriesEquals(newS, wantS) {
		t.Errorf("SelectRows.Swap() returned %v, want %v", newS, wantS)
	}
}

// // func TestAtLabel_singleIdx(t *testing.T) {
// // 	var tests = []struct {
// // 		input string
// // 		want  Series
// // 	}{
// // 		{"10", mustNew([]string{"hot"}, Idx([]int{10}))},
// // 		{"100", mustNew([]string{"dog"}, Idx([]int{100}))},
// // 		{"1000", mustNew([]string{"log"}, Idx([]int{1000}))},
// // 	}
// // 	for _, test := range tests {
// // 		s, _ := New([]string{"hot", "dog", "log"}, Idx([]int{10, 100, 1000}))
// // 		got := s.AtLabel(test.input)
// // 		if !seriesEquals(got, test.want) {
// // 			t.Errorf("Returned Series %v, want %v", got, test.want)
// // 		}
// // 	}
// // }

// func TestAt_fail(t *testing.T) {
// 	s, _ := New([]string{"hot", "dog", "log"}, Idx([]int{10, 100, 1000}))
// 	_, err := s.At(3)
// 	if err == nil {
// 		t.Error("Returned nil err, want out-of-range fail", err)
// 	}
// }

// // func TestAtLabel_fail(t *testing.T) {
// // 	s, _ := New([]string{"hot", "dog", "log"}, Idx([]int{10, 100, 1000}))
// // 	got := s.AtLabel("NotPresent")
// // 	if !seriesEquals(got, s) {
// // 		t.Errorf("Returned %v, want original series", got)
// // 	}
// // }

// func TestSelect_RowsmultiIdx(t *testing.T) {
// 	var tests = []struct {
// 		input int
// 		want  interface{}
// 	}{
// 		{0, mustNew([]string{"hot"}, Idx([]int{10}), Idx([]string{"A"}))},
// 		{1, mustNew([]string{"dog"}, Idx([]int{100}), Idx([]string{"B"}))},
// 		{2, mustNew([]string{"log"}, Idx([]int{1000}), Idx([]string{"C"}))},
// 	}
// 	for _, test := range tests {
// 		s, _ := New(
// 			[]string{"hot", "dog", "log"},
// 			Idx([]int{10, 100, 1000}),
// 			Idx([]string{"A", "B", "C"}),
// 		)
// 		got := s.At(test.input)
// 		if !seriesEquals(got, test.want) {
// 			t.Errorf("Returned Series %v, want %v", got, test.want)
// 		}
// 	}
// }

// // func TestAtLabel_multiIdx(t *testing.T) {
// // 	var tests = []struct {
// // 		input string
// // 		want  Series
// // 	}{
// // 		{"10", mustNew([]string{"hot"}, Idx([]int{10}), Idx([]string{"A"}))},
// // 		{"100", mustNew([]string{"dog"}, Idx([]int{100}), Idx([]string{"B"}))},
// // 		{"1000", mustNew([]string{"log"}, Idx([]int{1000}), Idx([]string{"C"}))},
// // 	}
// // 	for _, test := range tests {
// // 		s, _ := New(
// // 			[]string{"hot", "dog", "log"},
// // 			Idx([]int{10, 100, 1000}),
// // 			Idx([]string{"A", "B", "C"}),
// // 		)
// // 		got := s.AtLabel(test.input)
// // 		if !seriesEquals(got, test.want) {
// // 			t.Errorf("Returned Series %v, want %v", got, test.want)
// // 		}
// // 	}
// // }

// func TestSelect_Failures(t *testing.T) {
// 	var tests = []struct {
// 		options  []options.SelectionOption
// 		errorMsg string
// 	}{
// 		{[]options.SelectionOption{options.ByRows([]int{0, 3})}, "row indexing out range"},
// 		{[]options.SelectionOption{options.ByLabels([]string{"0", "3"})}, "label indexing out range"},
// 		{[]options.SelectionOption{options.ByLevels([]int{0, 2})}, "index level out range"},
// 		{[]options.SelectionOption{options.ByLevelNames([]string{"foo", "baz"})}, "index name not in index"},
// 		{[]options.SelectionOption{options.ByLabels([]string{"0"}), options.ByRows([]int{0})}, "multiple row selectors supplied"},
// 		{[]options.SelectionOption{options.ByLevelNames([]string{"foo"}), options.ByLevels([]int{0})}, "multiple index level selectors supplied"},
// 		{[]options.SelectionOption{options.ByLevels([]int{0, 1}), options.ByLabels([]string{"0"})}, "multiple levels plus labels supplied."},
// 	}
// 	for _, test := range tests {
// 		s, _ := New([]int{1, 5, 10}, Idx([]int{0, 1, 2}, options.Name("foo")), Idx([]string{"A", "B", "C"}, options.Name("bar")))
// 		sel := s.Select(test.options...)
// 		if !seriesEquals(sel.s, s) {
// 			t.Errorf("Select() returned %v, want return of underlying series due to %v", sel.s, test.errorMsg)
// 		}
// 	}
// }

// func TestSelect_Get(t *testing.T) {
// 	var tests = []struct {
// 		options    []options.SelectionOption
// 		wantSeries Series
// 		desc       string
// 	}{
// 		{
// 			[]options.SelectionOption{options.ByRows([]int{0})},
// 			mustNew([]int{0}, Idx([]int{0}, options.Name("foo")), Idx([]string{"A"}, options.Name("bar"))),
// 			"one row only",
// 		},
// 		{
// 			[]options.SelectionOption{options.ByRows([]int{0, 2})},
// 			mustNew([]int{0, 2}, Idx([]int{0, 2}, options.Name("foo")), Idx([]string{"A", "C"}, options.Name("bar"))),
// 			"multiple rows",
// 		},
// 		{
// 			[]options.SelectionOption{options.ByLabels([]string{"0"})},
// 			mustNew([]int{0}, Idx([]int{0}, options.Name("foo")), Idx([]string{"A"}, options.Name("bar"))),
// 			"one label only",
// 		},
// 		{
// 			[]options.SelectionOption{options.ByLabels([]string{"0", "2"})},
// 			mustNew([]int{0, 2}, Idx([]int{0, 2}, options.Name("foo")), Idx([]string{"A", "C"}, options.Name("bar"))),
// 			"multiple labels",
// 		},
// 		{
// 			[]options.SelectionOption{options.ByLevels([]int{0})},
// 			mustNew([]int{0, 1, 2}, Idx([]int{0, 1, 2}, options.Name("foo"))),
// 			"one index level only",
// 		},
// 		{
// 			[]options.SelectionOption{options.ByLevels([]int{0, 1})},
// 			mustNew([]int{0, 1, 2}, Idx([]int{0, 1, 2}, options.Name("foo")), Idx([]string{"A", "B", "C"}, options.Name("bar"))),
// 			"multiple index levels",
// 		},
// 		{
// 			[]options.SelectionOption{options.ByLevelNames([]string{"foo"})},
// 			mustNew([]int{0, 1, 2}, Idx([]int{0, 1, 2}, options.Name("foo"))),
// 			"index name only",
// 		},
// 		{
// 			[]options.SelectionOption{options.ByLevelNames([]string{"foo", "bar"})},
// 			mustNew([]int{0, 1, 2}, Idx([]int{0, 1, 2}, options.Name("foo")), Idx([]string{"A", "B", "C"}, options.Name("bar"))),
// 			"multiple index names",
// 		},
// 		{
// 			[]options.SelectionOption{options.ByRows([]int{0, 1}), options.ByLevels([]int{0})},
// 			mustNew([]int{0, 1}, Idx([]int{0, 1}, options.Name("foo"))),
// 			"rows and one index level",
// 		},
// 		{
// 			[]options.SelectionOption{options.ByRows([]int{0, 1}), options.ByLevels([]int{0, 1})},
// 			mustNew([]int{0, 1}, Idx([]int{0, 1}, options.Name("foo")), Idx([]string{"A", "B"}, options.Name("bar"))),
// 			"rows and multiple index levels",
// 		},
// 		{
// 			[]options.SelectionOption{options.ByRows([]int{0, 1}), options.ByLevelNames([]string{"foo"})},
// 			mustNew([]int{0, 1}, Idx([]int{0, 1}, options.Name("foo"))),
// 			"rows and one index name",
// 		},
// 		{
// 			[]options.SelectionOption{options.ByRows([]int{0, 1}), options.ByLevelNames([]string{"foo", "bar"})},
// 			mustNew([]int{0, 1}, Idx([]int{0, 1}, options.Name("foo")), Idx([]string{"A", "B"}, options.Name("bar"))),
// 			"rows and multiple index names",
// 		},
// 	}
// 	for _, test := range tests {
// 		s, _ := New([]int{0, 1, 2}, Idx([]int{0, 1, 2}, options.Name("foo")), Idx([]string{"A", "B", "C"}, options.Name("bar")))
// 		sel := s.Select(test.options...)
// 		newS, err := sel.Get()
// 		if err != nil {
// 			t.Errorf("selection.Get() returned %v when selecting %s", err, test.desc)
// 		}
// 		if !seriesEquals(newS, test.wantSeries) {
// 			t.Errorf("selection.Get() returned \n%v when selecting %s, want \n%v", newS, test.desc, test.wantSeries)
// 		}
// 	}
// }

// func TestSelect_Swap(t *testing.T) {
// 	var tests = []struct {
// 		options    []options.SelectionOption
// 		wantSeries Series
// 		desc       string
// 	}{
// 		{
// 			[]options.SelectionOption{options.ByRows([]int{0, 1})},
// 			mustNew([]int{1, 0, 2}, Idx([]int{1, 0, 2}, options.Name("foo")), Idx([]string{"B", "A", "C"}, options.Name("bar"))),
// 			"swap two rows by position",
// 		},
// 		{
// 			[]options.SelectionOption{options.ByRows([]int{1, 0})},
// 			mustNew([]int{1, 0, 2}, Idx([]int{1, 0, 2}, options.Name("foo")), Idx([]string{"B", "A", "C"}, options.Name("bar"))),
// 			"swap two rows by position - reverse arguments",
// 		},
// 		{
// 			[]options.SelectionOption{options.ByLabels([]string{"0", "1"})},
// 			mustNew([]int{1, 0, 2}, Idx([]int{1, 0, 2}, options.Name("foo")), Idx([]string{"B", "A", "C"}, options.Name("bar"))),
// 			"swap two rows by label",
// 		},
// 		{
// 			[]options.SelectionOption{options.ByLabels([]string{"1", "0"})},
// 			mustNew([]int{1, 0, 2}, Idx([]int{1, 0, 2}, options.Name("foo")), Idx([]string{"B", "A", "C"}, options.Name("bar"))),
// 			"swap two rows by label - reverse arguments",
// 		},
// 		{
// 			[]options.SelectionOption{options.ByLevels([]int{0, 1})},
// 			mustNew([]int{0, 1, 2}, Idx([]string{"A", "B", "C"}, options.Name("bar")), Idx([]int{0, 1, 2}, options.Name("foo"))),
// 			"swap two levels by position",
// 		},
// 		{
// 			[]options.SelectionOption{options.ByLevels([]int{1, 0})},
// 			mustNew([]int{0, 1, 2}, Idx([]string{"A", "B", "C"}, options.Name("bar")), Idx([]int{0, 1, 2}, options.Name("foo"))),
// 			"swap two levels by position - reverse arguments",
// 		},
// 		{
// 			[]options.SelectionOption{options.ByLevelNames([]string{"foo", "bar"})},
// 			mustNew([]int{0, 1, 2}, Idx([]string{"A", "B", "C"}, options.Name("bar")), Idx([]int{0, 1, 2}, options.Name("foo"))),
// 			"swap two levels by name",
// 		},
// 		{
// 			[]options.SelectionOption{options.ByLevelNames([]string{"bar", "foo"})},
// 			mustNew([]int{0, 1, 2}, Idx([]string{"A", "B", "C"}, options.Name("bar")), Idx([]int{0, 1, 2}, options.Name("foo"))),
// 			"swap two levels by name - reverse arguments",
// 		},
// 	}
// 	for _, test := range tests {
// 		s, _ := New([]int{0, 1, 2}, Idx([]int{0, 1, 2}, options.Name("foo")), Idx([]string{"A", "B", "C"}, options.Name("bar")))
// 		origS, _ := New([]int{0, 1, 2}, Idx([]int{0, 1, 2}, options.Name("foo")), Idx([]string{"A", "B", "C"}, options.Name("bar")))
// 		sel := s.Select(test.options...)
// 		newS, err := sel.Swap()
// 		if err != nil {
// 			t.Errorf("selection.Swap() returned %v when selecting %s", err, test.desc)
// 		}
// 		if !seriesEquals(newS, test.wantSeries) {
// 			t.Errorf("selection.Swap() returned \n%v when selecting %s, want \n%v", newS.index, test.desc, test.wantSeries.index)
// 		}
// 		if !seriesEquals(origS, s) {
// 			t.Errorf("selection.Swap() modifying Series in place instead of returning new")
// 		}
// 	}
// }

// func TestSelect_Set(t *testing.T) {
// 	s, _ := New([]int{0, 1, 2}, Idx([]int{0, 1, 2}, options.Name("foo")), Idx([]string{"A", "B", "C"}, options.Name("bar")))
// 	origS, _ := New([]int{0, 1, 2}, Idx([]int{0, 1, 2}, options.Name("foo")), Idx([]string{"A", "B", "C"}, options.Name("bar")))
// 	sel := s.Select(options.ByRows([]int{0, 1, 2}))
// 	newS, err := sel.Set(5)
// 	if err != nil {
// 		t.Errorf("selection.Set() %v", err)
// 	}
// 	if !seriesEquals(origS, s) {
// 		t.Errorf("selection.Set() modifying Series in place instead of returning new")
// 	}
// 	if newS.Element(0).Value != int64(5) || newS.Element(1).Value != int64(5) || newS.Element(2).Value != int64(5) {
// 		t.Errorf("selection.Set() did not set all values")
// 	}
// }
