package dataframe

import (
	"fmt"
	"log"
	"time"

	"github.com/ptiger10/pd/options"
)

func ExampleNew_empty_dataframe() {
	df := MustNew(nil)
	fmt.Println(df)
	// Output:
	// {Empty DataFrame}
}

func ExampleNew_float64() {
	df, err := New(
		[]interface{}{[]float64{0, 1.5}, []float64{2.5, 3}},
		Config{
			Index: []string{"foo", "bar"},
			Col:   []string{"baz", "qux"},
		})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//         baz   qux
	// foo    0.00  2.50
	// bar    1.50  3.00
	//
	// datatype: float64
}

func ExampleNew_datetime() {
	df, err := New([]interface{}{[]time.Time{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//                      0
	// 0    1/1/2019T00:00:00
	// 1    1/2/2019T00:00:00
	//
	// datatype: dateTime
}

func ExampleNew_no_interpolation() {
	df := MustNew([]interface{}{[]interface{}{"foo", "bar"}}, Config{Manual: true})
	fmt.Println(df)
	// Output:
	//        0
	// 0    foo
	// 1    bar
	//
	// datatype: interface
}

func ExampleNew_string_indexUnnamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar", Index: "baz"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//          0    1
	// baz    foo  bar
	//
	// datatype: string
	// name: foobar
}
func ExampleNew_string_indexNamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar", Index: "baz", IndexName: "qux"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	// qux
	//          0    1
	// baz    foo  bar
	//
	// datatype: string
	// name: foobar
}

func ExampleNew_string_indexRepeated() {
	df, err := New([]interface{}{[]string{"foo", "bar"}},
		Config{Name: "foobar", Index: []string{"baz", "baz"}, IndexName: "qux"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	// qux
	//          0
	// baz    foo
	//        bar
	//
	// datatype: string
	// name: foobar
}

func ExampleNew_string_indexRepeated_allowed() {
	options.SetDisplayRepeatedLabels(true)
	df, err := New([]interface{}{[]string{"foo", "bar"}},
		Config{Name: "foobar", Index: []string{"baz", "baz"}, IndexName: "qux"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	options.RestoreDefaults()
	// Output:
	// qux
	//          0
	// baz    foo
	// baz    bar
	//
	// datatype: string
	// name: foobar
}

func ExampleNew_string_multiIndexUnnamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar", MultiIndex: []interface{}{"baz", "corge"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//                0    1
	// baz corge    foo  bar
	//
	// datatype: string
	// name: foobar
}

func ExampleNew_string_multiIndexNamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar", MultiIndex: []interface{}{"baz", "corge"}, MultiIndexNames: []string{"qux", "quux"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	// qux  quux
	//                0    1
	// baz corge    foo  bar
	//
	// datatype: string
	// name: foobar
}

func ExampleNew_string_colsUnnamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar", Col: []string{"baz", "qux"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//      baz  qux
	// 0    foo  bar
	//
	// datatype: string
	// name: foobar
}

func ExampleNew_string_colsNamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar", Col: []string{"baz", "qux"}, ColName: "corge"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//      corge baz  qux
	// 0          foo  bar
	//
	// datatype: string
	// name: foobar
}

func ExampleNew_string_multicolUnnamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar", MultiCol: [][]string{{"quux", "quux"}, {"baz", "qux"}}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//      quux
	//       baz  qux
	// 0     foo  bar
	//
	// datatype: string
	// name: foobar
}

func ExampleNew_string_multicolNamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar", MultiCol: [][]string{{"quux", "quax"}, {"baz", "qux"}}, MultiColNames: []string{"corge", "grault"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//       corge quux  quax
	//      grault  baz   qux
	// 0            foo   bar
	//
	// datatype: string
	// name: foobar
}

func ExampleNew_string_multicolNamed_repeat() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar", MultiCol: [][]string{{"quux", "quux"}, {"baz", "qux"}}, MultiColNames: []string{"corge", "grault"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//       corge quux
	//      grault  baz  qux
	// 0            foo  bar
	//
	// datatype: string
	// name: foobar
}

