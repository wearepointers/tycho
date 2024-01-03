package utils

import "reflect"

func IsSlice(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice
}

func IsBool(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Bool
}
