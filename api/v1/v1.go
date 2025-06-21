package v1

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HandleSuccess(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	resp := Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
	ctx.JSON(http.StatusOK, resp)
}

func HandleCreated(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	resp := Response{
		Code:    0,
		Message: "created successfully",
		Data:    data,
	}
	ctx.JSON(http.StatusCreated, resp)
}

func HandleError(ctx *gin.Context, httpCode int, message string, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	resp := Response{
		Code:    httpCode,
		Message: message,
		Data:    data,
	}
	ctx.JSON(httpCode, resp)
}

func HandleBadRequest(ctx *gin.Context, message string) {
	HandleError(ctx, http.StatusBadRequest, message, nil)
}

func HandleNotFound(ctx *gin.Context, message string) {
	HandleError(ctx, http.StatusNotFound, message, nil)
}

func HandleInternalError(ctx *gin.Context, message string) {
	HandleError(ctx, http.StatusInternalServerError, message, nil)
}

func HandleValidationError(ctx *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var messages []string
		for _, ve := range validationErrors {
			switch ve.Tag() {
			case "required":
				messages = append(messages, ve.Field()+" é obrigatório")
			case "gt":
				messages = append(messages, ve.Field()+" deve ser maior que "+ve.Param())
			case "len":
				messages = append(messages, ve.Field()+" deve ter "+ve.Param()+" caracteres")
			case "brazilian_state":
				messages = append(messages, ve.Field()+" deve ser um estado brasileiro válido")
			case "uuid":
				messages = append(messages, ve.Field()+" deve ser um UUID válido")
			default:
				messages = append(messages, ve.Field()+" é inválido")
			}
		}
		HandleError(ctx, http.StatusBadRequest, strings.Join(messages, ", "), nil)
		return
	}

	HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
}

func HandleDatabaseError(ctx *gin.Context, err error, message string) {
	if errors.Is(err, sql.ErrNoRows) {
		HandleNotFound(ctx, message)
		return
	}
	HandleInternalError(ctx, message)
}
