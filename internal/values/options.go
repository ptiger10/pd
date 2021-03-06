package values

var displayValuesWhitespaceBuffer = 4
var displayColumnsWhitespaceBuffer = 2
var displayElementWhitespaceBuffer = 1
var displayIndexWhitespaceBuffer = 1
var multiColNameSeparator = " | "
var interpolationMaximum = 50
var interpolationThreshold = .80

// GetDisplayValuesWhitespaceBuffer returns displayValuesWhitespaceBuffer.
// displayValuesWhitespaceBuffer is an option when printing a Series or DataFrame.
// It is the number of spaces between the last level of index labels
// and the first collection of values. In a Series, there is only one collection of values.
// In a DataFrame, the first collection of values is the first Series.
//
// Default buffer: 4 spaces
func GetDisplayValuesWhitespaceBuffer() int {
	return displayValuesWhitespaceBuffer
}

// GetDisplayColumnsWhitespaceBuffer returns displayColumnsWhitespaceBuffer.
// displayColumnsWhitespaceBuffer is an option when printing a Series or DataFrame.
// It is the number of spaces between columns in a DataFrame.
//
// Default buffer: 2 spaces
func GetDisplayColumnsWhitespaceBuffer() int {
	return displayColumnsWhitespaceBuffer
}

// GetDisplayElementWhitespaceBuffer returns displayElementWhitespaceBuffer.
// DisplayElementWhitespaceBuffer is an option when printing an Element.
// It is the number of spaces between the last level of index labels and the first value.
//
// // Default buffer: 1 space
func GetDisplayElementWhitespaceBuffer() int {
	return displayElementWhitespaceBuffer
}

// GetDisplayIndexWhitespaceBuffer returns displayIndexWhitespaceBuffer.
// DisplayIndexWhitespaceBuffer is an option when printing a Series.
// It is the number of spaces between index labels. This applies only to a multi-level index.
//
// Default buffer: 1 space
func GetDisplayIndexWhitespaceBuffer() int {
	return displayIndexWhitespaceBuffer
}

// GetMultiColNameSeparator returns the multiColNameSeparator.
// The multiColNameSeparator separates col names whenever a multicol is concatenated together (e.g., into a Series name or index level name).
//
// Default: " | "
func GetMultiColNameSeparator() string {
	return multiColNameSeparator
}

// GetInterpolationMaximum returns the max number of records that will be checked during an interpolation check.
//
// Default: 50
func GetInterpolationMaximum() int {
	return interpolationMaximum
}

// GetInterpolationThreshold returns the ratio of type inclusion required for a dataType to be interpolated.
//
// Default: .80
func GetInterpolationThreshold() float64 {
	return interpolationThreshold
}
