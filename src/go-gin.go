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

	r.Run(":8000")
}

func (g GoGin) Test() {
	go_gin_ts1()
}