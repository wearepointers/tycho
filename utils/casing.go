package utils

import (
	"github.com/danielgtaylor/casing"
)

func ToSnakeCase(s string) string {
	return casing.Snake(s)
}

func ToCamelCase(s string) string {
	return casing.LowerCamel(s)
}

func ToPascalCase(s string) string {
	return casing.Camel(s)
}
