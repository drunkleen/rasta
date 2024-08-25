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

func NewResetPwd(repository *userrepository.ResetPwdRepository) *ResetPwdService {
	return &ResetPwdService{Repository: repository}
}

func (s *ResetPwdService) GenerateResetPwdAndSendEmail(userModel *usermodel.User, userId uuid.UUID) error {
	otpCode := auth.GenerateOtpCode(8)
	expTime := time.Now().Add(time.Duration(config.GetEnvEmailOTPExpiry()) * time.Second)

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

func (s *ResetPwdService) FindByUserId(id uuid.UUID) (*usermodel.ResetPwd, error) {
	topData, err := s.Repository.FindByUserId(id)
	if err != nil {
		log.Printf("Error finding reset password model: %v", err)
		return &usermodel.ResetPwd{}, errors.New(commonerrors.ErrInvalidUserId)
	}
	return topData, nil
}

func (s *ResetPwdService) FindByUserEmailIncludingResetPwd(email *string) (*usermodel.User, error) {
	user, err := s.Repository.FindByUserEmailIncludingResetPwd(email)
	if err != nil {
		log.Printf("Error finding reset password model: %v", err)
		return &usermodel.User{}, errors.New(commonerrors.ErrInvalidUserId)
	}
	return user, nil
}

func (s *ResetPwdService) FindByUserIdIncludingResetPwd(id *uuid.UUID) (*usermodel.User, error) {
	user, err := s.Repository.FindByUserIdIncludingResetPwd(id)
	if err != nil {
		log.Printf("Error finding reset password model: %v", err)
		return &usermodel.User{}, errors.New(commonerrors.ErrInvalidUserId)
	}
	return user, nil
}

func (s *ResetPwdService) Delete(id uuid.UUID) error {
	err := s.Repository.Delete(id)
	if err != nil {
		log.Printf("Error deleting reset password model: %v", err)
		return errors.New(commonerrors.ErrInvalidUserId)
	}
	return nil
}
