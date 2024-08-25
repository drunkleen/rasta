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

func NewOAuthService(repository *userrepository.OAuthRepository) *OAuthService {
	return &OAuthService{Repository: repository}
}

func (s *OAuthService) OAuthValidate(user *usermodel.User, oauth string) error {
	if auth.ValidateOTP(oauth, user.OAuth.Secret) {
		return nil
	}
	return errors.New(commonerrors.ErrInvalidOAuth)
}

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

func (s *OAuthService) UpdateOAuthEnabled(id uuid.UUID, oauthEnabled bool) error {
	return s.Repository.UpdateOAuthEnabled(id, oauthEnabled)
}

func (s *OAuthService) DeleteOAuth(id uuid.UUID) error {
	return s.Repository.DeleteOAuth(id)
}

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
