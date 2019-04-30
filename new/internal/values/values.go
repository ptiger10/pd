package values

type Values interface {
	Describe() string
	In([]int) interface{}
}
