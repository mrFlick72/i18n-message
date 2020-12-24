package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github/mrflick72/i18n-message/api"
	"github/mrflick72/i18n-message/internal/heath"
	"github/mrflick72/i18n-message/internal/message/repository"
	"github/mrflick72/i18n-message/internal/web"
	"sync"
)

func main() {
	// Creates an iris application without any middleware by default
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go applicationServer(wg)
	go actuatorServer(wg)

	wg.Wait()
}

func actuatorServer(wg *sync.WaitGroup) {
	app := newWebServer()
	endpoints := heath.HealthEndpoint{}
	endpoints.ResgisterEndpoints(app)
	app.Listen(":8081")
	wg.Done()
}

func applicationServer(wg *sync.WaitGroup) {
	app := newWebServer()
	messageRepository := configureMessageRepository()
	configureMessageEndpoints(messageRepository, app)

	// Listen and serve on 0.0.0.0:8080
	app.Listen(":8080")
	wg.Done()
}

func configureMessageRepository() repository.RestMessageRepository {
	messageRepository := repository.RestMessageRepository{
		Client:               web.New(),
		RepositoryServiceUrl: "http://local.onlyone-portal.com/repository-service",
		RegistrationName:     "i18n-service",
	}
	return messageRepository
}

func configureMessageEndpoints(messageRepository repository.RestMessageRepository, app *iris.Application) {
	endpoints := api.MessageEndpoints{
		Repository: &messageRepository,
	}
	endpoints.ResgisterEndpoint(app)
}

func newWebServer() *iris.Application {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	return app
}
