package oauthDTO

type Response struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Token   string `json:"oauth_token,omitempty"`
	OtpUrl  string `json:"oauth_url,omitempty"`
	Enabled bool   `json:"is_active"`
}

// ToOAuthResponse generates an OAuth response based on the provided parameters.
//
// Parameter message is the response message, token is the OAuth token, url is the OAuth URL, and isActive is a boolean indicating whether the response is active or not.
// Return type is a pointer to the Response struct.
func ToOAuthResponse(message, token, url string, isActive bool) *Response {
	return &Response{
		Status:  "success",
		Message: message,
		Token:   token,
		OtpUrl:  url,
		Enabled: isActive,
	}
}
