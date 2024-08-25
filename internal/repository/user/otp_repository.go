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

// NewOtpRepository returns a new instance of OtpRepository.
//
// Parameter:
// - db: a pointer to the gorm DB instance.
//
// Returns:
// - *OtpRepository
func NewOtpRepository(db *gorm.DB) *OtpRepository {
	return &OtpRepository{DB: db}
}

// Create creates a new OTP entry for the given user ID.
//
// It takes three parameters: userId, otpCode, and expTime.
// The userId is the unique identifier of the user,
// the otpCode is the one-time password to be stored,
// and the expTime is the time when the OTP expires.
//
// It returns an error if the operation fails.
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

// Delete deletes an OTP record by ID.
//
// Parameters:
// - id: the UUID of the OTP record to be deleted.
// Returns:
// - error: an error if the deletion fails.
func (r *OtpRepository) Delete(id uuid.UUID) error {
	err := r.DB.Where("user_id = ?", id).Delete(&usermodel.OtpEmail{}).Error
	if err != nil {
		log.Printf("failed to delete otp: %v", err)
		return err
	}
	return nil
}

// FindByUserId finds an OTP record by user ID.
//
// Parameters:
// - id: the UUID of the user.
//
// Returns:
// - *usermodel.OtpEmail: the OTP record associated with the user, or nil if not found.
// - error: an error if the query fails.
func (r *OtpRepository) FindByUserId(id uuid.UUID) (*usermodel.OtpEmail, error) {
	var otpEmail usermodel.OtpEmail
	err := r.DB.Where("user_id = ?", id).First(&otpEmail).Error
	return &otpEmail, err
}

// FindByUserIdIncludingOtp finds a user by ID and includes the associated OTP record.
//
// Parameters:
// - id: a pointer to the UUID of the user.
//
// Returns:
// - *usermodel.User
func (r *OtpRepository) FindByUserIdIncludingOtp(id *uuid.UUID) (*usermodel.User, error) {
	var user usermodel.User
	err := r.DB.Preload("OtpEmail").Where("id = ?", *id).First(&user).Error
	return &user, err
}

// FindByUserEmailIncludingOtp finds a user by email and includes the OTP record associated with the user.
//
// email is the email address of the user to find.
// Returns the user and OTP record if found, or an error if not found.
func (r *OtpRepository) FindByUserEmailIncludingOtp(email *string) (*usermodel.User, error) {
	var user usermodel.User
	err := r.DB.Preload("OtpEmail").Where("email = ?", *email).First(&user).Error
	return &user, err
}

// DeleteByUserId deletes an OTP record by user ID.
//
// id is the user ID of the OTP record to be deleted.
// Returns an error if the deletion fails.
func (r *OtpRepository) DeleteByUserId(id uuid.UUID) error {
	err := r.DB.Delete(&usermodel.OtpEmail{}, "user_id = ?", id).Error
	if err != nil {
		log.Printf("failed to delete otp: %v", err)
		return err
	}
	return nil
}
