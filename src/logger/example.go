package logger

import "go.uber.org/zap"

// 使用示例（仅供参考，不会被编译到主程序）
func ExampleUsage() {
	// 示例 1: 使用默认配置（开发环境，控制台输出）
	// 默认配置已在 init() 中初始化，无需手动调用

	// 示例 2: 自定义开发环境配置
	devConfig := &Config{
		Level:    "debug",
		Env:      "dev",
		FilePath: "", // 只输出到控制台
	}
	Init(devConfig)

	// 示例 3: 生产环境配置（JSON 格式 + 文件轮转）
	prodConfig := &Config{
		Level:      "info",
		Env:        "prod",
		FilePath:   "logs/app.log",
		MaxSize:    100,  // 100MB
		MaxBackups: 10,   // 保留 10 个备份
		MaxAge:     30,   // 保留 30 天
		Compress:   true, // 压缩旧日志
	}
	Init(prodConfig)

	// 示例 4: 结构化日志
	Info("用户登录",
		zap.String("username", "alice"),
		zap.Int("user_id", 123),
		zap.String("ip", "192.168.1.1"),
	)

	// 示例 5: 格式化日志（更方便但性能略低）
	Infof("用户 %s (ID: %d) 从 %s 登录", "alice", 123, "192.168.1.1")

	// 示例 6: 不同日志级别
	Debug("调试信息")
	Info("普通信息")
	Warn("警告信息")
	Error("错误信息")

	// 示例 7: 格式化日志
	Debugf("这是调试信息: %v", map[string]int{"count": 5})
	Infof("处理了 %d 个请求", 100)
	Warnf("磁盘使用率达到 %d%%", 85)
	Errorf("连接数据库失败: %s", "connection timeout")

	// 示例 8: 带预设字段的 logger
	userLogger := With(
		zap.String("service", "user-service"),
		zap.String("version", "1.0.0"),
	)
	userLogger.Info("服务启动")

	// 示例 9: 获取原始 logger
	rawLogger := GetLogger()
	rawLogger.Info("使用原始 logger")

	sugarLogger := GetSugar()
	sugarLogger.Infow("使用 sugar logger",
		"key1", "value1",
		"key2", 123,
	)

	// 程序退出前刷新日志
	defer Sync()
}
