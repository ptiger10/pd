package dataframe

import (
	"strings"
	"testing"
)

// Modify tests check both inplace and copy functionality in the same test, if both are available
func TestRename(t *testing.T) {
	df := MustNew([]interface{}{"foo"}, Config{Name: "baz"})
	want := "qux"
	df.Rename(want)
	got := df.Name()
	if got != want {
		t.Errorf("Rename() got %v, want %v", got, want)
	}
}

// func TestModify_Sort(t *testing.T) {
// 	testDate1 := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
// 	testDate2 := testDate1.Add(24 * time.Hour)
// 	testDate3 := testDate2.Add(24 * time.Hour)

// 	type args struct {
// 		asc bool
// 	}
// 	type want struct {
// 		df *DataFrame
// 	}
// 	var tests = []struct {
// 		name  string
// 		input *DataFrame
// 		args  args
// 		want  want
// 	}{
// 		{"float",
// 			MustNew([]interface{}{[]float64{3, 5, 1}}), args{true},
// 			want{MustNew([]interface{}{[]float64{1, 3, 5}}, Config{Index: []int{2, 0, 1}})}},
// 		{"float desc",
// 			MustNew([]interface{}{[]float64{3, 5, 1}}), args{false},
// 			want{MustNew([]interface{}{[]float64{5, 3, 1}}, Config{Index: []int{1, 0, 2}})}},

// 		{"int",
// 			MustNew([]int{3, 5, 1}), args{true},
// 			want{MustNew([]int{1, 3, 5}, Config{Index: []int{2, 0, 1}})}},
// 		{"int desc",
// 			MustNew([]int{3, 5, 1}), args{false},
// 			want{MustNew([]int{5, 3, 1}, Config{Index: []int{1, 0, 2}})}},

// 		{"string",
// 			MustNew([]string{"15", "3", "1"}), args{true},
// 			want{MustNew([]string{"1", "15", "3"}, Config{Index: []int{2, 0, 1}})}},
// 		{"string desc",
// 			MustNew([]string{"15", "3", "1"}), args{false},
// 			want{MustNew([]string{"3", "15", "1"}, Config{Index: []int{1, 0, 2}})}},

// 		{"bool",
// 			MustNew([]bool{false, true, false}), args{true},
// 			want{MustNew([]bool{false, false, true}, Config{Index: []int{0, 2, 1}})}},
// 		{"bool desc",
// 			MustNew([]bool{false, true, false}), args{false},
// 			want{MustNew([]bool{true, false, false}, Config{Index: []int{1, 0, 2}})}},

// 		{"datetime",
// 			MustNew([]time.Time{testDate2, testDate3, testDate1}), args{true},
// 			want{MustNew([]time.Time{testDate1, testDate2, testDate3}, Config{Index: []int{2, 0, 1}})}},
// 		{"datetime desc",
// 			MustNew([]time.Time{testDate2, testDate3, testDate1}), args{false},
// 			want{MustNew([]time.Time{testDate3, testDate2, testDate1}, Config{Index: []int{1, 0, 2}})}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := tt.input
// 			dfArchive := tt.input.Copy()
// 			s.InPlace.Sort(tt.args.asc)
// 			if !Equal(s, tt.want.df) {
// 				t.Errorf("InPlace.Sort() got %v, want %v", s, tt.want.df)
// 			}

// 			dfCopy := dfArchive.Sort(tt.args.asc)
// 			if !Equal(dfCopy, tt.want.df) {
// 				t.Errorf("DataFrame.Sort() got %v, want %v", dfCopy, tt.want.df)
// 			}
// 			if Equal(dfArchive, dfCopy) {
// 				t.Errorf("DataFrame.Sort() retained access to original, want copy")
// 			}
// 		})
// 	}
// }

