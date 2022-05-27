package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"saas/kernel/validator"
)

func Unauthorized(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Code:    40100,
		Message: "Unauthorized",
	})
}

func Forbidden(ctx *gin.Context) {
	ctx.JSON(http.StatusForbidden, Response{
		Code:    40400,
		Message: "Forbidden",
	})
}

func NotFound(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, Response{
		Code:    40400,
		Message: message,
	})
}

func FailByRequest(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, Response{
		Code:    40000,
		Message: validator.Translate(err),
	})
}

func FailByRequestWithMessage(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, Response{
		Code:    40000,
		Message: message,
	})
}

func FailByLogin(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Code:    40100,
		Message: "登陆失败",
	})
}

func Success(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Code:    20000,
		Message: "Success",
	})
}

func SuccessByData(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Response{
		Code:    20000,
		Message: "Success",
		Data:    data,
	})
}

func SuccessByList(ctx *gin.Context, list []any) {
	ctx.JSON(http.StatusOK, Responses{
		Code:    20000,
		Message: "Success",
		Data:    list,
	})
}

func SuccessByPaginate(ctx *gin.Context, data Paginate) {
	ctx.JSON(http.StatusOK, Response{
		Code:    20000,
		Message: "Success",
		Data:    data,
	})
}

func Fail(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, Response{
		Code:    60000,
		Message: message,
	})
}
