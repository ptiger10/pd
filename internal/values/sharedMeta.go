package values

// Generic methods for valueTypeValue type converters

func (val valueTypeValue) toFloat64() float64Value {
	return float64Value{}
}

func (val valueTypeValue) toInt64() int64Value {
	return int64Value{}
}

func (val valueTypeValue) toBool() boolValue {
	return boolValue{}
}

func (val valueTypeValue) toDateTime() dateTimeValue {
	return dateTimeValue{}
}

func (val interfaceValue) tovalueType() valueTypeValue {
	return valueTypeValue{}
}
