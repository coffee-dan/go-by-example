package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	app.LoadHTMLFiles("./views/index.tmpl")
	app.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"favicon": "😎",
		})
	})
	app.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	app.Run(":8080")
}
