package main

import (
	"hello-go/src/logger"

	"go.uber.org/zap"
)

// func test1() {
// 	a := Add(1, 2)
// 	fmt.Println(a)
// }

// func test2() {
// 	a, s := Test1(1, 2, "cjl")
// 	fmt.Println(a, s)
// }

// func test3() {
// 	t := Test{}
// 	t.Method0()
// 	t.Method1(10, "hello")
// 	r2 := t.Method2()
// 	fmt.Println(r2)
// 	r3a, r3b := t.Method3(5, "go")
// 	fmt.Println(r3a, r3b)
// }

// func test4() {
// 	u := User{Name: "Alice", Age: 20}

// 	greeting := u.Greet("Hello")
// 	fmt.Println(greeting)

// 	isAdult := u.IsAdult()
// 	fmt.Println(isAdult)
// }

func test5() {
	http := TestHttp{}
	http.Start()
}

func test6() {
	o := TestGoRoutine{}
	o.Test()
}

func mysqlTest() {
	// t := MyGormSQL{}
	// t.Test()
	// t := GoRedis{}
	// t := GoEtcd{}
	// t := GoKafka{}
	// t := GoRabbit{}
	// t := GoGin{}
	// t := GoGin1{}
	// t.Test()

	// 在 goroutine 中启动 RPC 服务器
	// s := RpcService{}
	// go s.Start()

	// // 等待服务器启动完成
	// time.Sleep(1 * time.Second)

	// // 启动 RPC 客户端
	// c := RpcClient{}
	// c.Start()

	t := TcpScan{}
	t.Start()
}

func main() {
	defer logger.Sync() // 确保在程序退出前刷新日志缓冲区

	// 可选：使用自定义配置初始化 logger
	// logConfig := &logger.Config{
	// 	Level:      "debug",
	// 	Env:        "prod",
	// 	FilePath:   "logs/app.log",
	// 	MaxSize:    100,
	// 	MaxBackups: 10,
	// 	MaxAge:     30,
	// 	Compress:   true,
	// }
	// logger.Init(logConfig)

	mysqlTest()

	// 演示结构化日志
	logger.Info("应用启动",
		zap.String("version", "1.0.0"),
		zap.Int("port", 8080),
	)
	logger.Infof("应用启动",
		zap.String("version", "1.0.0"),
		zap.Int("port", 8080),
	)

	// 演示格式化日志
	logger.Infof("处理了 100 个请求")
	logger.Infof("处理了 %d 个请求", 100)
	logger.Debug("调试信息")
	logger.Debugf("调试信息: %v", map[string]int{"count": 5})
	logger.Warnf("警告: CPU 使用率 %d%%", 85)
	logger.Errorf("错误: %s", "连接超时")
	// test5()
	// test4()
	// test3()
	// test2()
	// test1()
}
