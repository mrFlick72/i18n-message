package web

import "github.com/go-resty/resty/v2"

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
	client *resty.Client
}

func (r *RestWebClient) Get(request *WebRequest) (*WebResponse, error) {
	response, _ := r.client.R().Get(request.Url)

	return &WebResponse{
		Body:   string(response.Body()),
		Status: response.StatusCode(),
	}, nil
}

func New() WebClient {
	return &RestWebClient{
		client: resty.New(),
	}
}
