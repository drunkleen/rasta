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

func NewNewsletterController(newsletterService *newsletterservice.NewsletterService) *NewsletterController {
	return &NewsletterController{NewsletterService: newsletterService}
}

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
