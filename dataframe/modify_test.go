package dataframe

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/d4l3k/messagediff"
	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/series"
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

// func TestDataFrame_Modify_Sort(t *testing.T) {
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

func TestDataFrame_Modify_SetSpecial(t *testing.T) {
	type args struct {
		colLabel string
		s        *series.Series
	}
	type want struct {
		df *DataFrame
	}
	var tests = []struct {
		name  string
		input *DataFrame
		args  args
		want  want
	}{
		{name: "pass",
			input: MustNew([]interface{}{"foo"}),
			args:  args{colLabel: "0", s: series.MustNew("bar")},
			want:  want{df: MustNew([]interface{}{"bar"})},
		},
		{"fail: invalid column label",
			MustNew([]interface{}{"foo"}),
			args{colLabel: "100", s: series.MustNew("bar")},
			want{df: MustNew([]interface{}{"foo"})},
		},
		{"fail: series length does not match dataframe length",
			MustNew([]interface{}{"foo"}),
			args{colLabel: "0", s: series.MustNew([]string{"bar", "baz"})},
			want{df: MustNew([]interface{}{"foo"})},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			df := tt.input
			dfArchive := tt.input.Copy()
			df.InPlace.Set(tt.args.colLabel, tt.args.s)
			if !Equal(df, tt.want.df) {
				t.Errorf("InPlace.Set() got %v, want %v", df, tt.want.df)
			}
			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("InPlace.Set() returned no log message, want log due to fail")
				}
				buf.Reset()
			}

			dfCopy := dfArchive.Set(tt.args.colLabel, tt.args.s)
			if !Equal(dfCopy, tt.want.df) {
				t.Errorf("DataFrame.Set() got %v, want %v", dfCopy, tt.want.df)
			}
			if !strings.Contains(tt.name, "fail") {
				if Equal(dfArchive, dfCopy) {
					t.Errorf("DataFrame.Set() retained access to original, want copy")
				}
			}
			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("DataFrame.Set() returned no log message, want log due to fail")
				}
			}
		})
	}
}

func TestDataFrame_Modify_SwapRows(t *testing.T) {
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
			df := MustNew([]interface{}{[]string{"foo", "bar"}}, Config{Index: []int{0, 1}})
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

func TestDataFrame_Modify_SwapColumns(t *testing.T) {
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
		{"fail: i", args{2, 0}, true},
		{"fail: j", args{0, 2}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := MustNew([]interface{}{"foo", "bar"}, Config{Col: []string{"0", "1"}})
			dfArchive := df.Copy()
			want := MustNew([]interface{}{"bar", "foo"}, Config{Col: []string{"1", "0"}})

			dfCopy, err := dfArchive.SwapColumns(tt.args.i, tt.args.j)
			if (err != nil) != tt.wantErr {
				t.Errorf("DataFrame.SwapColumns() error = %v, want %v", err, tt.wantErr)
				return
			}

			// intentionally skip fail case
			if !strings.Contains(tt.name, "fail") {
				df.InPlace.SwapColumns(tt.args.i, tt.args.j)
				if !Equal(df, want) {
					t.Errorf("InPlace.SwapColumns() got %v, want %v", df.cols.Levels[0].DataType, want.cols.Levels[0].DataType)
				}
				if !Equal(dfCopy, want) {
					t.Errorf("DataFrame.SwapColumns() got %v, want %v", dfCopy, want)
				}
				if Equal(dfArchive, dfCopy) {
					t.Errorf("DataFrame.SwapColumns() retained access to original, want copy")
				}
			}
		})
	}
}

