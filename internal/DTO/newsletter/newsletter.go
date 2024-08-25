package newsletterDTO

type GenericResponse struct {
	Status  string      `json:"status"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type CreateNewsletterRequest struct {
	EmailText string `json:"email_text" binding:"required"`
	Limit     int    `json:"limit"`
}
