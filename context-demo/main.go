package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			fmt.Println("goroutine 1 done")
			return
		default:
			fmt.Println("goroutine 1")
		}
	}(ctx)

	go func(ctx context.Context) {
		time.Sleep(3 * time.Second)
		select {
		case <-ctx.Done():
			fmt.Println("goroutine 2 done")
			return
		default:
			fmt.Println("goroutine 2")
		}
	}(ctx)

	time.Sleep(2 * time.Second)
	cancel()
	fmt.Println("call cancal")
	// fmt.Println("ctx error", ctx.Err())
	// t, ok := ctx.Deadline()
	// fmt.Println("ctx error", t, ok)
	// fmt.Println("ctx error", ctx.Value("aa"))

	time.Sleep(2 * time.Second)
	fmt.Println(&ctx)

	ctx2 := context.WithValue(ctx, 111, "a")
	fmt.Println(&ctx2)

	ctx3, _ := context.WithDeadline(ctx, time.Now())
	fmt.Println(ctx3.Value(&ctx))
}
