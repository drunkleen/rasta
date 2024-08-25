package utils

import (
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func HashString(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CompareHashWithString(plainString, hashString string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashString), []byte(plainString))
	return err == nil
}

func PasswordValid(password string) bool {

	if len(password) < 8 {
		return false
	}
	if !strings.ContainsAny(password, "!@#$%^&*()_+`-=[]{}|;':\",./<>?") {
		return false
	}

	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") && !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return false
	}
	if strings.Contains(password, " ") || strings.Contains(password, "\t") {
		return false
	}
	return true
}

// EmailValidate checks if an email address is valid.
//
// Parameters:
// - email: the email address to be checked.
//
// Returns:
// - bool: true if the email address is valid, false otherwise.
func EmailValidate(email *string) bool {
	if *email == "" {
		return false
	}
	if strings.Contains(*email, "+") {
		return false
	}
	if !strings.Contains(*email, "@") {
		return false
	}
	if !strings.Contains(strings.Split(*email, "@")[1], ".") {
		return false
	}
	parts := strings.SplitN(*email, "@", 2)
	localPart := strings.ToLower(parts[0])
	lowerEmail := localPart + "@" + parts[1]
	*email = lowerEmail
	return true
}

func UsernameValid(username string) bool {
	if len(username) < 4 {
		return false
	}
	if strings.ContainsAny(username, "!@#$%^&*()+`-=[]{}|;':\",./<>?") {
		return false
	}
	if !strings.ContainsAny(username, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") && !strings.ContainsAny(username, "abcdefghijklmnopqrstuvwxyz") {
		return false
	}
	if strings.Contains(username, " ") || strings.Contains(username, "\t") {
		return false
	}
	return true
}
