package web

import "github.com/go-resty/resty/v2"

type Response struct {
	Body   string
	Status int
}

type Uri = string
type Headers = map[string]string

type Request struct {
	Url    Uri
	Header Headers
}

type Client interface {
	Get(request *Request) (*Response, error)
}

type RestWebClient struct {
	client *resty.Client
}

func (r *RestWebClient) Get(request *Request) (*Response, error) {
	client := r.client.R()
	for key, value := range request.Header {
		client.SetHeader(key, value)
	}

	response, _ := client.Get(request.Url)

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
