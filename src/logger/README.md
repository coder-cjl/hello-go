# Logger 使用文档

基于 zap 的生产级 Go 日志库，提供结构化日志、文件轮转、环境配置等功能。

## 功能特性

✅ **环境配置切换** - 开发环境（彩色控制台）/ 生产环境（JSON 格式）  
✅ **文件输出与轮转** - 自动切割日志文件，支持压缩  
✅ **多种日志方法** - 结构化日志 + 格式化日志  
✅ **可配置的日志级别** - debug, info, warn, error, fatal  
✅ **多输出目标** - 同时输出到控制台和文件  
✅ **调用位置追踪** - 自动记录调用文件和行号  

## 快速开始

### 1. 使用默认配置

```go
package main

import "hello-go/src/logger"

func main() {
    defer logger.Sync() // 程序退出前刷新日志缓冲区

    logger.Info("应用启动")
    logger.Debugf("处理了 %d 个请求", 100)
}
```

默认配置：
- 日志级别：info
- 环境：dev（彩色控制台输出）
- 输出：仅控制台

### 2. 生产环境配置

```go
package main

import "hello-go/src/logger"

func main() {
    defer logger.Sync()

    // 初始化生产环境配置
    config := &logger.Config{
        Level:      "info",
        Env:        "prod",        // 生产环境（JSON 格式）
        FilePath:   "logs/app.log", // 输出到文件
        MaxSize:    100,            // 单文件最大 100MB
        MaxBackups: 10,             // 保留 10 个备份
        MaxAge:     30,             // 保留 30 天
        Compress:   true,           // 压缩旧日志
    }
    logger.Init(config)

    logger.Info("应用启动")
}
```

### 3. 开发环境配置

```go
config := &logger.Config{
    Level:    "debug",
    Env:      "dev",    // 开发环境（彩色控制台）
    FilePath: "",       // 不输出到文件
}
logger.Init(config)
```

## 使用示例

### 结构化日志（推荐）

性能更好，适合生产环境：

```go
import "go.uber.org/zap"

// 基础用法
logger.Info("用户登录成功")
logger.Error("连接数据库失败")

// 带字段
logger.Info("用户登录",
    zap.String("username", "alice"),
    zap.Int("user_id", 123),
    zap.String("ip", "192.168.1.1"),
)

logger.Error("API 调用失败",
    zap.String("endpoint", "/api/users"),
    zap.Int("status_code", 500),
    zap.Duration("latency", time.Millisecond*150),
)
```

### 格式化日志

更方便，适合开发环境：

```go
// 基础用法
logger.Infof("处理了 %d 个请求", 100)
logger.Debugf("配置: %+v", config)
logger.Warnf("磁盘使用率达到 %d%%", 85)
logger.Errorf("连接失败: %s", err.Error())
```

### 所有日志级别

```go
logger.Debug("调试信息")         // 开发环境可见
logger.Info("普通信息")          // 常规日志
logger.Warn("警告信息")          // 需要注意的问题
logger.Error("错误信息")         // 错误但不影响运行
logger.Fatal("致命错误")         // 记录后退出程序

// 格式化版本
logger.Debugf("debug: %v", data)
logger.Infof("info: %s", msg)
logger.Warnf("warn: %d", count)
logger.Errorf("error: %v", err)
logger.Fatalf("fatal: %s", reason)
```

### 带预设字段的 Logger

```go
import "go.uber.org/zap"

// 创建带服务信息的 logger
serviceLogger := logger.With(
    zap.String("service", "user-service"),
    zap.String("version", "1.0.0"),
    zap.String("env", "prod"),
)

// 所有日志都会包含预设字段
serviceLogger.Info("服务启动")
// 输出: {"level":"info","service":"user-service","version":"1.0.0","env":"prod","msg":"服务启动"}

// 为特定用户创建 logger
userLogger := logger.With(
    zap.Int("user_id", 123),
    zap.String("username", "alice"),
)

userLogger.Info("执行操作")
userLogger.Error("操作失败")
```

