package usercontroller

import (
	userDTO "github.com/drunkleen/rasta/internal/DTO/user"
	commonerrors "github.com/drunkleen/rasta/internal/common/errors"
	"github.com/drunkleen/rasta/internal/common/utils"
	"github.com/drunkleen/rasta/internal/service/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type OtpController struct {
	OtpService  *userservice.OtpService
	UserService *userservice.UserService
}

// NewOtpController returns a new instance of the OtpController struct.
//
// Parameters:
// - otpService: a pointer to the userservice.OtpService object.
// - userService: a pointer to the userservice.UserService object.
//
// Returns a pointer to the OtpController struct.
func NewOtpController(otpService *userservice.OtpService, userService *userservice.UserService) *OtpController {
	return &OtpController{OtpService: otpService, UserService: userService}
}

// VerifyEmail godoc
// @Summary Verify Email with OTP
// @Description Verifies the user's email using the provided OTP. If successful, marks the email as verified and deletes the OTP.
// @Tags OTP
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param otp body map[string]string true "OTP code"
// @Success 200 {object} userDTO.GenericResponse "Email verified successfully"
// @Failure 400 {object} commonerrors.ErrorMap "Bad Request"
// @Failure 401 {object} commonerrors.ErrorMap "Unauthorized"
// @Failure 500 {object} commonerrors.ErrorMap "Internal Server Error"
// @Router /users/otp/{id}/verify [post]
func (c *OtpController) VerifyEmail(ctx *gin.Context) {
	var reqBody map[string]string
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(commonerrors.ErrInvalidRequestBody))
		return
	}
	otp, otpExists := reqBody["otp"]
	if !otpExists || otp == "" || len(otp) != 8 {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(commonerrors.ErrInvalidRequestBody))
		return
	}
	userId := uuid.MustParse(ctx.Param("id"))
	user, err := c.OtpService.FindByUserIdIncludingOtp(&userId)
	if err != nil || user.IsVerified {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(commonerrors.ErrInvalidUserId))
		return
	}
	if !utils.CompareHashWithString(otp, user.OtpEmail.Code) || time.Now().After(user.OtpEmail.Expiry) {
		ctx.JSON(http.StatusUnauthorized, commonerrors.NewErrorMap("invalid or expired otp"))
		return
	}
	err = c.UserService.MarkEmailAsVerified(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, commonerrors.NewErrorMap(err.Error()))
		return
	}
	err = c.OtpService.Delete(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, commonerrors.NewErrorMap(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, userDTO.GenericResponse{
		Status: "success",
		Data: struct {
			Message string `json:"message"`
		}{
			Message: "Email verified successfully",
		},
	})

}

// ResendOtp godoc
// @Summary Resend OTP to Email
// @Description Resends the OTP to the user's email for verification purposes.
// @Tags OTP
// @Accept  json
// @Produce  json
// @Param email body map[string]string true "User email"
// @Success 200 {object} userDTO.GenericResponse "OTP sent successfully"
// @Failure 400 {object} commonerrors.ErrorMap "Bad Request"
// @Failure 500 {object} commonerrors.ErrorMap "Internal Server Error"
// @Router /users/otp/resend [post]
func (c *OtpController) ResendOtp(ctx *gin.Context) {
	var reqBody map[string]string
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(commonerrors.ErrInvalidRequestBody))
		return
	}
	email := reqBody["email"]
	if !utils.EmailValidate(&email) {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(commonerrors.ErrInvalidRequestBody))
		return
	}
	user, err := c.OtpService.FindByUserEmailIncludingOtp(&email)
	if err != nil || user.IsVerified {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(commonerrors.ErrInvalidUserId))
		return
	}
	if user.IsVerified {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap("user already verified"))
		return
	}
	if err = c.OtpService.GenerateOtpAndSendEmail(user, user.Id); err != nil {
		ctx.JSON(http.StatusInternalServerError, commonerrors.NewErrorMap("failed to generate otp"))
		return
	}
	ctx.JSON(http.StatusOK, userDTO.GenericResponse{
		Status: "success",
		Data: struct {
			Message string    `json:"message"`
			Id      uuid.UUID `json:"id"`
		}{
			Message: "otp sent successfully to your email",
			Id:      user.Id,
		},
	})
}