func TestDataFrame_Modify_InsertRow(t *testing.T) {
	multi := MustNew([]interface{}{[]string{"foo"}}, Config{MultiIndex: []interface{}{"A", 1}})

	type args struct {
		row       int
		val       []interface{}
		idxLabels []interface{}
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
			args:  args{row: 0, val: []interface{}{"foo"}, idxLabels: []interface{}{"A"}},
			want:  want{df: MustNew([]interface{}{"foo"}, Config{Index: "A"}), err: false}},
		{"singleIndex",
			MustNew([]interface{}{[]string{"foo"}}, Config{Index: "A"}),
			args{0, []interface{}{"bar"}, []interface{}{"B"}},
			want{df: MustNew([]interface{}{[]string{"bar", "foo"}}, Config{Index: []string{"B", "A"}}), err: false}},
		{"no label provided, not default",
			MustNew([]interface{}{[]string{"foo"}}, Config{Index: "A"}),
			args{0, []interface{}{"bar"}, nil},
			want{df: MustNew([]interface{}{[]string{"bar", "foo"}}, Config{Index: []string{"NaN", "A"}}), err: false}},
		{"no label provided, default",
			MustNew([]interface{}{"foo"}),
			args{0, []interface{}{"bar"}, nil},
			want{df: MustNew([]interface{}{[]string{"bar", "foo"}}), err: false}},
		{"multiIndex",
			multi,
			args{1, []interface{}{"bar"}, []interface{}{"B", 2}},
			want{df: MustNew([]interface{}{[]string{"foo", "bar"}}, Config{MultiIndex: []interface{}{[]string{"A", "B"}, []int{1, 2}}}), err: false}},
		{"fail: exceeds index length",
			multi,
			args{1, []interface{}{"bar"}, []interface{}{"C", "D", "E"}},
			want{nil, true}},
		{"fail: wrong values length",
			multi,
			args{1, []interface{}{"bar", "baz"}, []interface{}{"C", 3}},
			want{nil, true}},
		{"fail: invalid position",
			multi,
			args{10, []interface{}{"bar"}, []interface{}{"C", 3}},
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
			err := df.InPlace.InsertRow(tt.args.row, tt.args.val, tt.args.idxLabels...)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.InsertRow() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(df, tt.want.df) {
					t.Errorf("InPlace.InsertRow() got %v, want %v", df, tt.want.df)
					diff, _ := messagediff.PrettyDiff(df, tt.want.df)
					fmt.Println(diff)
				}
			}

			dfCopy, err := dfArchive.InsertRow(tt.args.row, tt.args.val, tt.args.idxLabels...)
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

func TestDataFrame_Modify_InsertColumn(t *testing.T) {
	single := MustNew([]interface{}{[]string{"foo"}}, Config{Col: []string{"bar"}})
	multi := MustNew([]interface{}{[]string{"foo"}}, Config{MultiCol: [][]string{{"bar"}, {"baz"}}})
	type args struct {
		col       int
		val       interface{}
		colLabels []string
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
			args:  args{col: 0, val: "foo", colLabels: []string{"bar"}},
			want:  want{df: MustNew([]interface{}{"foo"}, Config{Col: []string{"bar"}}), err: false}},
		{"single col",
			single,
			args{0, "baz", []string{"qux"}},
			want{df: MustNew([]interface{}{"baz", "foo"}, Config{Col: []string{"qux", "bar"}}), err: false}},
		{"insert label into default range",
			MustNew([]interface{}{"foo", "bar"}),
			args{1, "baz", []string{"qux"}},
			want{df: MustNew([]interface{}{"foo", "baz", "bar"}, Config{Col: []string{"0", "qux", "1"}}), err: false}},
		{"no label provided, not default",
			single,
			args{0, "baz", nil},
			want{df: MustNew([]interface{}{"baz", "foo"}, Config{Col: []string{"NaN", "bar"}}), err: false}},
		{"no label provided, default",
			MustNew([]interface{}{[]string{"foo"}}),
			args{0, "baz", nil},
			want{df: MustNew([]interface{}{"baz", "foo"}), err: false}},
		{"multi col",
			multi,
			args{1, "corge", []string{"qux", "quux"}},
			want{df: MustNew([]interface{}{"foo", "corge"}, Config{MultiCol: [][]string{{"bar", "qux"}, {"baz", "quux"}}}), err: false}},
		{"fail: exceeds column length",
			multi,
			args{1, "bar", []string{"A", "B", "C"}},
			want{nil, true}},
		{"fail: wrong values length",
			multi,
			args{1, []string{"bar", "baz"}, nil},
			want{nil, true}},
		{"fail: invalid column insertion point",
			single,
			args{10, "bar", nil},
			want{nil, true}},
		{"fail: unsupported values",
			single,
			args{0, complex64(1), nil},
			want{nil, true}},
		{"fail: unsupported value in empty dataframe",
			newEmptyDataFrame(),
			args{0, complex64(1), nil},
			want{nil, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := tt.input.Copy()
			dfArchive := tt.input.Copy()
			err := df.InPlace.InsertColumn(tt.args.col, tt.args.val, tt.args.colLabels...)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.InsertColumn() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(df, tt.want.df) {
					t.Errorf("InPlace.InsertColumn() got %v, want %v", df, tt.want.df)
				}
				diff, _ := messagediff.PrettyDiff(df, tt.want.df)
				fmt.Println(diff)
			}

			dfCopy, err := dfArchive.InsertColumn(tt.args.col, tt.args.val, tt.args.colLabels...)
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

func TestDataFrame_Modify_AppendRow(t *testing.T) {
	type args struct {
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
		{"singleIndex",
			MustNew([]interface{}{[]string{"foo"}}, Config{Index: []int{1}}),
			args{val: []interface{}{"bar"}, idx: []interface{}{2}},
			want{df: MustNew([]interface{}{[]string{"foo", "bar"}}, Config{Index: []int{1, 2}}), err: false}},
		{"multiIndex",
			MustNew([]interface{}{[]string{"foo"}}, Config{MultiIndex: []interface{}{[]string{"A"}, []int{1}}}),
			args{[]interface{}{"bar"}, []interface{}{"B", 2}},
			want{MustNew([]interface{}{[]string{"foo", "bar"}}, Config{MultiIndex: []interface{}{[]string{"A", "B"}, []int{1, 2}}}), false}},
		{"fail multiIndex: excessive index values",
			MustNew([]interface{}{[]string{"foo"}}, Config{MultiIndex: []interface{}{[]string{"A"}, []int{1}}}),
			args{[]interface{}{"bar"}, []interface{}{"B", "C", "D"}},
			want{nil, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input
			dfArchive := tt.input.Copy()
			err := s.InPlace.AppendRow(tt.args.val, tt.args.idx...)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.AppendRow() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(s, tt.want.df) {
					t.Errorf("InPlace.AppendRow() got %v, want %v", s, tt.want.df)
				}
			}
			dfCopy, err := dfArchive.AppendRow(tt.args.val, tt.args.idx...)
			if (err != nil) != tt.want.err {
				t.Errorf("DataFrame.Append() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(dfCopy, tt.want.df) {
					t.Errorf("DataFrame.Append() got %v, want %v", dfCopy, tt.want.df)
				}
				if Equal(dfArchive, dfCopy) {
					t.Errorf("DataFrame.Append() retained access to original, want copy")
				}
			}
		})
	}
}

func TestDataFrame_Modify_AppendColumn(t *testing.T) {
	type args struct {
		val       interface{}
		colLabels []string
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
			args:  args{val: "foo", colLabels: []string{"bar"}},
			want:  want{df: MustNew([]interface{}{"foo"}, Config{Col: []string{"bar"}}), err: false}},
		{"single col",
			MustNew([]interface{}{"foo"}, Config{Col: []string{"bar"}}),
			args{"baz", []string{"qux"}},
			want{df: MustNew([]interface{}{"foo", "baz"}, Config{Col: []string{"bar", "qux"}}), err: false}},
		{"fail: exceed cols",
			MustNew([]interface{}{"foo"}, Config{Col: []string{"bar"}}),
			args{"baz", []string{"qux", "quux", "corge"}},
			want{nil, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := tt.input.Copy()
			dfArchive := tt.input.Copy()
			err := df.InPlace.AppendColumn(tt.args.val, tt.args.colLabels...)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.AppendColumn() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(df, tt.want.df) {
					t.Errorf("InPlace.AppendColumn() got %v, want %v", df, tt.want.df)
				}
			}

			dfCopy, err := dfArchive.AppendColumn(tt.args.val, tt.args.colLabels...)
			if (err != nil) != tt.want.err {
				t.Errorf("DataFrame.AppendColumn() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(dfCopy, tt.want.df) {
					t.Errorf("DataFrame.AppendColumn() got %v, want %v", dfCopy, tt.want.df)
				}
				if Equal(dfArchive, dfCopy) {
					t.Errorf("DataFrame.AppendColumn() retained access to original, want copy")
				}
			}
		})
	}
}

