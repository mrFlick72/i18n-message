package application

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kataras/iris/v12"
	"github/mrflick72/i18n-message/src/api"
	"github/mrflick72/i18n-message/src/configuration"
	"github/mrflick72/i18n-message/src/internal/message/listener"
	"github/mrflick72/i18n-message/src/internal/message/repository"
	"strconv"
	"sync"
	"time"
)

var manager = configuration.GetConfigurationManagerInstance()

func ConfigureMessageRepository() repository.S3MessageRepository {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	messageRepository := repository.S3MessageRepository{
		Client:    s3.New(sess),
		BuketName: manager.GetConfigFor("MESSAGES_BUKET"),
	}
	return messageRepository
}

func ConfigureMessageEndpoints(messageRepository repository.S3MessageRepository, app *iris.Application) {
	endpoints := api.MessageEndpoints{
		Repository: &messageRepository,
	}
	endpoints.RegisterEndpoint(app)
}

func DocumentUpdatesListener(wg *sync.WaitGroup) {
	timeout, _ := strconv.ParseInt(manager.GetConfigFor("SQS_TIMEOUT"), 10, 64)
	waitTimeSeconds, _ := strconv.ParseInt(manager.GetConfigFor("SQS_WAIT_TIME_SECONDS"), 10, 64)
	maxNumberOfMessages, _ := strconv.ParseInt(manager.GetConfigFor("SQS_MAX_NUMBER_OF_MESSAGES"), 10, 64)
	sleep, _ := time.ParseDuration(manager.GetConfigFor("SQS_LISTENER_PAUSE_TIMEOUT"))
	mapping := manager.GetStringMapFor("update-signals")

	listener.New(
		mapping,
		manager.GetConfigFor("SQS_QUEUE_URL"),
		timeout,
		waitTimeSeconds,
		maxNumberOfMessages,
		sleep,
	).Start(wg)
}
