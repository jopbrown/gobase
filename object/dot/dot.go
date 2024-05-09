package dot

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/jopbrown/gobase/errors"
)

type Tag string

func Get[V any](input any, key string) V {
	return GetWithTag[V](input, key, "json")
}

func TryGet[V any](input any, key string) (V, error) {
	return TryGetWithTag[V](input, key, "json")
}

func GetWithTag[V any](input any, key string, tagName string) V {
	return errors.Must1(TryGetWithTag[V](input, key, tagName))
}

func TryGetWithTag[V any](input any, key, tagName string) (V, error) {
	var none V
	keyPath := strings.Split(key, ".")
	value := reflect.ValueOf(input)
	k := value.Kind()
	if k == reflect.Ptr {
		value = value.Elem()
	}
	k = value.Kind()
	if k != reflect.Struct && k != reflect.Map {
		return none, errors.Errorf("only support map and struct, but got `%s`", k)
	}

	ret, err := deepGet(value, tagName, keyPath)
	if err != nil {
		return none, errors.ErrorAt(err)
	}

	i := ret.Interface()
	v, ok := i.(V)
	if !ok {
		return none, errors.Errorf("`%s` is not the type of `%T` but got `%T`", key, none, i)
	}

	return v, nil
}

func deepGet(value reflect.Value, tagName string, keyPath []string) (reflect.Value, error) {
	name := keyPath[0]
	arrName, arrlDX, isArrayNotation := deepGetArray(keyPath[0])

	if isArrayNotation {
		name = arrName
	}

	innerVal, err := getField(value, name, tagName)
	if err != nil {
		return reflect.Value{}, errors.ErrorAt(err)
	}
	inner := reflect.ValueOf(innerVal)

	if isArrayNotation {
		if k := inner.Kind(); k != reflect.Array && k != reflect.Slice {
			return reflect.Value{}, errors.Errorf("the field is not array or slice: %s", name)
		}
		if arrlDX < 0 || inner.Len() <= arrlDX {
			return reflect.Value{}, errors.Errorf("access array out of range(len=%d): %d", inner.Len(), arrlDX)
		}
		inner = inner.Index(arrlDX)
	}
	if len(keyPath) == 1 {
		return inner, nil
	}

	return deepGet(inner, tagName, keyPath[1:])
}

func getField(value reflect.Value, name, tagName string) (any, error) {
	k := value.Kind()
	switch k {
	case reflect.Struct, reflect.Ptr:
		return getStructField(value, name, tagName)
	case reflect.Map:
		return getMapField(value, name)
	case reflect.Interface:
		return getField(value.Elem(), name, tagName)
	default:
		return reflect.Value{}, errors.Errorf("parent's kind not of supported to get `%s`: %v", k, name)
	}
}

func getMapField(value reflect.Value, name string) (any, error) {
	iter := value.MapRange()
	for iter.Next() {
		key := iter.Key()
		if key.String() == name {
			return iter.Value().Interface(), nil
		}
	}

	return nil, errors.Errorf("map key not found: %s", name)
}

func getStructField(value reflect.Value, name, tagName string) (any, error) {
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	t := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := t.Field(i)
		var fieldName string

		tag := field.Tag.Get(tagName)
		if tagName != "" && tag != "" {
			fieldName = strings.Split(tag, ",")[0]
		} else {
			fieldName = field.Name
		}

		if fieldName == name {
			return value.Field(i).Interface(), nil
		}
	}
	return reflect.Value{}, errors.Errorf("struct field not found: %s", name)
}

var pattArrayNotation = regexp.MustCompile(`^(\w+)\[(-?\d+)\]$`)

func deepGetArray(key string) (string, int, bool) {
	match := pattArrayNotation.FindStringSubmatch(key)
	if len(match) == 0 {
		return "", 0, false
	}

	name := match[1]

	index, err := strconv.Atoi(match[2])
	if err != nil {
		return "", 0, false
	}
	return name, index, true
}
