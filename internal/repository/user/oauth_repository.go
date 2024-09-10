package userrepository

import (
	"errors"
	"github.com/drunkleen/rasta/internal/models/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
)

type OAuthRepository struct {
	DB *gorm.DB
}

// NewOAuthRepository creates a new instance of the OAuthRepository struct,
// initialized with the provided *gorm.DB instance.
//
// Parameters:
// - db: A pointer to a *gorm.DB instance.
//
// Returns:
// - A pointer to an OAuthRepository instance.
func NewOAuthRepository(db *gorm.DB) *OAuthRepository {
	return &OAuthRepository{DB: db}
}

// Create creates a new OAuth entry for the given user.
//
// Parameters:
//
//	user: the user for whom the OAuth entry is created
//	secret: the secret code for the OAuth entry
//
// Returns:
//
//	an error, if any
func (r *OAuthRepository) Create(user *usermodel.User, secret string) error {
	if user.Id == uuid.Nil {
		return errors.New("user ID is required")
	}
	if secret == "" {
		return errors.New("OAuth secret is required")
	}
	if err := r.DB.Where("user_id = ?", user.Id).First(&usermodel.OAuth{}).Error; err == nil {
		err = r.DeleteOAuth(user.Id)
		if err != nil {
			return err
		}
	}

	oauth := &usermodel.OAuth{
		UserId:  user.Id,
		Enabled: false,
		Secret:  secret,
	}
	if err := r.DB.Create(oauth).Error; err != nil {
		return err
	}
	return nil
}

// UpdateOAuthEnabled updates the OAuth enabled status for the given user ID.
//
// Parameters:
//
//	id: the UUID of the user
//	oauthEnabled: whether OAuth is enabled
//
// Returns:
//
//	an error, if any
func (r *OAuthRepository) UpdateOAuthEnabled(id uuid.UUID, oauthEnabled bool) error {
	err := r.DB.Model(&usermodel.OAuth{}).Where("user_id = ?", id).Update("enabled", oauthEnabled).Error
	if err != nil {
		log.Printf("failed to update oauth_enabled: %v", err)
		return err
	}
	return nil
}

// DeleteOAuth deletes the OAuth record associated with the given user ID.
//
// Parameters:
// - id: the UUID of the user
//
// Returns:
// - error: an error if the deletion fails, nil otherwise.
func (r *OAuthRepository) DeleteOAuth(id uuid.UUID) error {
	err := r.DB.Where("user_id = ?", id).Delete(&usermodel.OAuth{}).Error
	if err != nil {
		log.Printf("failed to delete oauth: %v", err)
		return err
	}
	return nil
}

// UpdateOAuthSecret updates the OAuth secret for the given user ID.
//
// Parameters:
//
//	id: the UUID of the user
//	oauthEnabled: whether OAuth is enabled
//	secret: the new secret code
//
// Returns:
//
//	an error, if any
func (r *OAuthRepository) UpdateOAuthSecret(id uuid.UUID, oauthEnabled bool, secret string) error {
	updates := map[string]interface{}{
		"enabled": oauthEnabled,
		"secret":  secret,
	}
	err := r.DB.Model(&usermodel.OAuth{}).Where("user_id = ?", id).Updates(updates).Error
	if err != nil {
		log.Printf("failed to update otp_enabled: %v", err)
		return err
	}
	return nil
}
