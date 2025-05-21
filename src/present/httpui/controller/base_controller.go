package controller

import (
	"TTCS/src/common/fault"
	"TTCS/src/core/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

const SuccessCode = 0

type Response struct {
	Err     string      `json:"err"`
	Key     string      `json:"key,omitempty"`
	Body    interface{} `json:"body"`
	Message string      `json:"message"`
}

type BaseController struct {
	validate *validator.Validate
}

func NewBaseController(validate *validator.Validate) *BaseController {
	return &BaseController{
		validate: validate,
	}
}

func (*BaseController) GetAuthUser(c *gin.Context) *domain.User {
	user, exist := c.Get("user")
	if !exist {
		return nil
	}
	domainUser, ok := user.(*domain.User)
	if !ok {
		return nil
	}
	return domainUser
}

func (*BaseController) ServeSuccessResponse(c *gin.Context, body interface{}) {

	c.JSON(http.StatusOK, Response{
		Body:    body,
		Message: "success",
	})
}

func (*BaseController) ServeErrResponse(c *gin.Context, err error, statusCodes ...int) {
	var statusCode int
	if len(statusCodes) > 0 {
		statusCode = statusCodes[0]
	} else {
		statusCode = fault.GetStatusCode(err)
	}

	errRes := Response{
		Err:     err.Error(),
		Key:     fault.GetKey(err),
		Message: fault.GetMessage(err),
		Body:    nil,
	}

	c.JSON(statusCode, errRes)
}

func (*BaseController) ServeRedirect(c *gin.Context, url string) {
	c.Redirect(http.StatusTemporaryRedirect, url)
}
