package conv

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ToDuration casts an interface to a time.Duration type.
func ToDuration(i interface{}) time.Duration {
	v, _ := ToDurationE(i)
	return v
}

// ToDurationE casts an interface to a time.Duration type.
func ToDurationE(i interface{}) (d time.Duration, err error) {
	switch s := i.(type) {
	case time.Duration:
		return s, nil
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		d = time.Duration(ToInt64(s))
		return
	case float32, float64:
		d = time.Duration(ToFloat64(s))
		return
	case string:
		if strings.ContainsAny(s, "nsuµmh") {
			d, err = time.ParseDuration(s)
		} else {
			d, err = time.ParseDuration(s + "ns")
		}
		return
	case json.Number:
		var v float64
		v, err = s.Float64()
		d = time.Duration(v)
		return
	default:
		err = fmt.Errorf("unable to cast %#v of type %T to Duration", i, i)
		return
	}
}
