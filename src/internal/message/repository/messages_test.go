package repository

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

var (
	baseUrl          = "http://localhost/repository-service"
	registrationName = "i18n-messages"
	application      = "AN_APPLICATION"
	language         = "it"
	expected         = map[string]string{"key1": "prop1", "key2": "prop2"}
	body             = "key1=prop1\nkey2=prop2"
)

func TestRestMessageRepository_Find(t *testing.T) {
	os.Setenv("AWS_ACCESS_KEY_ID", "XXXXXXXXXXXX")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "XXXXXXXXXXXX")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	lang := "it"
	buketName := "buket"
	filePath := "messages_it.properties"

	client := s3.New(sess, &aws.Config{Region: aws.String("us-east-1"), Endpoint: aws.String("http://localhost:4566/")})

	setUp(filePath, buketName, client, t)

	repository := S3MessageRepository{
		Client:    client,
		BuketName: buketName,
	}

	actual, err := repository.Find("AN_APPLICATION", lang, nil)
	fmt.Println("err")
	fmt.Println(err)

	assert.EqualValues(t, *actual, expected)

	tearDown(buketName, client)
}

func setUp(filePath string, buketName string, client *s3.S3, t *testing.T) {
	result, err := client.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(buketName)})
	if err != nil {
		fmt.Println("error during the buket creation")
		fmt.Println(err.Error())
		t.Fail()
	} else {
		fmt.Println(result)
	}

	// Open the file from the file path
	upFile, _ := os.Open(filePath)
	defer upFile.Close()

	// Get the file info
	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)

	_, err = client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(buketName),
		Key:           aws.String("AN_APPLICATION/messages_it.properties"),
		ACL:           aws.String("private"),
		Body:          bytes.NewReader(fileBuffer),
		ContentLength: aws.Int64(fileSize),
		ContentType:   aws.String(http.DetectContentType(fileBuffer)),
	})
	if err != nil {
		fmt.Println("error during the object upload")
		fmt.Println(err.Error())
		t.Fail()
	}
}

func tearDown(buketName string, client *s3.S3) {
	bucket, err := client.DeleteBucket(&s3.DeleteBucketInput{Bucket: aws.String(buketName)})
	if err != nil {
		fmt.Println("error")
		fmt.Println(err.Error())
	} else {
		fmt.Println(bucket)
	}
}

/*
func TestRestMessageRepository_Find_WithoutA_Defined_Language(t *testing.T) {
	client := new(MockedWebClientObject)

	client.On("Get", &web.Request{
		Url: serviceUrlFor(baseUrl, registrationName, application, ""),
	}).Return(&web.Response{
		Body:   body,
		Status: 200,
	})

	repository := RestMessageRepository{
		Client:               client,
		RepositoryServiceUrl: "http://localhost/repository-service",
		RegistrationName:     "i18n-messages",
	}

	actual, _ := repository.Find("AN_APPLICATION", "", nil)
	assert.EqualValues(t, *actual, expected)
}

func TestRestMessageRepository_Find_Wit_Fallback(t *testing.T) {
	client := new(MockedWebClientObject)

	repository := RestMessageRepository{
		Client:               client,
		RepositoryServiceUrl: baseUrl,
		RegistrationName:     registrationName,
	}

	client.On("Get", &web.Request{
		Url: serviceUrlFor(baseUrl, registrationName, application, language),
	}).Return(&web.Response{
		Status: 404,
	})

	client.On("Get", &web.Request{
		Url: serviceUrlFor(baseUrl, registrationName, application, ""),
	}).Return(&web.Response{
		Body:   body,
		Status: 200,
	})

	actual, _ := repository.Find(application, language, nil)
	assert.EqualValues(t, *actual, expected)
}

type MockedWebClientObject struct {
	mock.Mock
}

func (mock *MockedWebClientObject) Get(request *web.Request) (*web.Response, error) {
	called := mock.Called(request)
	return called.Get(0).(*web.Response), nil
}

func serviceUrlFor(baseUrl string, registrationUrl string, application string, language Language) string {
	if language != "" {
		return fmt.Sprintf("%s/documents/%s?path=%s&fileName=messages_%v&fileExt=properties",
			baseUrl, registrationUrl, application, language)
	} else {
		return fmt.Sprintf("%s/documents/%s?path=%s&fileName=messages&fileExt=properties",
			baseUrl, registrationUrl, application)
	}
}
*/
