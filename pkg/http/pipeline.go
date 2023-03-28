package http

import "net/http"

type pipeline struct {
	*http.Client
}

type Pipeline interface {
	Do(*Request) (*Response, error)
}

func NewPipeline() Pipeline {
	return pipeline{defaultHTTPClient}
}

func (p pipeline) Do(req *Request) (*Response, error) {
	resp, err := p.Client.Do(req.Request)
	if err != nil {
		return nil, err
	}
	return &Response{resp}, nil
}
