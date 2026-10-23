package validatex

import "net/http"

// ErrorResponse 定义通用错误响应
type ErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// NewErrorResponse 创建通用错误响应
func NewErrorResponse(code int, message string, data any) ErrorResponse {
	return ErrorResponse{
		Code: code,
		Msg:  message,
		Data: data,
	}
}

// NewBindErrorResponse 创建绑定错误响应
func NewBindErrorResponse(data, message string) ErrorResponse {
	return NewErrorResponse(http.StatusBadRequest, message, data)
}

// NewValidationErrorResponse 创建校验错误响应
func NewValidationErrorResponse(errors Errors, message string) ErrorResponse {
	return NewErrorResponse(http.StatusUnprocessableEntity, message, errors.ToMap())
}
