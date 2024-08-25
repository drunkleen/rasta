package userservice

import (
	"errors"
	userDTO "github.com/drunkleen/rasta/internal/DTO/user"
	commonerrors "github.com/drunkleen/rasta/internal/common/errors"
	"github.com/drunkleen/rasta/internal/common/utils"
	usermodel "github.com/drunkleen/rasta/internal/models/user"
	"github.com/drunkleen/rasta/internal/repository/user"
	"github.com/google/uuid"
	"log"
	"strings"
)

type UserService struct {
	Repository *userrepository.UserRepository
}

// NewUserService creates a new instance of the UserService struct.
//
// It takes a pointer to a UserRepository as a parameter and returns a pointer to a UserService.
func NewUserService(repository *userrepository.UserRepository) *UserService {
	return &UserService{Repository: repository}
}

// GetAllUsers fetches all users in the database.
//
// No parameters.
// Returns a pointer to a slice of userDTO.User and an error if any.
func (s *UserService) GetAllUsers() ([]userDTO.User, error) {
	dbUsers, err := s.Repository.GetAll()
	if err != nil {
		log.Printf("Error fetching all users: %v", err)
		return nil, errors.New(commonerrors.ErrInternalServer)
	}
	var respUsers = make([]userDTO.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		respUsers[i] = *userDTO.FromModelToUserResponse(&dbUser)
	}
	return respUsers, nil
}

// GetUsersWithPagination fetches a paginated list of users.
//
// The function takes two parameters: limit and page, both integers.
// The limit parameter specifies the number of users to be returned per page.
// The page parameter specifies the current page number.
// Returns a pointer to a slice of userDTO.User and an error.
// GetUsersWithPagination fetches a paginated list of users.
//
// The function takes two parameters: limit and page, both integers.
// The limit parameter specifies the number of users to be returned per page.
// The page parameter specifies the current page number.
// Returns a pointer to a slice of userDTO.User and an error.
func (s *UserService) GetUsersWithPagination(limit, page int) (*[]userDTO.User, error) {
	if limit == 0 {
		limit = 1
	}
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	dbUsers, err := s.Repository.GetLimited(offset, limit)
	if err != nil {
		log.Printf("Error fetching all users: %v", err)
		return nil, errors.New(commonerrors.ErrInternalServer)
	}
	if len(*dbUsers) == 0 {
		return &[]userDTO.User{}, nil
	}
	var respUsers = make([]userDTO.User, len(*dbUsers))
	for i, dbUser := range *dbUsers {
		respUsers[i] = *userDTO.FromModelToUserResponseForAdmins(&dbUser)
	}
	return &respUsers, nil
}

// GetAllUsersCount returns the total count of all users in the database.
//
// No parameters.
// Returns an int64 representing the total count of users and an error if any.
func (s *UserService) GetAllUsersCount() (int64, error) {
	dbUsersCount, err := s.Repository.CountUsers()
	if err != nil {
		log.Printf("Error fetching all users count: %v", err)
		return 0, errors.New(commonerrors.ErrInternalServer)
	}
	return dbUsersCount, nil
}

// FindById finds a user by their ID.
//
// The ID is used to search for a user in the database. If the user is found, it
// is returned. If the user is not found, an error is returned. The error is
// either an internal server error or a user not found error.
func (s *UserService) FindById(id uuid.UUID) (*usermodel.User, error) {
	dbUser, err := s.Repository.FindById(id)
	if err != nil {
		log.Printf("Error finding user by ID: %v", err)
		return &usermodel.User{}, errors.New(commonerrors.ErrInvalidUserId)
	}
	return &dbUser, nil
}

// FindByUsername finds a user by their username.
//
// The username is used to search for a user in the database. If the user is
// found, it is returned. If the user is not found, an error is returned. The
// error is either an internal server error or a user not found error.
func (s *UserService) FindByUsername(username string) (*userDTO.User, error) {
	dbUser, err := s.Repository.FindByUsername(username)
	if err != nil {
		log.Printf("Error finding user by username: %v", err)
		return &userDTO.User{}, errors.New(commonerrors.ErrUsernameNotExists)
	}
	return userDTO.FromModelToUserResponse(&dbUser), nil
}

// FindByEmail finds a user by their email.
//
// The email is used to search for a user in the database. If the user is found,
// it is returned. If the user is not found, an error is returned. The error is
// either an internal server error or a user not found error.
func (s *UserService) FindByEmail(email string) (usermodel.User, error) {
	dbUser, err := s.Repository.FindByEmail(email)
	if err != nil {
		log.Printf("Error finding user by email: %v", err)
		return usermodel.User{}, errors.New(commonerrors.ErrEmailNotExists)
	}
	return dbUser, nil
}

// FindByUsernameOrEmail finds a user by either their username or email.
//
// The username or email is used to search for a user in the database. If the
// user is found, it is returned. If the user is not found, an error is returned.
// The error is either an internal server error or a user not found error.
func (s *UserService) FindByUsernameOrEmail(username string) (usermodel.User, error) {
	dbUser, err := s.Repository.FindByUsernameOrEmail(username, username)
	if err != nil {
		log.Printf("Error finding user by email or username: %v", err)
		return usermodel.User{}, errors.New("email or username not exists")
	}
	return dbUser, nil
}

