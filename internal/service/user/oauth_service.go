package userservice

import (
	"errors"
	"github.com/drunkleen/rasta/internal/common/auth"
	commonerrors "github.com/drunkleen/rasta/internal/common/errors"
	"github.com/drunkleen/rasta/internal/models/user"
	"github.com/drunkleen/rasta/internal/repository/user"
	"github.com/google/uuid"
	"log"
)

type OAuthService struct {
	Repository *userrepository.OAuthRepository
}

// NewOAuthService creates a new instance of the OAuthService.
//
// It takes a pointer to the OAuthRepository as a parameter to initialize the OAuthService.
// It returns a pointer to the OAuthService.
func NewOAuthService(repository *userrepository.OAuthRepository) *OAuthService {
	return &OAuthService{Repository: repository}
}

// OAuthValidate validates the OAuth secret for a given user.
//
// It takes a user and an OAuth secret as parameters to identify the user and validate their OAuth secret.
// It returns an error if the validation fails.
func (s *OAuthService) OAuthValidate(user *usermodel.User, oauth string) error {
	if auth.ValidateOTP(oauth, user.OAuth.Secret) {
		return nil
	}
	return errors.New(commonerrors.ErrInvalidOAuth)
}

// GenerateOAuthSecret generates an OAuth secret for a given user.
//
// It takes a user as a parameter to identify the user and generate a new OAuth secret.
// It returns the new OAuth secret, the OAuth URL, and an error.
func (s *OAuthService) GenerateOAuthSecret(user *usermodel.User) (string, string, error) {
	secret, err := auth.CreateOAuth(user.Email)
	if err != nil {
		return "", "", errors.New(commonerrors.ErrInternalServer)
	}
	if err = s.Repository.Create(user, secret); err != nil {
		return "", "", err
	}
	return secret, auth.GenerateOAuthUrl(user.Email, secret), nil
}

// UpdateOAuthEnabled updates the OAuth enabled status for a given user.
//
// It takes a user ID and an OAuth enabled status as parameters to identify the user and update their OAuth status.
// It returns an error if the update fails.
func (s *OAuthService) UpdateOAuthEnabled(id uuid.UUID, oauthEnabled bool) error {
	return s.Repository.UpdateOAuthEnabled(id, oauthEnabled)
}

// DeleteOAuth deletes the OAuth secret for a given user.
//
// It takes a user ID as a parameter to identify the user.
// It returns an error if the deletion fails.
// UpdateOAuthSecret updates the OAuth secret for a given user.
//
// It takes an email and a user ID as parameters to identify the user and generate a new OAuth secret.
// It returns the new OAuth secret as a string and an error.
func (s *OAuthService) DeleteOAuth(id uuid.UUID) error {
	return s.Repository.DeleteOAuth(id)
}

// UpdateOAuthSecret updates the OAuth secret for a given user.
//
// It takes an email and a user ID as parameters to identify the user and generate a new OAuth secret.
// It returns the new OAuth secret as a string and an error.
func (s *OAuthService) UpdateOAuthSecret(email string, id uuid.UUID) (string, error) {
	oauthSecret, err := auth.CreateOAuth(email)
	if err != nil {
		return "", errors.New(commonerrors.ErrInternalServer)
	}
	err = s.Repository.UpdateOAuthSecret(id, false, oauthSecret)
	if err != nil {
		log.Printf("failed to update otp_secret: %v", err)
		return "", errors.New(commonerrors.ErrInternalServer)
	}
	return oauthSecret, nil
}
