package response

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data"`
}

func NewResponse(message string, success bool, data any) *Response {
	return &Response{
		Message: message,
		Success: success,
		Data:    data,
	}
}
