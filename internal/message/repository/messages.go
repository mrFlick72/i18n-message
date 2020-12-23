package repository

import (
	"bufio"
	"fmt"
	"github/mrflick72/i18n-message/internal/web"
	"strings"
)

type Language = string
type Message = map[string]string

type MessageRepository interface {
	Find(application string, language string) (Message, error)
}

type RestMessageRepository struct {
	client               web.WebClient
	repositoryServiceUrl string
	registrationName     string
}

func (repository *RestMessageRepository) Find(application string, language Language) (*Message, error) {
	result := make(Message)
	repositoryServiceUrl := repository.repositoryUrlFor(application, language)
	content := repository.contentFor(application, repositoryServiceUrl)

	repository.loadFrom(content, result)

	return &result, nil
}

func (repository *RestMessageRepository) contentFor(application string, repositoryServiceUrl string) string {
	webResponse, _ := repository.client.Get(web.WebRequest{Url: repositoryServiceUrl})
	if webResponse.Status == 404 {
		repositoryServiceUrl := repository.repositoryUrlFor(application, "")
		webResponse, _ = repository.client.Get(web.WebRequest{Url: repositoryServiceUrl})
	}
	content := webResponse.Body
	return content
}

func (repository *RestMessageRepository) loadFrom(content string, result Message) {
	reader := strings.NewReader(content)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		split := strings.Split(scanner.Text(), "=")
		result[split[0]] = split[1]
	}
}

func (repository *RestMessageRepository) repositoryUrlFor(application string, language Language) string {
	if language != "" {
		return fmt.Sprintf("%s/documents/%s?path=%s&fileName=messages_%v&fileExt=properties",
			repository.repositoryServiceUrl, repository.registrationName, application, language)
	} else {
		return fmt.Sprintf("%s/documents/%s?path=%s&fileName=messages&fileExt=properties",
			repository.repositoryServiceUrl, repository.registrationName, application)
	}
}
