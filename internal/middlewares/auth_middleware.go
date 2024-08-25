package middlewares

import (
	"errors"
	"github.com/drunkleen/rasta/internal/common/auth"
	commonerrors "github.com/drunkleen/rasta/internal/common/errors"
	usermodel "github.com/drunkleen/rasta/internal/models/user"
	userrepository "github.com/drunkleen/rasta/internal/repository/user"
	userservice "github.com/drunkleen/rasta/internal/service/user"
	"github.com/drunkleen/rasta/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

var (
	userRepository = userrepository.NewUserRepository(database.DB)
	userService    = userservice.NewUserService(userRepository)
)

// extractAndValidateToken extracts and validates a JWT token from the Authorization header.
//
// If the JWT token is empty, the function returns an error.
// If the token is invalid, the function returns an error.
// If the token is valid, the function returns the user ID and email associated with the token.
//
// Parameters:
// c *gin.Context is the gin context.
//
// Returns:
// uuid.UUID is the user ID associated with the token.
// string is the user email associated with the token.
// error is an error object that is returned if the token is invalid or empty.
func extractAndValidateToken(c *gin.Context) (uuid.UUID, string, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		return uuid.Nil, "", errors.New(commonerrors.ErrUnauthorizedToken)
	}
	userIdStr, userEmail, err := auth.ValidateJWTToken(token)
	if err != nil {
		return uuid.Nil, "", errors.New(commonerrors.ErrUnauthorizedToken)
	}
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return uuid.Nil, "", errors.New(commonerrors.ErrUnauthorizedToken)
	}
	return userId, userEmail, nil
}

// JWTAuthMiddleware authenticates a user by validating the JWT token in the Authorization header.
//
// Parameter c *gin.Context is the gin context.
//
// Returns None
func JWTAuthMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, commonerrors.NewErrorMap(commonerrors.ErrUnauthorizedToken))
		return
	}
	userId, userEmail, err := auth.ValidateJWTToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, commonerrors.NewErrorMap(commonerrors.ErrUnauthorizedToken))
		return
	}
	c.Set("userId", userId)
	c.Set("userEmail", userEmail)
	c.Next()
}

// AdminAuthMiddleware is a middleware function that authenticates and authorizes admin users.
//
// Parameters:
// c *gin.Context is the gin context.
//
// Returns:
// None
func AdminAuthMiddleware(c *gin.Context) {
	userId, userEmail, err := extractAndValidateToken(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, commonerrors.NewErrorMap(err.Error()))
		return
	}
	userModel, err := userService.FindById(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, commonerrors.NewErrorMap(err.Error()))
		return
	}
	if userModel.Account != usermodel.AccountTypeAdmin {
		c.AbortWithStatusJSON(http.StatusForbidden, commonerrors.NewErrorMap(commonerrors.ErrForbidden))
		return
	}
	c.Set("userId", userModel.Id)
	c.Set("userEmail", userEmail)
	c.Set("userModel", userModel)
	c.Next()
}