func TestModify_SwapRows(t *testing.T) {
	type args struct {
		i int
		j int
	}
	var tests = []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"0,1", args{0, 1}, false},
		{"1,0", args{1, 0}, false},
		{"fail i", args{2, 0}, true},
		{"fail j", args{0, 2}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := MustNew([]interface{}{[]string{"foo", "bar"}})
			dfArchive := df.Copy()
			want := MustNew([]interface{}{[]string{"bar", "foo"}}, Config{Index: []int{1, 0}})

			dfCopy, err := dfArchive.SwapRows(tt.args.i, tt.args.j)
			if (err != nil) != tt.wantErr {
				t.Errorf("DataFrame.Swap() error = %v, want %v", err, tt.wantErr)
				return
			}

			// intentionally skip fail case
			if !strings.Contains(tt.name, "fail") {
				df.InPlace.SwapRows(tt.args.i, tt.args.j)
				if !Equal(df, want) {
					t.Errorf("InPlace.Swap() got %v, want %v", df, want)
				}
				if !Equal(dfCopy, want) {
					t.Errorf("DataFrame.Swap() got %v, want %v", dfCopy, want)
				}
				if Equal(dfArchive, dfCopy) {
					t.Errorf("DataFrame.Sort() retained access to original, want copy")
				}
			}
		})
	}
}

func TestModify_InsertRow(t *testing.T) {
	multi := MustNew([]interface{}{[]string{"foo"}}, Config{MultiIndex: []interface{}{"A", 1}})
	misaligned := MustNew([]interface{}{[]string{"foo", "bar"}})
	misaligned.index.Levels[0].Labels.Drop(1)

	type args struct {
		pos int
		val []interface{}
		idx []interface{}
	}
	type want struct {
		df  *DataFrame
		err bool
	}
	var tests = []struct {
		name  string
		input *DataFrame
		args  args
		want  want
	}{
		{name: "emptySeries",
			input: newEmptyDataFrame(),
			args:  args{pos: 0, val: []interface{}{"foo"}, idx: []interface{}{"A"}},
			want:  want{df: MustNew([]interface{}{"foo"}, Config{Index: "A"}), err: false}},
		{"singleIndex",
			MustNew([]interface{}{[]string{"foo"}}, Config{Index: "A"}),
			args{0, []interface{}{"bar"}, []interface{}{"B"}},
			want{df: MustNew([]interface{}{[]string{"bar", "foo"}}, Config{Index: []string{"B", "A"}}), err: false}},
		{"multiIndex",
			multi,
			args{1, []interface{}{"bar"}, []interface{}{"B", 2}},
			want{df: MustNew([]interface{}{[]string{"foo", "bar"}}, Config{MultiIndex: []interface{}{[]string{"A", "B"}, []int{1, 2}}}), err: false}},
		{"fail: wrong index length",
			multi,
			args{1, []interface{}{"bar"}, []interface{}{"C"}},
			want{nil, true}},
		{"fail: invalid position",
			multi,
			args{10, []interface{}{"bar"}, []interface{}{"C", 3}},
			want{nil, true}},
		{"fail: misaligned df position",
			misaligned,
			args{0, []interface{}{"bar"}, []interface{}{"B"}},
			want{nil, true}},
		{"fail: unsupported index value",
			MustNew([]interface{}{[]string{"foo"}}, Config{Index: "A"}),
			args{0, []interface{}{"bar"}, []interface{}{complex64(1)}},
			want{nil, true}},
		{"fail: unsupported df value",
			MustNew([]interface{}{[]string{"foo"}}, Config{Index: "A"}),
			args{0, []interface{}{complex64(1)}, []interface{}{"A"}},
			want{nil, true}},
		{"fail: unsupported value inserting into empty df",
			newEmptyDataFrame(),
			args{0, []interface{}{complex64(1)}, []interface{}{"A"}},
			want{nil, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := tt.input.Copy()
			dfArchive := tt.input.Copy()
			err := df.InPlace.InsertRow(tt.args.pos, tt.args.val, tt.args.idx)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.InsertRow() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(df, tt.want.df) {
					t.Errorf("InPlace.InsertRow() got %v, want %v", df, tt.want.df)
				}
			}

			dfCopy, err := dfArchive.InsertRow(tt.args.pos, tt.args.val, tt.args.idx)
			if (err != nil) != tt.want.err {
				t.Errorf("DataFrame.Insert() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(dfCopy, tt.want.df) {
					t.Errorf("DataFrame.Insert() got %v, want %v", dfCopy, tt.want.df)
				}
				if Equal(dfArchive, dfCopy) {
					t.Errorf("DataFrame.Insert() retained access to original, want copy")
				}
			}
		})
	}
}

