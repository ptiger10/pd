package cmp

import "time"

// Gt - Greater Than (Numeric)
func Gt(comparison float64) func(float64) bool {
	return func(elem float64) bool {
		return elem > comparison
	}
}

// Gte - Greater Than or Equal To (Numeric)
func Gte(comparison float64) func(float64) bool {
	return func(elem float64) bool {
		return elem >= comparison
	}
}

// Lt - Less Than (Numeric)
func Lt(comparison float64) func(float64) bool {
	return func(elem float64) bool {
		return elem < comparison
	}
}

// Lte - Less Than or Equal To (Numeric)
func Lte(comparison float64) func(float64) bool {
	return func(elem float64) bool {
		return elem <= comparison
	}
}

// Eq - Equal To (Numeric)
func Eq(comparison float64) func(float64) bool {
	return func(elem float64) bool {
		return elem == comparison
	}
}

// Neq - Not Equal To (Numeric)
func Neq(comparison float64) func(float64) bool {
	return func(elem float64) bool {
		return elem != comparison
	}
}

// In - Contained within list (String)
func In(list []string) func(string) bool {
	return func(elem string) bool {
		for _, s := range list {
			if elem == s {
				return true
			}
		}
		return false
	}
}

// Nin - Not contained within list (String)
func Nin(list []string) func(string) bool {
	return func(elem string) bool {
		for _, s := range list {
			if elem == s {
				return false
			}
		}
		return true
	}
}

// True - True, non-null value (Bool)
func True() func(bool) bool {
	return func(elem bool) bool {
		return elem
	}
}

// False - False, non-null value (Bool)
func False() func(bool) bool {
	return func(elem bool) bool {
		return !elem
	}
}

// After - After a specific datetime (DateTime)
func After(t time.Time) func(time.Time) bool {
	return func(elem time.Time) bool {
		return elem.After(t)
	}
}

// Before - Before a specific datetime (DateTime)
func Before(t time.Time) func(time.Time) bool {
	return func(elem time.Time) bool {
		return elem.Before(t)
	}
}