func TestDataFrame_Modify_Set(t *testing.T) {
	type args struct {
		rowPositions int
		val          interface{}
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
		{name: "singleRow",
			input: MustNew([]interface{}{"foo"}), args: args{rowPositions: 0, val: "bar"},
			want: want{df: MustNew([]interface{}{"bar"}), err: false}},
		{"fail: invalid index singleRow",
			MustNew([]interface{}{"foo"}), args{1, "bar"},
			want{MustNew([]interface{}{"foo"}), true}},
		{"fail: unsupported value",
			MustNew([]interface{}{"foo"}), args{0, complex64(1)},
			want{MustNew([]interface{}{"foo"}), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input
			dfArchive := tt.input.Copy()
			err := s.InPlace.SetRow(tt.args.rowPositions, tt.args.val)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.SetRow() error = %v, want %v", err, tt.want.err)
				return
			}
			if !Equal(s, tt.want.df) {
				t.Errorf("InPlace.SetRow() got %v, want %v", s, tt.want.df)
			}

			dfCopy, err := dfArchive.SetRow(tt.args.rowPositions, tt.args.val)
			if (err != nil) != tt.want.err {
				t.Errorf("DataFrame.SetRow() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(dfCopy, tt.want.df) {
					t.Errorf("DataFrame.SetRow() got %v, want %v", dfCopy, tt.want.df)
				}
				if Equal(dfArchive, dfCopy) {
					t.Errorf("DataFrame.SetRow() retained access to original, want copy")
				}
			}
		})
	}
}

func TestDataFrame_Modify_SetRows(t *testing.T) {
	type args struct {
		rowPositions []int
		val          interface{}
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
		{"singleRow",
			MustNew([]interface{}{"foo"}), args{rowPositions: []int{0}, val: "bar"},
			want{df: MustNew([]interface{}{"bar"}), err: false}},
		{"multiRow",
			MustNew([]interface{}{[]string{"foo", "bar"}}), args{[]int{0, 1}, "baz"},
			want{MustNew([]interface{}{[]string{"baz", "baz"}}), false}},
		{"fail: singleRow",
			MustNew([]interface{}{"foo"}), args{rowPositions: []int{0}, val: complex64(1)},
			want{df: MustNew([]interface{}{"foo"}), err: true}},
		{"fail: invalid index singleRow",
			MustNew([]interface{}{"foo"}), args{[]int{1}, "bar"},
			want{MustNew([]interface{}{"foo"}), true}},
		{"fail: partial success on multiRow",
			MustNew([]interface{}{"foo"}), args{[]int{0, 2}, "bar"},
			want{MustNew([]interface{}{"foo"}), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := tt.input
			dfArchive := tt.input.Copy()
			err := df.InPlace.SetRows(tt.args.rowPositions, tt.args.val)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.Set() error = %v, want %v", err, tt.want.err)
				return
			}
			if !Equal(df, tt.want.df) {
				t.Errorf("InPlace.Set() got %v, want %v", df, tt.want.df)
			}

			dfCopy, err := dfArchive.SetRows(tt.args.rowPositions, tt.args.val)
			if (err != nil) != tt.want.err {
				t.Errorf("DataFrame.Set() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(dfCopy, tt.want.df) {
					t.Errorf("DataFrame.Set() got %v, want %v", dfCopy, tt.want.df)
				}
				if Equal(dfArchive, dfCopy) {
					t.Errorf("DataFrame.Set() retained access to original, want copy")
				}
			}
		})
	}
}

