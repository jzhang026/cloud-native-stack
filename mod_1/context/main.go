package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	baseCtx := context.Background()
	ctx := context.WithValue(baseCtx, "a", "b")
	go func(ctx context.Context) {
		fmt.Println(ctx.Value("a"))
	}(ctx)
	timeoutCtx, cancel := context.WithTimeout(baseCtx, time.Second)
	defer cancel()

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Second)
		for _ = range ticker.C {
			select {
			case <-ctx.Done():
				println("child process interrupted ...")
			default:
				println("select case `default` branch")
			}
		}
	}(timeoutCtx)

	time.Sleep(1 * time.Second)

	select {
	case <-timeoutCtx.Done():
		println("main process exit")
		time.Sleep(time.Second)
	}
}
