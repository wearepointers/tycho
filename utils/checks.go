package utils

func IsSlice[T any](v any) (bool, []T) {
	val, ok := v.([]T)
	return ok, val
}

func IsBool(v any) (bool, bool) {
	val, ok := v.(bool)
	return ok, val
}

func IsString(v any) (bool, string) {
	val, ok := v.(string)
	return ok, val
}

func IsFloat64(v any) (bool, float64) {
	val, ok := v.(float64)
	return ok, val
}

func IsFloat32(v any) (bool, float32) {
	val, ok := v.(float32)
	return ok, val
}

func IsInt(v any) (bool, int) {
	val, ok := v.(int)
	return ok, val
}

func IsInt64(v any) (bool, int64) {
	val, ok := v.(int64)
	return ok, val
}

func IsInt32(v any) (bool, int32) {
	val, ok := v.(int32)
	return ok, val
}

func IsInt16(v any) (bool, int16) {
	val, ok := v.(int16)
	return ok, val
}

func IsInt8(v any) (bool, int8) {
	val, ok := v.(int8)
	return ok, val
}
