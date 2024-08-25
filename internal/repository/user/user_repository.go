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

// NewUserRepository returns a new instance of UserRepository.
//
// Parameters:
// - db: the database connection to be used by the UserRepository.
//
// Returns:
// - *UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// GetAll returns all users from the database.
//
// No parameters are required.
// Returns a slice of usermodel.User and an error.
func (r *UserRepository) GetAll() ([]usermodel.User, error) {
	var users []usermodel.User
	err := r.DB.Find(&users).Error
	return users, err
}

// GetLimited returns a limited number of users from the database, based on the given offset and limit parameters.
//
// The offset parameter specifies the number of records to skip before starting to return records.
// The limit parameter specifies the maximum number of records to return.
//
// If there are no users found, an error is returned with the message "no users found".
func (r *UserRepository) GetLimited(offset, limit int) (*[]usermodel.User, error) {
	var users []usermodel.User
	err := r.DB.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		log.Printf("no users found: %v", err)
		return nil, errors.New("no users found")
	}
	return &users, nil
}

// CountUsers returns the total count of users in the database.
//
// It returns an error if the query fails.
func (r *UserRepository) CountUsers() (int64, error) {
	var count int64
	err := r.DB.Model(&usermodel.User{}).Count(&count).Error
	if err != nil {
		log.Printf("failed to count users: %v", err)
		return 0, errors.New("failed to count users")
	}
	return count, nil
}

// FindById finds a user by their id.
//
// Parameters:
// - id: the id of the user to find.
//
// Returns:
// - usermodel.User
// - error
func (r *UserRepository) FindById(id uuid.UUID) (usermodel.User, error) {
	var dbUser usermodel.User
	err := r.DB.Preload("OAuth").Where("id = ?", id).First(&dbUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("user not found: %v", err)
		return dbUser, errors.New("user not found")
	}
	return dbUser, nil
}

// FindByUsername finds a user by their username.
//
// Parameters:
// - username: the username to search for.
//
// Returns:
// - usermodel.User
func (r *UserRepository) FindByUsername(username string) (usermodel.User, error) {
	var dbUser usermodel.User
	err := r.DB.Preload("OAuth").Where("username = ?", username).First(&dbUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("user not found: %v", err)
		return dbUser, errors.New("user not found")
	}
	return dbUser, nil
}

// FindByEmail finds a user by their email.
//
// Parameters:
// - email: the email to search for.
//
// Returns:
// - usermodel.User
func (r *UserRepository) FindByEmail(email string) (usermodel.User, error) {
	var dbUser usermodel.User
	err := r.DB.Preload("OAuth").Where("email = ?", email).First(&dbUser).Error
	if err != nil {
		log.Printf("failed to find user: %v", err)
		return dbUser, errors.New("failed to find user")
	}
	return dbUser, nil
}

// FindByUsernameOrEmail finds a user by their username or email.
//
// Parameters:
// - username: the username to search for.
// - email: the email to search for.
//
// Returns:
// - usermodel.User
// FindByUsernameOrEmail finds a user by their username or email.
//
// Parameters:
// - username: the username to search for.
// - email: the email to search for.
//
// Returns:
// - usermodel.User
func (r *UserRepository) FindByUsernameOrEmail(username, email string) (usermodel.User, error) {
	var dbUser usermodel.User
	err := r.DB.Where("username = ?", username).Or("email = ?", email).First(&dbUser).Error
	if err != nil {
		log.Printf("failed to find user: %v", err)
		return dbUser, errors.New("failed to find user")
	}
	return dbUser, nil
}

// Create creates a new user in the UserRepository.
//
// Parameters:
// - user: the user to be created.
//
// Returns:
// - error: if the creation operation fails, an error is returned.
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

// Update updates a user in the UserRepository.
//
// Parameters:
// - user: the user to be updated.
//
// Returns:
// - error: if the update operation fails, an error is returned.
func (r *UserRepository) Update(user *usermodel.User) error {
	user.UpdatedAt = time.Now()
	err := r.DB.Save(user).Error
	if err != nil {
		log.Printf("failed to update user: %v", err)
		return errors.New("failed to update user")
	}
	return nil
}

// Delete permanently removes a user from the UserRepository.
//
// Parameters:
// - id: the unique identifier of the user to be deleted.
//
// Returns:
// - error: if the deletion operation fails, an error is returned.
func (r *UserRepository) Delete(id uuid.UUID) error {
	err := r.DB.Delete(&usermodel.User{}, "id = ?", id).Error
	if err != nil {
		log.Printf("failed to delete user: %v", err)
		return errors.New("failed to delete user")
	}
	return nil
}

// UpdateEmail updates the email of a user in the UserRepository.
//
// Parameters:
// - id: the unique identifier of the user.
// - email: the new email to be updated.
//
// Returns:
// - error: if the update operation fails, an error is returned.
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

// UpdatePassword updates the password of a user in the UserRepository.
//
// Parameters:
// - id: the unique identifier of the user.
// - password: the new password to be updated.
//
// Returns:
// - error: if the update operation fails, an error is returned.
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

// UpdateUsername updates the username of a user in the UserRepository.
//
// Parameters:
// - id: the unique identifier of the user.
// - username: the new username to be updated.
//
// Returns:
// - error: if the update operation fails, an error is returned.
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

// UpdateRegion updates the region field of a user in the database.
//
// id is the unique identifier of the user to update.
// region is the new value of the region field.
// Returns an error if the update operation fails.
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

// UpdateIsVerified updates the is_verified field of a user in the database.
//
// id is the unique identifier of the user to update.
// isVerified is the new value of the is_verified field.
// Returns an error if the update operation fails.
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

// UpdateIsDisabled updates the is_disabled field of the user with the given id.
// If a error occurred during the update, it will return the error.
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
