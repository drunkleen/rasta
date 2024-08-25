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

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetAll() ([]usermodel.User, error) {
	var users []usermodel.User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *UserRepository) GetLimited(offset, limit int) (*[]usermodel.User, error) {
	var users []usermodel.User
	err := r.DB.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		log.Printf("no users found: %v", err)
		return nil, errors.New("no users found")
	}
	return &users, nil
}

func (r *UserRepository) CountUsers() (int64, error) {
	var count int64
	err := r.DB.Model(&usermodel.User{}).Count(&count).Error
	if err != nil {
		log.Printf("failed to count users: %v", err)
		return 0, errors.New("failed to count users")
	}
	return count, nil
}

func (r *UserRepository) FindById(id uuid.UUID) (usermodel.User, error) {
	var dbUser usermodel.User
	err := r.DB.Preload("OAuth").Where("id = ?", id).First(&dbUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("user not found: %v", err)
		return dbUser, errors.New("user not found")
	}
	return dbUser, nil
}

func (r *UserRepository) FindByUsername(username string) (usermodel.User, error) {
	var dbUser usermodel.User
	err := r.DB.Preload("OAuth").Where("username = ?", username).First(&dbUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("user not found: %v", err)
		return dbUser, errors.New("user not found")
	}
	return dbUser, nil
}

func (r *UserRepository) FindByEmail(email string) (usermodel.User, error) {
	var dbUser usermodel.User
	err := r.DB.Preload("OAuth").Where("email = ?", email).First(&dbUser).Error
	if err != nil {
		log.Printf("failed to find user: %v", err)
		return dbUser, errors.New("failed to find user")
	}
	return dbUser, nil
}

func (r *UserRepository) FindByUsernameOrEmail(username, email string) (usermodel.User, error) {
	var dbUser usermodel.User
	err := r.DB.Where("username = ?", username).Or("email = ?", email).First(&dbUser).Error
	if err != nil {
		log.Printf("failed to find user: %v", err)
		return dbUser, errors.New("failed to find user")
	}
	return dbUser, nil
}

func (r *UserRepository) Create(user *usermodel.User) error {
	user.Id = uuid.New()
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	var err error
	user.Password, err = utils.HashString(user.Password)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		return errors.New("failed to hash password")
	}
	if err := r.DB.Create(user).Error; err != nil {
		log.Printf("failed to create user: %v", err)
		return errors.New("failed to create user")
	}
	return nil
}

func (r *UserRepository) Update(user *usermodel.User) error {
	user.UpdatedAt = time.Now()
	err := r.DB.Save(user).Error
	if err != nil {
		log.Printf("failed to update user: %v", err)
		return errors.New("failed to update user")
	}
	return nil
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	err := r.DB.Delete(&usermodel.User{}, "id = ?", id).Error
	if err != nil {
		log.Printf("failed to delete user: %v", err)
		return errors.New("failed to delete user")
	}
	return nil
}

func (r *UserRepository) UpdateEmail(id uuid.UUID, email string) error {
	updates := map[string]interface{}{
		"email":      email,
		"updated_at": time.Now(),
	}
	if err := r.DB.Model(&usermodel.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Printf("failed to update email: %v", err)
		return errors.New("failed to update email")
	}
	return nil
}

func (r *UserRepository) UpdatePassword(id uuid.UUID, password string) error {
	var err error
	password, err = utils.HashString(password)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		return errors.New("failed to hash password")
	}
	updates := map[string]interface{}{
		"password":   password,
		"updated_at": time.Now(),
	}
	err = r.DB.Model(&usermodel.User{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		log.Printf("failed to update password: %v", err)
		return errors.New("failed to update password")
	}
	return nil
}

func (r *UserRepository) UpdateUsername(id uuid.UUID, username string) error {
	updates := map[string]interface{}{
		"username":   username,
		"updated_at": time.Now(),
	}
	if err := r.DB.Model(&usermodel.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Printf("failed to update username: %v", err)
		return errors.New("failed to update username")
	}
	return nil
}

func (r *UserRepository) UpdateRegion(id uuid.UUID, region string) error {
	updates := map[string]interface{}{
		"region":     region,
		"updated_at": time.Now(),
	}
	if err := r.DB.Model(&usermodel.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Printf("failed to update region: %v", err)
		return errors.New("failed to update region")
	}
	return nil
}

func (r *UserRepository) UpdateIsVerified(id uuid.UUID, isVerified bool) error {
	updates := map[string]interface{}{
		"is_verified": isVerified,
		"updated_at":  time.Now(),
	}
	if err := r.DB.Model(&usermodel.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Printf("failed to update is_verified: %v", err)
		return errors.New("failed to update is_verified")
	}
	return nil
}

func (r *UserRepository) UpdateIsDisabled(id uuid.UUID, isDisabled bool) error {
	updates := map[string]interface{}{
		"is_disabled": isDisabled,
		"updated_at":  time.Now(),
	}
	if err := r.DB.Model(&usermodel.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Printf("failed to update is_disabled: %v", err)
		return errors.New("failed to update is_disabled")
	}
	return nil
}
