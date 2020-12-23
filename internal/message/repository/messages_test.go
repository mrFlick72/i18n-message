package repository

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github/mrflick72/i18n-message/internal/web"
	"testing"
)

var (
	baseUrl          = "http://localhost/repository-service"
	registrationName = "i18n-messages"
	application      = "AN_APPLICATION"
	language         = "it"
)

func TestRestMessageRepository_Find(t *testing.T) {
	lang := "it"
	client := &TestableWebClient{
		baseUrl:         "http://localhost/repository-service",
		registrationUrl: "i18n-messages",
		application:     "AN_APPLICATION",
		language:        &lang,
		status:          200,
	}
	repository := RestMessageRepository{
		client:               client,
		repositoryServiceUrl: "http://localhost/repository-service",
		registrationName:     "i18n-messages",
	}

	actual, _ := repository.Find("AN_APPLICATION", &lang)
	expected := map[string]string{"key1": "prop1"}

	assert.EqualValues(t, *actual, expected)
}

func TestRestMessageRepository_Find_WithoutA_Defined_Language(t *testing.T) {
	client := &TestableWebClient{
		baseUrl:         "http://localhost/repository-service",
		registrationUrl: "i18n-messages",
		application:     "AN_APPLICATION",
		status:          200,
	}
	repository := RestMessageRepository{
		client:               client,
		repositoryServiceUrl: "http://localhost/repository-service",
		registrationName:     "i18n-messages",
	}

	var language *string

	actual, _ := repository.Find("AN_APPLICATION", language)
	expected := map[string]string{"key1": "prop1"}

	assert.EqualValues(t, *actual, expected)
}

func TestRestMessageRepository_Find_Wit_Fallback(t *testing.T) {
	client := new(MockedWebClientObject)

	repository := RestMessageRepository{
		client:               client,
		repositoryServiceUrl: baseUrl,
		registrationName:     registrationName,
	}

	client.On("Get", web.WebRequest{
		Url: serviceUrlFor(baseUrl, registrationName, application, language),
	}).Return(web.WebResponse{
		Status: 404,
	})

	client.On("Get", web.WebRequest{
		Url: serviceUrlFor(baseUrl, registrationName, application, ""),
	}).Return(web.WebResponse{
		Body:   "{\"key1\":\"prop1\"}",
		Status: 200,
	})

	actual, _ := repository.Find(application, &language)
	expected := map[string]string{"key1": "prop1"}

	assert.EqualValues(t, *actual, expected)
}

type MockedWebClientObject struct {
	mock.Mock
}

func (mock *MockedWebClientObject) Get(request web.WebRequest) (web.WebResponse, error) {
	called := mock.Called(request)
	return called.Get(0).(web.WebResponse), nil
}

type TestableWebClient struct {
	baseUrl         string
	registrationUrl string
	application     string
	language        *string
	status          int
}

func (receiver *TestableWebClient) Get(request web.WebRequest) (web.WebResponse, error) {
	var expectedUrl string

	if receiver.language != nil {
		expectedUrl = fmt.Sprintf("%s/documents/%s?path=%s&fileName=messages_%v&fileExt=properties",
			receiver.baseUrl, receiver.registrationUrl, receiver.application, "it")
	} else {
		expectedUrl = fmt.Sprintf("%s/documents/%s?path=%s&fileName=messages&fileExt=properties",
			receiver.baseUrl, receiver.registrationUrl, receiver.application)
	}

	if request.Url == expectedUrl {
		return web.WebResponse{
			Body:   "{\"key1\":\"prop1\"}",
			Status: receiver.status,
		}, nil

	}
	return web.WebResponse{}, nil
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
