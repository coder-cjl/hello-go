# Elasticsearch 日志接入方案

## 方案一：Filebeat（推荐）

最常用、最稳定的方案，通过 Filebeat 读取日志文件并发送到 Elasticsearch。

### 1. 配置生产环境日志

```go
// main.go
func init() {
    logger.Init(&logger.Config{
        Level:      "info",
        Env:        "prod",              // JSON 格式
        Topic:      "luca-app",
        FilePath:   "/var/log/app/luca.log",  // 日志文件路径
        MaxSize:    100,
        MaxBackups: 10,
        MaxAge:     30,
        Compress:   false,  // Filebeat 处理时建议不压缩
    })
}
```

### 2. 安装 Filebeat

```bash
# macOS
brew install filebeat

# Linux (Debian/Ubuntu)
sudo apt-get install filebeat

# Docker
docker pull docker.elastic.co/beats/filebeat:8.11.0
```

### 3. 配置 Filebeat

创建 `filebeat.yml`：

```yaml
# 输入配置
filebeat.inputs:
  - type: log
    enabled: true
    paths:
      - /var/log/app/luca.log
    
    # JSON 日志解析
    json.keys_under_root: true
    json.add_error_key: true
    json.message_key: msg
    
    # 添加额外字段
    fields:
      app: luca
      env: production
    fields_under_root: true

# Elasticsearch 输出
output.elasticsearch:
  hosts: ["localhost:9200"]
  index: "luca-logs-%{+yyyy.MM.dd}"
  
  # 如果使用认证
  # username: "elastic"
  # password: "changeme"

# 索引模板
setup.template.name: "luca-logs"
setup.template.pattern: "luca-logs-*"
setup.ilm.enabled: false

# Kibana 配置（可选）
setup.kibana:
  host: "localhost:5601"
```

### 4. 启动 Filebeat

```bash
# 测试配置
filebeat test config -c filebeat.yml

# 启动
sudo filebeat -e -c filebeat.yml
```

---

## 方案二：Fluentd

更灵活的日志处理方案，支持复杂的过滤和转换。

### 1. 安装 Fluentd

```bash
# macOS
brew install fluentd

# Docker
docker pull fluent/fluentd:latest
```

### 2. 配置 Fluentd

创建 `fluent.conf`：

```xml
<source>
  @type tail
  path /var/log/app/luca.log
  pos_file /var/log/fluent/luca.log.pos
  tag luca.app
  
  <parse>
    @type json
    time_key ts
    time_format %s
  </parse>
</source>

<filter luca.app>
  @type record_transformer
  <record>
    hostname ${hostname}
    env production
  </record>
</filter>

<match luca.app>
  @type elasticsearch
  host localhost
  port 9200
  index_name luca-logs
  type_name _doc
  
  <buffer>
    flush_interval 10s
  </buffer>
</match>
```

---

## 方案三：直接写入 ES（不推荐生产环境）

适合测试环境，不推荐生产使用（性能和可靠性问题）。

### 1. 安装 ES Writer

```bash
go get github.com/olivere/elastic/v7
```

### 2. 创建 ES Writer

```go
// src/logger/es_writer.go
package logger

import (
    "context"
    "github.com/olivere/elastic/v7"
    "go.uber.org/zap/zapcore"
)

type ElasticsearchWriter struct {
    client *elastic.Client
    index  string
}

func NewElasticsearchWriter(url, index string) (*ElasticsearchWriter, error) {
    client, err := elastic.NewClient(
        elastic.SetURL(url),
        elastic.SetSniff(false),
    )
    if err != nil {
        return nil, err
    }
    
    return &ElasticsearchWriter{
        client: client,
        index:  index,
    }, nil
}

func (w *ElasticsearchWriter) Write(p []byte) (n int, err error) {
    _, err = w.client.Index().
        Index(w.index).
        BodyString(string(p)).
        Do(context.Background())
    
    if err != nil {
        return 0, err
    }
    return len(p), nil
}
```

### 3. 使用 ES Writer

