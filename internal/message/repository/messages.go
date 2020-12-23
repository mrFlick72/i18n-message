package repository

import (
	"encoding/json"
	"fmt"
	"github/mrflick72/i18n-message/internal/web"
)

type Language = string

type MessageRepository interface {
	Find(application string, language string) (map[string]string, error)
}

type RestMessageRepository struct {
	client               web.WebClient
	repositoryServiceUrl string
	registrationName     string
}

func (repository *RestMessageRepository) Find(application string, language *Language) (*map[string]string, error) {
	client := repository.client
	result := make(map[string]string)
	serviceUrl := repositoryUrlFor(application, language, repository)

	webResponse, _ := client.Get(web.WebRequest{Url: serviceUrl})
	content := webResponse.Body

	err := json.Unmarshal([]byte(content), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func repositoryUrlFor(application string, language *Language, repository *RestMessageRepository) string {
	if language != nil {
		return fmt.Sprintf("%s/documents/%s?path=%s&fileName=messages_%v&fileExt=properties",
			repository.repositoryServiceUrl, repository.registrationName, application, *language)
	} else {
		return fmt.Sprintf("%s/documents/%s?path=%s&fileName=messages&fileExt=properties",
			repository.repositoryServiceUrl, repository.registrationName, application)
	}
}
