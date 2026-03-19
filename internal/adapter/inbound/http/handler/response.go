package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/mobigen/golang-web-template/internal/adapter/inbound/http/dto"
	"github.com/mobigen/golang-web-template/internal/domain"
)

// httpStatusForAppError 도메인 에러 코드를 HTTP 상태 코드로 변환한다.
// HTTP 관심사는 어댑터 레이어에서만 다룬다.
func httpStatusForAppError(ae domain.AppError) int {
	switch {
	case ae.Code == domain.ErrNotFound, ae.Code == domain.ErrRouteNotFound:
		return http.StatusNotFound
	case ae.Code == domain.ErrMethodNotAllowed:
		return http.StatusMethodNotAllowed
	case ae.Code == domain.ErrAlreadyExists, ae.Code == domain.ErrAlreadyProcessed:
		return http.StatusConflict
	case ae.Code == domain.ErrInvalidStatusTransition:
		return http.StatusUnprocessableEntity
	case ae.Code >= 2000 && ae.Code < 3000:
		return http.StatusBadRequest
	case ae.Code == domain.ErrUnauthorized, ae.Code == domain.ErrTokenExpired, ae.Code == domain.ErrInvalidToken:
		return http.StatusUnauthorized
	case ae.Code == domain.ErrForbidden:
		return http.StatusForbidden
	case ae.Code >= 4000 && ae.Code < 5000:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}

// OK 성공 응답 (HTTP 200)
func OK[T any](c *echo.Context, data T) error {
	return c.JSON(http.StatusOK, dto.HTTPResponse[T]{
		IsSuccess: true,
		Code:      domain.Success,
		Message:   "success",
		Data:      data,
	})
}

// Fail 실패 응답. msg가 빈 문자열이면 ErrMessages에서 기본 메시지를 사용한다.
func Fail(c *echo.Context, httpStatus int, code int, msg string) error {
	if msg == "" {
		if m, ok := domain.ErrMessages[code]; ok {
			msg = m
		} else {
			msg = "unknown error"
		}
	}
	return c.JSON(httpStatus, dto.HTTPResponse[any]{
		IsSuccess: false,
		Code:      code,
		Message:   msg,
	})
}

// FailApp AppError 기반 실패 응답.
func FailApp(c *echo.Context, ae domain.AppError) error {
	httpStatus := httpStatusForAppError(ae)
	msg := ae.Message
	if msg == "" {
		if m, ok := domain.ErrMessages[ae.Code]; ok {
			msg = m
		} else {
			msg = "unknown error"
		}
	}
	return c.JSON(httpStatus, dto.HTTPResponse[any]{
		IsSuccess: false,
		Code:      ae.Code,
		Message:   msg,
	})
}

// Wrap AOP-style 헬퍼. (T, error)를 반환하는 핸들러 함수를 echo.HandlerFunc로 변환한다.
func Wrap[T any](h func(c *echo.Context) (T, error)) echo.HandlerFunc {
	return func(c *echo.Context) error {
		result, err := h(c)
		if err != nil {
			if ae, ok := err.(domain.AppError); ok {
				return FailApp(c, ae)
			}
			return Fail(c, http.StatusInternalServerError, domain.ErrInternalServer, "")
		}
		return OK(c, result)
	}
}
