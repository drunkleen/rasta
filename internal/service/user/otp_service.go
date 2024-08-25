package userservice

import (
	"errors"
	"github.com/drunkleen/rasta/config"
	"github.com/drunkleen/rasta/internal/common/auth"
	commonerrors "github.com/drunkleen/rasta/internal/common/errors"
	"github.com/drunkleen/rasta/internal/models/user"
	"github.com/drunkleen/rasta/internal/repository/user"
	emailPkg "github.com/drunkleen/rasta/pkg/email"
	"github.com/google/uuid"
	"log"
	"time"
)

type OtpService struct {
	Repository *userrepository.OtpRepository
}

// NewOtpService returns a new instance of the OtpService struct.
//
// Parameter repository is a pointer to the userrepository.OtpRepository object.
// Return type is a pointer to the OtpService struct.
func NewOtpService(repository *userrepository.OtpRepository) *OtpService {
	return &OtpService{Repository: repository}
}

// GenerateOtpAndSendEmail generates a new OTP code, saves it to the repository, and sends an email to the user with the OTP code.
//
// Parameter userModel is the usermodel.User object of the user to send the OTP to, and userId is the unique identifier of the user.
// Return type is an error object that is returned if any of the operations fail.
func (s *OtpService) GenerateOtpAndSendEmail(userModel *usermodel.User, userId uuid.UUID) error {
	otpCode := auth.GenerateOtpCode(8)
	expTime := time.Now().Add(time.Duration(config.GetEnvEmailOTPExpiry()) * time.Second)

	err := s.Repository.Create(userId, otpCode, expTime)
	if err != nil {
		return err
	}

	userModel.OtpEmail.Code = otpCode

	err = emailPkg.SendEmailVerify(userModel)
	if err != nil {
		log.Printf("Error sending email Otp: %v", err)
		_ = s.Repository.Delete(userId)
		return errors.New(commonerrors.ErrInternalServer)
	}
	return nil
}

// FindByUserId finds a user by its ID.
//
// id: The UUID of the user to find.
// Returns a pointer to the usermodel.OtpEmail struct and an error if any.
func (s *OtpService) FindByUserId(id uuid.UUID) (*usermodel.OtpEmail, error) {
	topData, err := s.Repository.FindByUserId(id)
	if err != nil {
		log.Printf("Error finding otp: %v", err)
		return &usermodel.OtpEmail{}, errors.New(commonerrors.ErrInvalidUserId)
	}
	return topData, nil
}

// FindByUserIdIncludingOtp finds a user by its ID including OTP.
//
// id: The UUID of the user to find.
// Returns a pointer to the usermodel.User struct and an error if any.
func (s *OtpService) FindByUserIdIncludingOtp(id *uuid.UUID) (*usermodel.User, error) {
	user, err := s.Repository.FindByUserIdIncludingOtp(id)
	if err != nil {
		log.Printf("Error finding otp: %v", err)
		return &usermodel.User{}, errors.New(commonerrors.ErrInvalidUserId)
	}
	return user, nil
}

// FindByUserEmailIncludingOtp finds a user by its email including OTP.
//
// email - the email of the user to find.
// *usermodel.User - the user found, or nil if none.
// error - an error if the user was not found.
func (s *OtpService) FindByUserEmailIncludingOtp(email *string) (*usermodel.User, error) {
	user, err := s.Repository.FindByUserEmailIncludingOtp(email)
	if err != nil {
		log.Printf("Error finding otp: %v", err)
		return &usermodel.User{}, errors.New(commonerrors.ErrInvalidUserId)
	}
	return user, nil
}

// Delete deletes an OTP entry by its ID.
//
// id - the UUID of the OTP entry to delete.
// error - an error if the deletion fails.
func (s *OtpService) Delete(id uuid.UUID) error {
	err := s.Repository.Delete(id)
	if err != nil {
		log.Printf("Error deleting otp: %v", err)
		return errors.New(commonerrors.ErrInvalidUserId)
	}
	return nil
}
