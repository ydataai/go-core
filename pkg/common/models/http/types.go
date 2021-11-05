package http

import "time"

// RESTRequest defines a struct for a request
type RESTRequest struct {
	URL         string
	Body        interface{}
	Headers     map[string]string
	QueryParams map[string]string
	Timeout     time.Duration
}
