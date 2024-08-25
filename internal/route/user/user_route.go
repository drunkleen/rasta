package userroute

import (
	"github.com/drunkleen/rasta/internal/controller/user"
	"github.com/drunkleen/rasta/internal/middlewares"
	"github.com/drunkleen/rasta/internal/repository/user"
	"github.com/drunkleen/rasta/internal/service/user"
	"github.com/drunkleen/rasta/pkg/database"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
	db := database.DB

	otpRepository := userrepository.NewOtpRepository(db)
	userRepository := userrepository.NewUserRepository(db)
	oauthRepository := userrepository.NewOAuthRepository(db)
	resetPwdRepositor := userrepository.NewResetPwdRepository(db)

	otpService := userservice.NewOtpService(otpRepository)
	userService := userservice.NewUserService(userRepository)
	oauthService := userservice.NewOAuthService(oauthRepository)
	resetPwdService := userservice.NewResetPwd(resetPwdRepositor)

	otpController := usercontroller.NewOtpController(otpService, userService)
	userController := usercontroller.NewUserController(userService, otpService)
	oauthController := usercontroller.NewOAuthController(oauthService, userService)
	resetPwdController := usercontroller.NewResetPwdController(resetPwdService, userService)

	userRoute := r.Group("/users")
	userRouteClosed := userRoute.Group("/")
	userRouteClosed.Use(middlewares.JWTAuthMiddleware)
	adminOnlyRoute := r.Group("/admin")
	adminOnlyRoute.Use(middlewares.AdminAuthMiddleware)

	registerOpenUserRoutes(userRoute, userController, resetPwdController)
	registerOpenOtpRoutes(userRoute, otpController)
	registerClosedUserRoutes(userRouteClosed, userController)
	registerClosedOAuthRoutes(userRouteClosed, oauthController)
	registerAdminRoutes(adminOnlyRoute, userController)
}

func registerOpenUserRoutes(r *gin.RouterGroup, userController *usercontroller.UserController, resetPwd *usercontroller.ResetPwdController) {
	r.POST("/login", userController.Login)
	r.POST("/signup", userController.Create)
	r.GET("/reset-password", resetPwd.Send)
	r.POST("/reset-password/:id/verify", resetPwd.VerifyAndResetPassword)
}

func registerOpenOtpRoutes(r *gin.RouterGroup, otpController *usercontroller.OtpController) {
	r.GET("/otp/resend", otpController.ResendOtp)
	r.POST("/otp/:id/verify", otpController.VerifyEmail)
}

func registerClosedUserRoutes(r *gin.RouterGroup, userController *usercontroller.UserController) {
	r.GET("/:username", userController.FindUserByUsername)
	r.GET("/:username/update-password", userController.UpdatePassword)
}

func registerClosedOAuthRoutes(r *gin.RouterGroup, oauthController *usercontroller.OAuthController) {
	r.GET("/oauth/generate", oauthController.GenerateOAuth)
	r.POST("/oauth/enable", oauthController.VerifyAndEnableOAuth)
	r.DELETE("/oauth/disable", oauthController.DisableOAuth)
}

func registerAdminRoutes(r *gin.RouterGroup, userController *usercontroller.UserController) {
	r.GET("/", userController.GetWithPagination)
	r.GET("/count", userController.GetAllUsersCount)
	r.GET("/id/:id", userController.FindUserByID)
}