func TestDataFrame_Modify_SetColumn(t *testing.T) {
	type args struct {
		col int
		val interface{}
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
		{"singleColumn",
			MustNew([]interface{}{[]string{"foo", "bar"}}), args{col: 0, val: []string{"baz", "qux"}},
			want{df: MustNew([]interface{}{[]string{"baz", "qux"}}), err: false}},
		{"fail: invalid col position",
			MustNew([]interface{}{"foo"}), args{col: 10, val: "bar"},
			want{df: MustNew([]interface{}{"foo"}), err: true}},
		{"fail: unsupported value",
			MustNew([]interface{}{"foo"}), args{col: 0, val: complex64(1)},
			want{df: MustNew([]interface{}{"foo"}), err: true}},
		{"fail: excessive val length",
			MustNew([]interface{}{"foo"}), args{col: 0, val: []string{"baz", "qux"}},
			want{df: MustNew([]interface{}{"foo"}), err: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := tt.input
			dfArchive := tt.input.Copy()
			err := df.InPlace.SetColumn(tt.args.col, tt.args.val)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.Set() error = %v, want %v", err, tt.want.err)
				return
			}
			if !Equal(df, tt.want.df) {
				t.Errorf("InPlace.Set() got %v, want %v", df, tt.want.df)
			}

			dfCopy, err := dfArchive.SetColumn(tt.args.col, tt.args.val)
			if (err != nil) != tt.want.err {
				t.Errorf("DataFrame.Set() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(dfCopy, tt.want.df) {
					t.Errorf("DataFrame.Set() got %v, want %v", dfCopy, tt.want.df)
				}
				if Equal(dfArchive, dfCopy) {
					t.Errorf("DataFrame.Set() retained access to original, want copy")
				}
			}
		})
	}
}

func TestDataFrame_Modify_SetColumns(t *testing.T) {
	type args struct {
		columnPositions []int
		val             interface{}
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
		{"singleColumn",
			MustNew([]interface{}{[]string{"foo", "bar"}, []string{"foo", "bar"}}), args{columnPositions: []int{0, 1}, val: []string{"baz", "qux"}},
			want{df: MustNew([]interface{}{[]string{"baz", "qux"}, []string{"baz", "qux"}}), err: false}},
		{"fail: invalid col position",
			MustNew([]interface{}{"foo"}), args{[]int{10}, "bar"},
			want{df: MustNew([]interface{}{"foo"}), err: true}},
		{"fail: unsupported value",
			MustNew([]interface{}{"foo"}), args{[]int{0}, complex64(1)},
			want{df: MustNew([]interface{}{"foo"}), err: true}},
		{"fail: excessive val length",
			MustNew([]interface{}{"foo"}), args{[]int{0}, []string{"baz", "qux"}},
			want{df: MustNew([]interface{}{"foo"}), err: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := tt.input
			dfArchive := tt.input.Copy()
			err := df.InPlace.SetColumns(tt.args.columnPositions, tt.args.val)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.Set() error = %v, want %v", err, tt.want.err)
				return
			}
			if !Equal(df, tt.want.df) {
				t.Errorf("InPlace.Set() got %v, want %v", df, tt.want.df)
			}

			dfCopy, err := dfArchive.SetColumns(tt.args.columnPositions, tt.args.val)
			if (err != nil) != tt.want.err {
				t.Errorf("DataFrame.Set() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(dfCopy, tt.want.df) {
					t.Errorf("DataFrame.Set() got %v, want %v", dfCopy, tt.want.df)
				}
				if Equal(dfArchive, dfCopy) {
					t.Errorf("DataFrame.Set() retained access to original, want copy")
				}
			}
		})
	}
}

func TestDataFrame_Modify_Drop(t *testing.T) {
	type args struct {
		rowPositions int
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
		{"drop to 0",
			MustNew([]interface{}{"foo"}), args{rowPositions: 0},
			want{df: newEmptyDataFrame(), err: false}},
		{"singleRow",
			MustNew([]interface{}{[]string{"foo", "bar", "baz"}}, Config{Index: []int{0, 1, 2}}), args{1},
			want{MustNew([]interface{}{[]string{"foo", "baz"}}, Config{Index: []int{0, 2}}), false}},
		{"fail: invalid index singleRow",
			MustNew([]interface{}{"foo"}), args{1},
			want{MustNew([]interface{}{"foo"}), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input
			dfArchive := tt.input.Copy()
			err := s.InPlace.DropRow(tt.args.rowPositions)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.DropRow() error = %v, want %v", err, tt.want.err)
				return
			}
			if !Equal(s, tt.want.df) {
				t.Errorf("InPlace.DropRow() got %v, want %v", s, tt.want.df)
			}

			dfCopy, err := dfArchive.DropRow(tt.args.rowPositions)
			if (err != nil) != tt.want.err {
				t.Errorf("DataFrame.DropRow() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(dfCopy, tt.want.df) {
					t.Errorf("DataFrame.DropRow() got %v, want %v", dfCopy, tt.want.df)
				}
				if Equal(dfArchive, dfCopy) {
					t.Errorf("DataFrame.DropRow() retained access to original, want copy")
				}
			}
		})
	}
}

