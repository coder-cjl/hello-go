package api

import (
	"hello-go/src/jwt"
	"hello-go/src/logger"
	"hello-go/src/mysql"
	"hello-go/src/redis"
	"log"
	"net/http"
	"strconv"
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

// 生成并设置 TraceID
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

type UserVo struct {
	ID       int64  `gorm:"column:id"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password_hash"`
}

func (UserVo) TableName() string {
	return "users"
}

type UserLoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type UserInfoVo struct {
	ID       int64  `gorm:"column:user_id" json:"user_id"`
	Nickname string `gorm:"column:nickname" json:"nickname"`
	Email    string `gorm:"column:email" json:"email"`
	Gender   string `gorm:"column:gender" json:"gender"`
}

func (UserInfoVo) TableName() string {
	return "user_info"
}

type UserApi struct{}

func (u *UserApi) RegisterRoutes(r *gin.Engine) {
	userGroup := r.Group("user")
	userGroup.POST("/login", loginHandle)
	userGroup.POST("/update-user-info", updateUserHandle)
}

// 生成 token 并缓存到 Redis
func generateAndCacheToken(user *UserVo, traceID string) (string, error) {
	token, err := jwt.GenerateToken(user.ID, user.Username, []string{"user"})
	if err != nil {
		logger.Errorf("[%s] 生成token失败: %v", traceID, err)
		return "", err
	}

	if err := redis.SetValue(strconv.FormatInt(user.ID, 10), token); err != nil {
		logger.Warnf("[%s] Redis存储token失败: %v", traceID, err)
		return "", err
	}

	return token, nil
}

// 从 Redis 获取有效的 token
func getCachedToken(userID int64, traceID string) (string, bool) {
	token, err := redis.GetValue(strconv.FormatInt(userID, 10))
	if err != nil {
		logger.Warnf("[%s] Redis获取token失败: %v", traceID, err)
		return "", false
	}

	if token == "" {
		return "", false
	}

	// 检查 token 是否过期
	expired, err := jwt.IsTokenExpired(token)
	if err != nil || expired {
		return "", false
	}

	return token, true
}

// 返回登录成功响应
func respondSuccess(c *gin.Context, token, traceID, username string) {
	c.JSON(http.StatusOK, ApiResponse{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    UserLoginResponse{AccessToken: token},
		TraceID: traceID,
	})
	logger.Infof("[%s] 用户登录成功: %s", traceID, username)
}

// 返回错误响应
func respondError(c *gin.Context, code int, message, traceID string) {
	c.JSON(code, ApiResponse{
		Code:    code,
		Message: message,
		TraceID: traceID,
	})
}

// 登录
func loginHandle(c *gin.Context) {
	traceID := c.GetString("TraceID")

	var reqParams UserRequest
	if err := c.ShouldBindJSON(&reqParams); err != nil {
		respondError(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), traceID)
		return
	}

	// 查询用户
	var user UserVo
	result := mysql.Db.Where("username = ? AND password_hash = ?", reqParams.Username, reqParams.Password).First(&user)

	if result.Error != nil {
		if result.RowsAffected == 0 {
			// 用户不存在，创建新用户
			user = UserVo{
				ID:       getSnowflakeNode().Generate().Int64(),
				Username: reqParams.Username,
				Password: reqParams.Password,
			}
			if err := mysql.Db.Create(&user).Error; err != nil {
				logger.Errorf("[%s] 创建用户失败: %v", traceID, err)
				respondError(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), traceID)
				return
			}
			logger.Infof("[%s] 创建新用户成功: %s", traceID, reqParams.Username)
		} else {
			logger.Errorf("[%s] 查询用户失败: %v", traceID, result.Error)
			respondError(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), traceID)
			return
		}
	} else {
		// 老用户，尝试从缓存获取 token
		if token, ok := getCachedToken(user.ID, traceID); ok {
			respondSuccess(c, token, traceID, reqParams.Username)
			return
		}
	}

	// 生成新 token
	token, err := generateAndCacheToken(&user, traceID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), traceID)
		return
	}

	respondSuccess(c, token, traceID, reqParams.Username)
}

// 更新用户信息
func updateUserHandle(c *gin.Context) {
	// 验证token是否合法有效
	traceID := c.GetString("TraceID")
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		respondError(c, http.StatusUnauthorized, "Missing Authorization header", traceID)
		return
	}
	logger.Debugf("[%s] Authorization header: %s", traceID, authHeader)

	tokenString := authHeader[len("Bearer "):]
	claims, err := jwt.GetClaimsFromToken(tokenString)
	if err != nil {
		respondError(c, http.StatusUnauthorized, "Invalid token", traceID)
		return
	}

	userID := claims.UserId

	// 绑定请求参数
	var userInfo UserInfoVo
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		respondError(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), traceID)
		return
	}

	// 更新用户信息
	userInfo.ID = userID
	if err := mysql.Db.Model(&UserInfoVo{}).Where("user_id = ?", userID).Updates(userInfo).Error; err != nil {
		logger.Errorf("[%s] 更新用户信息失败: %v", traceID, err)
		respondError(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), traceID)
		return
	}

	c.JSON(http.StatusOK, ApiResponse{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		TraceID: traceID,
	})
	logger.Infof("[%s] 用户信息更新成功: user_id=%d", traceID, userID)
}
