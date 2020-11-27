package cs3users

import (
	"reflect"
)

var (
	validKinds = []reflect.Kind{
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
	}
)

// verifies an autoincrement field kind on the target struct.
func isValidKind(k reflect.Kind) bool {
	for _, v := range validKinds {
		if k == v {
			return true
		}
	}
	return false
}

func getKind(i interface{}, field string) (reflect.Kind, error) {
	r := reflect.ValueOf(i)
	return reflect.Indirect(r).FieldByName(field).Kind(), nil
}
