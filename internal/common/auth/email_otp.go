package auth

import "math/rand"

// GenerateOtpCode generates a One Time Password (OTP) with a given length.
//
// The generated OTP is a string composed of characters from the charset
// "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ". The generated OTP is of the given length.
//
// The function uses the math/rand package to generate the OTP, and each call
// to GenerateOtpCode will generate a new OTP.
func GenerateOtpCode(length int) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
