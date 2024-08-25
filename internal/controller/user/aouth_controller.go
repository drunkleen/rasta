package usercontroller

import (
	oauthDTO "github.com/drunkleen/rasta/internal/DTO/oauth"
	commonerrors "github.com/drunkleen/rasta/internal/common/errors"
	"github.com/drunkleen/rasta/internal/service/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type OAuthController struct {
	OAuthService *userservice.OAuthService
	UserService  *userservice.UserService
}

func NewOAuthController(oauthService *userservice.OAuthService, userService *userservice.UserService) *OAuthController {
	return &OAuthController{OAuthService: oauthService, UserService: userService}
}

func (c *OAuthController) GenerateOAuth(ctx *gin.Context) {
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
	user, err := c.UserService.FindById(id)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		ctx.JSON(http.StatusUnauthorized, commonerrors.NewErrorMap(err.Error()))
		return
	}
	if user.OAuth.Enabled {
		ctx.JSON(http.StatusUnauthorized, commonerrors.NewErrorMap("OAuth is already enabled"))
		return
	}
	oauthSecret, oauthUrl, err := c.OAuthService.GenerateOAuthSecret(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, commonerrors.NewErrorMap(commonerrors.ErrInternalServer))
		return
	}
	ctx.JSON(http.StatusOK, oauthDTO.ToOAuthResponse("", oauthSecret, oauthUrl, false))
}

func (c *OAuthController) VerifyAndEnableOAuth(ctx *gin.Context) {
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
	var reqBody map[string]string
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(commonerrors.ErrInvalidRequestBody))
		return
	}
	oauth, exists := reqBody["oauth"]
	if !exists || oauth == "" {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(commonerrors.ErrInvalidRequestBody))
		return
	}
	user, err := c.UserService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, commonerrors.NewErrorMap(err.Error()))
		return
	}
	if user.OAuth.Enabled {
		ctx.JSON(http.StatusUnauthorized, commonerrors.NewErrorMap("OAuth is already enabled"))
		return
	}
	if err = c.OAuthService.OAuthValidate(user, oauth); err != nil {
		ctx.JSON(http.StatusUnauthorized, commonerrors.NewErrorMap("Invalid OAuth code"))
		return
	}
	if err = c.OAuthService.UpdateOAuthEnabled(user.Id, true); err != nil {
		ctx.JSON(http.StatusInternalServerError, commonerrors.NewErrorMap(commonerrors.ErrInternalServer))
		return
	}
	ctx.JSON(http.StatusOK, oauthDTO.ToOAuthResponse("Otp enabled", "", "", true))
}

func (c *OAuthController) DisableOAuth(ctx *gin.Context) {
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
	var reqBody map[string]string
	if err = ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(commonerrors.ErrInvalidRequestBody))
		return
	}
	otp, exists := reqBody["oauth"]
	if !exists || otp == "" {
		ctx.JSON(http.StatusBadRequest, commonerrors.NewErrorMap(commonerrors.ErrInvalidRequestBody))
		return
	}
	user, err := c.UserService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, commonerrors.NewErrorMap(err.Error()))
		return
	}
	if !user.OAuth.Enabled {
		ctx.JSON(http.StatusUnauthorized, commonerrors.NewErrorMap("OAuth is already disabled"))
		return
	}
	if err = c.OAuthService.OAuthValidate(user, otp); err != nil {
		ctx.JSON(http.StatusUnauthorized, commonerrors.NewErrorMap("Invalid OAuth code"))
		return
	}
	if err = c.OAuthService.DeleteOAuth(user.Id); err != nil {
		ctx.JSON(http.StatusInternalServerError, commonerrors.NewErrorMap(commonerrors.ErrInternalServer))
		return
	}
	ctx.JSON(http.StatusOK, oauthDTO.ToOAuthResponse("OAuth disabled", "", "", false))
}
