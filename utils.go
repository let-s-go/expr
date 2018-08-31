package expr

import (
	"fmt"
	"reflect"
	"strconv"
)

func isBool1(val interface{}) bool {
	return val != nil && reflect.TypeOf(val).Kind() == reflect.Bool
}

func toBool(val interface{}) (bool, bool) {
	if b, ok := val.(bool); ok {
		return b, true
	}
	if v, err := cast(val); err == nil {
		return v > 0, true
	}
	return false, false
}

func isText(val interface{}) bool {
	return val != nil && reflect.TypeOf(val).Kind() == reflect.String
}

func toText(val interface{}) string {
	return reflect.ValueOf(val).String()
}

func equal(left, right interface{}) bool {
	if l, err := cast(left); err == nil {
		if r, err := cast(right); err == nil {
			return l == r
		}
		return false
	}
	return reflect.DeepEqual(left, right)
}

func isNumber(val interface{}) bool {
	return val != nil && reflect.TypeOf(val).Kind() == reflect.Float64
}

func cast(v interface{}) (float64, error) {
	switch t := v.(type) {
	case float32:
		return float64(t), nil
	case float64:
		return t, nil
	case int:
		return float64(t), nil
	case int8:
		return float64(t), nil
	case int16:
		return float64(t), nil
	case int32:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case uint:
		return float64(t), nil
	case uint8:
		return float64(t), nil
	case uint16:
		return float64(t), nil
	case uint32:
		return float64(t), nil
	case uint64:
		return float64(t), nil
	case string:
		return strconv.ParseFloat(t, 64)
	}
	//	if v != nil {
	//		switch reflect.TypeOf(v).Kind() {
	//		case reflect.Float32, reflect.Float64:
	//			return v.(float64), nil

	//		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	//			return float64(reflect.ValueOf(v).Int()), nil

	//		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	//			return float64(reflect.ValueOf(v).Uint()), nil // TODO: Check if uint64 fits into float64.
	//		}
	//	}
	return 0, fmt.Errorf("can't cast %T to float64", v)
}

func canBeNumber(v interface{}) bool {
	if v != nil {
		return isNumberType(reflect.TypeOf(v))
	}
	return false
}

func extract(from interface{}, it string) (interface{}, error) {
	if from != nil {
		if m, ok := from.(map[string]interface{}); ok {
			return m[it], nil
		}
		switch reflect.TypeOf(from).Kind() {
		case reflect.Map:
			value := reflect.ValueOf(from).MapIndex(reflect.ValueOf(it))
			if value.IsValid() && value.CanInterface() {
				return value.Interface(), nil
			}
			return nil, nil
		case reflect.Struct:
			value := reflect.ValueOf(from).FieldByName(it)
			if value.IsValid() && value.CanInterface() {
				return value.Interface(), nil
			}
			return nil, nil
		case reflect.Ptr:
			value := reflect.ValueOf(from).Elem()
			if value.IsValid() && value.CanInterface() {
				return extract(value.Interface(), it)
			}
		}
	}
	return nil, fmt.Errorf("can't get %q from %T", it, from)
}

func extractIt(from interface{}, it interface{}) (interface{}, error) {
	if from != nil {
		if m, ok := from.(map[string]interface{}); ok {
			if k, ok := it.(string); ok {
				return m[k], nil
			}
		}
		switch reflect.TypeOf(from).Kind() {
		case reflect.Array, reflect.Slice, reflect.String:
			i, err := cast(it)
			if err != nil {
				return nil, err
			}

			value := reflect.ValueOf(from).Index(int(i))
			if value.IsValid() && value.CanInterface() {
				return value.Interface(), nil
			}
		case reflect.Map:
			value := reflect.ValueOf(from).MapIndex(reflect.ValueOf(it))
			if value.IsValid() && value.CanInterface() {
				return value.Interface(), nil
			}
			return nil, nil
		case reflect.Struct:
			value := reflect.ValueOf(from).FieldByName(reflect.ValueOf(it).String())
			if value.IsValid() && value.CanInterface() {
				return value.Interface(), nil
			}
			return nil, nil
		case reflect.Ptr:
			value := reflect.ValueOf(from).Elem()
			if value.IsValid() && value.CanInterface() {
				return extractIt(value.Interface(), it)
			}
		}
	}
	return nil, fmt.Errorf("can't get %q from %T", it, from)
}

func contains(needle interface{}, array interface{}) (bool, error) {
	if array != nil {
		value := reflect.ValueOf(array)
		switch reflect.TypeOf(array).Kind() {
		case reflect.Array, reflect.Slice:
			for i := 0; i < value.Len(); i++ {
				value := value.Index(i)
				if value.IsValid() && value.CanInterface() {
					if equal(value.Interface(), needle) {
						return true, nil
					}
				}
			}
			return false, nil
		}
		return false, fmt.Errorf("operator in not defined on %T", array)
	}
	return false, nil
}
