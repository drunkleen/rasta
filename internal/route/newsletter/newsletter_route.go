package newsletterroute

import (
	newslettercontroller "github.com/drunkleen/rasta/internal/controller/newsletter"
	"github.com/drunkleen/rasta/internal/middlewares"
	newsletterrepository "github.com/drunkleen/rasta/internal/repository/newsletter"
	newsletterservice "github.com/drunkleen/rasta/internal/service/newsletter"
	"github.com/drunkleen/rasta/pkg/database"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
	db := database.DB
	nlRepository := newsletterrepository.NewNewsletterRepository(db)
	nlService := newsletterservice.NewNewsletterService(nlRepository)
	nlController := newslettercontroller.NewNewsletterController(nlService)

	userRoute := r.Group("/users/newsletter")
	//userRoute.Use(middlewares.JWTAuthMiddleware)

	adminOnlyRoute := r.Group("/admin/newsletter")
	adminOnlyRoute.Use(middlewares.AdminAuthMiddleware)

	registerOpenRoutes(userRoute, nlController)
	registerAdminOnlyRoutes(adminOnlyRoute, nlController)
}

func registerOpenRoutes(r *gin.RouterGroup, newsletterController *newslettercontroller.NewsletterController) {
	r.POST("/subscribe", newsletterController.Subscribe)
	r.POST("/unsubscribe", newsletterController.Unsubscribe)
}
func registerAdminOnlyRoutes(r *gin.RouterGroup, newsletterController *newslettercontroller.NewsletterController) {
	r.GET("/subscribers", newsletterController.GetSubscribers)
	r.GET("/subscribers/count", newsletterController.GetSubscribersCount)
	r.GET("/unsubscribed/count", newsletterController.GetUnsubscribedCount)
	r.DELETE("/delete", newsletterController.DeleteSubscriber)
	r.POST("/send", newsletterController.SendNewsletterToEveryActiveParticipants)
}
