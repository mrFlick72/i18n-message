package repository

import "github/mrflick72/i18n-message/internal/web"

type MessageRepository interface {
	Find(application string, language string) (*map[string]string, error)
}

type RestMessageRepository struct {
	client               web.WebClient
	repositoryServiceUrl string
	registrationName     string
}

func (repository *RestMessageRepository) Find(application string, language string) (*map[string]string, error) {

	return nil, nil
}