// func TestModify_Append(t *testing.T) {
// 	type args struct {
// 		val interface{}
// 		idx []interface{}
// 	}
// 	type want struct {
// 		df  *DataFrame
// 		err bool
// 	}
// 	var tests = []struct {
// 		name  string
// 		input *DataFrame
// 		args  args
// 		want  want
// 	}{
// 		{"singleIndex",
// 			MustNew([]string{"foo"}, Config{Index: []int{1}}),
// 			args{val: "bar", idx: []interface{}{2}},
// 			want{df: MustNew([]string{"foo", "bar"}, Config{Index: []int{1, 2}}), err: false}},
// 		{"multiIndex",
// 			MustNew([]string{"foo"}, Config{MultiIndex: []interface{}{[]string{"A"}, []int{1}}}),
// 			args{"bar", []interface{}{"B", 2}},
// 			want{MustNew([]string{"foo", "bar"}, Config{MultiIndex: []interface{}{[]string{"A", "B"}, []int{1, 2}}}), false}},
// 		{"fail singleIndex: nil index values",
// 			MustNew([]string{"foo"}, Config{Index: []int{1}}),
// 			args{"bar", nil},
// 			want{nil, true}},
// 		{"fail multiIndex: insufficient index values",
// 			MustNew([]string{"foo"}, Config{MultiIndex: []interface{}{[]string{"A"}, []int{1}}}),
// 			args{"bar", []interface{}{"B"}},
// 			want{nil, true}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := tt.input
// 			dfArchive := tt.input.Copy()
// 			err := s.InPlace.Append(tt.args.val, tt.args.idx)
// 			if (err != nil) != tt.want.err {
// 				t.Errorf("InPlace.Append() error = %v, want %v", err, tt.want.err)
// 				return
// 			}
// 			if !strings.Contains(tt.name, "fail") {
// 				if !Equal(s, tt.want.df) {
// 					t.Errorf("InPlace.Append() got %v, want %v", s, tt.want.df)
// 				}
// 			}
// 			dfCopy, err := dfArchive.Append(tt.args.val, tt.args.idx)
// 			if (err != nil) != tt.want.err {
// 				t.Errorf("DataFrame.Append() error = %v, want %v", err, tt.want.err)
// 				return
// 			}
// 			if !strings.Contains(tt.name, "fail") {
// 				if !Equal(dfCopy, tt.want.df) {
// 					t.Errorf("DataFrame.Append() got %v, want %v", dfCopy, tt.want.df)
// 				}
// 				if Equal(dfArchive, dfCopy) {
// 					t.Errorf("DataFrame.Append() retained access to original, want copy")
// 				}
// 			}
// 		})
// 	}
// }

