package fluent

import (
	"fmt"
	"reflect"
	"strings"
)

func ConvertToValue(p interface{}, tagName string) interface{} {
	rv := toValue(p)
	switch rv.Kind() {
	case reflect.Struct:
		return converFromStruct(rv.Interface(), tagName)
	case reflect.Map:
		return convertFromMap(rv, tagName)
	case reflect.Slice:
		return convertFromSlice(rv, tagName)
	case reflect.Chan:
		return nil
	case reflect.Invalid:
		return nil
	default:
		return rv.Interface()
	}
}

func convertFromMap(rv reflect.Value, tagName string) interface{} {
	result := make(map[string]interface{})
	for _, key := range rv.MapKeys() {
		kv := rv.MapIndex(key)
		result[fmt.Sprint(key.Interface())] = ConvertToValue(kv.Interface(), tagName)
	}
	return result
}

func convertFromSlice(rv reflect.Value, tagName string) interface{} {
	var result []interface{}
	for i, max := 0, rv.Len(); i < max; i++ {
		result = append(result, ConvertToValue(rv.Index(i).Interface(), tagName))
	}
	return result
}

func converFromStruct(p interface{}, tagName string) interface{} {
	result := make(map[string]interface{})
	t := toType(p)
	values := toValue(p)
	for i, max := 0, t.NumField(); i < max; i++ {
		f := t.Field(i)
		if f.PkgPath != "" && !f.Anonymous {
			continue
		}
		tag, opts := parseTag(f, tagName)
		if tag == "-" {
			continue
		}

		v := values.Field(i)
		if opts.Has("omitempty") && isZero(v) {
			continue
		}
		name := getNameFromTag(f, tagName)
		result[name] = ConvertToValue(v.Interface(), TagName)
	}
	return result
}

func toValue(p interface{}) reflect.Value {
	v := reflect.ValueOf(p)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func toType(p interface{}) reflect.Type {
	t := reflect.ValueOf(p).Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func isZero(v reflect.Value) bool {
	zero := reflect.Zero(v.Type()).Interface()
	value := v.Interface()
	return reflect.DeepEqual(value, zero)
}

func getNameFromTag(f reflect.StructField, tagName string) string {
	tag, _ := parseTag(f, tagName)
	if tag != "" {
		return tag
	}
	return f.Name
}

func getTagValues(f reflect.StructField, tag string) string {
	return f.Tag.Get(tag)
}

func parseTag(f reflect.StructField, tag string) (string, options) {
	return splitTags(getTagValues(f, tag))
}

func splitTags(tags string) (string, options) {
	res := strings.Split(tags, ",")
	return res[0], res[1:]
}

type options []string

func (t options) Has(tag string) bool {
	for _, opt := range t {
		if opt == tag {
			return true
		}
	}
	return false
}