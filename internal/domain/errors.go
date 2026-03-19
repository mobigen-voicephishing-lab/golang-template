package domain

import (
	"net/http"
)

// AppError 도메인/비즈니스 에러. error 인터페이스 구현.
// Service 계층에서 생성하여 Controller로 전달한다.
type AppError struct {
	HttpStatus int    // HTTP 응답 상태 코드 (예: 404, 400, 409)
	Code       int    // 내부 에러 코드
	Message    string // 에러 메시지
}

func (e AppError) Error() string {
	return e.Message
}

// NewAppError AppError 생성 헬퍼.
func NewAppError(httpStatus int, code int, msg string) AppError {
	return AppError{HttpStatus: httpStatus, Code: code, Message: msg}
}

// NewNotFoundError 404 Not Found AppError 생성 shortcut
func NewNotFoundError(msg string) AppError {
	return AppError{HttpStatus: http.StatusNotFound, Code: 1001, Message: msg}
}

// NewBadRequestError 400 Bad Request AppError 생성 shortcut
func NewBadRequestError(code int, msg string) AppError {
	return AppError{HttpStatus: http.StatusBadRequest, Code: code, Message: msg}
}

// NewInternalError 500 Internal Server Error AppError 생성 shortcut
func NewInternalError(msg string) AppError {
	return AppError{HttpStatus: http.StatusInternalServerError, Code: 5000, Message: msg}
}
