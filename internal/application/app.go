package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/leetcode-golang-classroom/golang-sample-with-rate-limiter/internal/config"
	rateLimiter "github.com/leetcode-golang-classroom/golang-sample-with-rate-limiter/internal/service/rate_limiter"
	"golang.org/x/time/rate"
)

type App struct {
	appConfig *config.Config
	router    *http.ServeMux
}

// New - 建立 App 物件
func New(appConfig *config.Config) *App {
	router := http.NewServeMux()
	// create app instance
	app := &App{
		appConfig: appConfig,
		router:    router,
	}
	// setup route
	app.setupGreetingRoute()
	return app
}

// Start - 啟動 server
func (app *App) Start(ctx context.Context) error {
	// setup rateLimiter
	rateLimitMux := rateLimiter.NewRateLimiter()
	handler := rateLimitMux.RateLimiterMiddleware(app.router, rate.Limit(2), 10)
	// setup server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.appConfig.Port),
		Handler: handler,
	}
	log.Printf("starting server on %s", app.appConfig.Port)
	errCh := make(chan error, 1)
	var err error
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			errCh <- fmt.Errorf("failed to start server: %w", err)
		}
		CloseChannel(errCh)
	}()
	select {
	case err = <-errCh:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		log.Printf("stopping server, wait for 10 seconds to stop")
		defer cancel()
		return server.Shutdown(timeout)
	}
}

func CloseChannel(ch chan error) {
	if _, ok := <-ch; ok {
		close(ch)
	}
}