// func TestModify_Set(t *testing.T) {
// 	type args struct {
// 		rowPositions int
// 		val          interface{}
// 	}
// 	type want struct {
// 		df  *DataFrame
// 		err bool
// 	}
// 	var tests = []struct {
// 		name  string
// 		input *DataFrame
// 		args  args
// 		want  want
// 	}{
// 		{"singleRow",
// 			MustNew("foo"), args{rowPositions: 0, val: "bar"},
// 			want{df: MustNew("bar"), err: false}},
// 		{"fail: invalid index singleRow",
// 			MustNew("foo"), args{1, "bar"},
// 			want{MustNew("foo"), true}},
// 		{"fail: unsupported value",
// 			MustNew("foo"), args{0, complex64(1)},
// 			want{MustNew("foo"), true}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := tt.input
// 			dfArchive := tt.input.Copy()
// 			err := s.InPlace.Set(tt.args.rowPositions, tt.args.val)
// 			if (err != nil) != tt.want.err {
// 				t.Errorf("InPlace.Set() error = %v, want %v", err, tt.want.err)
// 				return
// 			}
// 			if !Equal(s, tt.want.df) {
// 				t.Errorf("InPlace.Set() got %v, want %v", s, tt.want.df)
// 			}

// 			dfCopy, err := dfArchive.Set(tt.args.rowPositions, tt.args.val)
// 			if (err != nil) != tt.want.err {
// 				t.Errorf("DataFrame.Set() error = %v, want %v", err, tt.want.err)
// 				return
// 			}
// 			if !strings.Contains(tt.name, "fail") {
// 				if !Equal(dfCopy, tt.want.df) {
// 					t.Errorf("DataFrame.Set() got %v, want %v", dfCopy, tt.want.df)
// 				}
// 				if Equal(dfArchive, dfCopy) {
// 					t.Errorf("DataFrame.Set() retained access to original, want copy")
// 				}
// 			}
// 		})
// 	}
// }

// func TestModify_SetRows(t *testing.T) {
// 	type args struct {
// 		rowPositions []int
// 		val          interface{}
// 	}
// 	type want struct {
// 		df  *DataFrame
// 		err bool
// 	}
// 	var tests = []struct {
// 		name  string
// 		input *DataFrame
// 		args  args
// 		want  want
// 	}{
// 		{"singleRow",
// 			MustNew("foo"), args{rowPositions: []int{0}, val: "bar"},
// 			want{df: MustNew("bar"), err: false}},
// 		{"multiRow",
// 			MustNew([]string{"foo", "bar"}), args{[]int{0, 1}, "baz"},
// 			want{MustNew([]string{"baz", "baz"}), false}},
// 		{"fail: singleRow",
// 			MustNew("foo"), args{rowPositions: []int{0}, val: complex64(1)},
// 			want{df: MustNew("foo"), err: true}},
// 		{"fail: invalid index singleRow",
// 			MustNew("foo"), args{[]int{1}, "bar"},
// 			want{MustNew("foo"), true}},
// 		{"fail: partial success on multiRow",
// 			MustNew("foo"), args{[]int{0, 2}, "bar"},
// 			want{MustNew("foo"), true}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := tt.input
// 			dfArchive := tt.input.Copy()
// 			err := s.InPlace.SetRows(tt.args.rowPositions, tt.args.val)
// 			if (err != nil) != tt.want.err {
// 				t.Errorf("InPlace.Set() error = %v, want %v", err, tt.want.err)
// 				return
// 			}
// 			if !Equal(s, tt.want.df) {
// 				t.Errorf("InPlace.Set() got %v, want %v", s, tt.want.df)
// 			}

// 			dfCopy, err := dfArchive.SetRows(tt.args.rowPositions, tt.args.val)
// 			if (err != nil) != tt.want.err {
// 				t.Errorf("DataFrame.Set() error = %v, want %v", err, tt.want.err)
// 				return
// 			}
// 			if !strings.Contains(tt.name, "fail") {
// 				if !Equal(dfCopy, tt.want.df) {
// 					t.Errorf("DataFrame.Set() got %v, want %v", dfCopy, tt.want.df)
// 				}
// 				if Equal(dfArchive, dfCopy) {
// 					t.Errorf("DataFrame.Set() retained access to original, want copy")
// 				}
// 			}
// 		})
// 	}
// }

