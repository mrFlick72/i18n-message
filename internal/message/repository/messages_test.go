package repository

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github/mrflick72/i18n-message/internal/web"
	"testing"
)

func TestRestMessageRepository_Find(t *testing.T) {
	lang := "it"
	client := &TestableWebClient{
		baseUrl:         "http://localhost/repository-service",
		registrationUrl: "i18n-messages",
		application:     "AN_APPLICATION",
		language:        &lang,
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

type TestableWebClient struct {
	baseUrl         string
	registrationUrl string
	application     string
	language        *string
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
			Body: "{\"key1\":\"prop1\"}",
		}, nil

	}
	return web.WebResponse{}, nil
}
