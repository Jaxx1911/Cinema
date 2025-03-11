package fault

import (
	"errors"
	"fmt"
)

// Danh sách tag lỗi có sẵn
const (
	TagBadRequest     = "BAD_REQUEST" // 400
	TagUnAuthorize    = "UnAuthorize" // 401
	TagNotFound       = "NOT_FOUND"
	TagInternalServer = "INTERNAL_SERVER" // 500
)

// Error FaultError chứa thông tin lỗi, caller, tag
type Error struct {
	Message string
	Err     error
	Tag     string
}

// Implement error interface
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s -> %v", e.Message, e.Err)
	}
	return fmt.Sprintf("%s", e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func Wrap(err error) *Error {
	return &Error{
		Message: err.Error(),
		Err:     err,
	}
}

func Wrapf(err error, format string, args ...interface{}) *Error {
	message := fmt.Sprintf(format, args...)

	return &Error{
		Message: message,
		Err:     err,
	}
}

func (e *Error) SetTag(tag string) *Error {
	e.Tag = tag
	return e
}

var TagMap = map[string]int{
	TagBadRequest:     400,
	TagUnAuthorize:    401,
	TagNotFound:       404,
	TagInternalServer: 500,
}

func GetStatusCode(err error) int {
	var faultErr *Error
	if errors.As(err, &faultErr) {
		if statusCode, exists := TagMap[faultErr.Tag]; exists {
			return statusCode
		}
	}
	return 500
}

func GetMessage(err error) string {
	var fault *Error
	if errors.As(err, &fault) {
		return fault.Message
	}
	return err.Error()
}

var (
	ErrBadRequest = Error{
		Message: "Bad request",
		Err:     errors.New("bad request"),
		Tag:     TagBadRequest,
	} // 400
	ErrUnauthenticated = Error{
		Message: "UnAuthorized",
		Err:     errors.New("UnAuthorized"),
		Tag:     TagUnAuthorize,
	}
	ErrNotFound = Error{
		Message: "Not Found",
		Err:     errors.New("not found"),
		Tag:     TagNotFound,
	}
	ErrInternalServer = Error{
		Message: "Internal Server",
		Err:     errors.New("internal server error"),
		Tag:     TagInternalServer,
	}
)
