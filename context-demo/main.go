package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func httpDo(ctx context.Context) error {
	// Run the HTTP request in a goroutine and pass the response to f.
	c := make(chan error, 2)
	go func() {
		fmt.Println("aaa")
		time.Sleep(1 * time.Second)
		c <- errors.New("some error")
	}()
	select {
	case <-ctx.Done():
		<-c // Wait for f to return.
		return ctx.Err()
	case err := <-c:
		return err
	}
}

func main() {
	ctx := context.Background()
	// fmt.Println(fmt.Sprintf("%p", ctx))
	ctx1, cancel := context.WithCancel(ctx)

	ctx2 := context.WithValue(ctx1, "foo", "bar")
	ctx3 := context.WithValue(ctx2, "bar", "bar")
	ctx4, cancel := context.WithCancel(ctx3)

	cancelCtx, cancel := context.WithCancel(ctx4)
	defer cancel()
	// fmt.Println(cancelCtx)
	fmt.Println(httpDo(cancelCtx))

	// ctx11 := context.Background()
	// fmt.Println(fmt.Sprintf("%p", ctx11))

	// ctx1 := context.WithValue(ctx, "foo", "bar1")
	// fmt.Println(fmt.Sprintf("%p", ctx1))
	// fmt.Println(ctx1, ctx1.Value("foo"))

	// ctx2 := context.WithValue(ctx1, "bar", "bar2")
	// fmt.Println(fmt.Sprintf("%p", ctx2))
	// fmt.Println(ctx2, ctx2.Value("foo"))

	// ctx3 := context.WithValue(ctx2, "bar", "bar3")
	// fmt.Println(fmt.Sprintf("%p", ctx3))
	// fmt.Println(ctx3, ctx3.Value("foo"))

	// time.Sleep(100 * time.Millisecond)
	// ctx, cancel := context.WithCancel(context.Background())

	// go func(ctx context.Context) {
	// 	time.Sleep(1 * time.Second)
	// 	select {
	// 	case <-ctx.Done():
	// 		fmt.Println("goroutine 1 done")
	// 		return
	// 	default:
	// 		fmt.Println("goroutine 1")
	// 	}
	// }(ctx)

	// go func(ctx context.Context) {
	// 	time.Sleep(3 * time.Second)
	// 	select {
	// 	case <-ctx.Done():
	// 		fmt.Println("goroutine 2 done")
	// 		return
	// 	default:
	// 		fmt.Println("goroutine 2")
	// 	}
	// }(ctx)

	// time.Sleep(2 * time.Second)
	// cancel()
	// fmt.Println("call cancal")
	// // fmt.Println("ctx error", ctx.Err())
	// // t, ok := ctx.Deadline()
	// // fmt.Println("ctx error", t, ok)
	// // fmt.Println("ctx error", ctx.Value("aa"))

	// time.Sleep(2 * time.Second)
	// fmt.Println(&ctx)
}
