package conv

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

var errNegativeNotAllowed = errors.New("unable to cast negative value")

// ToUint casts an interface to a uint type.
func ToUint(i interface{}) uint {
	v, _ := ToUintE(i)
	return v
}

// ToUint64 casts an interface to a uint64 type.
func ToUint64(i interface{}) uint64 {
	v, _ := ToUint64E(i)
	return v
}

// ToUint32 casts an interface to a uint32 type.
func ToUint32(i interface{}) uint32 {
	v, _ := ToUint32E(i)
	return v
}

// ToUint16 casts an interface to a uint16 type.
func ToUint16(i interface{}) uint16 {
	v, _ := ToUint16E(i)
	return v
}

// ToUint8 casts an interface to a uint8 type.
func ToUint8(i interface{}) uint8 {
	v, _ := ToUint8E(i)
	return v
}

// ToUintE casts an interface to a uint type.
func ToUintE(i interface{}) (uint, error) {
	intv, ok := toInt(i)
	if ok {
		if intv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(intv), nil
	}

	switch s := i.(type) {
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if err == nil {
			if v < 0 {
				return 0, errNegativeNotAllowed
			}
			return uint(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint", i, i)
	case json.Number:
		return ToUintE(string(s))
	case int64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(s), nil
	case int32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(s), nil
	case int16:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(s), nil
	case int8:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(s), nil
	case uint:
		return s, nil
	case uint64:
		return uint(s), nil
	case uint32:
		return uint(s), nil
	case uint16:
		return uint(s), nil
	case uint8:
		return uint(s), nil
	case float64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(s), nil
	case float32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(s), nil
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint", i, i)
	}
}

// ToUint64E casts an interface to a uint64 type.
func ToUint64E(i interface{}) (uint64, error) {
	intv, ok := toInt(i)
	if ok {
		if intv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(intv), nil
	}

	switch s := i.(type) {
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if err == nil {
			if v < 0 {
				return 0, errNegativeNotAllowed
			}
			return uint64(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint64", i, i)
	case json.Number:
		return ToUint64E(string(s))
	case int64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(s), nil
	case int32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(s), nil
	case int16:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(s), nil
	case int8:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(s), nil
	case uint:
		return uint64(s), nil
	case uint64:
		return s, nil
	case uint32:
		return uint64(s), nil
	case uint16:
		return uint64(s), nil
	case uint8:
		return uint64(s), nil
	case float32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(s), nil
	case float64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(s), nil
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint64", i, i)
	}
}

// ToUint32E casts an interface to a uint32 type.
func ToUint32E(i interface{}) (uint32, error) {
	intv, ok := toInt(i)
	if ok {
		if intv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(intv), nil
	}

	switch s := i.(type) {
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if err == nil {
			if v < 0 {
				return 0, errNegativeNotAllowed
			}
			return uint32(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint32", i, i)
	case json.Number:
		return ToUint32E(string(s))
	case int64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(s), nil
	case int32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(s), nil
	case int16:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(s), nil
	case int8:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(s), nil
	case uint:
		return uint32(s), nil
	case uint64:
		return uint32(s), nil
	case uint32:
		return s, nil
	case uint16:
		return uint32(s), nil
	case uint8:
		return uint32(s), nil
	case float64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(s), nil
	case float32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(s), nil
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint32", i, i)
	}
}

// ToUint16E casts an interface to a uint16 type.
func ToUint16E(i interface{}) (uint16, error) {
	intv, ok := toInt(i)
	if ok {
		if intv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(intv), nil
	}

	switch s := i.(type) {
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if err == nil {
			if v < 0 {
				return 0, errNegativeNotAllowed
			}
			return uint16(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint16", i, i)
	case json.Number:
		return ToUint16E(string(s))
	case int64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(s), nil
	case int32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(s), nil
	case int16:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(s), nil
	case int8:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(s), nil
	case uint:
		return uint16(s), nil
	case uint64:
		return uint16(s), nil
	case uint32:
		return uint16(s), nil
	case uint16:
		return s, nil
	case uint8:
		return uint16(s), nil
	case float64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(s), nil
	case float32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(s), nil
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint16", i, i)
	}
}

// ToUint8E casts an interface to a uint type.
func ToUint8E(i interface{}) (uint8, error) {
	intv, ok := toInt(i)
	if ok {
		if intv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(intv), nil
	}

	switch s := i.(type) {
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if err == nil {
			if v < 0 {
				return 0, errNegativeNotAllowed
			}
			return uint8(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint8", i, i)
	case json.Number:
		return ToUint8E(string(s))
	case int64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(s), nil
	case int32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(s), nil
	case int16:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(s), nil
	case int8:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(s), nil
	case uint:
		return uint8(s), nil
	case uint64:
		return uint8(s), nil
	case uint32:
		return uint8(s), nil
	case uint16:
		return uint8(s), nil
	case uint8:
		return s, nil
	case float64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(s), nil
	case float32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(s), nil
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint8", i, i)
	}
}
