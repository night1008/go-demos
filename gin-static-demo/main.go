package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// if Allow DirectoryIndex
	//r.Use(static.Serve("/", static.LocalFile("/tmp", true)))
	// set prefix
	//r.Use(static.Serve("/static", static.LocalFile("/tmp", true)))

	r.StaticFS("/static", http.Dir("web"))

	// TODO 为什么这种情况，前端静态资源不能启用缓存呢
	// r.StaticFS("/static", AssetFile())

	// r.Use(static.Serve("/", static.LocalFile("web", true)))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"pong": "pong"})
	})

	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
