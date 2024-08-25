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

func NewOAuthRepository(db *gorm.DB) *OAuthRepository {
	return &OAuthRepository{DB: db}
}

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

func (r *OAuthRepository) UpdateOAuthEnabled(id uuid.UUID, oauthEnabled bool) error {
	err := r.DB.Model(&usermodel.OAuth{}).Where("id = ?", id).Update("enabled", oauthEnabled).Error
	if err != nil {
		log.Printf("failed to update oauth_enabled: %v", err)
		return err
	}
	return nil
}

func (r *OAuthRepository) DeleteOAuth(id uuid.UUID) error {
	err := r.DB.Where("user_id = ?", id).Delete(&usermodel.OAuth{}).Error
	if err != nil {
		log.Printf("failed to delete oauth: %v", err)
		return err
	}
	return nil
}

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
