package pd

import (
	"fmt"
	"log"
)

func ExampleSeries_defaultIndex() {
	s, err := Series([]string{"foo", "bar", "baz"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
	// Output:
	// 0    foo
	// 1    bar
	// 2    baz
	// datatype: string
}

func ExampleDataFrame_default() {
	df, err := DataFrame([]interface{}{[]string{"foo", "bar", "baz"}, []int{7, 11, 19}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	// Output:
	//        0  1
	// 0    foo  7
	// 1    bar 11
	// 2    baz 19
}
