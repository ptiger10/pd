package options

var defaultOptions = struct {
	DisplayIndexMaxWidth          int
	DisplayValuesWhitespaceBuffer int
	DisplayIndexWhitespaceBuffer  int
	DisplayFloatPrecision         int
	DisplayRepeatedIndexLabels    bool
	DisplayStringNullFiller       string
	StringNullValues              []string
	LogWarnings                   bool
}{
	DisplayIndexMaxWidth,
	DisplayValuesWhitespaceBuffer,
	DisplayIndexWhitespaceBuffer,
	DisplayFloatPrecision,
	DisplayRepeatedIndexLabels,
	DisplayStringNullFiller,
	StringNullValues,
	LogWarnings,
}

// RestoreDefaults resets options back to their default setting
func RestoreDefaults() {
	SetDisplayIndexMaxWidth(defaultOptions.DisplayIndexMaxWidth)
	SetDisplayValuesWhitespaceBuffer(defaultOptions.DisplayValuesWhitespaceBuffer)
	SetDisplayIndexWhitespaceBuffer(DisplayIndexWhitespaceBuffer)
	SetDisplayFloatPrecision(DisplayFloatPrecision)
	SetDisplayRepeatedIndexLabels(DisplayRepeatedIndexLabels)
	SetDisplayStringNullFiller(DisplayStringNullFiller)
	SetStringNullValues(StringNullValues)
	SetLogWarnings(LogWarnings)
}

// DisplayIndexMaxWidth is an option when printing a Series.
// It is the widest allowable character width for an index label.
// If a label is longer than the max, it will be elided at the end.
var DisplayIndexMaxWidth = 25

// DisplayValuesWhitespaceBuffer is an option when printing a Series.
// It is the number of spaces between the last level of index labels
// and the first collection of values. In a Series, there is only one collection of values.
// In a DataFrame, the first collection of values is the first Series.
var DisplayValuesWhitespaceBuffer = 4

// DisplayIndexWhitespaceBuffer is an option when printing a Series.
// It is the number of spaces between index labels. This applies only to a multi-level index.
var DisplayIndexWhitespaceBuffer = 1

// DisplayFloatPrecision is an option when printing a Series.
// It is the number of decimal points in floating point values and index labels.
var DisplayFloatPrecision = 2

// DisplayRepeatedIndexLabels is an option when printing a Series.
// If true, all index labels will be shown, like so:
//
// A 0    foo
//
// B 0    bar
//
// C 1    baz
//
// If false, repeated index labels in the same level will be excluded, like so:
//
// A 0    foo
//
// B ... bar
//
// C 1    baz
//
// NB: ellipsis not included in actual printing
var DisplayRepeatedIndexLabels = false

// DisplayStringNullFiller is an option when printing a Series.
// It is how null string values are represented.
var DisplayStringNullFiller = "NaN"

// StringNullValues is an option when constructing or converting a Series.
// It is the list of string values that are considered null.
var StringNullValues = []string{"NaN", "n/a", "N/A", "", "nil"}

// LogWarnings is an option when executing functions within this module.
// If true, non-returned errors are logged to stderr.
// This is relevant for many common exploratory methods, which are often chained together and therefore not designed to return an error value.
var LogWarnings = true

// SetDisplayIndexMaxWidth sets DisplayIndexMaxWidth to n characters.
func SetDisplayIndexMaxWidth(n int) {
	DisplayIndexMaxWidth = n
}

// SetDisplayValuesWhitespaceBuffer sets DisplayValuesWhitespaceBuffer to n spaces.
func SetDisplayValuesWhitespaceBuffer(n int) {
	DisplayValuesWhitespaceBuffer = n
}

// SetDisplayIndexWhitespaceBuffer sets DisplayIndexWhitespaceBuffer to n spaces.
func SetDisplayIndexWhitespaceBuffer(n int) {
	DisplayIndexWhitespaceBuffer = n
}

// SetDisplayFloatPrecision sets DisplayFloatPrecision to n decimal places.
func SetDisplayFloatPrecision(n int) {
	DisplayFloatPrecision = n
}

// SetDisplayStringNullFiller sets DisplayStringNullFiller to "s".
func SetDisplayStringNullFiller(s string) {
	DisplayStringNullFiller = s
}

// SetDisplayRepeatedIndexLabels sets DisplayRepeatedIndexLabels to boolean.
func SetDisplayRepeatedIndexLabels(boolean bool) {
	DisplayRepeatedIndexLabels = boolean
}

// SetStringNullValues sets StringNullValues to include only those items contained in nullList.
func SetStringNullValues(nullList []string) {
	StringNullValues = nullList
}

// SetLogWarnings sets LogWarnings to boolean.
func SetLogWarnings(boolean bool) {
	LogWarnings = boolean
}
