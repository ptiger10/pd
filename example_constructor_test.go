package pd

// // Index -------------------------------------------------------------
// // Int ---------------------------------------------------------------
// var arrayInt = []int{1, -1, 0}
// var sInt = Series(arrayInt)

// func ExampleSeries_int() {
// 	fmt.Print(sInt)
// 	// Output:
// 	// 0    1
// 	// 1    -1
// 	// 2    0
// 	// dtype: int
// }
// func ExampleSeries_int_values() {
// 	fmt.Print(sInt.Values)
// 	// Output:
// 	// [1 -1 0]
// }

// func ExampleSeries_int_index() {
// 	var sInt = Series(arrayInt)
// 	fmt.Print(sInt.Index)
// 	// Output:
// 	// [0 1 2] dtype: int
// }

// func ExampleSeries_int_at() {
// 	fmt.Print(sInt.At(1))
// 	// Output:
// 	// -1
// }

// // Custom String ---------------------------------------------------------------
// var sStr = Series(arrayInt, []interface{}{"positive", "negative", "neutralButLong"})

// func ExampleSeries_str() {
// 	fmt.Print(sStr)
// 	// Output:
// 	// positive          1
// 	// negative          -1
// 	// neutralButLong    0
// 	// dtype: int
// }

// func ExampleSeries_str_index() {
// 	fmt.Print(sStr.Index)
// 	// Output:
// 	// [positive negative neutralButLong] dtype: string
// }

// func ExampleSeries_str_at_1() {
// 	fmt.Print(sStr.At(1))
// 	// Output: -1
// }

// func ExampleSeries_str_at_2() {
// 	fmt.Print(sStr.At("negative"))
// 	// Output: -1
// }

// // Custom Int ---------------------------------------------------------------
// var sCustomInt = Series(arrayInt, []interface{}{3, 4, 5})

// func ExampleSeries_int_idx_int() {
// 	fmt.Print(sCustomInt)
// 	// Output:
// 	// 3    1
// 	// 4    -1
// 	// 5    0
// 	// dtype: int
// }

// func ExampleSeries_int_idx_int_index() {
// 	fmt.Print(sCustomInt.Index)
// 	// Output:
// 	// [3 4 5] dtype: int
// }

// func ExampleSeries_int_idx_int_at() {
// 	fmt.Print(sCustomInt.At(4))
// 	// Output: -1
// }
