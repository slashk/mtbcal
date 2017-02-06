package actions

import (
	"reflect"
	"time"
)

// ConvertFormDate sets a new data html parser for Schema
func ConvertFormDate(value string) reflect.Value {
	s, _ := time.Parse("2006-01-_2", value)
	return reflect.ValueOf(s)
}
