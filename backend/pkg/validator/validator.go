package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func IsValidPassword(password string) (bool, error) {
	if len(password) < 8 {
		return false, errors.New("password must be at least 8 characters")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	switch false {
	case hasUpper:
		return false, errors.New("password must contain at least one uppercase letter")
	case hasLower:
		return false, errors.New("password must contain at least one lowercase letter")
	case hasSpecial:
		return false, errors.New("password must contain at least one special character")
	}

	return true, nil
}

func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func IsValidRole(role string) bool {
	validRoles := []string{"donor", "charity_staff", "admin"}
	role = strings.ToLower(role)
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

func IsValidStatus(status string) bool {
	validStatuses := []string{"pending", "approved", "rejected", "distributed"}
	status = strings.ToLower(status)
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func ValidateRequired(fields map[string]string) error {
	for fieldName, fieldValue := range fields {
		if IsEmpty(fieldValue) {
			return fmt.Errorf("%s is required", fieldName)
		}
	}
	return nil
}
