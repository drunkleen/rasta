package userrepository

import (
	"errors"
	"github.com/drunkleen/rasta/internal/common/utils"
	usermodel "github.com/drunkleen/rasta/internal/models/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"time"
)

type ResetPwdRepository struct {
	DB *gorm.DB
}

func NewResetPwdRepository(db *gorm.DB) *ResetPwdRepository {
	return &ResetPwdRepository{DB: db}
}

func (r *ResetPwdRepository) Create(userId uuid.UUID, otpCode string, expTime time.Time) error {
	if userId == uuid.Nil {
		return errors.New("user ID is required")
	}
	if otpCode == "" {
		return errors.New("otp code is required")
	}
	if err := r.DB.Where("user_id = ?", userId).First(&usermodel.ResetPwd{}).Error; err == nil {
		err = r.Delete(userId)
		if err != nil {
			return err
		}
	}

	hashedOtpCode, err := utils.HashString(otpCode)
	if err != nil {
		return err
	}

	resetPwdModel := &usermodel.ResetPwd{
		UserId: userId,
		Code:   hashedOtpCode,
		Expiry: expTime,
	}
	if err = r.DB.Create(resetPwdModel).Error; err != nil {
		return err
	}
	return nil
}

func (r *ResetPwdRepository) Delete(id uuid.UUID) error {
	err := r.DB.Where("user_id = ?", id).Delete(&usermodel.ResetPwd{}).Error
	if err != nil {
		log.Printf("failed to delete otp: %v", err)
		return err
	}
	return nil
}

func (r *ResetPwdRepository) FindByUserId(id uuid.UUID) (*usermodel.ResetPwd, error) {
	var forgetPasswordModel usermodel.ResetPwd
	err := r.DB.Where("user_id = ?", id).First(&forgetPasswordModel).Error
	return &forgetPasswordModel, err
}

func (r *ResetPwdRepository) FindByUserEmailIncludingResetPwd(email *string) (*usermodel.User, error) {
	var user usermodel.User
	err := r.DB.Preload("ResetPwd").Where("email = ?", *email).First(&user).Error
	return &user, err
}

func (r *ResetPwdRepository) FindByUserIdIncludingResetPwd(id *uuid.UUID) (*usermodel.User, error) {
	var user usermodel.User
	err := r.DB.Preload("ResetPwd").Where("id = ?", *id).First(&user).Error
	return &user, err
}

func (r *ResetPwdRepository) DeleteByUserId(id uuid.UUID) error {
	err := r.DB.Delete(&usermodel.ResetPwd{}, "user_id = ?", id).Error
	if err != nil {
		log.Printf("failed to delete otp: %v", err)
		return err
	}
	return nil
}
