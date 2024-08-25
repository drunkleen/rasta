package newslettercontroller

import (
	newsletterDTO "github.com/drunkleen/rasta/internal/DTO/newsletter"
	commonerrors "github.com/drunkleen/rasta/internal/common/errors"
	"github.com/drunkleen/rasta/internal/common/utils"
	newsletterservice "github.com/drunkleen/rasta/internal/service/newsletter"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NewsletterController struct {
	NewsletterService *newsletterservice.NewsletterService
}

// NewNewsletterController creates a new instance of NewsletterController
//
// It takes a pointer to a newsletterservice.NewsletterService as a parameter to
// initialize the NewsletterController.
// It returns a pointer to the NewsletterController.
func NewNewsletterController(newsletterService *newsletterservice.NewsletterService) *NewsletterController {
	return &NewsletterController{NewsletterService: newsletterService}
}

// Subscribe godoc
// @Summary Subscribe to Newsletter
// @Description Subscribes the user to the newsletter with the provided email address.
// @Tags Newsletter
// @Accept  json
// @Produce  json
// @Param email body map[string]string true "User email"
// @Success 201 {object} newsletterDTO.GenericResponse "Successfully subscribed to newsletter"
// @Failure 400 {object} commonerrors.ErrorMap "Invalid request body"
// @Failure 406 {object} commonerrors.ErrorMap "Email already subscribed"
// @Failure 500 {object} commonerrors.ErrorMap "Internal Server Error"
// @Router /newsletter/subscribe [post]
func (c *NewsletterController) Subscribe(ctx *gin.Context) {
	var reqBody map[string]string
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, commonerrors.InvalidRequestBodyError())
		return
	}

	email, exists := reqBody["email"]
	if !exists || !utils.EmailValidate(&email) {
		ctx.JSON(http.StatusBadRequest, commonerrors.InvalidRequestBodyError())
		return
	}

	subscriber, err := c.NewsletterService.FindByEmail(&email)
	if err != nil {
		if err := c.NewsletterService.Create(&email); err != nil {
			ctx.JSON(http.StatusInternalServerError, commonerrors.InternalServerError())
			return
		}
		ctx.JSON(http.StatusCreated, newsletterDTO.GenericResponse{
			Status:  "success",
			Message: "Successfully subscribed for newsletter",
		})
	}

	if (*subscriber).IsActive {
		ctx.JSON(http.StatusNotAcceptable, commonerrors.EmailAlreadyExistsError())
		return
	}

	if err = c.NewsletterService.UpdateActiveStatus(&email, true); err != nil {
		ctx.JSON(http.StatusInternalServerError, commonerrors.InternalServerError())
	}
	ctx.JSON(http.StatusCreated, newsletterDTO.GenericResponse{
		Status:  "success",
		Message: "Successfully subscribed for newsletter",
	})

}

// Unsubscribe godoc
// @Summary Unsubscribe from Newsletter
// @Description Unsubscribes the user from the newsletter using the provided email address.
// @Tags Newsletter
// @Accept  json
// @Produce  json
// @Param email body map[string]string true "User email"
// @Success 200 {object} newsletterDTO.GenericResponse "Successfully unsubscribed from newsletter"
// @Failure 400 {object} commonerrors.ErrorMap "Invalid request body"
// @Failure 406 {object} commonerrors.ErrorMap "Email not subscribed"
// @Failure 500 {object} commonerrors.ErrorMap "Internal Server Error"
// @Router /newsletter/unsubscribe [post]
func (c *NewsletterController) Unsubscribe(ctx *gin.Context) {
	var reqBody map[string]string
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, commonerrors.InvalidRequestBodyError())
		return
	}

	email, exists := reqBody["email"]
	if !exists || !utils.EmailValidate(&email) {
		ctx.JSON(http.StatusBadRequest, commonerrors.InvalidRequestBodyError())
		return
	}

	subscriber, err := c.NewsletterService.FindByEmail(&email)
	if err != nil {
		ctx.JSON(http.StatusCreated, commonerrors.EmailNotExistsError())
	}

	if !(*subscriber).IsActive {
		ctx.JSON(http.StatusNotAcceptable, commonerrors.EmailNotExistsError())
		return
	}

	if err = c.NewsletterService.UpdateActiveStatus(&email, false); err != nil {
		ctx.JSON(http.StatusInternalServerError, commonerrors.InternalServerError())
	}
	ctx.JSON(http.StatusCreated, newsletterDTO.GenericResponse{
		Status:  "success",
		Message: "Successfully unsubscribed from newsletter",
	})
}

