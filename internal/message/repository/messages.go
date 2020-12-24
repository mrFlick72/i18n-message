package repository

import (
	"bufio"
	"fmt"
	"github/mrflick72/i18n-message/internal/strings/utils"
	"github/mrflick72/i18n-message/internal/web"
	"strings"
)

type Language = string
type Message = map[string]string

type MessageRepository interface {
	Find(application string, language string) (*Message, error)
}

type RestMessageRepository struct {
	Client               web.WebClient
	RepositoryServiceUrl string
	RegistrationName     string
}

func (repository *RestMessageRepository) Find(application string, language Language) (*Message, error) {
	result := make(Message)
	repositoryServiceUrl := repository.repositoryUrlFor(application, language)
	fmt.Println("repositoryServiceUrl")
	fmt.Println(repositoryServiceUrl)
	content := repository.contentFor(application, repositoryServiceUrl)

	repository.loadFrom(content, result)

	return &result, nil
}

func (repository *RestMessageRepository) contentFor(application string, repositoryServiceUrl string) string {
	webResponse, _ := repository.Client.Get(&web.WebRequest{Url: repositoryServiceUrl})
	if webResponse.Status == 404 {
		repositoryServiceUrl := repository.repositoryUrlFor(application, "")
		webResponse, _ = repository.Client.Get(&web.WebRequest{Url: repositoryServiceUrl})
	}
	content := webResponse.Body
	return content
}

func (repository *RestMessageRepository) loadFrom(content string, result Message) {
	fmt.Println("content")
	fmt.Println(content)
	reader := strings.NewReader(content)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		split := strings.Split(scanner.Text(), "=")
		if len(split) == 2 {
			result[utils.TrimStace(split[0])] = utils.TrimStace(split[1])
		}
	}
}

func (repository *RestMessageRepository) repositoryUrlFor(application string, language Language) string {
	if language != "" {
		return fmt.Sprintf("%s/documents/%s?path=%s&fileName=messages_%v&fileExt=properties",
			repository.RepositoryServiceUrl, repository.RegistrationName, application, language)
	} else {
		return fmt.Sprintf("%s/documents/%s?path=%s&fileName=messages&fileExt=properties",
			repository.RepositoryServiceUrl, repository.RegistrationName, application)
	}
}
