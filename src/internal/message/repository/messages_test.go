package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	application          = "AN_APPLICATION"
	language             = "it"
	notAvailableLanguage = "en"
	expected             = map[string]string{"key1": "prop1", "key2": "prop2"}
	buketName            = "/bucket"

	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	client = s3.New(sess, &aws.Config{Region: aws.String("us-east-1"), Endpoint: aws.String("http://localhost:4566")})

	repository = S3MessageRepository{
		Client:    client,
		BuketName: buketName,
	}
)

func TestRestMessageRepository_Find(t *testing.T) {
	actual, _ := repository.Find(application, language, map[string]string{})
	assert.EqualValues(t, *actual, expected)
}

func TestRestMessageRepository_Find_WithoutA_Defined_Language(t *testing.T) {
	actual, _ := repository.Find(application, "", nil)
	assert.EqualValues(t, *actual, expected)
}

func TestRestMessageRepository_Find_With_Fallback(t *testing.T) {
	actual, _ := repository.Find(application, notAvailableLanguage, nil)
	assert.EqualValues(t, *actual, expected)
}
