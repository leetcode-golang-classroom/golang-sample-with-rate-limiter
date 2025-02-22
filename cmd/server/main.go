package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/leetcode-golang-classroom/golang-sample-with-rate-limiter/internal/application"
	"github.com/leetcode-golang-classroom/golang-sample-with-rate-limiter/internal/config"
)

func main() {
	// 建立 application instance
	app := application.New(config.AppConfig)
	// 設定中斷訊號監聽
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt,
		syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	err := app.Start(ctx)
	if err != nil {
		log.Println("failed to start app:", err)
	}
}
