package api

import (
	"github.com/gin-gonic/gin"
)

type AppApi struct{}

func (a AppApi) Start() {
	r := gin.Default()
	r.Use(TraceIDMiddleware())

	// 注册用户相关路由
	userApi := &UserApi{}
	userApi.RegisterRoutes(r)

	r.Run(":8000")
}