func TestDataFrame_Modify_DropRows(t *testing.T) {
	type args struct {
		rowPositions []int
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
		{"drop to 0",
			MustNew([]interface{}{"foo"}), args{rowPositions: []int{0}},
			want{df: newEmptyDataFrame(), err: false}},
		{"singleRow",
			MustNew([]interface{}{[]string{"foo", "bar", "baz"}}, Config{Index: []int{0, 1, 2}}), args{[]int{1}},
			want{MustNew([]interface{}{[]string{"foo", "baz"}}, Config{Index: []int{0, 2}}), false}},
		{"multiRow",
			MustNew([]interface{}{[]string{"foo", "bar", "baz"}}, Config{Index: []int{0, 1, 2}}), args{[]int{1, 2}},
			want{MustNew([]interface{}{[]string{"foo"}}, Config{Index: []int{0}}), false}},
		{"multiRow reverse",
			MustNew([]interface{}{[]string{"foo", "bar", "baz"}}, Config{Index: []int{0, 1, 2}}), args{[]int{2, 1}},
			want{MustNew([]interface{}{[]string{"foo"}}, Config{Index: []int{0}}), false}},
		{"fail: invalid index singleRow",
			MustNew([]interface{}{"foo"}), args{[]int{1}},
			want{MustNew([]interface{}{"foo"}), true}},
		{"fail: partial success on multiRow",
			MustNew([]interface{}{[]string{"foo", "bar"}}), args{[]int{0, 2}},
			want{MustNew([]interface{}{[]string{"foo", "bar"}}), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input
			dfArchive := tt.input.Copy()
			err := s.InPlace.DropRows(tt.args.rowPositions)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.DropRows() error = %v, want %v", err, tt.want.err)
			}
			if !Equal(s, tt.want.df) {
				t.Errorf("InPlace.DropRows() got %v, want %v", s, tt.want.df)
			}

			dfCopy, err := dfArchive.DropRows(tt.args.rowPositions)
			if (err != nil) != tt.want.err {
				t.Errorf("DataFrame.DropRows() error = %v, want %v", err, tt.want.err)
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(dfCopy, tt.want.df) {
					t.Errorf("DataFrame.DropRows() got %v, want %v", dfCopy, tt.want.df)
				}
				if Equal(dfArchive, dfCopy) {
					t.Errorf("DataFrame.DropRows() retained access to original, want copy")
				}
			}
		})
	}
}

func TestDataFrame_Modify_DropColumn(t *testing.T) {
	type args struct {
		col int
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
		{"drop to 0",
			MustNew([]interface{}{"foo"}), args{col: 0},
			want{df: newEmptyDataFrame(), err: false}},
		{"multiple cols",
			MustNew([]interface{}{"foo", "bar"}, Config{Col: []string{"0", "1"}}), args{1},
			want{MustNew([]interface{}{"foo"}, Config{Col: []string{"0"}}), false}},
		{"fail: invalid col",
			MustNew([]interface{}{"foo"}), args{10},
			want{MustNew([]interface{}{"foo"}), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := tt.input
			dfArchive := tt.input.Copy()
			err := df.InPlace.DropColumn(tt.args.col)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.DropColumn() error = %v, want %v", err, tt.want.err)
			}
			if !Equal(df, tt.want.df) {
				t.Errorf("InPlace.DropColumn() got %v, want %v", df, tt.want.df)
				diff, _ := messagediff.PrettyDiff(df, tt.want.df)
				fmt.Println(diff)
			}

			dfCopy, err := dfArchive.DropColumn(tt.args.col)
			if (err != nil) != tt.want.err {
				t.Errorf("DataFrame.DropColumn() error = %v, want %v", err, tt.want.err)
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(dfCopy, tt.want.df) {
					t.Errorf("DataFrame.DropColumn() got %v, want %v", dfCopy, tt.want.df)
				}
				if Equal(dfArchive, dfCopy) {
					t.Errorf("DataFrame.DropColumn() retained access to original, want copy")
				}
			}
		})
	}
}

func TestDataFrame_Modify_DropColumns(t *testing.T) {
	type args struct {
		col []int
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
		{name: "multiple cols",
			input: MustNew([]interface{}{"foo", "bar", "baz"}, Config{Col: []string{"0", "1", "2"}}),
			args:  args{[]int{0, 2}},
			want:  want{MustNew([]interface{}{"bar"}, Config{Col: []string{"1"}}), false}},
		{"fail: invalid col",
			MustNew([]interface{}{"foo"}), args{[]int{10}},
			want{MustNew([]interface{}{"foo"}), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := tt.input
			dfArchive := tt.input.Copy()
			err := df.InPlace.DropColumns(tt.args.col)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.DropColumns() error = %v, want %v", err, tt.want.err)
			}
			if !Equal(df, tt.want.df) {
				t.Errorf("InPlace.DropColumns() got %v, want %v", df, tt.want.df)
				diff, _ := messagediff.PrettyDiff(df, tt.want.df)
				fmt.Println(diff)
			}

			dfCopy, err := dfArchive.DropColumns(tt.args.col)
			if (err != nil) != tt.want.err {
				t.Errorf("DataFrame.DropColumns() error = %v, want %v", err, tt.want.err)
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(dfCopy, tt.want.df) {
					t.Errorf("DataFrame.DropColumns() got %v, want %v", dfCopy, tt.want.df)
				}
				if Equal(dfArchive, dfCopy) {
					t.Errorf("DataFrame.DropColumns() retained access to original, want copy")
				}
			}
		})
	}
}