func ExampleNew_string_indexNamed_colsUnnamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar",
			Index: "baz", IndexName: "corge",
			Col: []string{"quux", "qux"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	// corge
	//          quux  qux
	//   baz     foo  bar
	//
	// datatype: string
	// name: foobar
}

func ExampleNew_string_indexNamed_colsNamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar",
			Index: "baz", IndexName: "corge",
			Col: []string{"quux", "qux"}, ColName: "quuz"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	// corge
	//          quuz quux  qux
	//   baz          foo  bar
	//
	// datatype: string
	// name: foobar
}

func ExampleNew_string_multiindexNamed_colsNamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar",
			MultiIndex: []interface{}{"baz", "garply"}, MultiIndexNames: []string{"corge", "grault"},
			Col: []string{"quux", "qux"}, ColName: "quuz"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	// corge grault
	//                 quuz quux  qux
	//   baz garply          foo  bar
	//
	// datatype: string
	// name: foobar
}

func ExampleNew_string_multiindexNamed_multicolNamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar",
			MultiIndex: []interface{}{"baz", "garply"}, MultiIndexNames: []string{"grault", "corge"},
			MultiCol: [][]string{{"fred", "fred"}, {"quux", "qux"}}, MultiColNames: []string{"waldo", "quuz"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	// grault  corge
	//                  waldo fred
	//                   quuz quux  qux
	//    baz garply           foo  bar
	//
	// datatype: string
	// name: foobar
}

