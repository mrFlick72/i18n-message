package repository

import (
	json "encoding/json"
	"fmt"
	"github/mrflick72/i18n-message/internal/web"
)

type MessageRepository interface {
	Find(application string, language string) (map[string]string, error)
}

type RestMessageRepository struct {
	client               web.WebClient
	repositoryServiceUrl string
	registrationName     string
}

func (repository *RestMessageRepository) Find(application string, language string) (*map[string]string, error) {
	client := repository.client
	serviceUrl := fmt.Sprintf("%s/documents/%s?path=%s&fileName=messages_it&fileExt=properties",
		repository.repositoryServiceUrl, repository.registrationName, application)

	result := make(map[string]string)
	webResponse, _ := client.Get(web.WebRequest{Url: serviceUrl})
	content := webResponse.Body

	err := json.Unmarshal([]byte(content), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
