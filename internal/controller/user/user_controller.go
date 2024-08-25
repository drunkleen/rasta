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

func NewUserController(userService *userservice.UserService, otpService *userservice.OtpService) *UserController {
	return &UserController{UserService: userService, OtpService: otpService}
}

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
