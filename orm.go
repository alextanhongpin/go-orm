package orm

import (
	"reflect"

	_ "github.com/lib/pq"
)

func getFields[T any](t T) []any {
	v := reflect.Indirect(reflect.ValueOf(t))

	fields := make([]any, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		fields[i] = f.Addr().Interface()
	}

	return fields
}

func getType(v any) string {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func getColumns(t any, tag string) []string {
	v := reflect.Indirect(reflect.ValueOf(t)).Type()

	fields := make([]string, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		name := f.Tag.Get(tag)
		if name == "" {
			name = f.Name
		}
		fields[i] = name
	}

	return fields
}
