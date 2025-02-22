package application

import "github.com/leetcode-golang-classroom/golang-sample-with-rate-limiter/internal/service/greeting"

func (app *App) setupGreetingRoute() {
	greetingRoute := greeting.NewRoute()
	greetingRoute.SetupRoute(app.router)
}
