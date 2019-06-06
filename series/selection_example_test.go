package series

// func ExampleSeries_Select_all() {
// 	s, _ := New([]int{0, 1, 2}, Idx([]int{0, 1, 2}, opt.Name("foo")), Idx([]string{"A", "B", "C"}, opt.Name("bar")))
// 	sel := s.Select()
// 	fmt.Println(sel)
// 	// Output:
// 	// Selection Object. To print underlying Series, call .Get()
// 	// DerivedIntent: Select All
// 	// Rows: [0 1 2]
// 	// Levels: []
// 	// Error: false
// }

// func ExampleSeries_Select_levels() {
// 	s, _ := New([]int{0, 1, 2}, Idx([]int{0, 1, 2}, opt.Name("foo")), Idx([]string{"A", "B", "C"}, opt.Name("bar")))
// 	sel := s.Select(opt.ByLevels([]int{0}))
// 	fmt.Println(sel)
// 	// Output:
// 	// Selection Object. To print underlying Series, call .Get()
// 	// DerivedIntent: Select Levels
// 	// Rows: []
// 	// Levels: [0]
// 	// Error: false
// }

// func ExampleSeries_Select_rows() {
// 	s, _ := New([]int{0, 1, 2}, Idx([]int{0, 1, 2}, opt.Name("foo")), Idx([]string{"A", "B", "C"}, opt.Name("bar")))
// 	sel := s.Select(opt.ByRows([]int{0}))
// 	fmt.Println(sel)
// 	// Output:
// 	// Selection Object. To print underlying Series, call .Get()
// 	// DerivedIntent: Select Rows
// 	// Rows: [0]
// 	// Levels: []
// 	// Error: false
// }

// func ExampleSeries_Select_xs() {
// 	s, _ := New([]int{0, 1, 2}, Idx([]int{0, 1, 2}, opt.Name("foo")), Idx([]string{"A", "B", "C"}, opt.Name("bar")))
// 	sel := s.Select(opt.ByRows([]int{0}), opt.ByLevels([]int{0}))
// 	fmt.Println(sel)
// 	// Output:
// 	// Selection Object. To print underlying Series, call .Get()
// 	// DerivedIntent: Select Cross-Section
// 	// Rows: [0]
// 	// Levels: [0]
// 	// Error: false
// }
