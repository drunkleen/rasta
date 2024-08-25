package userDTO

import (
	"errors"
	oauthDTO "github.com/drunkleen/rasta/internal/DTO/oauth"
	"github.com/drunkleen/rasta/internal/models/user"
	"time"

	"github.com/google/uuid"
)

type GenericResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

type UserCreate struct {
	FirstName string               `json:"first_name" binding:"required"`
	LastName  string               `json:"last_name" binding:"required"`
	Username  string               `json:"username" binding:"required"`
	Email     string               `json:"email" binding:"required"`
	Password  string               `json:"password" binding:"required"`
	Region    usermodel.RegionType `json:"region" binding:"required"`
}

// UserCreateToModel converts a UserCreate DTO to a usermodel.User, ready to be
// inserted into the database.
func (u *UserCreate) UserCreateToModel() *usermodel.User {
	return &usermodel.User{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Region:    u.Region,
	}
}

type UserCreateResponse struct {
	Id         uuid.UUID             `json:"id"`
	FirstName  string                `json:"first_name"`
	LastName   string                `json:"last_name"`
	Username   string                `json:"username"`
	Email      string                `json:"email"`
	Account    usermodel.AccountType `json:"account"`
	Region     usermodel.RegionType  `json:"country"`
	IsVerified bool                  `json:"is_verified"`
	CreatedAt  time.Time             `json:"created_at"`
}

// UserCreateResponse returns a new UserCreateResponse instance.
//
// It takes no parameters.
// Returns a pointer to a UserCreateResponse struct.
func (u *UserCreate) UserCreateResponse() *UserCreateResponse {
	return &UserCreateResponse{
		Id:         uuid.New(),
		Username:   u.Username,
		Email:      u.Email,
		Region:     u.Region,
		IsVerified: false,
		CreatedAt:  time.Now(),
	}
}

// UserCreateResponseToModel converts a UserCreateResponse DTO to a usermodel.User, ready to be
// inserted into the database.
//
// It takes no parameters.
// Returns a pointer to a usermodel.User struct.
func (u *UserCreate) UserCreateResponseToModel() *usermodel.User {
	return &usermodel.User{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Region:    u.Region,
	}
}

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	OTP      string `json:"otp" binding:"-"`
}

type User struct {
	Id         uuid.UUID             `json:"id,omitempty"`
	FirstName  string                `json:"first_name,omitempty"`
	LastName   string                `json:"last_name,omitempty"`
	Username   string                `json:"username,omitempty"`
	Email      string                `json:"email,omitempty"`
	Account    usermodel.AccountType `json:"account,omitempty"`
	Region     usermodel.RegionType  `json:"country,omitempty"`
	IsVerified bool                  `json:"is_verified,omitempty"`
	IsDisabled bool                  `json:"is_disabled,omitempty"`
	CreatedAt  time.Time             `json:"created_at,omitempty"`
	UpdatedAt  time.Time             `json:"updated_at,omitempty"`
	OAuth      oauthDTO.Response     `json:"oauth,omitempty"`
}

// FromModelToUserResponse converts a usermodel.User to a User DTO.
//
// It takes a pointer to a usermodel.User struct as a parameter.
// Returns a pointer to a User struct.
func FromModelToUserResponse(user *usermodel.User) *User {
	return &User{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Account:   user.Account,
		Region:    user.Region,
		CreatedAt: user.CreatedAt,
	}
}

// FromModelToUserResponseForAdmins converts a usermodel.User to a User DTO for admins.
//
// It takes a pointer to a usermodel.User struct as a parameter.
// Returns a pointer to a User struct.
func FromModelToUserResponseForAdmins(user *usermodel.User) *User {
	return &User{
		Id:         user.Id,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Username:   user.Username,
		Email:      user.Email,
		Account:    user.Account,
		Region:     user.Region,
		IsVerified: user.IsVerified,
		IsDisabled: user.IsDisabled,
		CreatedAt:  user.CreatedAt,
		OAuth: oauthDTO.Response{
			Enabled: user.OAuth.Enabled,
		},
	}
}

type LoginResponse struct {
	Status string `json:"status"`
	User   *User  `json:"user"`
	Token  string `json:"token"`
}

// FromModelToUserLoginResponse converts a usermodel.User to a LoginResponse DTO.
//
// It takes a pointer to a usermodel.User struct and a token string as parameters.
// Returns a pointer to a LoginResponse struct.
func FromModelToUserLoginResponse(user *usermodel.User, token string) *LoginResponse {
	return &LoginResponse{
		Status: "success",
		User:   FromModelToUserResponse(user),
		Token:  token,
	}
}

type UpdatePassword struct {
	OldPassword  string `json:"old_password" binding:"required"`
	NewPassword1 string `json:"new_password1" binding:"required"`
	NewPassword2 string `json:"new_password2" binding:"required"`
}

// Validate checks if the new password and its confirmation match.
//
// It takes no parameters, but uses the fields of the UpdatePassword struct.
// Returns an error if the passwords do not match, nil otherwise.
func (u *UpdatePassword) Validate() error {
	if u.NewPassword1 != u.NewPassword2 {
		return errors.New("passwords do not match")
	}
	return nil
}

type ResetPassword struct {
	Otp          string `json:"otp" binding:"required"`
	NewPassword1 string `json:"new_password1" binding:"required"`
	NewPassword2 string `json:"new_password2" binding:"required"`
}

// Validate checks if the new password and its confirmation match.
//
// It takes no parameters, but uses the fields of the ResetPassword struct.
// Returns an error if the passwords do not match, nil otherwise.
func (u *ResetPassword) Validate() error {
	if u.NewPassword1 != u.NewPassword2 {
		return errors.New("passwords do not match")
	}
	return nil
}
