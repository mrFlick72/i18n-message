package listener

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github/mrflick72/i18n-message/src/internal/logging"
	"sync"
	"time"
)

var logger = logging.GetLoggerInstance()

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

func New(queueMapping map[string]string, queueURL string, timeout int64, waitTimeSeconds int64, maxNumberOfMessages int64, sleep time.Duration) *UpdateSignalsListener {
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
	for {
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
	if msgErr == nil {
		listener.dispatchMessageTo(msgResult, client)
	} else {
		listener.logger.LogErrorFor(fmt.Sprintf("error in receiving message error is: %v", msgErr))
	}
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
	return msgResult, msgErr
}

func (listener *UpdateSignalsListener) dispatchMessageTo(msgResult *sqs.ReceiveMessageOutput, client *sqs.SQS) {
	if len(msgResult.Messages) != 0 {
		listener.processMessages(msgResult, client)
	} else {
		listener.logger.LogDebugFor("empty message")
	}
}

func (listener *UpdateSignalsListener) processMessages(msgResult *sqs.ReceiveMessageOutput, client *sqs.SQS) {
	listener.logger.LogInfoFor("msgResult")
	listener.logger.LogInfoFor(msgResult)

	for _, message := range msgResult.Messages {
		applicationUpdateSignal, err := listener.getApplicationDataFrom(message)
		fireUpdateEventToError := listener.fireUpdateEventTo(*applicationUpdateSignal, err, client)
		listener.deleteConsumedMessage(*message, fireUpdateEventToError, client)
	}
}

func (listener *UpdateSignalsListener) getApplicationDataFrom(message *sqs.Message) (*I18nApplicationUpdateSignal, error) {
	listener.logger.LogDebugFor("message.Body")
	listener.logger.LogDebugFor(*message.Body)

	applicationUpdateSignal := I18nApplicationUpdateSignal{}
	err := json.Unmarshal([]byte(*message.Body), &applicationUpdateSignal)

	if err != nil {
		listener.logger.LogDebugFor(fmt.Sprintf("error during unmarshalling error message: %v", err))
		return nil, err
	} else {
		listener.logger.LogInfoFor(fmt.Sprintf("application data to update: %v", applicationUpdateSignal))
	}
	return &applicationUpdateSignal, err
}

func (listener *UpdateSignalsListener) fireUpdateEventTo(application I18nApplicationUpdateSignal, err error, client *sqs.SQS) error {
	if err == nil {
		message := "signal for messages updates from i18n-messages service"
		queue := listener.queueMapping[application.Name]
		input := sqs.SendMessageInput{
			MessageBody: &message,
			QueueUrl:    &queue,
		}
		_, err = client.SendMessage(&input)

		if err != nil {
			listener.logger.LogErrorFor(fmt.Sprintf("error during update signal send. Error message: %v", err))
		} else {
			listener.logger.LogInfoFor(fmt.Sprintf("message sent:  %v", input))
		}
	}
	return err
}

func (listener *UpdateSignalsListener) deleteConsumedMessage(message sqs.Message, err error, client *sqs.SQS) {
	if err == nil {
		listener.logger.LogInfoFor(
			fmt.Sprintf("start to delete message with ReceiptHandle: %v \n and quque: %v", *message.ReceiptHandle, listener.queueURL))
		_, err = client.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      &listener.queueURL,
			ReceiptHandle: message.ReceiptHandle,
		})
	}

	if err != nil {
		listener.logger.LogErrorFor(fmt.Sprintf("Error to delete consumed message error message: %v", err))
	}
}

type I18nApplicationUpdateSignal struct {
	Name string `json:"path"`
}