// Create creates a new user.
//
// The user is created with the provided userDto, and the password is checked
// for complexity. If the password is not complex enough, an error is returned.
// The user is also checked for uniqueness by their username and email. If
// either the username or email is already in use, an error is returned.
// Finally, the user is created in the database, and the created user is returned.
// If the creation fails, an error is returned.
func (s *UserService) Create(userDto *userDTO.UserCreate) (*usermodel.User, error) {
	userModel := userDto.UserCreateResponseToModel()

	if !utils.PasswordValid(userModel.Password) {
		return &usermodel.User{}, errors.New(commonerrors.ErrPasswordTooWeak)
	}
	userModel.Username = strings.ToLower(userModel.Username)
	if _, err := s.Repository.FindByUsername(userModel.Username); err == nil {
		return &usermodel.User{}, errors.New(commonerrors.ErrUsernameAlreadyExists)
	}
	if _, err := s.Repository.FindByEmail(userModel.Email); err == nil {
		return &usermodel.User{}, errors.New(commonerrors.ErrEmailAlreadyExists)
	}
	if !utils.EmailValidate(&userModel.Email) {
		return &usermodel.User{}, errors.New(commonerrors.ErrInvalidEmail)
	}
	if !utils.UsernameValid(userModel.Username) {
		return &usermodel.User{}, errors.New(commonerrors.ErrInvalidUsername)
	}

	err := s.Repository.Create(userModel)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return &usermodel.User{}, errors.New(commonerrors.ErrInternalServer)
	}

	return userModel, nil
}

// Login authenticates a user by their username or email and password.
//
// usernameOrEmail is the username or email of the user to authenticate.
// password is the password of the user to authenticate.
// Returns the authenticated user and an error if authentication fails.
func (s *UserService) Login(usernameOrEmail, password string) (usermodel.User, error) {
	data := strings.ToLower(usernameOrEmail)
	dbUser, err := s.Repository.FindByUsernameOrEmail(data, usernameOrEmail)
	if err != nil {
		log.Println("Error finding user: ", err)
		return usermodel.User{}, errors.New(commonerrors.ErrInvalidCredentials)
	}
	if !utils.CompareHashWithString(password, dbUser.Password) {
		return usermodel.User{}, errors.New(commonerrors.ErrInvalidCredentials)
	}
	return dbUser, nil
}

// Update updates a user.
//
// user is the user to update.
// Returns an error if the update operation fails.
func (s *UserService) Update(user *usermodel.User) error {
	return s.Repository.Update(user)
}

// Delete deletes a user by ID.
//
// id is the unique identifier of the user to delete.
// Returns an error if the deletion operation fails.
func (s *UserService) Delete(id uuid.UUID) error {
	err := s.Repository.Delete(id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return errors.New(commonerrors.ErrInvalidUserId)
	}
	return nil
}

// UpdateEmail updates the email address associated with a user.
//
// id is the unique identifier of the user, and email is the new email address to associate with the user.
// Returns an error if the update operation fails.
func (s *UserService) UpdateEmail(id uuid.UUID, email string) error {
	return s.Repository.UpdateEmail(id, email)
}

// UpdatePassword updates the password associated with a user.
//
// id is the unique identifier of the user, and newPassword is the new password to associate with the user.
// Returns an error if the update operation fails.
func (s *UserService) UpdatePassword(id uuid.UUID, newPassword string) error {
	if !utils.PasswordValid(newPassword) {
		return errors.New(commonerrors.ErrPasswordTooWeak)
	}
	userModel, err := s.Repository.FindById(id)
	if err != nil {
		log.Printf("Error finding user by ID: %v", err)
		return errors.New(commonerrors.ErrInvalidUserId)
	}
	if !utils.CompareHashWithString(newPassword, userModel.Password) {
		return errors.New(commonerrors.ErrInvalidCredentials)
	}
	return s.Repository.UpdatePassword(id, newPassword)
}

// ResetPassword resets the password associated with a user.
//
// id is the unique identifier of the user, and newPassword is the new password to associate with the user.
// Returns an error if the update operation fails.
func (s *UserService) ResetPassword(id uuid.UUID, newPassword string) error {
	if !utils.PasswordValid(newPassword) {
		return errors.New(commonerrors.ErrPasswordTooWeak)
	}
	return s.Repository.UpdatePassword(id, newPassword)
}

// UpdateUsername updates the username of a user.
//
// id is the unique identifier of the user, and username is the new username to associate with the user.
// Returns an error if the update operation fails.
func (s *UserService) UpdateUsername(id uuid.UUID, username string) error {
	if !utils.UsernameValid(username) {
		return errors.New(commonerrors.ErrInvalidUsername)
	}
	return s.Repository.UpdateUsername(id, username)
}

// UpdateRegion updates the region of a user.
//
// id is the unique identifier of the user, and region is the name of the region to associate with the user.
// Returns an error if the update operation fails.
func (s *UserService) UpdateRegion(id uuid.UUID, country string) error {
	return s.Repository.UpdateRegion(id, country)
}

// MarkEmailAsVerified marks the email address associated with the user as verified.
//
// id is the unique identifier of the user.
// Returns an error if the update operation fails.
func (s *UserService) MarkEmailAsVerified(id uuid.UUID) error {
	return s.Repository.UpdateIsVerified(id, true)
}

// UpdateIsDisabled updates the disabled status of a user.
//
// id is the unique identifier of the user, and isDisabled is a boolean indicating whether the user should be disabled.
// Returns an error if the update operation fails.
// UpdateIsDisabled updates the disabled status of a user.
//
// id is the unique identifier of the user, and isDisabled is a boolean indicating whether the user should be disabled.
// Returns an error if the update operation fails.
func (s *UserService) UpdateIsDisabled(id uuid.UUID, isDisabled bool) error {
	return s.Repository.UpdateIsDisabled(id, isDisabled)
}
