package web

type WebResponse struct {
	Body   string
	Status int
}

type Uri = string

type WebRequest struct {
	Url Uri
}

type WebClient interface {
	Get(request *WebRequest) (*WebResponse, error)
}

type RestWebClient struct {
}

func (r *RestWebClient) Get(request *WebRequest) (*WebResponse, error) {
	panic("implement me")
}
