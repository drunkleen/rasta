package oauthDTO

type Response struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Token   string `json:"oauth_token,omitempty"`
	OtpUrl  string `json:"oauth_url,omitempty"`
	Enabled bool   `json:"is_active"`
}

func ToOAuthResponse(message, token, url string, isActive bool) *Response {
	return &Response{
		Status:  "success",
		Message: message,
		Token:   token,
		OtpUrl:  url,
		Enabled: isActive,
	}
}
