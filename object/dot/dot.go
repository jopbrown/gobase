package dot

import (
	"fmt"
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
	fields, err := fetchFields(key, input, tagName)
	if err != nil {
		return none, errors.ErrorAt(err)
	}

	i := fields[len(fields)-1].Inner.Interface()
	v, ok := i.(V)
	if !ok {
		return none, errors.Errorf("`%s` is not the type of `%T` but got `%T`", key, none, i)
	}

	return v, nil
}

func Set[V any](input any, key string, v V) {
	SetWithTag[V](input, key, "json", v)
}

func TrySet[V any](input any, key string, v V) error {
	return TrySetWithTag[V](input, key, "json", v)
}

func SetWithTag[V any](input any, key string, tagName string, v V) {
	errors.Must(TrySetWithTag[V](input, key, tagName, v))
}

func TrySetWithTag[V any](input any, key, tagName string, v V) error {
	fields, err := fetchFields(key, input, tagName)
	if err != nil {
		return errors.ErrorAt(err)
	}

	parent, child := fields[len(fields)-2], fields[len(fields)-1]
	pv, pk := parent.Inner, parent.Inner.Kind()
	if pk == reflect.Map {
		iter := pv.MapRange()
		for iter.Next() {
			if iter.Key().String() == child.Name {
				pv.SetMapIndex(iter.Key(), reflect.ValueOf(v))
				break
			}
		}
		return nil
	}

	if pk == reflect.Slice || pk == reflect.Array {
		err = setFieldValue(child.Outter, v)
		if err != nil {
			return errors.ErrorAt(err)
		}
		return nil
	}

	err = setFieldValue(child.Inner, v)
	if err != nil {
		return errors.ErrorAt(err)
	}

	return nil
}

type reflectField struct {
	Name   string
	Outter reflect.Value
	Inner  reflect.Value
}

func fetchFields(key string, input any, tagName string) ([]*reflectField, error) {
	keyElems := splitKey(key)
	rv := reflect.ValueOf(input)
	rv = getInnerElem(rv)

	k := rv.Kind()
	if k != reflect.Struct && k != reflect.Map && k != reflect.Array && k != reflect.Slice {
		return nil, errors.Errorf("top level only support map, struct, array and slice, but got `%s`", k)
	}

	fields, err := fetchElemFields(rv, tagName, keyElems)
	if err != nil {
		return nil, errors.ErrorAt(err)
	}

	return fields, nil
}

func fetchElemFields(rv reflect.Value, tagName string, keyElems []string) ([]*reflectField, error) {
	fields := make([]*reflectField, 0, len(keyElems)+1)
	fields = append(fields, &reflectField{Name: "", Outter: rv, Inner: getInnerElem(rv)})
	for i, keyElem := range keyElems {
		outter, inner, err := findField(rv, keyElem, tagName)
		if err != nil {
			return nil, errors.ErrorAtf(err, "unable to fetch field: %s", strings.Join(keyElems[:i+1], ""))
		}
		fields = append(fields, &reflectField{Name: keyElem, Outter: outter, Inner: inner})
		rv = inner
	}

	return fields, nil
}

func findField(rv reflect.Value, name, tagName string) (reflect.Value, reflect.Value, error) {
	k := rv.Kind()
	switch k {
	case reflect.Struct:
		return findStructField(rv, name, tagName)
	case reflect.Map:
		return findMapField(rv, name)
	case reflect.Array, reflect.Slice:
		return findArrayField(rv, name)
	default:
		return reflect.Value{}, reflect.Value{}, errors.Errorf("parent's kind not of supported to get `%s`: %v", k, name)
	}
}

func findArrayField(rv reflect.Value, name string) (reflect.Value, reflect.Value, error) {
	idx := -1
	if n, _ := fmt.Sscanf(name, "[%d]", &idx); n != 1 {
		return reflect.Value{}, reflect.Value{}, errors.Errorf("the key is notarray notation: %s", name)
	}
	if idx < 0 || rv.Len() <= idx {
		return reflect.Value{}, reflect.Value{}, errors.Errorf("access array out of range(len=%d): %d", rv.Len(), idx)
	}
	v := rv.Index(idx)
	return v, getInnerElem(v), nil
}

func findMapField(rv reflect.Value, name string) (reflect.Value, reflect.Value, error) {
	iter := rv.MapRange()
	for iter.Next() {
		key := iter.Key()
		if key.String() == name {
			v := iter.Value()
			return v, getInnerElem(v), nil
		}
	}

	return reflect.Value{}, reflect.Value{}, errors.Errorf("map key not found: %s", name)
}

func findStructField(rv reflect.Value, name, tagName string) (reflect.Value, reflect.Value, error) {
	t := rv.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		var fieldName string

		tag := field.Tag.Get(tagName)
		if tagName != "" && tag != "" {
			fieldName = strings.Split(tag, ",")[0]
		} else {
			fieldName = field.Name
		}

		if fieldName == name {
			v := rv.Field(i)
			return v, getInnerElem(v), nil
		}
	}
	return reflect.Value{}, reflect.Value{}, errors.Errorf("struct field not found: %s", name)
}

func setFieldValue(rv reflect.Value, v any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Error(r)
		}
	}()

	rv.Set(reflect.ValueOf(v))
	return
}

func splitKey(key string) []string {
	dotSplitKeyElems := strings.Split(key, ".")
	keyElems := make([]string, 0, len(dotSplitKeyElems)+3)
	for _, dotkey := range dotSplitKeyElems {
		name, idx, ok := parseArrayElem(dotkey)
		if ok {
			if name != "" {
				keyElems = append(keyElems, name)
			}
			keyElems = append(keyElems, fmt.Sprintf("[%d]", idx))
		} else {
			keyElems = append(keyElems, dotkey)
		}
	}

	return keyElems
}

var pattArrayNotation = regexp.MustCompile(`^(\w+)?\[(-?\d+)\]$`)

func parseArrayElem(keyElem string) (string, int, bool) {
	match := pattArrayNotation.FindStringSubmatch(keyElem)
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

func getInnerElem(rv reflect.Value) reflect.Value {
	k := rv.Kind()
	if k == reflect.Ptr || k == reflect.Interface {
		return getInnerElem(rv.Elem())
	}

	return rv
}
