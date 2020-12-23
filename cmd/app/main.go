package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github/mrflick72/i18n-message/api"
	"github/mrflick72/i18n-message/internal/message/repository"
	"github/mrflick72/i18n-message/internal/web"
)

func main() {
	// Creates an iris application without any middleware by default
	app := iris.New()

	// Global middleware using `UseRouter`.
	//    "github.com/kataras/iris/v12/middleware/recover"
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	app.Use(recover.New())
	app.Use(logger.New())

	messageRepository := repository.RestMessageRepository{
		Client:               web.New(),
		RepositoryServiceUrl: "http://local.onlyone-portal.com/repository-service",
		RegistrationName:     "i18n-service",
	}
	endpoints := api.MessageEndpoints{
		Repository: &messageRepository,
	}

	// Per route middleware, you can add as many as you desire.
	app.Get("/messages/{application}", endpoints.GetMessagesEndpoint)

	// Listen and serve on 0.0.0.0:8080
	app.Listen(":8080")
}
