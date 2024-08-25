package usercontroller

import (
	"fmt"
	userDTO "github.com/drunkleen/rasta/internal/DTO/user"
	"github.com/drunkleen/rasta/internal/common/auth"
	commonerrors "github.com/drunkleen/rasta/internal/common/errors"
	"github.com/drunkleen/rasta/internal/service/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
)

type UserController struct {
	UserService  *userservice.UserService
	OAuthService *userservice.OAuthService
	OtpService   *userservice.OtpService
}

// NewUserController creates a new instance of the UserController.
//
// userService is the UserService instance to be used by the UserController.
// otpService is the OtpService instance to be used by the UserController.
// Returns a pointer to the newly created UserController instance.
func NewUserController(userService *userservice.UserService, otpService *userservice.OtpService) *UserController {
	return &UserController{UserService: userService, OtpService: otpService}
}

// GetWithPagination godoc
// @Summary Get users with pagination
// @Description Get a list of users with pagination support
// @Tags Users
// @Accept  json
// @Produce  json
// @Param limit query int false "Number of users per page" default(10)
// @Param page query int false "Page number" default(1)
// @Success 200 {object} userDTO.GenericResponse
// @Failure 500 {object} userDTO.GenericResponse
// @Router /admin/users [get]
func (c *UserController) GetWithPagination(ctx *gin.Context) {
	limitStr := ctx.Query("limit")
	pageStr := ctx.Query("page")
	limit := 10
	page := 1

	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	users, err := c.UserService.GetUsersWithPagination(limit, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			userDTO.GenericResponse{
				Status: "error",
				Error:  err.Error(),
			},
		)
		return
	}
	ctx.JSON(http.StatusOK, userDTO.GenericResponse{
		Status: "success",
		Data:   users,
	})
}

// GetAllUsersCount godoc
// @Summary Get the total number of users
// @Description Retrieve the total count of users in the system
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {object} userDTO.GenericResponse
// @Failure 500 {object} userDTO.GenericResponse
// @Router /admin/users/count [get]
func (c *UserController) GetAllUsersCount(ctx *gin.Context) {
	count, err := c.UserService.GetAllUsersCount()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			commonerrors.NewErrorMap(err.Error()),
		)
		return
	}
	ctx.JSON(http.StatusOK, userDTO.GenericResponse{
		Status: "success",
		Data: struct {
			UserCount int64 `json:"user_count"`
		}{
			UserCount: count,
		},
	})
}

// FindUserByID godoc
// @Summary Get user by ID
// @Description Retrieve user details by their ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} userDTO.GenericResponse
// @Failure 404 {object} userDTO.GenericResponse
// @Router /admin/users/id/{id} [get]
func (c *UserController) FindUserByID(ctx *gin.Context) {
	userId := uuid.MustParse(ctx.Param("id"))
	user, err := c.UserService.FindById(userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound,
			commonerrors.NewErrorMap(err.Error()),
		)
		return
	}
	ctx.JSON(http.StatusOK, userDTO.GenericResponse{
		Status: "success",
		Data:   user,
	})
}

// FindUserByUsername godoc
// @Summary Get user by username
// @Description Retrieve user details by their username
// @Tags Users
// @Accept  json
// @Produce  json
// @Param username path string true "Username"
// @Success 200 {object} userDTO.GenericResponse
// @Failure 404 {object} userDTO.GenericResponse
// @Router /users/{username} [get]
func (c *UserController) FindUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")
	user, err := c.UserService.FindByUsername(username)
	if err != nil {
		ctx.JSON(http.StatusNotFound,
			commonerrors.NewErrorMap(err.Error()),
		)
		return
	}
	ctx.JSON(http.StatusOK, userDTO.GenericResponse{
		Status: "success",
		Data:   user,
	})
}