// func TestModify_Drop(t *testing.T) {
// 	type args struct {
// 		rowPositions int
// 	}
// 	type want struct {
// 		df  *DataFrame
// 		err bool
// 	}
// 	var tests = []struct {
// 		name  string
// 		input *DataFrame
// 		args  args
// 		want  want
// 	}{
// 		{"drop to 0",
// 			MustNew("foo"), args{rowPositions: 0},
// 			want{df: newEmptyDataFrame(), err: false}},
// 		{"singleRow",
// 			MustNew([]string{"foo", "bar", "baz"}), args{1},
// 			want{MustNew([]string{"foo", "baz"}, Config{Index: []int{0, 2}}), false}},
// 		{"fail: invalid index singleRow",
// 			MustNew("foo"), args{1},
// 			want{MustNew("foo"), true}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := tt.input
// 			dfArchive := tt.input.Copy()
// 			err := s.InPlace.Drop(tt.args.rowPositions)
// 			if (err != nil) != tt.want.err {
// 				t.Errorf("InPlace.Drop() error = %v, want %v", err, tt.want.err)
// 				return
// 			}
// 			if !Equal(s, tt.want.df) {
// 				t.Errorf("InPlace.Drop() got %v, want %v", s, tt.want.df)
// 			}

// 			dfCopy, err := dfArchive.Drop(tt.args.rowPositions)
// 			if (err != nil) != tt.want.err {
// 				t.Errorf("DataFrame.Drop() error = %v, want %v", err, tt.want.err)
// 				return
// 			}
// 			if !strings.Contains(tt.name, "fail") {
// 				if !Equal(dfCopy, tt.want.df) {
// 					t.Errorf("DataFrame.Drop() got %v, want %v", dfCopy, tt.want.df)
// 				}
// 				if Equal(dfArchive, dfCopy) {
// 					t.Errorf("DataFrame.Drop() retained access to original, want copy")
// 				}
// 			}
// 		})
// 	}
// }

// func TestModify_DropRows(t *testing.T) {
// 	type args struct {
// 		rowPositions []int
// 	}
// 	type want struct {
// 		df  *DataFrame
// 		err bool
// 	}
// 	var tests = []struct {
// 		name  string
// 		input *DataFrame
// 		args  args
// 		want  want
// 	}{
// 		{"drop to 0",
// 			MustNew("foo"), args{rowPositions: []int{0}},
// 			want{df: newEmptyDataFrame(), err: false}},
// 		{"singleRow",
// 			MustNew([]string{"foo", "bar", "baz"}), args{[]int{1}},
// 			want{MustNew([]string{"foo", "baz"}, Config{Index: []int{0, 2}}), false}},
// 		{"multiRow",
// 			MustNew([]string{"foo", "bar", "baz"}), args{[]int{1, 2}},
// 			want{MustNew([]string{"foo"}, Config{Index: []int{0}}), false}},
// 		{"multiRow reverse",
// 			MustNew([]string{"foo", "bar", "baz"}), args{[]int{2, 1}},
// 			want{MustNew([]string{"foo"}, Config{Index: []int{0}}), false}},
// 		{"fail: invalid index singleRow",
// 			MustNew("foo"), args{[]int{1}},
// 			want{MustNew("foo"), true}},
// 		{"fail: partial success on multiRow",
// 			MustNew([]string{"foo", "bar"}), args{[]int{0, 2}},
// 			want{MustNew([]string{"foo", "bar"}), true}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := tt.input
// 			dfArchive := tt.input.Copy()
// 			err := s.InPlace.DropRows(tt.args.rowPositions)
// 			if (err != nil) != tt.want.err {
// 				t.Errorf("InPlace.Drop() error = %v, want %v", err, tt.want.err)
// 				return
// 			}
// 			if !Equal(s, tt.want.df) {
// 				t.Errorf("InPlace.Drop() got %v, want %v", s, tt.want.df)
// 			}