func ExampleNew_float64_indexNamed_multicolNamed() {
	df, err := New([]interface{}{"qux", "waldo"},
		Config{
			Name:  "foobar",
			Index: "foo", IndexName: "grault",
			MultiCol: [][]string{{"quux", "quux"}, {"bar", "baz"}}, MultiColNames: []string{"quuz", "garply"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	// grault
	//             quuz quux
	//           garply  bar    baz
	//    foo            qux  waldo
	//
	// datatype: string
	// name: foobar
}

func ExampleNew_float64_colsNamed_repeat_resume() {
	df := MustNew([]interface{}{"qux", "bar", "fred"},
		Config{Name: "foobar", Col: []string{"waldo", "waldo", "foo"}})
	fmt.Println(df)
	// Output:
	//      waldo        foo
	// 0      qux  bar  fred
	//
	// datatype: string
	// name: foobar
}

func ExampleNew_float64_colsNamed_repeat_allowed() {
	options.SetDisplayRepeatedLabels(true)
	df := MustNew([]interface{}{"qux", "bar", "fred"},
		Config{Name: "foobar", Col: []string{"waldo", "waldo", "foo"}})
	fmt.Println(df)
	options.RestoreDefaults()
	// Output:
	//      waldo  waldo   foo
	// 0      qux    bar  fred
	//
	// datatype: string
	// name: foobar
}

func ExampleNew_maxWidth_index() {
	options.SetDisplayMaxWidth(10)
	df := MustNew([]interface{}{[]string{"foo", "bar"}}, Config{Index: []string{"This is a very long index row. Very long indeed.", "qux"}, IndexName: "baz"})
	fmt.Println(df)
	options.RestoreDefaults()
	// Output:
	//        baz
	//                 0
	// This is...    foo
	//        qux    bar
	//
	// datatype: string
}

func ExampleNew_maxWidth_value() {
	s := MustNew([]interface{}{[]string{"This is a very long value row. Very long indeed.", "foo"}})
	fmt.Println(s)
	// Output:
	//                                        0
	// 0    This is a very long value row. V...
	// 1                                    foo
	//
	// datatype: string
}

func ExampleNew_exceed_maxRows_even() {
	options.SetDisplayMaxRows(2)
	s := MustNew([]interface{}{[]float64{0, 1, 2, 3, 4}})
	fmt.Println(s)
	options.RestoreDefaults()
	// Output:
	//         0
	// 0    0.00
	// ...
	// 4    4.00
	//
	// datatype: float64
}

func ExampleNew_exceed_maxRows_odd() {
	options.SetDisplayMaxRows(3)
	s := MustNew([]interface{}{[]float64{0, 1, 2, 3, 4}})
	fmt.Println(s)
	options.RestoreDefaults()

	// Output:
	//         0
	// 0    0.00
	// 1    1.00
	// ...
	// 4    4.00
	//
	// datatype: float64
}

func ExampleNew_exceed_maxColumns_even() {
	options.SetDisplayMaxColumns(4)
	s := MustNew([]interface{}{0, 1, 2, 3, 4})
	fmt.Println(s)
	options.RestoreDefaults()

	// Output:
	//      0  1  ...  3  4
	// 0    0  1       3  4
	//
	// datatype: int64
}

func ExampleNew_exceed_maxColumns_odd() {
	options.SetDisplayMaxColumns(3)
	s := MustNew([]interface{}{0, 1, 2, 3, 4})
	fmt.Println(s)
	options.RestoreDefaults()

	// Output:
	//      0  1  ...  4
	// 0    0  1       4
	//
	// datatype: int64
}

func ExampleDataFrame_Col() {
	df, err := New([]interface{}{[]float64{1, 3, 5}, []string{"foo", "bar", "baz"}}, Config{Col: []string{"qux", "corge"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df.Col("corge"))
	// Output:
	// 0    foo
	// 1    bar
	// 2    baz
	//
	// datatype: string
	// name: corge
}

// Selects the first column with this label from the first level
func ExampleDataFrame_multiCol_col() {
	df, err := New([]interface{}{[]int{1, 3, 5}, []string{"foo", "bar", "baz"}},
		Config{MultiCol: [][]string{{"qux", "qux"}, {"quux", "quuz"}}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df.Col("qux"))
	// Output:
	// 0    1
	// 1    3
	// 2    5
	//
	// datatype: int64
	// name: qux | quux
}

// [START additional structs]

func ExampleDataFrame_empty_df() {
	df := newEmptyDataFrame()
	fmt.Println(df)
	// Output:
	// {Empty DataFrame}
}

func ExampleRow_valid_printer() {
	df := MustNew([]interface{}{"foo", 5, true, ""})
	fmt.Println(df.Row(0))
	// Output:
	// 	   Values: [foo 5 true NaN]
	//     IsNull: [false false false true]
	// ValueTypes: [string int64 bool string]
	//     Labels: [0]
	// LabelTypes: [int64]
}

func ExampleInPlace_method_list() {
	df := MustNew([]interface{}{"foo"})
	fmt.Println(df.InPlace)
	// Output:
	// {InPlace DataFrame Handler}
	// Methods:
	// AppendCol
	// AppendRow
	// Convert
	// DropCol
	// DropCols
	// DropDuplicates
	// DropNull
	// DropRow
	// DropRows
	// InsertCol
	// InsertRow
	// Len
	// ResetIndex
	// Set
	// SetCol
	// SetCols
	// SetIndex
	// SetRow
	// SetRows
	// String
	// SubsetColumns
	// SubsetRows
	// SwapColumns
	// SwapRows
	// ToBool
	// ToDateTime
	// ToFloat64
	// ToInt64
	// ToInterface
	// ToString
}

func ExampleIndex_valid_printer() {
	df := MustNew([]interface{}{[]string{"foo", "bar", "baz"}})
	fmt.Println(df.Index)
	// Output:
	// {DataFrame Index | Len: 3, NumLevels: 1}
}

func ExampleColumns_valid_printer() {
	df := MustNew([]interface{}{[]string{"foo", "bar", "baz"}})
	fmt.Println(df.Columns)
	// Output:
	// {DataFrame Columns | NumCols: 1, NumLevels: 1}
}
func ExampleGrouping_method_list() {
	s := MustNew(
		[]interface{}{[]string{"foo", "bar", "baz"}},
		Config{MultiIndex: []interface{}{[]int{0, 0, 1}, []int{100, 100, 101}}})
	g := s.GroupByIndex()
	fmt.Println(g)
	// Output:
	// {DataFrame Grouping | NumGroups: 2, Groups: [0 | 100, 1 | 101]}
}
