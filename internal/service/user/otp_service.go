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

func NewOtpService(repository *userrepository.OtpRepository) *OtpService {
	return &OtpService{Repository: repository}
}

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

func (s *OtpService) FindByUserId(id uuid.UUID) (*usermodel.OtpEmail, error) {
	topData, err := s.Repository.FindByUserId(id)
	if err != nil {
		log.Printf("Error finding otp: %v", err)
		return &usermodel.OtpEmail{}, errors.New(commonerrors.ErrInvalidUserId)
	}
	return topData, nil
}

func (s *OtpService) FindByUserIdIncludingOtp(id *uuid.UUID) (*usermodel.User, error) {
	user, err := s.Repository.FindByUserIdIncludingOtp(id)
	if err != nil {
		log.Printf("Error finding otp: %v", err)
		return &usermodel.User{}, errors.New(commonerrors.ErrInvalidUserId)
	}
	return user, nil
}

func (s *OtpService) FindByUserEmailIncludingOtp(email *string) (*usermodel.User, error) {
	user, err := s.Repository.FindByUserEmailIncludingOtp(email)
	if err != nil {
		log.Printf("Error finding otp: %v", err)
		return &usermodel.User{}, errors.New(commonerrors.ErrInvalidUserId)
	}
	return user, nil
}

func (s *OtpService) Delete(id uuid.UUID) error {
	err := s.Repository.Delete(id)
	if err != nil {
		log.Printf("Error deleting otp: %v", err)
		return errors.New(commonerrors.ErrInvalidUserId)
	}
	return nil
}
