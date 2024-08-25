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

func NewOtpController(otpService *userservice.OtpService, userService *userservice.UserService) *OtpController {
	return &OtpController{OtpService: otpService, UserService: userService}
}

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
