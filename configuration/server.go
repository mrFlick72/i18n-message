package configuration

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github/mrflick72/i18n-message/internal/heath"
	"sync"
)

func newWebServer() *iris.Application {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	return app
}

func NewActuatorServer(wg *sync.WaitGroup) {
	app := newWebServer()
	endpoints := heath.HealthEndpoint{}
	endpoints.ResgisterEndpoints(app)
	app.Listen(":8081")
	wg.Done()
}

func NewApplicationServer(wg *sync.WaitGroup) {
	app := newWebServer()
	messageRepository := ConfigureMessageRepository()
	ConfigureMessageEndpoints(messageRepository, app)

	// Listen and serve on 0.0.0.0:8080
	app.Listen(":8080")
	wg.Done()
}
