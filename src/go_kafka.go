package main

import (
	"github.com/IBM/sarama"
)

type GoKafka struct{}

func go_kafka_ts1() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true

	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = "web_log"
	msg.Value = sarama.StringEncoder("this is a test log")
	// 连接kafka
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		Log.Error("producer closed, err:", err)
		return
	}
	defer client.Close()
	// 发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		Log.Error("send msg failed, err:", err)
		return
	}

	Log.Debugf("pid:%v offset:%v\n", pid, offset)
}

func go_kafka_ts2() {
	Log.Debug("Start")

	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		Log.Error("failed to start consumer, err:", err)
		return
	}

	partitionList, err := consumer.Partitions("web_log")
	if err != nil {
		Log.Error("failed to get the list of partitions, err:", err)
		return
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition("web_log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			Log.Error("failed to start consumer for partition %d, err:", partition, err)
			return
		}

		defer pc.AsyncClose()

		// 异步从每个分区消费消息
		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				Log.Debugf("Partition:%d Offset:%d Key:%s Value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
		}(pc)
	}
}

func (g GoKafka) Test() {
	// go_kafka_ts1()
	go_kafka_ts2()
}
