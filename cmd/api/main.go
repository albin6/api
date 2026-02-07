package main

import (
	"github.com/gin-gonic/gin"
	"github.com/albin6/api/pkg/database"
)

func main() {
	database.NewPostgresDB()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	r.Run(":8080")
}