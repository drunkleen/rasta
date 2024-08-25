package userrepository

import (
	"errors"
	"github.com/drunkleen/rasta/internal/common/utils"
	"github.com/drunkleen/rasta/internal/models/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"time"
)

type OtpRepository struct {
	DB *gorm.DB
}

func NewOtpRepository(db *gorm.DB) *OtpRepository {
	return &OtpRepository{DB: db}
}

func (r *OtpRepository) Create(userId uuid.UUID, otpCode string, expTime time.Time) error {
	if userId == uuid.Nil {
		return errors.New("user ID is required")
	}
	if otpCode == "" {
		return errors.New("otp code is required")
	}
	if err := r.DB.Where("user_id = ?", userId).First(&usermodel.OtpEmail{}).Error; err == nil {
		err = r.Delete(userId)
		if err != nil {
			return err
		}
	}
	hashedOtpCode, err := utils.HashString(otpCode)
	if err != nil {
		return err
	}

	topEmail := &usermodel.OtpEmail{
		UserId: userId,
		Code:   hashedOtpCode,
		Expiry: expTime,
	}
	if err = r.DB.Create(topEmail).Error; err != nil {
		return err
	}
	return nil
}

func (r *OtpRepository) Delete(id uuid.UUID) error {
	err := r.DB.Where("user_id = ?", id).Delete(&usermodel.OtpEmail{}).Error
	if err != nil {
		log.Printf("failed to delete otp: %v", err)
		return err
	}
	return nil
}

func (r *OtpRepository) FindByUserId(id uuid.UUID) (*usermodel.OtpEmail, error) {
	var otpEmail usermodel.OtpEmail
	err := r.DB.Where("user_id = ?", id).First(&otpEmail).Error
	return &otpEmail, err
}

func (r *OtpRepository) FindByUserIdIncludingOtp(id *uuid.UUID) (*usermodel.User, error) {
	var user usermodel.User
	err := r.DB.Preload("OtpEmail").Where("id = ?", *id).First(&user).Error
	return &user, err
}

func (r *OtpRepository) FindByUserEmailIncludingOtp(email *string) (*usermodel.User, error) {
	var user usermodel.User
	err := r.DB.Preload("OtpEmail").Where("email = ?", *email).First(&user).Error
	return &user, err
}

func (r *OtpRepository) DeleteByUserId(id uuid.UUID) error {
	err := r.DB.Delete(&usermodel.OtpEmail{}, "user_id = ?", id).Error
	if err != nil {
		log.Printf("failed to delete otp: %v", err)
		return err
	}
	return nil
}
