package options

import (
	"reflect"
	"testing"
)

func TestSetting(t *testing.T) {
	SetDisplayIndexMaxWidth(15)
	if DisplayIndexMaxWidth != 15 {
		t.Error("Unable to set DisplayIndexMaxWidth")
	}
	SetDisplayIndexWhitespaceBuffer(5)
	if DisplayIndexWhitespaceBuffer != 5 {
		t.Error("Unable to set DisplayIndexWhitespaceBuffer")
	}
	SetDisplayValuesWhitespaceBuffer(10)
	if DisplayValuesWhitespaceBuffer != 10 {
		t.Error("Unable to set DisplayValuesWhitespaceBuffer")
	}
	SetDisplayFloatPrecision(10)
	if DisplayFloatPrecision != 10 {
		t.Error("Unable to set DisplayFloatPrecision")
	}

	SetDisplayRepeatedIndexLabels(true)
	if DisplayRepeatedIndexLabels != true {
		t.Error("Unable to set DisplayRepeatedIndexLabels")
	}

	SetDisplayStringNullFiller("Nothing")
	if DisplayStringNullFiller != "Nothing" {
		t.Error("Unable to set DisplayStringNullFiller")
	}

	SetStringNullValues([]string{"Nada", "Nothing"})
	if !reflect.DeepEqual(StringNullValues, []string{"Nada", "Nothing"}) {
		t.Error("Unable to set StringNullValues")
	}
	SetLogWarnings(false)
	if LogWarnings != false {
		t.Error("Unable to set LogWarnings")
	}
	RestoreDefaults()
	if DisplayIndexMaxWidth != 25 {
		t.Error("Unable to restore default for DisplayIndexMaxWidth")
	}

}
