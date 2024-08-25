package usercontroller

import (
	userDTO "github.com/drunkleen/rasta/internal/DTO/user"
	commonerrors "github.com/drunkleen/rasta/internal/common/errors"
	"github.com/drunkleen/rasta/internal/common/utils"
	userservice "github.com/drunkleen/rasta/internal/service/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type ResetPwdController struct {
	ResetPwdService *userservice.ResetPwdService
	UserService     *userservice.UserService
}

// NewResetPwdController returns a new instance of ResetPwdController.
//
// It takes two parameters: resetPwdService and userService, both pointers to services used for password reset and user management respectively.
// Returns a pointer to a ResetPwdController.
func NewResetPwdController(resetPwdService *userservice.ResetPwdService, userService *userservice.UserService) *ResetPwdController {
	return &ResetPwdController{ResetPwdService: resetPwdService, UserService: userService}
}

// VerifyAndResetPassword godoc
// @Summary Verify OTP and Reset Password
// @Description Verifies the provided OTP and, if valid, allows the user to reset their password.
// @Tags Password Reset
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param ResetPassword body userDTO.ResetPassword true "Password reset request body"
// @Success 200 {object} userDTO.GenericResponse "Password reset successfully"
// @Failure 400 {object} commonerrors.ErrorMap "Bad Request"
// @Failure 401 {object} commonerrors.ErrorMap "Unauthorized"
// @Failure 406 {object} commonerrors.ErrorMap "Not Acceptable"
// @Failure 500 {object} commonerrors.ErrorMap "Internal Server Error"
// @Router /users/reset-password/{id}/verify [post]
func (c *ResetPwdController) VerifyAndResetPassword(ctx *gin.Context) {
	var ResetPassword userDTO.ResetPassword
	if err := ctx.ShouldBindJSON(&ResetPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(commonerrors.ErrInvalidRequestBody))
		return
	}
	if err := ResetPassword.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(commonerrors.ErrPasswordsNotMatch))
	}
	if !utils.PasswordValid(ResetPassword.NewPassword1) {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(commonerrors.ErrPasswordTooWeak))
		return
	}
	userId := uuid.MustParse(ctx.Param("id"))
	user, err := c.ResetPwdService.FindByUserIdIncludingResetPwd(&userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(commonerrors.ErrInvalidUserId))
		return
	}
	if !utils.CompareHashWithString(ResetPassword.Otp, user.ResetPwd.Code) || time.Now().After(user.ResetPwd.Expiry) {
		ctx.JSON(http.StatusUnauthorized, commonerrors.NewErrorMap("invalid or expired otp"))
		return
	}
	err = c.UserService.ResetPassword(userId, ResetPassword.NewPassword1)
	if err != nil {
		ctx.JSON(http.StatusNotAcceptable, commonerrors.NewErrorMap(err.Error()))
		return
	}
	err = c.ResetPwdService.Delete(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, commonerrors.NewErrorMap(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, userDTO.GenericResponse{
		Status: "success",
		Data: struct {
			Message string `json:"message"`
		}{
			Message: "password reset successfully",
		},
	})
}

// Send godoc
// @Summary Send Password Reset Code
// @Description Generates a password reset code and sends it to the user's email.
// @Tags Password Reset
// @Accept  json
// @Produce  json
// @Param email body map[string]string true "User email"
// @Success 200 {object} userDTO.GenericResponse "Password reset code sent successfully"
// @Failure 400 {object} commonerrors.ErrorMap "Bad Request"
// @Failure 500 {object} commonerrors.ErrorMap "Internal Server Error"
// @Router /users/reset-password [get]
func (c *ResetPwdController) Send(ctx *gin.Context) {
	var reqBody map[string]string
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(commonerrors.ErrInvalidRequestBody))
		return
	}
	userEmail := reqBody["email"]
	if !utils.EmailValidate(&userEmail) {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(commonerrors.ErrInvalidRequestBody))
		return
	}
	user, err := c.ResetPwdService.FindByUserEmailIncludingResetPwd(&userEmail)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(commonerrors.ErrInvalidRequestBody))
		return
	}
	if err = c.ResetPwdService.GenerateResetPwdAndSendEmail(user, user.Id); err != nil {
		ctx.JSON(http.StatusInternalServerError, commonerrors.NewErrorMap("Failed to generate password reset code"))
		return
	}
	ctx.JSON(http.StatusOK, userDTO.GenericResponse{
		Status: "success",
		Data: struct {
			Message string    `json:"message"`
			Id      uuid.UUID `json:"id"`
		}{
			Message: "password reset code sent successfully to your email",
			Id:      user.Id,
		},
	})
}
