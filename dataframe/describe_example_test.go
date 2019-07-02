package dataframe

import (
	"fmt"
	"log"
)

func ExampleColSlice() {
	df, _ := New([]interface{}{[]string{"bar"}, []string{"baz", "qux"}})
	fmt.Println(df)
	//Output:
}

func ExampleNew_float64() {
	df, err := New(
		[]interface{}{[]float64{0, 1.5}, []float64{2.5, 3}},
		Config{
			Index: []string{"foo", "bar"},
			Cols:  []interface{}{"baz", "qux"},
		})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//         baz  qux
	// foo       0  2.5
	// bar     1.5    3
	//
	// datatype: float64
}

func ExampleNew_string_indexUnnamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar", Index: "baz"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//          0   1
	// baz    foo bar
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
	//          0   1
	// baz    foo bar
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
	//                0   1
	// baz corge    foo bar
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
	//                0   1
	// baz corge    foo bar
	// datatype: string
	// name: foobar
}

func ExampleNew_string_colsUnnamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar", Cols: []interface{}{"baz", "qux"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//      baz qux
	// 0    foo bar
	// datatype: string
	// name: foobar
}

func ExampleNew_string_colsNamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar", Cols: []interface{}{"baz", "qux"}, ColsName: "corge"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//      corge baz qux
	// 0          foo bar
	// datatype: string
	// name: foobar
}

func ExampleNew_string_multicolUnnamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar", MultiCol: [][]interface{}{{"quux", "quux"}, {"baz", "qux"}}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//      quux
	//       baz qux
	// 0     foo bar
	// datatype: string
	// name: foobar
}

func ExampleNew_string_multicolNamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar", MultiCol: [][]interface{}{{"quux", "quax"}, {"baz", "qux"}}, MultiColNames: []string{"corge", "grault"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//       corge quux quax
	//      grault  baz  qux
	// 0            foo  bar
	// datatype: string
	// name: foobar
}

func ExampleNew_string_multicolNamed_repeat() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar", MultiCol: [][]interface{}{{"quux", "quux"}, {"baz", "qux"}}, MultiColNames: []string{"corge", "grault"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//       corge quux
	//      grault  baz qux
	// 0            foo bar
	// datatype: string
	// name: foobar
}

func ExampleNew_string_indexNamed_colsUnnamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar",
			Index: "baz", IndexName: "corge",
			Cols: []interface{}{"quux", "qux"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	// corge
	//          quux qux
	//   baz     foo bar
	// datatype: string
	// name: foobar
}

func ExampleNew_string_indexNamed_colsNamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar",
			Index: "baz", IndexName: "corge",
			Cols: []interface{}{"quux", "qux"}, ColsName: "quuz"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	// corge
	//          quuz quux qux
	//   baz          foo bar
	// datatype: string
	// name: foobar
}

func ExampleNew_string_multiindexNamed_colsNamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar",
			MultiIndex: []interface{}{"baz", "garply"}, MultiIndexNames: []string{"corge", "grault"},
			Cols: []interface{}{"quux", "qux"}, ColsName: "quuz"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	// corge grault
	//                 quuz quux qux
	//   baz garply          foo bar
	// datatype: string
	// name: foobar
}

func ExampleNew_string_multiindexNamed_multicolNamed() {
	df, err := New([]interface{}{"foo", "bar"},
		Config{Name: "foobar",
			MultiIndex: []interface{}{"baz", "garply"}, MultiIndexNames: []string{"grault", "corge"},
			MultiCol: [][]interface{}{{"fred", "fred"}, {"quux", "qux"}}, MultiColNames: []string{"waldo", "quuz"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	// grault  corge
	//                  waldo fred
	//                   quuz quux qux
	//    baz garply           foo bar
	// datatype: string
	// name: foobar
}

func ExampleNew_float64_indexNamed_multicolNamed() {
	df, err := New([]interface{}{"qux", "waldo"},
		Config{
			Name:  "foobar",
			Index: "foo", IndexName: "grault",
			MultiCol: [][]interface{}{{"quux", "quux"}, {"bar", "baz"}}, MultiColNames: []string{"quuz", "garply"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	// grault
	//             quuz quux
	//           garply  bar   baz
	//    foo            qux waldo
	// datatype: string
	// name: foobar
}

func ExampleNew_float64_colsNamed_repeat_resume() {
	df := MustNew([]interface{}{"qux", "waldo", "fred"},
		Config{Name: "foobar", Cols: []interface{}{"quux", "quux", "foo"}})
	fmt.Println(df)
	// Output:
	//      quux        foo
	// 0     qux waldo fred
	// datatype: string
	// name: foobar
}

func ExampleDataFrame_Col() {
	df, err := New([]interface{}{[]float64{1, 3, 5}, []string{"foo", "bar", "baz"}}, Config{Cols: []interface{}{"qux", "corge"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df.Col("corge"))
	// Output:
	// 0    foo
	// 1    bar
	// 2    baz
	// datatype: string
	// name: corge
}

// Selects the first column with this label from the first level
func ExampleDataFrame_multiCol_col() {
	df, err := New([]interface{}{[]float64{1, 3, 5}, []string{"foo", "bar", "baz"}},
		Config{MultiCol: [][]interface{}{{"qux", "qux"}, {"quux", "quuz"}}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df.Col("qux"))
	// Output:
	// 0    1
	// 1    3
	// 2    5
	// datatype: float64
	// name: qux | quux
}

func ExampleDataFrame_subset() {
	df, err := New([]interface{}{[]float64{1, 3, 5}, []string{"foo", "bar", "baz"}},
		Config{MultiCol: [][]interface{}{{"qux", "qux"}, {"quux", "quuz"}}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df.Subset([]int{0}))
	// Output:
	//       qux
	//      quux quuz
	// 0       1  foo
	// datatype: float64
	//  <nil>
}
