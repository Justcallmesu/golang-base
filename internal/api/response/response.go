package response

type BaseResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type Response struct {
	BaseResponse
	Data any `json:"data"`
}

type ErrorResponse struct {
	BaseResponse

	Errors any `json:"errors,omitempty"`
}

type APIValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewResponse(message string, success bool, data any) *Response {
	return &Response{
		BaseResponse: BaseResponse{
			Message: message,
			Success: success,
		},
		Data: data,
	}
}

func NewErrorResponse(message string, errors any) *ErrorResponse {
	return &ErrorResponse{
		BaseResponse: BaseResponse{
			Message: message,
			Success: false,
		},
		Errors: errors,
	}
}
