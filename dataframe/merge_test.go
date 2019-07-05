package dataframe

import "testing"

func TestMerge_appendDataFrameRow(t *testing.T) {
	type args struct {
		df2 *DataFrame
	}
	tests := []struct {
		name  string
		input *DataFrame
		args  args
		want  *DataFrame
	}{
		{name: "empty", input: newEmptyDataFrame(),
			args: args{df2: MustNew([]interface{}{"foo"})},
			want: MustNew([]interface{}{"foo"})},
		{"same datatype", MustNew([]interface{}{"foo"}, Config{Index: "1"}),
			args{MustNew([]interface{}{"bar"}, Config{Index: "2"})},
			MustNew([]interface{}{[]string{"foo", "bar"}}, Config{Index: []string{"1", "2"}})},
		{"different datatype", MustNew([]interface{}{"foo"}, Config{Index: "1"}),
			args{MustNew([]interface{}{10}, Config{Index: "2"})},
			MustNew([]interface{}{[]string{"foo", "10"}}, Config{Index: []string{"1", "2"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := tt.input.Copy()
			df.InPlace.appendDataFrameRow(tt.args.df2)
			if !Equal(df, tt.want) {
				t.Errorf("InPlace.appendDataFrameRow() = %v, want %v", df, tt.want)
			}
		})
	}
}

func TestMerge_appendDataFrameColumn(t *testing.T) {
	type args struct {
		df2 *DataFrame
	}
	type want struct {
		df  *DataFrame
		err bool
	}
	tests := []struct {
		name  string
		input *DataFrame
		args  args
		want  want
	}{
		{name: "empty", input: newEmptyDataFrame(),
			args: args{df2: MustNew([]interface{}{"foo"})},
			want: want{df: MustNew([]interface{}{"foo"}), err: false}},
		{"pass", MustNew([]interface{}{"foo"}, Config{Col: []string{"1"}}),
			args{MustNew([]interface{}{"bar"}, Config{Col: []string{"2"}})},
			want{MustNew([]interface{}{"foo", "bar"}, Config{Col: []string{"1", "2"}}), false}},
		// fix to append multiple columns in order
		// {"fail: too many columns in df2", MustNew([]interface{}{"foo"}, Config{Col: []string{"1"}}),
		// 	args{MustNew([]interface{}{"bar", "baz"}, Config{Col: []string{"1", "2"}})},
		// 	want{MustNew([]interface{}{"foo"}, Config{Col: []string{"1"}}), false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := tt.input.Copy()
			err := df.InPlace.appendDataFrameColumn(tt.args.df2)
			if (err != nil) != tt.want.err {
				t.Errorf("DataFrame.SwapColumns() error = %v, want %v", err, tt.want.err)
				return
			}
			if !Equal(df, tt.want.df) {
				t.Errorf("InPlace.appendDataFrameColumn() = %v, want %v", df, tt.want.df)
			}
		})
	}
}
