package listener

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"sync"
	"time"
)

type UpdateSignalsListener struct {
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
	wg.Done()
}

func (listener *UpdateSignalsListener) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)
	for true {
		listener.onMessage(svc)
		time.Sleep(listener.sleep)
	}
}

func (listener *UpdateSignalsListener) onMessage(client *sqs.SQS) {
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

	if msgErr != nil {
		fmt.Printf("error in receiving message error is: %v", msgErr)
	}
	// business logic
	fmt.Print("msgResult")
	fmt.Print(msgResult)
}