// 			dfCopy, err := dfArchive.DropRows(tt.args.rowPositions)
// 			if (err != nil) != tt.want.err {
// 				t.Errorf("DataFrame.Drop() error = %v, want %v", err, tt.want.err)
// 				return
// 			}
// 			if !strings.Contains(tt.name, "fail") {
// 				if !Equal(dfCopy, tt.want.df) {
// 					t.Errorf("DataFrame.Drop() got %v, want %v", dfCopy, tt.want.df)
// 				}
// 				if Equal(dfArchive, dfCopy) {
// 					t.Errorf("DataFrame.Drop() retained access to original, want copy")
// 				}
// 			}
// 		})
// 	}
// }

// func TestModify_DropDuplicates(t *testing.T) {
// 	var tests = []struct {
// 		name  string
// 		input *DataFrame
// 		want  *DataFrame
// 	}{
// 		{"single",
// 			MustNew([]string{"foo", "foo", "bar"}, Config{
// 				Index: []int{0, 0, 1}}),
// 			MustNew([]string{"foo", "bar"}, Config{Index: []int{0, 1}})},
// 		{"multi",
// 			MustNew([]string{"foo", "foo", "bar"}, Config{
// 				MultiIndex: []interface{}{[]int{0, 0, 1}, []int{2, 2, 3}}}),
// 			MustNew([]string{"foo", "bar"}, Config{MultiIndex: []interface{}{[]int{0, 1}, []int{2, 3}}})},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := tt.input
// 			dfArchive := tt.input.Copy()
// 			s.InPlace.DropDuplicates()
// 			if !Equal(s, tt.want) {
// 				t.Errorf("InPlace.DropDuplicates() got %v, want %v", s, tt.want)
// 			}

// 			dfCopy := dfArchive.DropDuplicates()
// 			if !strings.Contains(tt.name, "fail") {
// 				if !Equal(dfCopy, tt.want) {
// 					t.Errorf("DataFrame.DropDuplicates() got %v, want %v", dfCopy, tt.want)
// 				}
// 				if Equal(dfArchive, dfCopy) {
// 					t.Errorf("DataFrame.DropDuplicates() retained access to original, want copy")
// 				}
// 			}
// 		})
// 	}
// }

// func TestModify_DropNull(t *testing.T) {
// 	type want struct {
// 		df *DataFrame
// 	}
// 	var tests = []struct {
// 		name  string
// 		input *DataFrame
// 		want  want
// 	}{
// 		{"control: no null rows",
// 			MustNew("foo"),
// 			want{df: MustNew("foo")}},
// 		{"null row",
// 			MustNew([]string{"foo", "", "baz"}),
// 			want{MustNew([]string{"foo", "baz"}, Config{Index: []int{0, 2}})}},
// 		{"null row reverse",
// 			MustNew([]string{"baz", "", "foo"}),
// 			want{MustNew([]string{"baz", "foo"}, Config{Index: []int{0, 2}})}},
// 		{"all null rows",
// 			MustNew([]string{"", ""}),
// 			want{newEmptyDataFrame()}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := tt.input
// 			dfArchive := tt.input.Copy()
// 			s.InPlace.DropNull()
// 			if !Equal(s, tt.want.df) {
// 				t.Errorf("InPlace.DropNull() got %v, want %v", s, tt.want.df)
// 			}

// 			dfCopy := dfArchive.DropNull()
// 			if !Equal(dfCopy, tt.want.df) {
// 				t.Errorf("DataFrame.DropNull() got %v, want %v", dfCopy, tt.want.df)
// 			}
// 			if !strings.Contains(tt.name, "control") {
// 				if Equal(dfArchive, dfCopy) {
// 					t.Errorf("DataFrame.DropNull() retained access to original, want copy")
// 				}
// 			}
// 		})
// 	}
// }

