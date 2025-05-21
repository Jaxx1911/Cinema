package fault

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

const (
	TagBadRequest     = "BAD_REQUEST" // 400
	TagUnAuthorize    = "UnAuthorize" // 401
	TagNotFound       = "NOT_FOUND"
	TagInternalServer = "INTERNAL_SERVER" // 500
	TagDuplicate      = "DUPLICATE"
	TagForbidden      = "FORBIDDEN"
)

const (
	KeyAuth     = "auth"
	KeyUser     = "user"
	KeyMovie    = "movie"
	KeyShowtime = "showtime"
	KeyOtp      = "otp"
	KeyRoom     = "room"
	KeyOrder    = "order"
	KeyPayment  = "payment"
	KeyCombo    = "combo"
	KeyCinema   = "cinema"
	KeyTicket   = "ticket"
	KeyDb       = "DB"
)

type Error struct {
	Message string
	Err     error
	Tag     string
	Key     string
}

// Implement error interface
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("error: %v, key: %s", e.Err, e.Key)
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

func (e *Error) SetKey(key string) *Error {
	e.Key = key
	return e
}

var TagMap = map[string]int{
	TagBadRequest:     400,
	TagUnAuthorize:    401,
	TagForbidden:      403,
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

func GetKey(err error) string {
	var fault *Error
	if errors.As(err, &fault) {
		return fault.Key
	}
	return ""
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
	ErrDBNotFound = Error{
		Message: "Not Found",
		Err:     gorm.ErrRecordNotFound,
		Tag:     TagNotFound,
	}
	ErrInternalServer = Error{
		Message: "Internal Server",
		Err:     errors.New("internal server error"),
		Tag:     TagInternalServer,
	}
)
