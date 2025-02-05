package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go-gin-demo/global"
	"go-gin-demo/utils"
	"go.uber.org/zap"
	"reflect"
)

type BaseApi struct {
	Ctx    *gin.Context
	Errors error
	Logger *zap.SugaredLogger
}

func NewBaseApi() BaseApi {
	return BaseApi{
		Logger: global.Logger,
	}
}

type BuildRequestOption struct {
	Ctx     *gin.Context
	DTO     any
	BindUri bool
	BindAll bool
}

func (m *BaseApi) AddError(errNew error) {
	m.Errors = utils.AppendError(m.Errors, errNew)
}

func (m *BaseApi) GetError() error {
	return m.Errors
}

func (m *BaseApi) BuildRequest(option BuildRequestOption) *BaseApi {
	var errResult error
	m.Ctx = option.Ctx
	if option.DTO != nil {
		if option.BindAll || option.BindUri {
			errResult = utils.AppendError(errResult, m.Ctx.ShouldBindUri(option.DTO))
		}

		if option.BindAll || !option.BindUri {
			errResult = utils.AppendError(errResult, m.Ctx.ShouldBindJSON(option.DTO))
		}

		if errResult != nil {
			errResult = m.ParseValidateErrors(errResult, option.DTO)
			m.AddError(errResult)
			m.Fail(ResponseJson{
				Msg: m.GetError().Error(),
			})
		}

	}
	return m
}

func (m *BaseApi) ParseValidateErrors(errs error, target any) error {
	var errResult error

	errValidator, ok := errs.(validator.ValidationErrors)
	if !ok {
		return errs
	}
	fields := reflect.TypeOf(target).Elem()
	for _, fieldErr := range errValidator {
		field, _ := fields.FieldByName(fieldErr.Field())
		errMessageTag := fmt.Sprintf("%s_err", fieldErr.Tag())
		errMessage := field.Tag.Get(errMessageTag)
		if errMessage == "" {
			errMessage = field.Tag.Get("message")
		}
		if errMessage == "" {
			errMessage = fmt.Sprintf("%s: %s Error", fieldErr.Field(), fieldErr.Tag())
		}

		errResult = utils.AppendError(errResult, errors.New(errMessage))
	}
	return errResult
}

func (m *BaseApi) Fail(resp ResponseJson) {
	Fail(m.Ctx, resp)
}

func (m *BaseApi) OK(resp ResponseJson) {
	OK(m.Ctx, resp)
}

func (m *BaseApi) ServerFail(resp ResponseJson) {
	ServerFail(m.Ctx, resp)
}
