package conv

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ToBool casts an interface to a bool type.
func ToBool(i interface{}) bool {
	v, _ := ToBoolE(i)
	return v
}

// ToBoolE casts an interface to a bool type.
func ToBoolE(i interface{}) (bool, error) {
	switch b := i.(type) {
	case bool:
		return b, nil
	case nil:
		return false, nil
	case int:
		if i.(int) != 0 {
			return true, nil
		}
		return false, nil
	case string:
		return strconv.ParseBool(i.(string))
	case json.Number:
		v, err := ToInt64E(b)
		if err == nil {
			return v != 0, nil
		}
		return false, fmt.Errorf("unable to cast %#v of type %T to bool", i, i)
	default:
		return false, fmt.Errorf("unable to cast %#v of type %T to bool", i, i)
	}
}
