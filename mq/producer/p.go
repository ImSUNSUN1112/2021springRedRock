package main

import (
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

type user struct {
	name string
	passwd string
}

func main(){

	//声明路由
	router := gin.Default()

	//客户端连接消息队列
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	//获取通道
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln(err)
	}

	//队列声明
	//name:队列名，同名进，同名出
	//其他的字段基本上都是一些关于持久化与否的，自动删除之类的
	q, err := ch.QueueDeclare("sunsun",false,false,false,false,nil)
	if err != nil {
		log.Fatalln(err)
	}

	router.GET("/user", func(ctx *gin.Context) {
		//获取注册信息
		var u user
		u.name=ctx.Request.FormValue("name")
		u.passwd=ctx.Request.FormValue("passwd")

		//只有一个body字段,所以我采用拼接的方法传递数据
		msg:=u.name+","+u.passwd

		//发送消息
		//翻查publishing发现有AppId这个字段，很适合用来做模块分化
		err = ch.Publish("", q.Name, false, false, amqp.Publishing{
			AppId: "SignUp",
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
		if err != nil {
			log.Fatalln(err)
			ctx.String(http.StatusBadRequest,"MQ写入失败")
			return
		}
		ctx.String(http.StatusOK,"MQ写入成功")
	})


	router.Run(":8080")
}

