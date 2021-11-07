package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	srv := &http.Server{
		Addr:    ":8082",
		Handler: r,
	}

	eg, ctx := errgroup.WithContext(context.Background())
	log.Println(ctx)

	eg.Go(func() error {
		log.Println("ListenAndServe")
		return srv.ListenAndServe()
	})

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			log.Println("errgroup exit...")
		}

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		time.Sleep(3 * time.Second)
		log.Println("shutting down server...")
		return srv.Shutdown(shutdownCtx)
	})

	eg.Go(func() error {
		quit := make(chan os.Signal, 0)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Reset()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-quit:
			return fmt.Errorf("get os signal: %v", sig)
		}
	})

	if err := eg.Wait(); err != nil {
		log.Println(err)
	}
}