### 获取原始 Logger

```go
// 获取 *zap.Logger
rawLogger := logger.GetLogger()
rawLogger.Info("使用原始 logger")

// 获取 *zap.SugaredLogger
sugar := logger.GetSugar()
sugar.Infow("使用 sugar logger",
    "key1", "value1",
    "key2", 123,
)
```

## 配置说明

```go
type Config struct {
    Level      string // 日志级别: "debug", "info", "warn", "error", "fatal"
    Env        string // 环境: "dev" (开发), "prod" (生产)
    FilePath   string // 日志文件路径，为空则只输出到控制台
    MaxSize    int    // 单个日志文件最大 MB，默认 100
    MaxBackups int    // 保留的旧日志文件数量，默认 10
    MaxAge     int    // 保留的旧日志文件最大天数，默认 30
    Compress   bool   // 是否压缩旧日志，默认 false
}
```

## 输出格式

### 开发环境 (Env="dev")

彩色控制台输出，易读：
```
2026-02-06T17:27:06.686+0800    INFO    src/main.go:93  应用启动        {"version": "1.0.0", "port": 8080}
2026-02-06T17:27:06.687+0800    WARN    src/main.go:101 警告: CPU 使用率 85%
```

### 生产环境 (Env="prod")

JSON 格式，便于日志收集和分析：
```json
{"level":"info","ts":1770369704.7217789,"caller":"src/main.go:78","msg":"应用启动","version":"1.0.0","port":8080}
{"level":"warn","ts":1770369704.721877,"caller":"src/main.go:79","msg":"警告: CPU 使用率 85%"}
```

## 日志轮转

当日志文件达到 `MaxSize` 时，自动创建新文件：
```
logs/app.log           # 当前日志
logs/app-2026-02-06.log.gz  # 压缩的旧日志
logs/app-2026-02-05.log.gz
```

## 最佳实践

1. **程序启动时初始化配置**
```go
func init() {
    config := loadConfigFromEnv() // 从环境变量或配置文件加载
    logger.Init(config)
}
```

2. **始终使用 defer Sync()**
```go
func main() {
    defer logger.Sync()
    // ...
}
```

3. **生产环境使用结构化日志**
```go
// ✅ 推荐
logger.Info("请求处理完成",
    zap.String("method", "POST"),
    zap.Int("status", 200),
    zap.Duration("latency", latency),
)

// ❌ 不推荐（性能较差）
logger.Infof("请求处理完成 method=%s status=%d latency=%v", "POST", 200, latency)
```

4. **错误日志包含上下文**
```go
logger.Error("数据库查询失败",
    zap.Error(err),
    zap.String("query", sql),
    zap.String("table", "users"),
)
```

5. **避免在循环中频繁记录日志**
```go
// ❌ 不推荐
for _, item := range items {
    logger.Debug("处理项目", zap.Any("item", item))
}

// ✅ 推荐
logger.Info("开始批量处理", zap.Int("count", len(items)))
// 处理...
logger.Info("批量处理完成", zap.Int("success", success), zap.Int("failed", failed))
```

## 常见问题

**Q: 如何在不同包中使用 logger？**
```go
import "hello-go/src/logger"

func someFunction() {
    logger.Info("在其他包中使用")
}
```

**Q: 如何动态修改日志级别？**
```go
// 重新初始化
newConfig := &logger.Config{
    Level: "debug",
    Env:   "dev",
}
logger.Init(newConfig)
```

**Q: 日志文件保存在哪里？**  
A: 由 `FilePath` 配置项决定，可以使用相对路径（相对于程序运行目录）或绝对路径。

**Q: 如何关闭某些日志输出？**  
A: 提高日志级别。例如设置 `Level: "warn"` 将不会输出 debug 和 info 级别的日志。
