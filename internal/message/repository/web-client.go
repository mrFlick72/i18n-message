package repository

type WebClient interface {
	Get()
}

type RestWebClient struct {
}

func (r *RestWebClient) Get() {
	panic("implement me")
}