func TestRow_hash(t *testing.T) {
	df := MustNew([]interface{}{"foo"})
	r := df.Row(0)
	got := r.hash()
	want := "30fce7113467b7e11a683e8d764529f6a23fdb0b"
	if got != want {
		t.Errorf("Row.hash() got %v, want %v", got, want)
	}
}

func TestDataFrame_Modify_DropDuplicates(t *testing.T) {
	var tests = []struct {
		name  string
		input *DataFrame
		want  *DataFrame
	}{
		{"single",
			MustNew([]interface{}{[]string{"foo", "foo", "bar"}}, Config{
				Index: []int{0, 0, 1}}),
			MustNew([]interface{}{[]string{"foo", "bar"}}, Config{Index: []int{0, 1}})},
		{"multi",
			MustNew([]interface{}{[]string{"foo", "foo", "bar"}}, Config{
				MultiIndex: []interface{}{[]int{0, 0, 1}, []int{2, 2, 3}}}),
			MustNew([]interface{}{[]string{"foo", "bar"}}, Config{MultiIndex: []interface{}{[]int{0, 1}, []int{2, 3}}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input
			dfArchive := tt.input.Copy()
			s.InPlace.DropDuplicates()
			if !Equal(s, tt.want) {
				t.Errorf("InPlace.DropDuplicates() got %v, want %v", s, tt.want)
			}

			dfCopy := dfArchive.DropDuplicates()
			if !strings.Contains(tt.name, "fail") {
				if !Equal(dfCopy, tt.want) {
					t.Errorf("DataFrame.DropDuplicates() got %v, want %v", dfCopy, tt.want)
				}
				if Equal(dfArchive, dfCopy) {
					t.Errorf("DataFrame.DropDuplicates() retained access to original, want copy")
				}
			}
		})
	}
}

func TestDataFrame_Modify_DropNull(t *testing.T) {
	type args struct {
		cols []int
	}
	type want struct {
		df *DataFrame
	}
	var tests = []struct {
		name  string
		input *DataFrame
		args  args
		want  want
	}{
		{name: "control: no null rows",
			input: MustNew([]interface{}{"foo"}),
			args:  args{cols: []int{}},
			want:  want{df: MustNew([]interface{}{"foo"})}},
		{"null row",
			MustNew([]interface{}{[]string{"foo", "", "baz"}}, Config{Index: []int{0, 1, 2}}),
			args{nil},
			want{MustNew([]interface{}{[]string{"foo", "baz"}}, Config{Index: []int{0, 2}})}},
		{"null row reverse",
			MustNew([]interface{}{[]string{"baz", "", "foo"}}, Config{Index: []int{0, 1, 2}}),
			args{nil},
			want{MustNew([]interface{}{[]string{"baz", "foo"}}, Config{Index: []int{0, 2}})}},
		{"all null rows",
			MustNew([]interface{}{[]string{"", ""}}, Config{Index: []int{0, 1}}),
			args{nil},
			want{newEmptyDataFrame()}},
		{"null in first column, second row",
			MustNew([]interface{}{[]string{"baz", "", "foo"}, []int{1, 2, 3}}, Config{Index: []int{0, 1, 2}}),
			args{[]int{0}},
			want{MustNew([]interface{}{[]string{"baz", "foo"}, []int{1, 3}}, Config{Index: []int{0, 2}})}},
		{"control: null in first column, but second column selected",
			MustNew([]interface{}{[]string{"baz", "", "foo"}, []int{1, 2, 3}}),
			args{[]int{1}},
			want{MustNew([]interface{}{[]string{"baz", "", "foo"}, []int{1, 2, 3}})}},
		{"control fail: invalid column selected",
			MustNew([]interface{}{"foo"}),
			args{[]int{10}},
			want{MustNew([]interface{}{"foo"})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			df := tt.input
			dfArchive := tt.input.Copy()
			df.InPlace.DropNull(tt.args.cols...)
			if !Equal(df, tt.want.df) {
				t.Errorf("InPlace.DropNull() got %v, want %v", df, tt.want.df)
			}
			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("InPlace.DropNull() returned no log message, want log due to fail")
				}
				buf.Reset()
			}

			dfCopy := dfArchive.DropNull(tt.args.cols...)
			if !Equal(dfCopy, tt.want.df) {
				t.Errorf("DataFrame.DropNull() got %v, want %v", dfCopy, tt.want.df)
			}
			if !strings.Contains(tt.name, "control") {
				if Equal(dfArchive, dfCopy) {
					t.Errorf("DataFrame.DropNull() retained access to original, want copy")
				}
			}
			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("DataFrame.DropNull() returned no log message, want log due to fail")
				}
			}
		})
	}
}