// DeleteSubscriber godoc
// @Summary Delete Subscriber
// @Description Deletes a subscriber from the newsletter system using the provided email address.
// @Tags Newsletter
// @Accept  json
// @Produce  json
// @Param email body map[string]string true "User email"
// @Success 200 {object} newsletterDTO.GenericResponse "Successfully deleted subscriber"
// @Failure 400 {object} commonerrors.ErrorMap "Invalid request body"
// @Failure 500 {object} commonerrors.GenericResponseError "Internal Server Error"
// @Router /newsletter/delete [delete]
func (c *NewsletterController) DeleteSubscriber(ctx *gin.Context) {
	var reqBody map[string]string
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, commonerrors.InvalidRequestBodyError())
		return
	}
	email, exists := reqBody["email"]
	if !exists || !utils.EmailValidate(&email) {
		ctx.JSON(http.StatusBadRequest, commonerrors.InvalidRequestBodyError())
		return
	}
	if err := c.NewsletterService.DeleteByEmail(&email); err != nil {
		ctx.JSON(http.StatusInternalServerError, commonerrors.GenericResponseError{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, newsletterDTO.GenericResponse{
		Status:  "success",
		Message: "Successfully deleted subscriber",
	})
}

// GetSubscribers godoc
// @Summary Get Active Subscribers
// @Description Retrieves a list of all active newsletter subscribers.
// @Tags Newsletter
// @Produce  json
// @Success 200 {object} newsletterDTO.GenericResponse "Successfully fetched subscribers"
// @Failure 500 {object} commonerrors.ErrorMap "Internal Server Error"
// @Router /newsletter/subscribers [get]
func (c *NewsletterController) GetSubscribers(ctx *gin.Context) {
	subscribers, err := c.NewsletterService.FindAllActive()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, commonerrors.InternalServerError())
		return
	}
	ctx.JSON(http.StatusOK, newsletterDTO.GenericResponse{
		Status:  "success",
		Message: "Successfully fetched subscribers",
		Data:    subscribers,
	})
}

// GetSubscribersCount godoc
// @Summary Get Active Subscribers Count
// @Description Retrieves the count of active newsletter subscribers.
// @Tags Newsletter
// @Produce  json
// @Success 200 {object} newsletterDTO.GenericResponse "Successfully fetched subscribers count"
// @Failure 500 {object} commonerrors.ErrorMap "Internal Server Error"
// @Router /newsletter/subscribers/count [get]
func (c *NewsletterController) GetSubscribersCount(ctx *gin.Context) {
	count, err := c.NewsletterService.CountActiveSubscribers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, commonerrors.InternalServerError())
		return
	}
	ctx.JSON(http.StatusOK, newsletterDTO.GenericResponse{
		Status:  "success",
		Message: "Successfully fetched subscribers count",
		Data: struct {
			SubscribersCount int64 `json:"subscribers_count"`
		}{
			SubscribersCount: count,
		},
	})
}

// GetUnsubscribedCount godoc
// @Summary Get Unsubscribed Count
// @Description Retrieves the count of unsubscribed users from the newsletter.
// @Tags Newsletter
// @Produce  json
// @Success 200 {object} newsletterDTO.GenericResponse "Successfully fetched unsubscribed count"
// @Failure 500 {object} commonerrors.ErrorMap "Internal Server Error"
// @Router /newsletter/unsubscribed/count [get]
func (c *NewsletterController) GetUnsubscribedCount(ctx *gin.Context) {
	count, err := c.NewsletterService.CountInactiveSubscribers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, commonerrors.InternalServerError())
		return
	}
	ctx.JSON(http.StatusOK, newsletterDTO.GenericResponse{
		Status:  "success",
		Message: "Successfully fetched unsubscribed count",
		Data: struct {
			UnsubscribedCount int64 `json:"unsubscribed_count"`
		}{
			UnsubscribedCount: count,
		},
	})
}

// SendNewsletterToEveryActiveParticipants godoc
// @Summary Send Newsletter to Active Subscribers
// @Description Sends the newsletter email to all active subscribers.
// @Tags Newsletter
// @Accept  json
// @Produce  json
// @Param newsletter body newsletterDTO.CreateNewsletterRequest true "Newsletter content and limit"
// @Success 200 {object} newsletterDTO.GenericResponse "Successfully sent newsletter to all active participants"
// @Failure 400 {object} commonerrors.ErrorMap "Invalid request body"
// @Failure 500 {object} commonerrors.GenericResponseError "Internal Server Error"
// @Router /newsletter/send [post]
func (c *NewsletterController) SendNewsletterToEveryActiveParticipants(ctx *gin.Context) {
	var newsletterReq newsletterDTO.CreateNewsletterRequest
	if err := ctx.ShouldBindJSON(&newsletterReq); err != nil {
		ctx.JSON(http.StatusBadRequest, commonerrors.InvalidRequestBodyError())
		return
	}
	if newsletterReq.Limit < 10 {
		newsletterReq.Limit = 10
	}
	err := c.NewsletterService.SendNewslettersEmail(&newsletterReq.EmailText, newsletterReq.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, commonerrors.GenericResponseError{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, newsletterDTO.GenericResponse{
		Status:  "success",
		Message: "Successfully sent newsletter to every active participants",
	})
}
