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
	toStop              bool
	queueURL            string
	timeout             int64
	maxNumberOfMessages int64
	sleep               time.Duration
}

func New(queueURL string, timeout int64, maxNumberOfMessages int64, sleep time.Duration) *UpdateSignalsListener {
	return &UpdateSignalsListener{
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
	listener.dispatchMessageTo(msgErr, msgResult)
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

func (listener *UpdateSignalsListener) dispatchMessageTo(msgErr error, msgResult *sqs.ReceiveMessageOutput) {
	if msgErr != nil {
		fmt.Printf("error in receiving message error is: %v", msgErr)
	}
	applicationNameQuery, _ := jsonpath.Prepare("$.application.value")
	applicationName, err := applicationNameQuery(msgResult)

	//business logic
	fmt.Print(msgResult)
	fmt.Printf("error during executing jsonpath query: %v", err)
	fmt.Print(applicationName)

}
