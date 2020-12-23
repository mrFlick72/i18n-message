package repository

type MessageRepository interface {
	Find(application string, language string) (*map[string]string, error)
}

type RestMessageRepository struct {
	client               WebClient
	repositoryServiceUrl string
	registrationName     string
}

func (m *RestMessageRepository) Find(application string, language string) (*map[string]string, error) {

	return nil, nil
}