// func TestModifyInPlace_DatatypeConversion(t *testing.T) {
// 	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
// 	epochDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
// 	type args struct {
// 		To (func(InPlace))
// 	}
// 	type want struct {
// 		df       *DataFrame
// 		datatype options.DataType
// 	}
// 	var tests = []struct {
// 		name string
// 		args args
// 		want want
// 	}{
// 		{"float", args{(InPlace).ToFloat64}, want{MustNew([]interface{}{[]float64{1.5, 1.0, 1.0, 0, 1.5566688e+18}), options.Float64}},
// 		{"int", args{(InPlace).ToInt64}, want{MustNew([]int64{1, 1, 1, 0, 1.5566688e+18}), options.Int64}},
// 		{"string", args{(InPlace).ToString}, want{MustNew([]string{"1.5", "1", "1", "false", "2019-05-01 00:00:00 +0000 UTC"}), options.String}},
// 		{"bool", args{(InPlace).ToBool}, want{MustNew([]bool{true, true, true, false, true}), options.Bool}},
// 		{"datetime", args{(InPlace).ToDateTime}, want{MustNew([]time.Time{epochDate, epochDate, time.Time{}, epochDate, testDate}), options.DateTime}},
// 		{"control", args{(InPlace).ToInterface}, want{MustNew([]interface{}{1.5, 1, "1", false, testDate}), options.Interface}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := MustNew([]interface{}{1.5, 1, "1", false, testDate})
// 			tt.args.To(s.InPlace)
// 			if !Equal(s, tt.want.df) {
// 				t.Errorf("InPlace.To... got %v, want %v", s, tt.want.df)
// 			}
// 			if s.datatype != tt.want.datatype {
// 				t.Errorf("InPlace.To... got datatype %v, want %v", s.datatype, tt.want.datatype)
// 			}
// 		})
// 	}
// }

// func TestModify_DatatypeConversion(t *testing.T) {
// 	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
// 	epochDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
// 	type args struct {
// 		To (func(*DataFrame) *DataFrame)
// 	}
// 	type want struct {
// 		df       *DataFrame
// 		datatype options.DataType
// 	}
// 	var tests = []struct {
// 		name string
// 		args args
// 		want want
// 	}{
// 		{"float", args{(*DataFrame).ToFloat64}, want{MustNew([]interface{}{[]float64{1.5, 1.0, 1.0, 0, 1.5566688e+18}), options.Float64}},
// 		{"int", args{(*DataFrame).ToInt64}, want{MustNew([]int64{1, 1, 1, 0, 1.5566688e+18}), options.Int64}},
// 		{"string", args{(*DataFrame).ToString}, want{MustNew([]string{"1.5", "1", "1", "false", "2019-05-01 00:00:00 +0000 UTC"}), options.String}},
// 		{"bool", args{(*DataFrame).ToBool}, want{MustNew([]bool{true, true, true, false, true}), options.Bool}},
// 		{"datetime", args{(*DataFrame).ToDateTime}, want{MustNew([]time.Time{epochDate, epochDate, time.Time{}, epochDate, testDate}), options.DateTime}},
// 		{"control: interface", args{(*DataFrame).ToInterface}, want{MustNew([]interface{}{1.5, 1, "1", false, testDate}), options.Interface}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := MustNew([]interface{}{1.5, 1, "1", false, testDate})
// 			got := tt.args.To(s)
// 			if !Equal(got, tt.want.df) {
// 				t.Errorf("DataFrame.To... got %v, want %v", got, tt.want.df)
// 			}
// 			if got.datatype != tt.want.datatype {
// 				t.Errorf("DataFrame.To... got datatype %v, want %v", got.datatype, tt.want.datatype)
// 			}
// 			if !strings.Contains(tt.name, "control") {
// 				if s.DataType() == got.DataType() {
// 					t.Errorf("DataFrame.To... retained access to original, want copy")
// 				}
// 			}
// 		})
// 	}
// }
