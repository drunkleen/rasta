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

// JWTAuthMiddleware is a middleware function that authenticates a user
// using a JSON Web Token (JWT) in the Authorization header.
//
// This function takes a pointer to a gin.Context as its sole argument.
// It first retrieves the JWT from the Authorization header and aborts the
// request if it is empty. It then validates the JWT using the
// `auth.ValidateJWTToken` function and aborts the request if validation
// fails. If validation succeeds, the function sets the user ID and email
// in the gin.Context and calls c.Next() to proceed to the next handler.
//
// The function does not return anything.
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

// AdminAuthMiddleware is a middleware function that authenticates and authorizes a user
// as an admin.
//
// This function takes a pointer to a gin.Context as its sole argument.
// It first extracts and validates the JWT token from the Authorization header.
// If the token is invalid, the function aborts the request with a 401 Unauthorized status.
// If the token is valid, the function retrieves the user information from the database.
// If the user is not an admin, the function aborts the request with a 403 Forbidden status.
// If the user is an admin, the function sets the user ID and email in the gin.Context
// and proceeds to the next handler.
//
// The function does not return anything.
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
