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

// NewResetPwdRepository creates a new ResetPwdRepository given a pointer to a GORM
// DB connection.
//
// The new repository is returned as a pointer to a ResetPwdRepository.
func NewResetPwdRepository(db *gorm.DB) *ResetPwdRepository {
	return &ResetPwdRepository{DB: db}
}

// Create creates a new reset password entry in the database.
//
// userId is the ID of the user for which the reset password entry is to be created.
// otpCode is the one-time password code to be stored.
// expTime is the expiration time after which the reset password entry is no longer valid.
//
// It returns an error if the creation fails.
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

// Delete deletes a reset password entry by user ID.
//
// id is the user ID of the reset password entry to be deleted.
// Returns an error if the deletion fails.
func (r *ResetPwdRepository) Delete(id uuid.UUID) error {
	err := r.DB.Where("user_id = ?", id).Delete(&usermodel.ResetPwd{}).Error
	if err != nil {
		log.Printf("failed to delete otp: %v", err)
		return err
	}
	return nil
}

// FindByUserId retrieves a reset password by user ID.
//
// id is the ID of the reset password to be retrieved.
// Returns a pointer to the reset password model and an error if the retrieval fails.
func (r *ResetPwdRepository) FindByUserId(id uuid.UUID) (*usermodel.ResetPwd, error) {
	var forgetPasswordModel usermodel.ResetPwd
	err := r.DB.Where("user_id = ?", id).First(&forgetPasswordModel).Error
	return &forgetPasswordModel, err
}

// FindByUserEmailIncludingResetPwd retrieves a user by their email, including their reset password information.
//
// email is the email of the user to be retrieved.
// Returns a pointer to the user model and an error if the retrieval fails.
func (r *ResetPwdRepository) FindByUserEmailIncludingResetPwd(email *string) (*usermodel.User, error) {
	var user usermodel.User
	err := r.DB.Preload("ResetPwd").Where("email = ?", *email).First(&user).Error
	return &user, err
}

// FindByUserIdIncludingResetPwd retrieves a user by their ID, including their reset password information.
//
// id is the ID of the user to be retrieved.
// Returns a pointer to the user model and an error if the retrieval fails.
func (r *ResetPwdRepository) FindByUserIdIncludingResetPwd(id *uuid.UUID) (*usermodel.User, error) {
	var user usermodel.User
	err := r.DB.Preload("ResetPwd").Where("id = ?", *id).First(&user).Error
	return &user, err
}

// DeleteByUserId deletes a reset password entry by user ID.
//
// id is the user ID of the reset password entry to be deleted.
// Returns an error if the deletion fails.
func (r *ResetPwdRepository) DeleteByUserId(id uuid.UUID) error {
	err := r.DB.Delete(&usermodel.ResetPwd{}, "user_id = ?", id).Error
	if err != nil {
		log.Printf("failed to delete otp: %v", err)
		return err
	}
	return nil
}
