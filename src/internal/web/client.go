package web

import "github.com/go-resty/resty/v2"

type Response struct {
	Body   string
	Status int
}

type Uri = string

type Request struct {
	Url Uri
}

type Client interface {
	Get(request *Request) (*Response, error)
}

type RestWebClient struct {
	client *resty.Client
}

func (r *RestWebClient) Get(request *Request) (*Response, error) {
	response, _ := r.client.R().Get(request.Url)

	return &Response{
		Body:   string(response.Body()),
		Status: response.StatusCode(),
	}, nil
}

func New() Client {
	return &RestWebClient{
		client: resty.New(),
	}
}
