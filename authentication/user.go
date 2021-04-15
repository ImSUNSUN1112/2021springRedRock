package controller

import (
	"ZhiHu/auth"
	"ZhiHu/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignUp(ctx *gin.Context){
	ok:=service.SignUp(ctx)

	if ok!=true{
		ctx.String(http.StatusBadRequest,"注册失败")
		fmt.Println("注册失败")
	}else {
		ctx.String(http.StatusOK,"注册成功")
		fmt.Println("注册成功")
	}

}

func SignIn(ctx *gin.Context){
	ok,username,telephone:=service.SignIn(ctx)
	if ok!=true{
		ctx.String(http.StatusBadRequest,"登录失败")
		fmt.Println("登录失败")
	}else {
		//登录成功,以用户名与手机号生成token
		ctx.JSON(http.StatusOK,gin.H{
			"token":auth.GetJWT(username,telephone).Token,
			"message":"登录成功",

		})
		fmt.Println("登录成功")
	}
}

func TestToken(ctx *gin.Context){
	ctx.String(http.StatusOK,"看来你通过token验证了")
}
