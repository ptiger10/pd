package series

import "fmt"

func ExampleSeries_string_defaultIndex() {
	s, _ := New([]string{"Joe", "Jamy", "", "Jenny"})
	fmt.Println(s)

	// Output:
	// 0 Joe
	// 1 Jamy
	// 2
	// 3 Jenny
}

func ExampleSeries_string_customIndex() {
	s, _ := New([]string{"Joe", "Jamy", "", "Jenny"}, Index([]int{100, 101, 102, 103}))
	fmt.Println(s)

	// Output:
	// 100 Joe
	// 101 Jamy
	// 102
	// 103 Jenny
}

func ExampleSeries_string_multiIndex() {
	s, _ := New([]string{"Joe", "Jamy", "", "Jenny"}, Index([]int{0, 1, 2, 3}), Index([]int{100, 101, 102, 103}))
	fmt.Println(s)

	// Output:
	// 0 100 Joe
	// 1 101 Jamy
	// 2 102
	// 3 103 Jenny
}
