package commonerrors

type ErrorMap struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewErrorMap(message string) ErrorMap {
	return ErrorMap{
		Status:  "error",
		Message: message,
	}
}

type GenericResponseError struct {
	Status  string      `json:"status"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func InvalidRequestBodyError() *GenericResponseError {
	return &GenericResponseError{
		Status:  "error",
		Message: ErrInvalidRequestBody,
	}
}

func EmailAlreadyExistsError() *GenericResponseError {
	return &GenericResponseError{
		Status:  "error",
		Message: ErrEmailAlreadyExists,
	}
}

func EmailNotExistsError() *GenericResponseError {
	return &GenericResponseError{
		Status:  "error",
		Message: ErrEmailNotExists,
	}
}

func InternalServerError() *GenericResponseError {
	return &GenericResponseError{
		Status:  "error",
		Message: ErrInternalServer,
	}
}
