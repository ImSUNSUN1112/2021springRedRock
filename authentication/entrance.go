package cmd

import (
	"ZhiHu/controller"
	"ZhiHu/mid"
	"github.com/gin-gonic/gin"
)

func Entrance() {

	//启动路由
	router := gin.Default()

	//添加中间件解决跨域问题
	router.Use(mid.Cors())

	routerGroup := router.Group("ZhiHu")
	{
		/*
		测试类
		 */
		routerGroup.POST("/test/token",mid.Token(),controller.TestToken)

		/*
		用户信息类
		*/
		//注册接口
		routerGroup.POST("/user/signUp",controller.SignUp)
		//登录接口
		//登录成功返回对应token
		routerGroup.POST("/user/signIn",controller.SignIn)

	}

	router.Run()
}
