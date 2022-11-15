package application

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github/mrflick72/i18n-message/src/internal/heath"
	"github/mrflick72/i18n-message/src/internal/web"
	"github/mrflick72/i18n-message/src/middleware/security"
	"strings"
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
	app.Listen(manager.GetConfigFor("ACTUATOR_PORT"))
	wg.Done()
}

func NewApplicationServer(wg *sync.WaitGroup) {
	app := newWebServer()
	security.SetUpOAuth2(app, security.Jwk{
		Url:    manager.GetConfigFor("security.jwk-uri"),
		Client: web.New(),
	}, strings.Split(manager.GetConfigFor("security.allowed-authority"), ","))
	messageRepository := ConfigureMessageRepository()
	ConfigureMessageEndpoints(messageRepository, app)

	// Listen and serve on 0.0.0.0:8080
	app.Listen(manager.GetConfigFor("PRODUCTION_PORT"))
	wg.Done()
}
