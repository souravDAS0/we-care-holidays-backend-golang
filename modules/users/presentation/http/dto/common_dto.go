package dto

import (
	"net/mail"
	"regexp"
	"strings"
)

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidPhone(phone string) bool {
	// Basic phone validation - adjust regex as needed
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	return phoneRegex.MatchString(strings.ReplaceAll(phone, " ", ""))
}
