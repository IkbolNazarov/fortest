package validator

import (
	"fmt"
	"regexp"
	"strings"
)

func ValidatePhoneNumber(phoneNumber string) error {

	re := regexp.MustCompile(`^992?[0-9]{9}$`)

	if len(strings.Trim(phoneNumber, " ")) == 0 {
		return fmt.Errorf("phonenumber is null")
	}

	if !re.MatchString(phoneNumber) {
		return fmt.Errorf("phone number wrong format")
	}

	return nil
}
