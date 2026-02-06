package main

import (
	"strconv"
	"time"
)

type GoRabbit struct{}

func go_rabbit_ts1() {
	rabbitmq := NewRabbitMQSimple("kuteng")
	rabbitmq.PublishSimple("Hello RabbitMQ 2!")
	Log.Info("Message published to RabbitMQ")
}

func go_rabbit_ts2() {
	rev := NewRabbitMQSimple("kuteng")
	rev.ConsumeSimple()
	Log.Info("Started consuming messages from RabbitMQ")
}

func go_rabbit_ts3() {
	r3 := NewRabbitMQSimple("kuteng")
	for i := 0; i <= 100; i++ {
		r3.PublishSimple("Hello RabbitMQ 3 ! " + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		Log.Debug(i)
	}
}

func go_rabbit_ts4() {
	r4 := NewRabbitMQPubSub("newProducer")
	for i := 0; i <= 100; i++ {
		r4.PublishPub("订阅模式生产第" +
			strconv.Itoa(i) + "条" + "数据")
		Log.Debug("订阅模式生产第" +
			strconv.Itoa(i) + "条" + "数据")
		time.Sleep(1 * time.Second)
	}
}

func go_rabbit_ts5() {
	r5 := NewRabbitMQPubSub("newProducer")
	r5.RecieveSub()
	Log.Info("订阅模式开始消费数据")
}

func go_rabbit_ts6() {
	k1 := NewRabbitMQRouting("kuteng", "kuteng_one")
	K2 := NewRabbitMQRouting("kuteng", "kuteng_two")

	for i := 0; i <= 100; i++ {
		k1.PublishRouting("hello kuteng one!" + strconv.Itoa(i))
		K2.PublishRouting("hello kuteng two!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		Log.Debug(i)
	}
}

func go_rabbit_ts7() {
	r7 := NewRabbitMQRouting("kuteng", "kuteng_one")
	r7.RecieveRouting()
	Log.Info("路由模式开始消费数据")
}

func (g GoRabbit) Test() {
	// go_rabbit_ts1()
	go_rabbit_ts7()
	// go_rabbit_ts3()
}
