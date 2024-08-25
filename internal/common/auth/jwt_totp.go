package auth

import (
	"fmt"
	"github.com/drunkleen/rasta/config"
	"github.com/pquerna/otp/totp"
	"log"
)

// CreateOAuth generates a TOTP token for the given email.
//
// It takes an email as a parameter and returns the generated TOTP token as a string and an error.
func CreateOAuth(email string) (string, error) {
	// func Create(email string) {
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      config.GetJwtIssuer(),
		AccountName: email,
	})
	if err != nil {
		log.Println(err)
	}
	return secret.Secret(), nil
}

// GenerateOAuthUrl generates a TOTP URL for the given email and secret code.
//
// The function takes an email and a secret code as parameters.
// It returns a string representing the TOTP URL.
func GenerateOAuthUrl(email string, SecretCode string) string {
	issuer := config.GetJwtIssuer()
	return fmt.Sprintf(
		"otpauth://totp/%s:%s?secret=%s&issuer=%s",
		issuer, email, SecretCode, issuer,
	)
}

// ValidateOTP validates a given TOTP passcode against a secret.
//
// Parameter UserTOTPPassCode is the TOTP passcode to be validated, and secret is the secret code used for validation.
// Return type is a boolean indicating whether the TOTP passcode is valid.
func ValidateOTP(UserTOTPPassCode, secret string) bool {
	isValid := totp.Validate(UserTOTPPassCode, secret)
	if !isValid {
		log.Println("Invalid TOTP code")
		return false
	}
	return true
}
