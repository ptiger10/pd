package series

import (
	"fmt"
)

func ExampleSeries_Select_rows() {
	s, _ := New(
		[]int{0, 1, 2},
		IndexLevel{Labels: []int{0, 1, 2}, Name: "foo"},
		IndexLevel{Labels: []string{"A", "B", "C"}, Name: "bar"},
	)
	sel := s.Select.ByRows([]int{0, 2})
	fmt.Println(sel)
	// Output:
	// Selection{rows: [0 2], levels: [], swappable: true, hasError: false}
}

// func ExampleSeries_Select_levels() {
// 	s, _ := New([]int{0, 1, 2}, Idx([]int{0, 1, 2}, options.Name("foo")), Idx([]string{"A", "B", "C"}, options.Name("bar")))
// 	sel := s.Select(options.ByLevels([]int{0}))
// 	fmt.Println(sel)
// 	// Output:
// 	// Selection Object. To print underlying Series, call .Get()
// 	// DerivedIntent: Select Levels
// 	// Rows: []
// 	// Levels: [0]
// 	// Error: false
// }

// func ExampleSeries_Select_rows() {
// 	s, _ := New([]int{0, 1, 2}, Idx([]int{0, 1, 2}, options.Name("foo")), Idx([]string{"A", "B", "C"}, options.Name("bar")))
// 	sel := s.Select(options.ByRows([]int{0}))
// 	fmt.Println(sel)
// 	// Output:
// 	// Selection Object. To print underlying Series, call .Get()
// 	// DerivedIntent: Select Rows
// 	// Rows: [0]
// 	// Levels: []
// 	// Error: false
// }

// func ExampleSeries_Select_xs() {
// 	s, _ := New([]int{0, 1, 2}, Idx([]int{0, 1, 2}, options.Name("foo")), Idx([]string{"A", "B", "C"}, options.Name("bar")))
// 	sel := s.Select(options.ByRows([]int{0}), options.ByLevels([]int{0}))
// 	fmt.Println(sel)
// 	// Output:
// 	// Selection Object. To print underlying Series, call .Get()
// 	// DerivedIntent: Select Cross-Section
// 	// Rows: [0]
// 	// Levels: [0]
// 	// Error: false
// }
