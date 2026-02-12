package validator

import (
	"strings"
	"unicode/utf8"
	"regexp"
)


type Validator struct {
	NonFieldErrors []string
	FieldErrors map[string]string
}

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")


func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddNonFieldError(message string) {
    v.NonFieldErrors = append(v.NonFieldErrors, message)
}


func NotBlank(str string) bool {
	return strings.TrimSpace(str) != ""
}

func MaxChar(str string, n int) bool {
	return utf8.RuneCountInString(str) < n
}

func MinChar(str string, n int) bool {
	return utf8.RuneCountInString(str) > n
}

func Matches(email string, rx *regexp.Regexp) bool {
	return rx.MatchString(email)
}

func Accept_values(expires int, value... int) bool {
	for r:= range value {
		if expires == value[r] {
			return true
		}
	}
	return false
}


func SamePassword(pass1, pass2 string) (bool) {
	return pass1 == pass2
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
