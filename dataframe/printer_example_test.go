package dataframe

import (
	"fmt"
	"log"

	"github.com/ptiger10/pd/series"
)

func ExampleNew_float64() {
	df, err := New([]interface{}{[]float64{0, 1}}, series.Idx([]string{"foo", "bar"}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//        0
	// foo    0
	// bar    1
	//
}

func ExampleNewWithConfig_float64() {
	df, err := NewWithConfig(Config{Columns: []string{"foobar"}}, []interface{}{[]float64{0, 1}}, series.Idx([]string{"foo", "bar"}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//        foobar
	// foo         0
	// bar         1
	//
}
