package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type Person1 struct {
	Name string `form:"name" json:"name" binding:"required"`
	Age  int    `form:"age" json:"age" binding:"required,gt=10"`
}

type GoGin1 struct{}

func loginHandler(c *gin.Context) {
	var login Login

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if login.User == "admin" && login.Password == "password" {
		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
	}
}

func loginHandler2(c *gin.Context) {
	var login Login

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := c.Cookie("token")
	if err != nil {
		cookie = "NotSet"
		c.SetCookie("token", "123456", 3600, "/", "localhost", false, true)
	}

	Log.Info("Cookie value:", cookie)

	if login.User == "admin" && login.Password == "password" {
		c.JSON(http.StatusOK, login)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
	}
}

func loginAsync(c *gin.Context) {
	var login Login

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func() {
		time.Sleep(5 * time.Second)
		Log.Info("User attempted to login:", login.User)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Login attempt received"})
}

func loginSync(c *gin.Context) {
	var login Login

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	time.Sleep(5 * time.Second)
	Log.Info("User attempted to login:", login.User)
	c.JSON(http.StatusOK, gin.H{"message": "Login attempt received"})
}

func middleWare1() gin.HandlerFunc {
	return func(c *gin.Context) {
		Log.Info("请求开始")
		c.Next()
		status := c.Writer.Status()
		Log.Info("请求状态码:", status)
		Log.Info("请求结束")
	}
}

func person1Handler(c *gin.Context) {
	var person Person1

	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Person data received", "person": person})
}

func (g GoGin1) Test() {
	r := gin.Default()
	r.Use(middleWare1())
	r.POST("/login", loginHandler)
	r.POST("/login2", loginHandler2)
	r.POST("/login_async", loginAsync)
	r.POST("/login_sync", loginSync)
	r.POST("/person1", person1Handler)
	r.Run(":8000")
}
