package sms

import (
	"errors"
	"regexp"
)

var validPhone = regexp.MustCompile(`^\+63[0-9]{10}$`)

func ValidatePhone(phone string) error {
	if !validPhone.MatchString(phone) {
		return errors.New("phone number format must be +63XXXXXXXXXX")
	}
	return nil
}
