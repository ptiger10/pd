package series

import "testing"

// func TestMain(m *testing.M) {
// 	opt.SetLogWarnings(false)
// 	exitCode := m.Run()
// 	opt.RestoreDefaults()
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
	sOrig, _ := New([]string{"hot", "dog", "log"}, Idx([]int{10, 100, 1000}))
	_, err := s.At(0)
	if err != nil {
		t.Errorf("Returned %v want nil", err)
	}
	_, err = s.At(1)
	if err != nil {
		t.Errorf("Returned %v want nil", err)
	}
	if !seriesEquals(s, sOrig) {
		t.Errorf("s.At() modified s, want fresh copy")
	}
}

func TestSelect_Rows_Drop(t *testing.T) {
	s, _ := New([]string{"hot", "dog", "log"}, Idx([]int{10, 100, 1000}))
	newS, err := s.Select.Rows([]int{0, 2}).Drop()
	if err != nil {
		t.Errorf("Select.Rows.Drop(): %v", err)
	}
	wantS, _ := New([]string{"dog"}, Idx([]int{100}))
	if !seriesEquals(newS, wantS) {
		t.Errorf("Select.Rows.Drop() returned %v, want %v", newS, wantS)
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
// 		options  []opt.SelectionOption
// 		errorMsg string
// 	}{
// 		{[]opt.SelectionOption{opt.ByRows([]int{0, 3})}, "row indexing out range"},
// 		{[]opt.SelectionOption{opt.ByLabels([]string{"0", "3"})}, "label indexing out range"},
// 		{[]opt.SelectionOption{opt.ByLevels([]int{0, 2})}, "index level out range"},
// 		{[]opt.SelectionOption{opt.ByLevelNames([]string{"foo", "baz"})}, "index name not in index"},
// 		{[]opt.SelectionOption{opt.ByLabels([]string{"0"}), opt.ByRows([]int{0})}, "multiple row selectors supplied"},
// 		{[]opt.SelectionOption{opt.ByLevelNames([]string{"foo"}), opt.ByLevels([]int{0})}, "multiple index level selectors supplied"},
// 		{[]opt.SelectionOption{opt.ByLevels([]int{0, 1}), opt.ByLabels([]string{"0"})}, "multiple levels plus labels supplied."},
// 	}
// 	for _, test := range tests {
// 		s, _ := New([]int{1, 5, 10}, Idx([]int{0, 1, 2}, opt.Name("foo")), Idx([]string{"A", "B", "C"}, opt.Name("bar")))
// 		sel := s.Select(test.options...)
// 		if !seriesEquals(sel.s, s) {
// 			t.Errorf("Select() returned %v, want return of underlying series due to %v", sel.s, test.errorMsg)
// 		}
// 	}
// }

// func TestSelect_Get(t *testing.T) {
// 	var tests = []struct {
// 		options    []opt.SelectionOption
// 		wantSeries Series
// 		desc       string
// 	}{
// 		{
// 			[]opt.SelectionOption{opt.ByRows([]int{0})},
// 			mustNew([]int{0}, Idx([]int{0}, opt.Name("foo")), Idx([]string{"A"}, opt.Name("bar"))),
// 			"one row only",
// 		},
// 		{
// 			[]opt.SelectionOption{opt.ByRows([]int{0, 2})},
// 			mustNew([]int{0, 2}, Idx([]int{0, 2}, opt.Name("foo")), Idx([]string{"A", "C"}, opt.Name("bar"))),
// 			"multiple rows",
// 		},
// 		{
// 			[]opt.SelectionOption{opt.ByLabels([]string{"0"})},
// 			mustNew([]int{0}, Idx([]int{0}, opt.Name("foo")), Idx([]string{"A"}, opt.Name("bar"))),
// 			"one label only",
// 		},
// 		{
// 			[]opt.SelectionOption{opt.ByLabels([]string{"0", "2"})},
// 			mustNew([]int{0, 2}, Idx([]int{0, 2}, opt.Name("foo")), Idx([]string{"A", "C"}, opt.Name("bar"))),
// 			"multiple labels",
// 		},
// 		{
// 			[]opt.SelectionOption{opt.ByLevels([]int{0})},
// 			mustNew([]int{0, 1, 2}, Idx([]int{0, 1, 2}, opt.Name("foo"))),
// 			"one index level only",
// 		},
// 		{
// 			[]opt.SelectionOption{opt.ByLevels([]int{0, 1})},
// 			mustNew([]int{0, 1, 2}, Idx([]int{0, 1, 2}, opt.Name("foo")), Idx([]string{"A", "B", "C"}, opt.Name("bar"))),
// 			"multiple index levels",
// 		},
// 		{
// 			[]opt.SelectionOption{opt.ByLevelNames([]string{"foo"})},
// 			mustNew([]int{0, 1, 2}, Idx([]int{0, 1, 2}, opt.Name("foo"))),
// 			"index name only",
// 		},
// 		{
// 			[]opt.SelectionOption{opt.ByLevelNames([]string{"foo", "bar"})},
// 			mustNew([]int{0, 1, 2}, Idx([]int{0, 1, 2}, opt.Name("foo")), Idx([]string{"A", "B", "C"}, opt.Name("bar"))),
// 			"multiple index names",
// 		},
// 		{
// 			[]opt.SelectionOption{opt.ByRows([]int{0, 1}), opt.ByLevels([]int{0})},
// 			mustNew([]int{0, 1}, Idx([]int{0, 1}, opt.Name("foo"))),
// 			"rows and one index level",
// 		},
// 		{
// 			[]opt.SelectionOption{opt.ByRows([]int{0, 1}), opt.ByLevels([]int{0, 1})},
// 			mustNew([]int{0, 1}, Idx([]int{0, 1}, opt.Name("foo")), Idx([]string{"A", "B"}, opt.Name("bar"))),
// 			"rows and multiple index levels",
// 		},
// 		{
// 			[]opt.SelectionOption{opt.ByRows([]int{0, 1}), opt.ByLevelNames([]string{"foo"})},
// 			mustNew([]int{0, 1}, Idx([]int{0, 1}, opt.Name("foo"))),
// 			"rows and one index name",
// 		},
// 		{
// 			[]opt.SelectionOption{opt.ByRows([]int{0, 1}), opt.ByLevelNames([]string{"foo", "bar"})},
// 			mustNew([]int{0, 1}, Idx([]int{0, 1}, opt.Name("foo")), Idx([]string{"A", "B"}, opt.Name("bar"))),
// 			"rows and multiple index names",
// 		},
// 	}
// 	for _, test := range tests {
// 		s, _ := New([]int{0, 1, 2}, Idx([]int{0, 1, 2}, opt.Name("foo")), Idx([]string{"A", "B", "C"}, opt.Name("bar")))
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
// 		options    []opt.SelectionOption
// 		wantSeries Series
// 		desc       string
// 	}{
// 		{
// 			[]opt.SelectionOption{opt.ByRows([]int{0, 1})},
// 			mustNew([]int{1, 0, 2}, Idx([]int{1, 0, 2}, opt.Name("foo")), Idx([]string{"B", "A", "C"}, opt.Name("bar"))),
// 			"swap two rows by position",
// 		},
// 		{
// 			[]opt.SelectionOption{opt.ByRows([]int{1, 0})},
// 			mustNew([]int{1, 0, 2}, Idx([]int{1, 0, 2}, opt.Name("foo")), Idx([]string{"B", "A", "C"}, opt.Name("bar"))),
// 			"swap two rows by position - reverse arguments",
// 		},
// 		{
// 			[]opt.SelectionOption{opt.ByLabels([]string{"0", "1"})},
// 			mustNew([]int{1, 0, 2}, Idx([]int{1, 0, 2}, opt.Name("foo")), Idx([]string{"B", "A", "C"}, opt.Name("bar"))),
// 			"swap two rows by label",
// 		},
// 		{
// 			[]opt.SelectionOption{opt.ByLabels([]string{"1", "0"})},
// 			mustNew([]int{1, 0, 2}, Idx([]int{1, 0, 2}, opt.Name("foo")), Idx([]string{"B", "A", "C"}, opt.Name("bar"))),
// 			"swap two rows by label - reverse arguments",
// 		},
// 		{
// 			[]opt.SelectionOption{opt.ByLevels([]int{0, 1})},
// 			mustNew([]int{0, 1, 2}, Idx([]string{"A", "B", "C"}, opt.Name("bar")), Idx([]int{0, 1, 2}, opt.Name("foo"))),
// 			"swap two levels by position",
// 		},
// 		{
// 			[]opt.SelectionOption{opt.ByLevels([]int{1, 0})},
// 			mustNew([]int{0, 1, 2}, Idx([]string{"A", "B", "C"}, opt.Name("bar")), Idx([]int{0, 1, 2}, opt.Name("foo"))),
// 			"swap two levels by position - reverse arguments",
// 		},
// 		{
// 			[]opt.SelectionOption{opt.ByLevelNames([]string{"foo", "bar"})},
// 			mustNew([]int{0, 1, 2}, Idx([]string{"A", "B", "C"}, opt.Name("bar")), Idx([]int{0, 1, 2}, opt.Name("foo"))),
// 			"swap two levels by name",
// 		},
// 		{
// 			[]opt.SelectionOption{opt.ByLevelNames([]string{"bar", "foo"})},
// 			mustNew([]int{0, 1, 2}, Idx([]string{"A", "B", "C"}, opt.Name("bar")), Idx([]int{0, 1, 2}, opt.Name("foo"))),
// 			"swap two levels by name - reverse arguments",
// 		},
// 	}
// 	for _, test := range tests {
// 		s, _ := New([]int{0, 1, 2}, Idx([]int{0, 1, 2}, opt.Name("foo")), Idx([]string{"A", "B", "C"}, opt.Name("bar")))
// 		origS, _ := New([]int{0, 1, 2}, Idx([]int{0, 1, 2}, opt.Name("foo")), Idx([]string{"A", "B", "C"}, opt.Name("bar")))
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
// 	s, _ := New([]int{0, 1, 2}, Idx([]int{0, 1, 2}, opt.Name("foo")), Idx([]string{"A", "B", "C"}, opt.Name("bar")))
// 	origS, _ := New([]int{0, 1, 2}, Idx([]int{0, 1, 2}, opt.Name("foo")), Idx([]string{"A", "B", "C"}, opt.Name("bar")))
// 	sel := s.Select(opt.ByRows([]int{0, 1, 2}))
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
