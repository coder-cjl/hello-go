package api

import (
	"hello-go/src/jwt"
	"hello-go/src/logger"
	"hello-go/src/redis"
	"log"
	"net/http"
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	snowflakeNode *snowflake.Node
	snowflakeOnce sync.Once
)

// 初始化 snowflake node（全局单例）
func getSnowflakeNode() *snowflake.Node {
	snowflakeOnce.Do(func() {
		var err error
		snowflakeNode, err = snowflake.NewNode(1)
		if err != nil {
			log.Fatalf("failed to create snowflake node: %v", err)
		}
	})
	return snowflakeNode
}

// TraceIDMiddleware 生成并设置 TraceID
func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Set("TraceID", traceID)
		c.Header("X-Trace-ID", traceID)
		c.Next()
	}
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type UserApi struct{}

func (u *UserApi) RegisterRoutes(r *gin.Engine) {
	userGroup := r.Group("user")
	userGroup.POST("/login", loginHandle)
}

// 登录
func loginHandle(c *gin.Context) {
	traceID := c.GetString("TraceID")

	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			TraceID: traceID,
		})
		logger.Warnf("请求参数错误", "traceID", traceID)
		return
	}

	if req.Username == "admin" && req.Password == "password" {
		node := getSnowflakeNode()
		userId := node.Generate().Int64()

		logger.Infof("[%s] 用户登录成功，生成用户ID: %d", traceID, userId)

		token, err := jwt.GenerateToken(uint(userId))
		if err != nil {
			logger.Errorf("[%s] 生成token失败: %v", traceID, err)
			c.JSON(http.StatusInternalServerError, ApiResponse{
				Code:    http.StatusInternalServerError,
				Message: "服务器内部错误",
				TraceID: traceID,
			})
			return
		}

		c.JSON(http.StatusOK, ApiResponse{
			Code:    http.StatusOK,
			Message: "登录成功",
			Data: UserLoginResponse{
				AccessToken: token,
			},
			TraceID: traceID,
		})

		// 将token 存入到redis
		err = redis.SetValue("userId", token)
		if err != nil {
			logger.Errorf("[%s] Redis存储token失败: %v", traceID, err)
		}

		logger.Infof("[%s] 用户登录成功: %s", traceID, req.Username)
	} else {
		c.JSON(http.StatusUnauthorized, ApiResponse{
			Code:    http.StatusUnauthorized,
			Message: "用户名或密码错误",
			TraceID: traceID,
		})
		logger.Warnf("[%s] 用户登录失败: %s", traceID, req.Username)
	}
}
