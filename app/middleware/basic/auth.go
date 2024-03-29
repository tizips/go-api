package basic

import (
	"github.com/gin-gonic/gin"
	"saas/app/constant"
	"saas/app/service/basic"
	"saas/kernel/authorize"
	"saas/kernel/response"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if !authorize.Check(ctx) {
			ctx.Abort()
			response.Unauthorized(ctx)
			return
		}

		claims := authorize.Jwt(ctx)

		if !claims.VerifyIssuer(constant.ContextAdmin, true) {
			ctx.Abort()
			response.Unauthorized(ctx)
			return
		}

		if !basic.CheckJwt(ctx, constant.ContextAdmin, *claims) {
			ctx.Abort()
			response.Unauthorized(ctx)
			return
		}

		ctx.Next()
	}
}
