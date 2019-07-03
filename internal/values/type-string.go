package values

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/ptiger10/pd/options"
)

// [START Constructor Functions]

func isNullString(s string) bool {
	nullStrings := options.GetStringNullValues()
	for _, ns := range nullStrings {
		if strings.TrimSpace(s) == ns {
			return true
		}
	}
	return false
}

// newString creates an stringValue from atomic string value
func newString(val string) stringValue {
	if isNullString(val) {
		return stringValue{options.GetDisplayStringNullFiller(), true}
	}
	return stringValue{val, false}
}

func (vals *stringValues) Less(i, j int) bool {
	if (*vals)[i].v < (*vals)[j].v {
		return true
	}
	return false
}

// [END Constructor Functions]

// [START Converters]

// toFloat converts stringValue to float64Value
//
// "1": 1.0, Null: NaN
func (val stringValue) toFloat64() float64Value {
	f, err := strconv.ParseFloat(val.v, 64)
	if math.IsNaN(f) || err != nil {
		return float64Value{math.NaN(), true}
	}
	return float64Value{f, false}
}

// toInt converts stringValue to int64Value
//
// "1": 1, null: NaN
func (val stringValue) toInt64() int64Value {
	if val.null {
		return int64Value{0, true}
	}
	f, err := strconv.ParseFloat(val.v, 64)
	if err != nil {
		return int64Value{0, true}
	}
	return int64Value{int64(f), false}
}

func (val stringValue) toString() stringValue {
	if isNullString(val.v) || val.null {
		return stringValue{options.GetDisplayStringNullFiller(), true}
	}
	return stringValue{fmt.Sprint(val.v), false}
}

// toBool converts stringValue to boolValue
//
// null: false; notnull: true
func (val stringValue) toBool() boolValue {
	if val.null {
		return boolValue{false, true}
	}
	return boolValue{true, false}
}

// toDateTime converts stringValue to dateTimeValue using an external parse library
//
// Jan 1 2019: 2019-01-01 00:00:00
//
// Acceptable DateTime string formats
/*
   	"May 8, 2009 5:57:51 PM",
   	"oct 7, 1970",
   	"oct 7, '70",
   	"oct. 7, 1970",
   	"oct. 7, 70",
   	"Mon Jan  2 15:04:05 2006",
   	"Mon Jan  2 15:04:05 MST 2006",
   	"Mon Jan 02 15:04:05 -0700 2006",
   	"Monday, 02-Jan-06 15:04:05 MST",
   	"Mon, 02 Jan 2006 15:04:05 MST",
   	"Tue, 11 Jul 2017 16:28:13 +0200 (CEST)",
   	"Mon, 02 Jan 2006 15:04:05 -0700",
   	"Thu, 4 Jan 2018 17:53:36 +0000",
   	"Mon Aug 10 15:44:11 UTC+0100 2015",
   	"Fri Jul 03 2015 18:04:07 GMT+0100 (GMT Daylight Time)",
   	"September 17, 2012 10:09am",
   	"September 17, 2012 at 10:09am PST-08",
   	"September 17, 2012, 10:10:09",
   	"October 7, 1970",
   	"October 7th, 1970",
   	"12 Feb 2006, 19:17",
   	"12 Feb 2006 19:17",
   	"7 oct 70",
   	"7 oct 1970",
   	"03 February 2013",
   	"1 July 2013",
   	"2013-Feb-03",
   	//   mm/dd/yy
   	"3/31/2014",
   	"03/31/2014",
   	"08/21/71",
   	"8/1/71",
   	"4/8/2014 22:05",
   	"04/08/2014 22:05",
   	"4/8/14 22:05",
   	"04/2/2014 03:00:51",
   	"8/8/1965 12:00:00 AM",
   	"8/8/1965 01:00:01 PM",
   	"8/8/1965 01:00 PM",
   	"8/8/1965 1:00 PM",
   	"8/8/1965 12:00 AM",
   	"4/02/2014 03:00:51",
   	"03/19/2012 10:11:59",
   	"03/19/2012 10:11:59.3186369",
   	// yyyy/mm/dd
   	"2014/3/31",
   	"2014/03/31",
   	"2014/4/8 22:05",
   	"2014/04/08 22:05",
   	"2014/04/2 03:00:51",
   	"2014/4/02 03:00:51",
   	"2012/03/19 10:11:59",
   	"2012/03/19 10:11:59.3186369",
   	// Chinese
   	"2014年04月08日",
   	//   yyyy-mm-ddThh
   	"2006-01-02T15:04:05+0000",
   	"2009-08-12T22:15:09-07:00",
   	"2009-08-12T22:15:09",
   	"2009-08-12T22:15:09Z",
   	//   yyyy-mm-dd hh:mm:ss
   	"2014-04-26 17:24:37.3186369",
   	"2012-08-03 18:31:59.257000000",
   	"2014-04-26 17:24:37.123",
   	"2013-04-01 22:43",
   	"2013-04-01 22:43:22",
   	"2014-12-16 06:20:00 UTC",
   	"2014-12-16 06:20:00 GMT",
   	"2014-04-26 05:24:37 PM",
   	"2014-04-26 13:13:43 +0800",
   	"2014-04-26 13:13:43 +0800 +08",
   	"2014-04-26 13:13:44 +09:00",
   	"2012-08-03 18:31:59.257000000 +0000 UTC",
   	"2015-09-30 18:48:56.35272715 +0000 UTC",
   	"2015-02-18 00:12:00 +0000 GMT",
   	"2015-02-18 00:12:00 +0000 UTC",
   	"2015-02-08 03:02:00 +0300 MSK m=+0.000000001",
   	"2015-02-08 03:02:00.001 +0300 MSK m=+0.000000001",
   	"2017-07-19 03:21:51+00:00",
   	"2014-04-26",
   	"2014-04",
   	"2014",
   	"2014-05-11 08:20:13,787",
   	// mm.dd.yy
   	"3.31.2014",
   	"03.31.2014",
   	"08.21.71",
   	"2014.03",
   	"2014.03.30",
   	//  yyyymmdd and similar
   	"20140601",
   	"20140722105203",
   	// unix seconds, ms, micro, nano
   	"1332151919",
   	"1384216367189",
   	"1384216367111222",
   	"1384216367111222333",
   }
*/
func (val stringValue) toDateTime() dateTimeValue {
	if val.null {
		return dateTimeValue{time.Time{}, true}
	}
	t, err := dateparse.ParseAny(val.v)
	if err != nil {
		return dateTimeValue{time.Time{}, true}
	}
	return dateTimeValue{t, false}
}

// [END Converters]
