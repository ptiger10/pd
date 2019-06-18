package dataframe

import (
	"fmt"
	"log"
)

func ExampleNew_float64() {
	df, err := New([]interface{}{[]float64{0, 1.5}, []float64{2.5, 3}}, Config{Index: []string{"foo", "bar"}, Cols: []interface{}{"baz", "qux"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//        baz qux
	// foo      0 2.5
	// bar    1.5   3
	// datatype: float64
	//
}

func ExampleNew_float64_multiCol() {
	df, err := New([]interface{}{[]float64{0, 1}, []float64{2, 3}}, Config{Index: []string{"foo", "bar"}, Cols: []interface{}{"baz", "qux"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//	     year     2019
	//       class     baz qux
	// index
	//   foo             0   2
	//   bar             1   3
	// datatype: float64
	//
}

func ExampleNew_mixed() {
	df, err := New([]interface{}{[]float64{1, 3, 5}, []string{"foo", "bar", "baz"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//      0   1
	// 0    1 foo
	// 1    3 bar
	// 2    5 baz
	// datatype: mixed
}
