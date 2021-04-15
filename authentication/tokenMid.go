package mid

import (
	"ZhiHu/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Token()gin.HandlerFunc{
	return func(ctx *gin.Context) {
		tokenHeader:=ctx.GetHeader("Authorization")
		_,err:=auth.CheckJWT(tokenHeader)
		if err!=nil{
			ctx.String(http.StatusBadRequest,"token验证失败")
			ctx.Abort()
		}else {
			ctx.Next()
		}
	}
}