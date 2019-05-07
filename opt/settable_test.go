package opt

import (
	"reflect"
	"testing"
)

func TestSetting(t *testing.T) {
	SetDisplayIndexMaxWidth(15)
	if GetDisplayIndexMaxWidth() != 15 {
		t.Error("Unable to set/get DisplayIndexMaxWidth")
	}
	SetDisplayIndexWhitespaceBuffer(5)
	if GetDisplayIndexWhitespaceBuffer() != 5 {
		t.Error("Unable to set/get DisplayIndexWhitespaceBuffer")
	}
	SetDisplayValuesWhitespaceBuffer(10)
	if GetDisplayValuesWhitespaceBuffer() != 10 {
		t.Error("Unable to set/get DisplayValuesWhitespaceBuffer")
	}
	SetDisplayFloatPrecision(10)
	if GetDisplayFloatPrecision() != 10 {
		t.Error("Unable to set/get DisplayFloatPrecision")
	}

	SetDisplayRepeatedIndexLabels(true)
	if GetDisplayRepeatedIndexLabels() != true {
		t.Error("Unable to set/get DisplayRepeatedIndexLabels")
	}

	SetDisplayStringNullFiller("Nothing")
	if GetDisplayStringNullFiller() != "Nothing" {
		t.Error("Unable to set/get DisplayStringNullFiller")
	}

	SetStringNullValues([]string{"Nada", "Nothing"})
	if !reflect.DeepEqual(GetStringNullValues(), []string{"Nada", "Nothing"}) {
		t.Error("Unable to set/get StringNullValues")
	}
	SetLogWarnings(false)
	if GetLogWarnings() != false {
		t.Error("Unable to set/get LogWarnings")
	}
	RestoreDefaults()
	if GetDisplayIndexMaxWidth() != 25 {
		t.Error("Unable to restore default for DisplayIndexMaxWidth")
	}

}
