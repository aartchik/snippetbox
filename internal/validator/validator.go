package validator

import (
	"strings"
	"unicode/utf8"
)


type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	if len(v.FieldErrors) > 0 {
		return false
	}
	return true
}

func NotBlank(str string) bool {
	return strings.TrimSpace(str) != ""
}

func MaxChar(str string, n int) bool {
	return utf8.RuneCountInString(str) < n
}

func Accept_values(expires int, value... int) bool {
	for r:= range value {
		if expires == value[r] {
			return true
		}
	}
	return false
}



func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldMap(key, message)
	}
}

func (v *Validator) AddFieldMap(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
	v.FieldErrors[key] = message
    }
}