func TestDataFrame_Modify_SetIndex(t *testing.T) {
	df := MustNew([]interface{}{"foo", "bar"}, Config{Index: []int{0}, Col: []string{"0", "1"}})
	type args struct {
		col int
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
		{name: "pass",
			input: df,
			args:  args{col: 0},
			want: want{df: MustNew([]interface{}{"bar"}, Config{
				Col:        []string{"1"},
				MultiIndex: []interface{}{[]string{"foo"}, []int{0}}, MultiIndexNames: []string{"0", ""}}),
				err: false,
			},
		},
		{"control: single column",
			MustNew([]interface{}{"foo"}),
			args{col: 0},
			want{MustNew([]interface{}{"foo"}), false},
		},
		{"fail: invalid col",
			df,
			args{col: 10},
			want{df, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := tt.input.Copy()
			dfArchive := tt.input.Copy()
			err := df.InPlace.SetIndex(tt.args.col)
			if !Equal(df, tt.want.df) {
				t.Errorf("InPlace.SetIndex() got %v, want %v", df, tt.want.df)
			}
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.SetIndex() error = %v, want %v", err, tt.want.err)
			}

			dfCopy, err := dfArchive.SetIndex(tt.args.col)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.SetIndex() error = %v, want %v", err, tt.want.err)
			}
			if !Equal(dfCopy, tt.want.df) {
				t.Errorf("DataFrame.SetIndex() got %v, want %v", dfCopy, tt.want.df)
			}

			if !strings.Contains(tt.name, "fail") {
				if !Equal(dfCopy, tt.want.df) {
					t.Errorf("DataFrame.SetIndex() got %v, want %v", dfCopy, tt.want.df)
				}
				if !strings.Contains(tt.name, "control") {
					if Equal(dfArchive, dfCopy) {
						t.Errorf("DataFrame.SetIndex() retained access to original, want copy")
					}
				}
			}
		})
	}
}

func TestDataFrame_Modify_ResetIndex(t *testing.T) {
	df := MustNew([]interface{}{"foo", "bar"}, Config{MultiIndex: []interface{}{[]int{0}, []string{"baz"}}, Col: []string{"0", "1"}})
	type args struct {
		level int
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
		{name: "pass",
			input: df,
			args:  args{level: 1},
			want: want{df: MustNew([]interface{}{"foo", "bar", "baz"}, Config{
				Col:   []string{"0", "1", ""},
				Index: []int{0}}),
				err: false,
			},
		},
		{"replace last index level with default int level",
			MustNew([]interface{}{"foo"}, Config{Index: "bar", IndexName: "qux", Col: []string{"baz"}}),
			args{level: 0},
			want{MustNew([]interface{}{"foo", "bar"}, Config{Col: []string{"baz", "qux"}}), false},
		},
		{"fail: invalid level",
			df,
			args{level: 10},
			want{df, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := tt.input.Copy()
			dfArchive := tt.input.Copy()
			err := df.InPlace.ResetIndex(tt.args.level)
			if !Equal(df, tt.want.df) {
				t.Errorf("InPlace.ResetIndex() got %v, want %v", df, tt.want.df)
			}
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.ResetIndex() error = %v, want %v", err, tt.want.err)
			}

			dfCopy, err := dfArchive.ResetIndex(tt.args.level)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.ResetIndex() error = %v, want %v", err, tt.want.err)
			}
			if !Equal(dfCopy, tt.want.df) {
				t.Errorf("DataFrame.ResetIndex() got %v, want %v", dfCopy, tt.want.df)
			}

			if !strings.Contains(tt.name, "fail") {
				if !Equal(dfCopy, tt.want.df) {
					t.Errorf("DataFrame.ResetIndex() got %v, want %v", dfCopy, tt.want.df)
				}
				if Equal(dfArchive, dfCopy) {
					t.Errorf("DataFrame.ResetIndex() retained access to original, want copy")
				}
			}
		})
	}
}

