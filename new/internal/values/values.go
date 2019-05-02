package values

// The Values interface is the primary means of handling a collection of values
// Thes same interface and value types are used for both Series values and Index labels
type Values interface {
	// returns all values regardless of null status
	All() []interface{}
	Describe() string
	In([]int) Values
	Len() int
	// Element returns a value element at an integer position in the form []interface{} {val, null}
	Element(int) []interface{}

	ToFloat() Values
	ToInt() Values
	ToString() Values
	ToBool() Values
	ToDateTime() Values
}
