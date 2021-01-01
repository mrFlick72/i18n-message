package listener

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	//"github.com/yalp/jsonpath"
	"github/mrflick72/i18n-message/src/internal/logging"
	"sync"
	"time"
)

type UpdateSignalsListener struct {
	queueMapping        map[string]string
	toStop              bool
	queueURL            string
	timeout             int64
	waitTimeSeconds     int64
	maxNumberOfMessages int64
	sleep               time.Duration
	logger              *logging.Logger
}

func New(queueMapping map[string]string, queueURL string, timeout int64, waitTimeSeconds int64, maxNumberOfMessages int64, sleep time.Duration, logger *logging.Logger) *UpdateSignalsListener {
	return &UpdateSignalsListener{
		queueMapping:        queueMapping,
		queueURL:            queueURL,
		timeout:             timeout,
		waitTimeSeconds:     waitTimeSeconds,
		maxNumberOfMessages: maxNumberOfMessages,
		sleep:               sleep,
		logger:              logger,
	}
}

func (listener *UpdateSignalsListener) Stop(wg *sync.WaitGroup) {
	listener.toStop = true
	wg.Done()
}

func (listener *UpdateSignalsListener) Start(wg *sync.WaitGroup) {
	listener.toStop = false
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)
	for true {
		if listener.toStop == true {
			break
		}

		listener.onMessage(svc)
		time.Sleep(listener.sleep)
	}
	wg.Done()
}

func (listener *UpdateSignalsListener) onMessage(client *sqs.SQS) {
	msgResult, msgErr := listener.receiveFrom(client)
	listener.dispatchMessageTo(msgErr, msgResult, client)
}

func (listener *UpdateSignalsListener) receiveFrom(client *sqs.SQS) (*sqs.ReceiveMessageOutput, error) {
	msgResult, msgErr := client.ReceiveMessage(
		&sqs.ReceiveMessageInput{
			AttributeNames: []*string{
				aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
			},
			MessageAttributeNames: []*string{
				aws.String(sqs.QueueAttributeNameAll),
			},
			QueueUrl:            &(listener).queueURL,
			MaxNumberOfMessages: &(listener).maxNumberOfMessages,
			VisibilityTimeout:   &(listener).timeout,
			WaitTimeSeconds:     &(listener).waitTimeSeconds,
		},
	)
	listener.logger.LogInfoFor(fmt.Sprintf("messge: %v", msgResult))
	fmt.Printf("messge: %v", msgResult)
	return msgResult, msgErr
}

func (listener *UpdateSignalsListener) dispatchMessageTo(msgErr error, msgResult *sqs.ReceiveMessageOutput, client *sqs.SQS) {
	if len(msgResult.Messages) != 0 {

		listener.logger.LogInfoFor("msgResult")
		listener.logger.LogInfoFor(msgResult)
		if msgErr != nil {
			listener.logger.LogErrorFor(fmt.Sprintf("error in receiving message error is: %v", msgErr))
			return
		}

		for _, message := range msgResult.Messages {
			listener.logger.LogDebugFor("message.Body")
			listener.logger.LogDebugFor(*message.Body)

			applicationUpdateSignal := I18nApplicationUpdateSignal{}
			err := json.Unmarshal([]byte(*message.Body), &applicationUpdateSignal)

			if err != nil {
				listener.logger.LogDebugFor(fmt.Sprintf("error during unmarshalling error message: %v", err))
				return
			} else {
				listener.logger.LogDebugFor(fmt.Sprintf("application data to update: %v", applicationUpdateSignal))
			}

			//business logic
			listener.fireUpdateEventTo(applicationUpdateSignal, err, client)
		}
	} else {
		listener.logger.LogDebugFor("empty message")
	}
}

func (listener *UpdateSignalsListener) fireUpdateEventTo(application I18nApplicationUpdateSignal, err error, client *sqs.SQS) {
	message := "signal for messages updates from i18n-messages service"
	queue := listener.queueMapping[application.Path]
	_, err = client.SendMessage(
		&sqs.SendMessageInput{
			MessageBody: &message,
			QueueUrl:    &queue,
		})
	if err != nil {
		listener.logger.LogErrorFor(fmt.Sprintf("error during update signal send. Error message: %v", err))
	}
}

type I18nApplicationUpdateSignal struct {
	Path string
}
