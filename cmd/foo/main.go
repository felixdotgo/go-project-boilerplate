package main

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	engine := gin.New()
	engine.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	engine.Use(ginzap.RecoveryWithZap(logger, true))
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	err = engine.Run(":8000")
	if err != nil {
		panic(err)
	}
}
