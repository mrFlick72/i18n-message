package listener

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/yalp/jsonpath"
	"sync"
	"time"
)

type UpdateSignalsListener struct {
	queueMapping        map[string]string
	toStop              bool
	queueURL            string
	timeout             int64
	maxNumberOfMessages int64
	sleep               time.Duration
}

func New(queueMapping map[string]string, queueURL string, timeout int64, maxNumberOfMessages int64, sleep time.Duration) *UpdateSignalsListener {
	return &UpdateSignalsListener{
		queueMapping:        queueMapping,
		queueURL:            queueURL,
		timeout:             timeout,
		maxNumberOfMessages: maxNumberOfMessages,
		sleep:               sleep,
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
		},
	)
	return msgResult, msgErr
}

func (listener *UpdateSignalsListener) dispatchMessageTo(msgErr error, msgResult *sqs.ReceiveMessageOutput, client *sqs.SQS) {
	if msgErr != nil {
		fmt.Printf("error in receiving message error is: %v", msgErr)
	}
	applicationNameQuery, err := jsonpath.Prepare("$.application.value")
	if err != nil {
		fmt.Printf("error during jsonpath query preparation error message: %v", err)
		return
	}

	applicationName, err := applicationNameQuery(msgResult)
	if err != nil {
		fmt.Printf("error during jsonpath query evaluation error message: %v", err)
		return
	}

	//business logic
	message := "signal for messages updates from i18n-messages service"
	queue := listener.queueMapping[applicationName.(string)]
	_, err = client.SendMessage(
		&sqs.SendMessageInput{
			MessageBody: &message,
			QueueUrl:    &queue,
		})
	fmt.Printf("error during update signal send. Error message: %v", err)
}
