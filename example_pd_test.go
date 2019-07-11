package pd

import (
	"fmt"
)

func ExampleSeries_defaultIndex() {
	s, _ := Series([]string{"foo", "bar", "baz"})
	fmt.Println(s)
	// Output:
	// 0    foo
	// 1    bar
	// 2    baz
	//
	// datatype: string
}

func ExampleDataFrame_default() {
	df, _ := DataFrame([]interface{}{[]string{"foo", "bar", "baz"}, []int{7, 11, 19}})
	fmt.Println(df)
	// Output:
	//        0   1
	// 0    foo   7
	// 1    bar  11
	// 2    baz  19
}
