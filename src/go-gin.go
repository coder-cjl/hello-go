package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GoGin struct{}

func go_gin_ts1() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, Gin!")
	})

	r.GET("/greet", greetHandler)
	r.GET("/action/:action", actionHandler)
	v2 := r.Group("/v2")
	{
		v2.GET("/user", userHandler)
	}

	r.Run(":8000")
}

func userHandler(c *gin.Context) {
	name := c.DefaultQuery("name", "枯藤")
	c.JSON(http.StatusOK, gin.H{
		"name": name,
	})
}

func greetHandler(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		name = "Guest"
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, " + name + "!",
	})
}

func actionHandler(c *gin.Context) {
	action := c.Param("action")
	c.JSON(http.StatusOK, gin.H{
		"action": action,
		"status": "completed",
	})
}

func (g GoGin) Test() {
	go_gin_ts1()
}