func TestDataFrame_ModifyInPlace_DatatypeConversion(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	epochDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	singleRow := MustNew([]interface{}{1.5, 1, "1", false, testDate})
	singleColumn := MustNew([]interface{}{[]interface{}{1.5, 1, "1", false, testDate}})
	type args struct {
		To (func(InPlace))
	}
	type want struct {
		df       *DataFrame
		datatype options.DataType
	}
	var tests = []struct {
		name  string
		input *DataFrame
		args  args
		want  want
	}{
		{name: "float row", input: singleRow, args: args{(InPlace).ToFloat64}, want: want{MustNew([]interface{}{1.5, 1.0, 1.0, 0.0, 1.5566688e+18}), options.Float64}},
		{"float col", singleColumn, args{(InPlace).ToFloat64}, want{MustNew([]interface{}{[]float64{1.5, 1.0, 1.0, 0.0, 1.5566688e+18}}), options.Float64}},
		{"int row", singleRow, args{(InPlace).ToInt64}, want{MustNew([]interface{}{int64(1), int64(1), int64(1), int64(0), int64(1.5566688e+18)}), options.Int64}},
		{"int col", singleColumn, args{(InPlace).ToInt64}, want{MustNew([]interface{}{[]int64{1, 1, 1, 0, 1.5566688e+18}}), options.Int64}},
		{"string row", singleRow, args{(InPlace).ToString}, want{MustNew([]interface{}{"1.5", "1", "1", "false", "2019-05-01 00:00:00 +0000 UTC"}), options.String}},
		{"string col", singleColumn, args{(InPlace).ToString}, want{MustNew([]interface{}{[]string{"1.5", "1", "1", "false", "2019-05-01 00:00:00 +0000 UTC"}}), options.String}},
		{"bool row", singleRow, args{(InPlace).ToBool}, want{MustNew([]interface{}{true, true, true, false, true}), options.Bool}},
		{"bool col", singleColumn, args{(InPlace).ToBool}, want{MustNew([]interface{}{[]bool{true, true, true, false, true}}), options.Bool}},
		{"datetime row", singleRow, args{(InPlace).ToDateTime}, want{MustNew([]interface{}{epochDate, epochDate, time.Time{}, epochDate, testDate}), options.DateTime}},
		{"datetime col", singleColumn, args{(InPlace).ToDateTime}, want{MustNew([]interface{}{[]time.Time{epochDate, epochDate, time.Time{}, epochDate, testDate}}), options.DateTime}},
		{"control row", singleRow, args{(InPlace).ToInterface}, want{MustNew([]interface{}{1.5, 1, "1", false, testDate}, Config{DataType: options.Interface}), options.Interface}},
		{"control col", singleColumn, args{(InPlace).ToInterface}, want{MustNew([]interface{}{[]interface{}{1.5, 1, "1", false, testDate}}), options.Interface}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := tt.input.Copy()
			tt.args.To(df.InPlace)
			if !Equal(df, tt.want.df) {
				t.Errorf("InPlace.To... got %v, want %v", df, tt.want.df)
			}
			if df.dataType() != tt.want.datatype.String() {
				t.Errorf("InPlace.To... got datatype %v, want %v", df.dataType(), tt.want.datatype.String())
			}
		})
	}
}

func TestDataFrame_Modify_DatatypeConversion(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	epochDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	singleRow := MustNew([]interface{}{1.5, 1, "1", false, testDate})
	singleColumn := MustNew([]interface{}{[]interface{}{1.5, 1, "1", false, testDate}})
	type args struct {
		To (func(*DataFrame) *DataFrame)
	}
	type want struct {
		df       *DataFrame
		datatype options.DataType
	}
	var tests = []struct {
		name  string
		input *DataFrame
		args  args
		want  want
	}{
		{name: "float row", input: singleRow, args: args{(*DataFrame).ToFloat64}, want: want{MustNew([]interface{}{1.5, 1.0, 1.0, 0.0, 1.5566688e+18}), options.Float64}},
		{"float col", singleColumn, args{(*DataFrame).ToFloat64}, want{MustNew([]interface{}{[]float64{1.5, 1.0, 1.0, 0.0, 1.5566688e+18}}), options.Float64}},
		{"int row", singleRow, args{(*DataFrame).ToInt64}, want{MustNew([]interface{}{int64(1), int64(1), int64(1), int64(0), int64(1.5566688e+18)}), options.Int64}},
		{"int col", singleColumn, args{(*DataFrame).ToInt64}, want{MustNew([]interface{}{[]int64{1, 1, 1, 0, 1.5566688e+18}}), options.Int64}},
		{"string row", singleRow, args{(*DataFrame).ToString}, want{MustNew([]interface{}{"1.5", "1", "1", "false", "2019-05-01 00:00:00 +0000 UTC"}), options.String}},
		{"string col", singleColumn, args{(*DataFrame).ToString}, want{MustNew([]interface{}{[]string{"1.5", "1", "1", "false", "2019-05-01 00:00:00 +0000 UTC"}}), options.String}},
		{"bool row", singleRow, args{(*DataFrame).ToBool}, want{MustNew([]interface{}{true, true, true, false, true}), options.Bool}},
		{"bool col", singleColumn, args{(*DataFrame).ToBool}, want{MustNew([]interface{}{[]bool{true, true, true, false, true}}), options.Bool}},
		{"datetime row", singleRow, args{(*DataFrame).ToDateTime}, want{MustNew([]interface{}{epochDate, epochDate, time.Time{}, epochDate, testDate}), options.DateTime}},
		{"datetime col", singleColumn, args{(*DataFrame).ToDateTime}, want{MustNew([]interface{}{[]time.Time{epochDate, epochDate, time.Time{}, epochDate, testDate}}), options.DateTime}},
		{"control row", singleRow, args{(*DataFrame).ToInterface}, want{MustNew([]interface{}{1.5, 1, "1", false, testDate}, Config{DataType: options.Interface}), options.Interface}},
		{"control col", singleColumn, args{(*DataFrame).ToInterface}, want{MustNew([]interface{}{[]interface{}{1.5, 1, "1", false, testDate}}), options.Interface}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := tt.input.Copy()
			dfArchive := df.Copy()
			got := tt.args.To(df)
			if !Equal(got, tt.want.df) {
				t.Errorf("DataFrame.To... got %v, want %v", got, tt.want.df)
			}
			if got.dataType() != tt.want.datatype.String() {
				t.Errorf("DataFrame.To... got datatype %v, want %v", got.dataType(), tt.want.datatype.String())
			}
			if !strings.Contains(tt.name, "control") {
				if Equal(got, dfArchive) {
					t.Errorf("DataFrame.To... retained access to original, want copy")
				}
			}
		})
	}
}
