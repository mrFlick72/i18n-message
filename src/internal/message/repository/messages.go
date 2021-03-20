package repository

import (
	"bufio"
	"fmt"
	"github/mrflick72/i18n-message/src/internal/strings/utils"
	"github/mrflick72/i18n-message/src/internal/web"
	"strings"
)

var (
	patternWithLanguage    = "%s/documents/%s?path=%s&fileName=messages_%v&fileExt=properties"
	patternWithoutLanguage = "%s/documents/%s?path=%s&fileName=messages&fileExt=properties"
	acceptableLength       = 2
	propertyKeyPosition    = 0
	propertyValuePosition  = 1
	translationNotFound    = 404
	noLanguage             = ""
	spliteratorCharacter   = "="
)

type Language = string
type Message = map[string]string

type MessageRepository interface {
	Find(application string, language string, context map[string]string) (*Message, error)
}

type RestMessageRepository struct {
	Client               web.Client
	RepositoryServiceUrl string
	RegistrationName     string
}

func (repository *RestMessageRepository) Find(application string, language string, context map[string]string) (*Message, error) {
	result := make(Message)
	repositoryServiceUrl := repository.repositoryUrlFor(application, language)
	content := repository.contentFor(application, repositoryServiceUrl, context)

	repository.loadFrom(content, result)

	return &result, nil
}

func (repository *RestMessageRepository) contentFor(application string, repositoryServiceUrl string, context map[string]string) string {
	webResponse, _ := repository.Client.Get(&web.Request{Url: repositoryServiceUrl, Header: context})
	if webResponse.Status == translationNotFound {
		repositoryServiceUrl := repository.repositoryUrlFor(application, noLanguage)
		webResponse, _ = repository.Client.Get(&web.Request{Url: repositoryServiceUrl})
	}
	content := webResponse.Body
	return content
}

func (repository *RestMessageRepository) loadFrom(content string, result Message) {
	reader := strings.NewReader(content)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		split := strings.Split(scanner.Text(), spliteratorCharacter)
		if len(split) == acceptableLength {
			result[utils.TrimStace(split[propertyKeyPosition])] = utils.TrimStace(split[propertyValuePosition])
		}
	}
}

func (repository *RestMessageRepository) repositoryUrlFor(application string, language Language) string {
	if language != noLanguage {
		return fmt.Sprintf(patternWithLanguage,
			repository.RepositoryServiceUrl, repository.RegistrationName, application, language)
	} else {
		return fmt.Sprintf(patternWithoutLanguage,
			repository.RepositoryServiceUrl, repository.RegistrationName, application)
	}
}
