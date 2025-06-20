package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func HandleValidationError(ctx *gin.Context, message string) {
	HandleError(ctx, http.StatusUnprocessableEntity, message, nil)
}
