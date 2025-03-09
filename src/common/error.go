package common

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type ErrorCode int

func (e ErrorCode) Error() string {
	return strconv.Itoa(int(e))
}

const (
	ErrorCodeBadRequest   ErrorCode = 400
	ErrorCodeNotFound     ErrorCode = 404
	ErrorCodeUnauthorized ErrorCode = 401
	ErrorCodeInternal     ErrorCode = 500
)

type Error struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Detail     string    `json:"detail"`
	HTTPStatus int       `json:"http_status"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("code:[%d], message:[%s], detail:[%s]", e.Code, e.Message, e.Detail)
}

func (e *Error) GetHttpStatus() int {
	return e.HTTPStatus
}

func (e *Error) GetCode() ErrorCode {
	return e.Code
}

func (e *Error) GetMessage() string {
	return e.Message
}

func (e *Error) SetCode(code ErrorCode) *Error {
	e.Code = code
	return e
}

func (e *Error) SetMessage(message string) *Error {
	e.Message = message
	return e
}

func (e *Error) SetDetail(detail string) *Error {
	e.Detail = detail
	return e
}

func (e *Error) GetDetail() string {
	return e.Detail
}

var (
	ErrUnauthorized = func(ctx context.Context) *Error {
		return &Error{
			Code:       ErrorCodeUnauthorized,
			Message:    DefaultUnauthorizedMessage,
			HTTPStatus: http.StatusUnauthorized,
		}
	}

	ErrBadRequest = func(ctx context.Context) *Error {
		return &Error{
			Code:       ErrorCodeBadRequest,
			Message:    DefaultBadRequestMessage,
			HTTPStatus: http.StatusBadRequest,
		}
	}

	ErrNotFound = func(ctx context.Context) *Error {
		return &Error{
			Code:       ErrorCodeNotFound,
			Message:    DefaultResourceMessage,
			HTTPStatus: http.StatusNotFound,
		}
	}

	ErrInternal = func(ctx context.Context, detail string) *Error {
		return &Error{
			Code:       ErrorCodeInternal,
			Message:    DefaultInternalErrorMessage,
			HTTPStatus: http.StatusInternalServerError,
			Detail:     detail,
		}
	}
)

const (
	DefaultInternalErrorMessage = "Something has gone wrong, please contact admin"
	DefaultResourceMessage      = "resource not found"
	DefaultBadRequestMessage    = "Invalid request"
	DefaultUnauthorizedMessage  = "Token invalid"
)

func IsInternalError(err *Error) bool {
	if err == nil {
		return false
	}
	if err.GetHttpStatus() >= http.StatusInternalServerError {
		return true
	}
	return false
}

func IsNotFoundError(err *Error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err.GetCode(), ErrorCodeNotFound)
}