```go
// 在 Init 函数中添加 ES core
esWriter, err := NewElasticsearchWriter("http://localhost:9200", "luca-logs")
if err != nil {
    return err
}

esCore := zapcore.NewCore(
    zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
    zapcore.AddSync(esWriter),
    level,
)
cores = append(cores, esCore)
```

---

## Elasticsearch 索引映射优化

为了更好地在 ES 中搜索和分析，建议定义索引模板：

```bash
# 创建索引模板
curl -X PUT "localhost:9200/_index_template/luca-logs-template" \
-H 'Content-Type: application/json' -d'
{
  "index_patterns": ["luca-logs-*"],
  "template": {
    "settings": {
      "number_of_shards": 1,
      "number_of_replicas": 1,
      "index.lifecycle.name": "luca-logs-policy"
    },
    "mappings": {
      "properties": {
        "@timestamp": { "type": "date" },
        "ts": { "type": "date" },
        "level": { "type": "keyword" },
        "topic": { "type": "keyword" },
        "caller": { "type": "text" },
        "msg": { 
          "type": "text",
          "fields": {
            "keyword": { "type": "keyword" }
          }
        },
        "error": { "type": "text" },
        "version": { "type": "keyword" },
        "user_id": { "type": "keyword" },
        "port": { "type": "integer" }
      }
    }
  }
}'
```

---

## 优化建议

### 1. 添加更多上下文字段

```go
// 在 Init 中添加全局字段
log = zap.New(core,
    zap.AddCaller(),
    zap.AddCallerSkip(1),
    zap.AddStacktrace(zapcore.ErrorLevel),
    zap.Fields(
        zap.String("host", hostname),
        zap.String("env", cfg.Env),
    ),
).Named(cfg.Topic)
```

### 2. 使用时间戳字段名 `@timestamp`

便于 ES 和 Kibana 自动识别：

```go
encoderConfig.TimeKey = "@timestamp"
```

### 3. 添加请求 ID 追踪

```go
// 为每个请求生成唯一 ID
requestLogger := logger.With(
    zap.String("request_id", uuid.New().String()),
)
```

---

## Kibana 查询示例

### 1. 按级别过滤

```
level: "error" AND topic: "LUCA"
```

### 2. 按时间范围

```
@timestamp: [now-1h TO now] AND level: ("error" OR "warn")
```

### 3. 全文搜索

```
msg: "MySQL" AND level: "error"
```

### 4. 聚合分析

在 Kibana Visualize 中创建：
- 按 `topic` 分组的日志量
- 按 `level` 的错误分布
- 响应时间趋势（如果有 duration 字段）

---

## Docker Compose 完整示例

```yaml
version: '3.8'

services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.11.0
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
    ports:
      - "9200:9200"
    volumes:
      - es-data:/usr/share/elasticsearch/data

  kibana:
    image: docker.elastic.co/kibana/kibana:8.11.0
    ports:
      - "5601:5601"
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
    depends_on:
      - elasticsearch

  filebeat:
    image: docker.elastic.co/beats/filebeat:8.11.0
    user: root
    volumes:
      - ./filebeat.yml:/usr/share/filebeat/filebeat.yml:ro
      - /var/log/app:/var/log/app:ro
      - filebeat-data:/usr/share/filebeat/data
    depends_on:
      - elasticsearch

  luca-app:
    build: .
    volumes:
      - /var/log/app:/var/log/app
    depends_on:
      - elasticsearch

volumes:
  es-data:
  filebeat-data:
```

---

## 总结

✅ **推荐方案：Filebeat + Elasticsearch + Kibana**

**优点：**
- 稳定可靠
- 低耦合（应用不需要知道 ES）
- Filebeat 自动处理重试、缓冲
- 支持多种输出（ES, Kafka, Logstash 等）

**缺点：**
- 需要额外部署 Filebeat

**部署步骤：**
1. 配置生产环境 logger（JSON 格式）
2. 确保日志写入文件
3. 安装配置 Filebeat
4. Filebeat 自动发送日志到 ES
5. 在 Kibana 中查看和分析

你的当前日志配置已经完全兼容这套方案！🎉
