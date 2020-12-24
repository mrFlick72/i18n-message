package configuration

import (
	"github.com/kataras/iris/v12"
	"github/mrflick72/i18n-message/api"
	"github/mrflick72/i18n-message/internal/message/listener"
	"github/mrflick72/i18n-message/internal/message/repository"
	"github/mrflick72/i18n-message/internal/web"
	"os"
	"strconv"
	"sync"
	"time"
)

func ConfigureMessageRepository() repository.RestMessageRepository {
	messageRepository := repository.RestMessageRepository{
		Client:               web.New(),
		RepositoryServiceUrl: os.Getenv("REPOSITORY_SERVICE_URL"),
		RegistrationName:     os.Getenv("REGISTRATION_NAME"),
	}
	return messageRepository
}

func ConfigureMessageEndpoints(messageRepository repository.RestMessageRepository, app *iris.Application) {
	endpoints := api.MessageEndpoints{
		Repository: &messageRepository,
	}
	endpoints.ResgisterEndpoint(app)
}

func DocumentUpdatesListener(wg *sync.WaitGroup) {
	timeout, _ := strconv.ParseInt(os.Getenv("SQS_TIMEOUT"), 10, 64)
	maxNumberOfMessages, _ := strconv.ParseInt(os.Getenv("SQS_MAX_NUMBER_OF_MESSAGES"), 10, 64)
	sleep, _ := time.ParseDuration(os.Getenv("SQS_LISTENER_PAUSE_TIMEOUT"))

	listener.New(
		os.Getenv("SQS_QUEUE_URL"),
		timeout,
		maxNumberOfMessages,
		sleep,
	).Start(wg)
}
