package models

import (
	"strings"

	"github.com/serenize/snaker"
)

// This is lifted from
// https://github.com/markbates/validate/blob/master/validators/time_is_before_time.go

// // TimeIsBeforeOrEqualTime is used in the validator below
// type TimeIsBeforeOrEqualTime struct {
// 	FirstName  string
// 	FirstTime  time.Time
// 	SecondName string
// 	SecondTime time.Time
// }
//
// // IsValid makes sure the first time is before or equal to the second time
// func (v *TimeIsBeforeOrEqualTime) IsValid(errors *validate.Errors) {
// 	if v.FirstTime.UnixNano() > v.SecondTime.UnixNano() {
// 		errors.Add(GenerateKey(v.FirstName), fmt.Sprintf("%s must be before %s.", v.FirstName, v.SecondName))
// 	}
// }

// CustomKeys is used to generate keys
var CustomKeys = map[string]string{}

// GenerateKey creates a validator error key
func GenerateKey(s string) string {
	key := CustomKeys[s]
	if key != "" {
		return key
	}
	key = strings.Replace(s, " ", "", -1)
	key = strings.Replace(key, "-", "", -1)
	key = snaker.CamelToSnake(key)
	return key
}
