package series

// Len returns the length of the Series (including null values)
func (s Series) Len() int {
	return s.Values.Len()
}
