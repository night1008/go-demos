package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.PUT("/put", func(c *gin.Context) {
		fmt.Println("===> put request start")
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
		defer cancel()
	Loop:
		for {
			select {
			case <-ctx.Done():
				break Loop
			case <-time.After(10 * time.Second):
				fmt.Println("elapse 10 seconds")
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
		fmt.Println("===> put request end")
	})

	r.GET("/a", func(c *gin.Context) {
		c.File("index.html")
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
