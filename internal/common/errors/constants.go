package commonerrors

const (
	ErrUserNotFound          = "user not found"
	ErrUnauthorizedToken     = "unauthorized, invalid token"
	ErrUnauthorizedExpToken  = "unauthorized, expired token"
	ErrForbidden             = "forbidden"
	ErrUserNotVerified       = "user not verified"
	ErrInvalidCredentials    = "invalid credentials"
	ErrInvalidOAuth          = "invalid one-time password"
	ErrInvalidUserId         = "invalid user ID"
	ErrEmailAlreadyExists    = "email already exists"
	ErrEmailNotExists        = "email not exists"
	ErrInvalidEmail          = "invalid email address"
	ErrUsernameAlreadyExists = "username already exists"
	ErrUsernameNotExists     = "username not exists"
	ErrInvalidUsername       = "username must be at least 4 characters long and contain only letters and numbers"
	ErrInvalidRequestBody    = "invalid request body"
	ErrPasswordTooWeak       = "password too weak. must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character"
	ErrPasswordsNotMatch     = "password do not match"
	ErrInternalServer        = "internal server error"
)
