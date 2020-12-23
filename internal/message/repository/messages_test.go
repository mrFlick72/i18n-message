package repository

import (
	"fmt"
	"github/mrflick72/i18n-message/internal/web"
	"testing"
)

func TestRestMessageRepository_Find(t *testing.T) {
	client := &TestableWebClient{}
	repository := RestMessageRepository{
		client:               client,
		repositoryServiceUrl: "http://localhost/repository-service",
		registrationName:     "i18n-messages",
	}

	messages, _ := repository.Find("AN_APPLICATION", "it")
	expected := map[string]string{"key1": "value1"}
	assertThatContentIs(t, messages, &expected, fmt.Sprintf("actual is %v while expected is %v", messages, expected))
}

func assertThatContentIs(t *testing.T, actual *map[string]string, expected *map[string]string, message string) {
	if actual != expected {
		t.Log(message)
		t.Fail()
	}
}

func assertThatNoErrorFor(t *testing.T, err error, errorMessage string) {
	if err != nil {
		t.Log(errorMessage)
		t.Fail()
	}
}

type TestableWebClient struct {
	baseUrl         string
	registrationUrl string
	application     string
}

func (receiver *TestableWebClient) Get(request web.WebRequest) (web.WebResponse, error) {
	expectedUrl := fmt.Sprintf("%s/documents/%s?path=%s&fileName=messages_it&fileExt=properties",
		receiver.baseUrl, receiver.registrationUrl, receiver.application)

	if request.Url == expectedUrl {
		return web.WebResponse{
			Body: "{'key1':'prop1'}",
		}, nil

	}
	return web.WebResponse{}, nil
}
