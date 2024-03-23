package cqlx_util

import (
	"reflect"

	"openmyth/messgener/util/database"
)

const (
	PART_KEY_TAG = "partkey"
	SORT_KEY_TAG = "sortkey"
)

// Fields takes a generic type T from database.DB and returns the partKeys and sortKeys.
//
// It takes a parameter e of type T and returns two slices of strings.
func Fields[T database.Entity](e T) (partKeys []string, sortKeys []string) {

	v := reflect.ValueOf(e).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		partKeys = append(partKeys, field.Tag.Get(PART_KEY_TAG))
		sortKeys = append(sortKeys, field.Tag.Get(SORT_KEY_TAG))
	}

	return partKeys, sortKeys
}
