package userservice

import (
	"errors"
	"github.com/drunkleen/rasta/config"
	"github.com/drunkleen/rasta/internal/common/auth"
	commonerrors "github.com/drunkleen/rasta/internal/common/errors"
	usermodel "github.com/drunkleen/rasta/internal/models/user"
	userrepository "github.com/drunkleen/rasta/internal/repository/user"
	emailPkg "github.com/drunkleen/rasta/pkg/email"
	"github.com/google/uuid"
	"log"
	"time"
)

type ResetPwdService struct {
	Repository *userrepository.ResetPwdRepository
}

// NewResetPwd creates a new ResetPwdService.
//
// Parameters:
//   - repository: The ResetPwdRepository to use.
//
// Returns:
//   - *ResetPwdService: The created ResetPwdService.
func NewResetPwd(repository *userrepository.ResetPwdRepository) *ResetPwdService {
	return &ResetPwdService{Repository: repository}
}

// GenerateResetPwdAndSendEmail generates a reset password OTP and sends it to the user via email.
//
// It takes a user model and a user ID as parameters and returns an error.
//
// The OTP is generated using the GenerateOtpCode function in the auth package.
// The OTP is valid for the amount of time specified in the env variable EMAIL_OTP_EXPIRY.
// The generated OTP is stored in the repository.
// The user model is then updated with the generated OTP.
// The email is sent using the SendEmailResetPassword function in the email package.
// If the email cannot be sent, the generated OTP is deleted from the repository and an error is returned.
func (s *ResetPwdService) GenerateResetPwdAndSendEmail(userModel *usermodel.User, userId uuid.UUID) error {
	otpCode := auth.GenerateOtpCode(8)
	expTime := time.Now().Add(time.Duration(config.GetEnvEmailOTPExpiry()) * time.Second)
	// FindByUserId retrieves a ResetPwd model by its User ID.
	//
	// Parameters:
	// - id: the UUID of the User.
	//
	// Returns:
	// - *usermodel.ResetPwd: the ResetPwd model if found, or nil if not found.
	// - error: an error if the ResetPwd model cannot be retrieved.

	err := s.Repository.Create(userId, otpCode, expTime)
	if err != nil {
		return err
	}

	userModel.ResetPwd.Code = otpCode

	err = emailPkg.SendEmailResetPassword(userModel)
	if err != nil {
		log.Printf("Error sending email reset password model: %v", err)
		_ = s.Repository.Delete(userId)
		return errors.New(commonerrors.ErrInternalServer)
	}
	return nil
}

// FindByUserId retrieves a ResetPwd model by its User ID.
//
// Parameters:
// - id: the UUID of the User.
//
// Returns:
// - *usermodel.ResetPwd: the ResetPwd model if found, or nil if not found.
// - error: an error if the ResetPwd model cannot be retrieved.
// FindByUserId retrieves a ResetPwd model by its User ID.
//
// Parameters:
// - id: the UUID of the User.
//
// Returns:
// - *usermodel.ResetPwd: the ResetPwd model if found, or nil if not found.
// - error: an error if the ResetPwd model cannot be retrieved.
func (s *ResetPwdService) FindByUserId(id uuid.UUID) (*usermodel.ResetPwd, error) {
	topData, err := s.Repository.FindByUserId(id)
	if err != nil {
		log.Printf("Error finding reset password model: %v", err)
		return &usermodel.ResetPwd{}, errors.New(commonerrors.ErrInvalidUserId)
	}
	return topData, nil
}

// FindByUserEmailIncludingResetPwd retrieves a user model including reset password information by email.
//
// email - The email address of the user to be retrieved.
// Returns a user model and an error.
func (s *ResetPwdService) FindByUserEmailIncludingResetPwd(email *string) (*usermodel.User, error) {
	user, err := s.Repository.FindByUserEmailIncludingResetPwd(email)
	if err != nil {
		log.Printf("Error finding reset password model: %v", err)
		return &usermodel.User{}, errors.New(commonerrors.ErrInvalidUserId)
	}
	return user, nil
}

// FindByUserIdIncludingResetPwd retrieves a user model including reset password information by user ID.
//
// id - The unique identifier of the user.
// Returns a user model and an error.
func (s *ResetPwdService) FindByUserIdIncludingResetPwd(id *uuid.UUID) (*usermodel.User, error) {
	user, err := s.Repository.FindByUserIdIncludingResetPwd(id)
	if err != nil {
		log.Printf("Error finding reset password model: %v", err)
		return &usermodel.User{}, errors.New(commonerrors.ErrInvalidUserId)
	}
	return user, nil
}

// Delete deletes a reset password model from the repository.
//
// It takes an id of type uuid.UUID as a parameter and returns an error.
func (s *ResetPwdService) Delete(id uuid.UUID) error {
	err := s.Repository.Delete(id)
	if err != nil {
		log.Printf("Error deleting reset password model: %v", err)
		return errors.New(commonerrors.ErrInvalidUserId)
	}
	return nil
}