// Create godoc
// @Summary Create a new user
// @Description Create a new user account and send a verification OTP email
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body userDTO.UserCreate true "User creation payload"
// @Success 200 {object} userDTO.GenericResponse
// @Failure 400 {object} userDTO.GenericResponse
// @Failure 500 {object} userDTO.GenericResponse
// @Router /users/signup [post]
func (c *UserController) Create(ctx *gin.Context) {
	var user userDTO.UserCreate
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest,
			commonerrors.NewErrorMap(commonerrors.ErrInvalidRequestBody),
		)
		return
	}
	newUser, err := c.UserService.Create(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(err.Error()))
		return
	}
	err = c.OtpService.GenerateOtpAndSendEmail(newUser, newUser.Id)
	if err != nil {
		return
	}
	jwtToken, err := auth.GenerateJWTToken(newUser.Email, fmt.Sprintf("%v", newUser.Id))
	if err != nil {
		log.Printf("failed to generate JWT token: %v", err)
		ctx.JSON(http.StatusInternalServerError,
			commonerrors.NewErrorMap(commonerrors.ErrInternalServer),
		)
		return
	}
	ctx.JSON(http.StatusOK, userDTO.GenericResponse{
		Status: "success",
		Data:   userDTO.FromModelToUserLoginResponse(newUser, jwtToken),
	})
}

// Delete godoc
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} userDTO.GenericResponse
// @Failure 500 {object} userDTO.GenericResponse
// @Router /admin/users/id/{id} [delete]
func (c *UserController) Delete(ctx *gin.Context) {
	userId := uuid.MustParse(ctx.Param("id"))
	err := c.UserService.Delete(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			commonerrors.NewErrorMap(err.Error()),
		)
		return
	}
	ctx.JSON(http.StatusOK, userDTO.GenericResponse{
		Status: "success",
		Data: struct {
			Message string `json:"message"`
		}{
			Message: "user deleted successfully",
		},
	})
}

// Login godoc
// @Summary User login
// @Description Authenticates a user and returns a JWT token
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body userDTO.UserLogin true "User login payload"
// @Success 202 {object} userDTO.LoginResponse
// @Failure 401 {object} userDTO.GenericResponse
// @Failure 500 {object} userDTO.GenericResponse
// @Router /users/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var user userDTO.UserLogin
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusUnauthorized,
			commonerrors.NewErrorMap(commonerrors.ErrInvalidRequestBody),
		)
		return
	}
	dbUser, err := c.UserService.Login(user.Username, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized,
			commonerrors.NewErrorMap(err.Error()),
		)
		return
	}
	if dbUser.IsVerified == false {
		ctx.JSON(http.StatusUnauthorized,
			commonerrors.NewErrorMap(commonerrors.ErrUserNotVerified),
		)
		return
	}
	if dbUser.OAuth.Enabled {
		if err = c.OAuthService.OAuthValidate(&dbUser, user.OTP); err != nil {
			ctx.JSON(http.StatusUnauthorized,
				commonerrors.NewErrorMap(err.Error()),
			)
			return
		}
	}
	jwtToken, err := auth.GenerateJWTToken(dbUser.Email, fmt.Sprintf("%v", dbUser.Id))
	if err != nil {
		log.Printf("failed to generate JWT token: %v", err)
		ctx.JSON(http.StatusInternalServerError,
			commonerrors.NewErrorMap(commonerrors.ErrInternalServer),
		)
		return
	}
	ctx.JSON(http.StatusAccepted, userDTO.FromModelToUserLoginResponse(&dbUser, jwtToken))
}

// UpdatePassword godoc
// @Summary Update user password
// @Description Updates the password for the currently authenticated user
// @Tags Users
// @Accept  json
// @Produce  json
// @Param updatePassword body userDTO.UpdatePassword true "Password update payload"
// @Success 200 {object} userDTO.GenericResponse
// @Failure 400 {object} userDTO.GenericResponse
// @Failure 500 {object} userDTO.GenericResponse
// @Router /users/{username}/update-password [put]
func (c *UserController) UpdatePassword(ctx *gin.Context) {
	var updatePassword userDTO.UpdatePassword
	if err := ctx.ShouldBindJSON(&updatePassword); err != nil {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(commonerrors.ErrInvalidRequestBody))
	}
	if err := updatePassword.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(commonerrors.ErrPasswordsNotMatch))
	}

	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, commonerrors.NewErrorMap(commonerrors.ErrInternalServer))
		return
	}
	userIdStr, ok := userId.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, commonerrors.NewErrorMap(commonerrors.ErrInternalServer))
		return
	}
	id, err := uuid.Parse(userIdStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, commonerrors.NewErrorMap(commonerrors.ErrInternalServer))
		return
	}
	err = c.UserService.UpdatePassword(id, updatePassword.NewPassword1)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, commonerrors.NewErrorMap(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, userDTO.GenericResponse{
		Status: "success",
		Data: struct {
			Message string `json:"message"`
		}{
			Message: "Password updated successfully",
		},
	})

}
