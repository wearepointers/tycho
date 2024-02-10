package utils

import (
	"errors"
	"fmt"
	"reflect"
	"time"
	"unsafe"

	"github.com/viant/xunsafe"
)

func StructToMap(obj interface{}) (map[string]any, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Struct && v.Kind() != reflect.Ptr {
		return nil, errors.New("input must be a struct or pointer to a struct")
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	numField := t.NumField()

	m := make(map[string]any, numField)
	for i := 0; i < numField; i++ {
		m[t.Field(i).Tag.Get("boil")] = v.Field(i).Interface()
	}

	return m, nil
}

type reflectT struct {
	Tp  reflect.Type
	Ptr unsafe.Pointer
}

func (o reflectT) getUnsafeField(name string) *xunsafe.Field {
	return xunsafe.FieldByName(o.Tp, name)
}
func (o reflectT) getUnsafeFieldByTag(name, tag string) *xunsafe.Field {
	if o.Tp.Kind() == reflect.Ptr {
		o.Tp = o.Tp.Elem()
	}

	numField := o.Tp.NumField()
	for i := 0; i < numField; i++ {
		field := o.Tp.Field(i)
		if field.Tag.Get(tag) == name {
			return xunsafe.NewField(field)
		}
	}
	return nil
}

func (o reflectT) GetUnsafeFieldAsString(name string) (string, error) {
	unsafeField := o.getUnsafeField(name)
	if unsafeField == nil {
		return "", fmt.Errorf("field %s not found", name)
	}
	switch v := unsafeField.Interface(o.Ptr).(type) {
	case string:
		return v, nil
	case time.Time:
		return v.Format(postgresTimestamptzTimeFormat), nil
	default:
		return "UNKNOWN_TYPE", nil
	}
}

func (o reflectT) GetUnsafeFieldByTagAsString(name, tag string) (string, error) {
	unsafeField := o.getUnsafeFieldByTag(name, tag)
	if unsafeField == nil {
		return "", fmt.Errorf("field %s not found", name)
	}

	switch v := unsafeField.Interface(o.Ptr).(type) {
	case string:
		return v, nil
	case time.Time:
		return v.Format(postgresTimestamptzTimeFormat), nil
	default:
		return "", fmt.Errorf("field %s type not found: %v", name, v)
	}
}

func NewReflectObject[T any](obj *T) reflectT {
	return reflectT{reflect.TypeOf(obj), unsafe.Pointer(obj)}
}
