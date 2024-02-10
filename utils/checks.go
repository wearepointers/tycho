package utils

func IsSlice[T any](v any) (bool, []T) {
	val, ok := v.([]T)
	return ok, val
}

func IsBool(v interface{}) (bool, bool) {
	val, ok := v.(bool)
	return ok, val
}
