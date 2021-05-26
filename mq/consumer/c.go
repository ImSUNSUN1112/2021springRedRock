package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func main(){
	//服务端连接消息队列
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

	//同样的声明队列
	//name:队列名，同名进，同名出
	//其他的字段基本上都是一些关于持久化与否的，自动删除之类的
	q, err := ch.QueueDeclare("sunsun", false, false, false, false, nil)
	if err != nil {
		log.Fatalln(err)
	}

	//设置每次从消息队列中获取的任务量
	err = ch.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		log.Fatalln(err)
	}

	//从消息队列中取出信息
	//consumer:
	//An empty string will cause,the library to generate a unique identity.  The consumer identity will be,included in every Delivery in the ConsumerTag field
	//起划分区域的作用
	//autoAck:自动确认与否
	//exclusive:是否独占队列
	//noLocal:RabbitMQ不支持nollocal标志,所以为什么这个方法会有这个字段呢?
	//noWait:当noWait为true时，不要等待服务器确认请求,立即开始交付。如果不可能消费，信道,异常将被引发，通道将被关闭
	//个人理解是会覆盖前面的autoAck字段
	//args:可选参数
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalln(err)
	}

	//我对这里的理解是相当于一个select的作用，监听这个管道，有相关队列的信息就进行读取
	for v:=range msgs{
		if v.AppId=="SignUp"{
			//输出读取到的注册信息,应该按name，passwd的格式
			//后序的注册服务就可以采用按,分割然后进行
			fmt.Printf("body=%s\n",v.Body)
			//读取到信息以后手动确认
			v.Ack(true)
		}
	}
}
