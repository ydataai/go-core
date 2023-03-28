package http

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	*http.Response
}

func (r *Response) DecodeJSON(into interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(into); err != nil {
		return err
	}
	return nil
}

// HasStatusCode returns true if the Response's status code is one of the specified values.
func (r *Response) HasStatusCode(statusCodes ...int) bool {
	if r == nil {
		return false
	}

	for _, sc := range statusCodes {
		if r.StatusCode == sc {
			return true
		}
	}
	return false
}
