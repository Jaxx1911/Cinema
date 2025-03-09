package controller

import (
	"TTCS/src/common"
	"TTCS/src/common/log"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type BaseController struct {
	validate *validator.Validate
}

func NewBaseController(validate *validator.Validate) *BaseController {
	return &BaseController{
		validate: validate,
	}
}

func (b *BaseController) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func (b *BaseController) Error(c *gin.Context, err *common.Error) {
	c.JSON(err.GetHttpStatus(), err)
}

func (b *BaseController) Redirect(c *gin.Context, url string) {
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (b *BaseController) BindAndValidateRequest(c *gin.Context, req interface{}) *common.Error {
	if err := c.BindUri(req); err != nil {
		log.Warn(c, "bind request err, err:[%s]", err)
		return common.ErrBadRequest(c).SetDetail(err.Error())
	}
	if err := c.Bind(req); err != nil {
		log.Warn(c, "bind request err, err:[%s]", err)
		return common.ErrBadRequest(c).SetDetail(err.Error())
	}
	return b.ValidateRequest(c, req)
}

func (b *BaseController) ValidateRequest(ctx context.Context, req interface{}) *common.Error {
	err := b.validate.Struct(req)

	if err != nil {
		var errs validator.ValidationErrors
		if !errors.As(err, &errs) {
			log.Error(ctx, "Cannot parse validate error: %+v", err)
			return common.ErrInternal(ctx, "ValidateFailed").SetDetail(err.Error())
		}
		var filedErrors []string
		for _, errValidate := range errs {
			log.Debug(ctx, "field invalid, err:[%s]", errValidate.Field())
			filedErrors = append(filedErrors, errValidate.Error())
		}
		str := strings.Join(filedErrors, ",")
		log.Warn(ctx, "invalid request, err:[%s]", err.Error())
		return common.ErrBadRequest(ctx).SetDetail(fmt.Sprintf("field invalidate [%s]", str))
	}
	return nil
}
