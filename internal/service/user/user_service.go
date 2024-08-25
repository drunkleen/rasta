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

// NewUserService returns a new instance of UserService.
//
// Parameters:
// - repository: a pointer to the UserRepository instance.
//
// Returns:
// - *UserService
func NewUserService(repository *userrepository.UserRepository) *UserService {
	return &UserService{Repository: repository}
}

// GetAllUsers retrieves all users from the database.
//
// No parameters.
// Returns a slice of models.User and an error.
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

// GetUsersWithPagination retrieves a limited number of users from the database without pagination.
//
// Parameters:
// - limit: the maximum number of users to retrieve.
// - page: the page number for the limited results.
//
// Returns:
// - []userDTO.UserCreateResponseDTO: a slice of user response DTOs.
// - error: an error if the operation fails.
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

// GetAllUsersCount retrieves the total number of users in the database.
//
// It calls the CountUsers method of the UserRepository to get the count of all users.
// If there is an error, it logs the error and returns an error with the ErrInternalServer code.
// Otherwise, it returns the count of users and no error.
//
// Returns:
// - int64: the total number of users in the database.
// - error: an error if there was an issue retrieving the count of users.
func (s *UserService) GetAllUsersCount() (int64, error) {
	dbUsersCount, err := s.Repository.CountUsers()
	if err != nil {
		log.Printf("Error fetching all users count: %v", err)
		return 0, errors.New(commonerrors.ErrInternalServer)
	}
	return dbUsersCount, nil
}

// FindById retrieves a user by their ID from the database.
//
// id is the unique identifier of the user to be retrieved.
// Returns a models.User and an error.
func (s *UserService) FindById(id uuid.UUID) (*usermodel.User, error) {
	dbUser, err := s.Repository.FindById(id)
	if err != nil {
		log.Printf("Error finding user by ID: %v", err)
		return &usermodel.User{}, errors.New(commonerrors.ErrInvalidUserId)
	}
	return &dbUser, nil
}

// FindByUsername retrieves a user from the database by their username.
//
// username is the username of the user to be retrieved.
// Returns a models.User and an error.
func (s *UserService) FindByUsername(username string) (*userDTO.User, error) {
	dbUser, err := s.Repository.FindByUsername(username)
	if err != nil {
		log.Printf("Error finding user by username: %v", err)
		return &userDTO.User{}, errors.New(commonerrors.ErrUsernameNotExists)
	}
	return userDTO.FromModelToUserResponse(&dbUser), nil
}

// FindByEmail retrieves a user from the database by their email.
//
// Parameters:
// - email: the email of the user to be retrieved.
//
// Returns:
// - models.User
// - error
func (s *UserService) FindByEmail(email string) (usermodel.User, error) {
	dbUser, err := s.Repository.FindByEmail(email)
	if err != nil {
		log.Printf("Error finding user by email: %v", err)
		return usermodel.User{}, errors.New(commonerrors.ErrEmailNotExists)
	}
	return dbUser, nil
}

// FindByUsernameOrEmail retrieves a user from the database by their username or email.
//
// Parameters:
// - username: the username or email of the user to be retrieved.
//
// Returns:
// - models.User
// - error
func (s *UserService) FindByUsernameOrEmail(username string) (usermodel.User, error) {
	dbUser, err := s.Repository.FindByUsernameOrEmail(username, username)
	if err != nil {
		log.Printf("Error finding user by email or username: %v", err)
		return usermodel.User{}, errors.New("email or username not exists")
	}
	return dbUser, nil
}

// Create creates a new user in the database.
//
// Parameters:
// - user: the user to be created.
//
// Returns:
// - models.User
// - error
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

// Login logs a user in and returns the user if the credentials are valid.
//
// Parameters:
// - usernameOrEmail: the username or email of the user to be logged in.
// - password: the password of the user to be logged in.
//
// Returns:
// - models.User
// - error
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

// Update updates a user in the database.
//
// Parameters:
// - user: the user to be updated.
//
// Returns:
// - error
func (s *UserService) Update(user *usermodel.User) error {
	return s.Repository.Update(user)
}

// Delete deletes a user by their ID from the database.
//
// Parameters:
// - id: the ID of the user to be deleted.
// Returns:
// - error
func (s *UserService) Delete(id uuid.UUID) error {
	err := s.Repository.Delete(id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return errors.New(commonerrors.ErrInvalidUserId)
	}
	return nil
}

// UpdateEmail updates the email of a user in the database.
//
// Parameters:
// - id: the ID of the user to be updated.
// - email: the new email of the user.
// Returns:
// - error
func (s *UserService) UpdateEmail(id uuid.UUID, email string) error {
	return s.Repository.UpdateEmail(id, email)
}

// UpdatePassword updates the password of a user in the database.
//
// Parameters:
// - id: the ID of the user to be updated.
// - password: the new password of the user.
// Returns:
// - error
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

func (s *UserService) ResetPassword(id uuid.UUID, newPassword string) error {
	if !utils.PasswordValid(newPassword) {
		return errors.New(commonerrors.ErrPasswordTooWeak)
	}
	return s.Repository.UpdatePassword(id, newPassword)
}

// UpdateUsername updates the username of a user in the database.
//
// Parameters:
// - id: the ID of the user to be updated.
// - username: the new username of the user.
// Returns:
// - error
func (s *UserService) UpdateUsername(id uuid.UUID, username string) error {
	if !utils.UsernameValid(username) {
		return errors.New(commonerrors.ErrInvalidUsername)
	}
	return s.Repository.UpdateUsername(id, username)
}

// UpdateRegion updates the country of a user in the database.
//
// Parameters:
// - id: the ID of the user to be updated.
// - country: the new country of the user.
// Returns:
// - error
func (s *UserService) UpdateRegion(id uuid.UUID, country string) error {
	return s.Repository.UpdateRegion(id, country)
}

// MarkEmailAsVerified updates the verified status of a user in the database.
//
// Parameters:
// - id: the ID of the user to be updated.
// - isVerified: the new verified status of the user.
// Returns:
// - error
func (s *UserService) MarkEmailAsVerified(id uuid.UUID) error {
	return s.Repository.UpdateIsVerified(id, true)
}

// UpdateIsDisabled updates the disabled status of a user in the database.
//
// Parameters:
// - id: the ID of the user to be updated.
// - isDisabled: the new disabled status of the user.
// Returns:
// - error
func (s *UserService) UpdateIsDisabled(id uuid.UUID, isDisabled bool) error {
	return s.Repository.UpdateIsDisabled(id, isDisabled)
}
