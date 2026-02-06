package main

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var etcdClient *clientv3.Client

// func init() {
// 	var err error
// 	etcdClient, err = clientv3.New(clientv3.Config{
// 		Endpoints:   []string{"localhost:2379"},
// 		DialTimeout: 5 * time.Second,
// 	})
// 	if err != nil {
// 		Log.Error("Failed to connect to etcd:", err)
// 		return
// 	}
// 	Log.Info("Connected to etcd successfully")
// }

func go_etcd_ts1() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = etcdClient.Put(ctx, "sample_key", "sample_value-1")
	cancel()
	if err != nil {
		Log.Error("Failed to put key-value pair:", err)
		return
	}
	Log.Info("Successfully put key-value pair to etcd")
}

func go_etcd_ts2() {
	rch := etcdClient.Watch(context.Background(), "sample_key")
	for wresp := range rch {
		for _, ev := range wresp.Events {
			Log.Debug(string(ev.Kv.Key), " : ", string(ev.Kv.Value))
		}
	}
	Log.Debug("Watch channel closed")
}

type GoEtcd struct{}

func (g GoEtcd) Test() {
	// 在后台 goroutine 中监听变化
	go go_etcd_ts2()

	// 等待 watch 启动
	time.Sleep(100 * time.Millisecond)

	// 触发变化
	go_etcd_ts1()

	// 等待一段时间观察输出
	time.Sleep(2 * time.Second)
}
