package orm

import (
	"fmt"
	"reflect"
	"strings"
)

type DataMapper[T any] struct {
	tableName  string
	tableAlias string
	columns    []string

	// The bit index of the included columns.
	colbits int
}

func NewDataMapper[T any](name, alias, tag string) *DataMapper[T] {
	var t T
	columns, colbits := getDataMapperColumns(t, tag)

	return &DataMapper[T]{
		tableName:  name,
		tableAlias: alias,
		columns:    columns,
		colbits:    colbits,
	}
}

func (m *DataMapper[T]) SelectName() string {
	return fmt.Sprintf("%s %s", m.tableName, m.tableAlias)
}

func (m *DataMapper[T]) Table() string {
	return m.tableName
}

func (m *DataMapper[T]) Alias() string {
	return m.tableAlias
}

func (m *DataMapper[T]) ColumnAlias(col string) string {
	return fmt.Sprintf("%s.%s", m.tableAlias, col)
}

func (m *DataMapper[T]) Columns() string {
	result := make([]string, len(m.columns))
	for i, c := range m.columns {
		result[i] = m.ColumnAlias(c)
	}

	return strings.Join(result, ", ")
}

func (m *DataMapper[T]) Fields(t any) []any {
	v := reflect.Indirect(reflect.ValueOf(t))

	fields := make([]any, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		b := 1 << i
		if m.colbits&b == b {
			fields = append(fields, f.Addr().Interface())
		}
	}

	return fields
}

func getDataMapperColumns(t any, tag string) ([]string, int) {
	v := reflect.Indirect(reflect.ValueOf(t)).Type()
	bit := 0

	fields := make([]string, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		name := v.Field(i).Tag.Get(tag)
		if name == "" || name == "-" {
			continue
		}
		bit |= 1 << i
		fields = append(fields, name)
	}

	return fields, bit
}
