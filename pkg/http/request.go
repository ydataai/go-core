package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type Request struct {
	*http.Request
}

// NewRequest creates a new http.Request with the specified input.
func NewRequest(ctx context.Context, httpMethod string, endpoint string) (*Request, error) {
	req, err := http.NewRequestWithContext(ctx, httpMethod, endpoint, nil)
	if err != nil {
		return nil, err
	}
	if req.URL.Host == "" {
		return nil, errors.New("no Host in request URL")
	}
	if !(req.URL.Scheme == "http" || req.URL.Scheme == "https") {
		return nil, fmt.Errorf("unsupported protocol scheme %s", req.URL.Scheme)
	}
	return &Request{req}, nil
}

func (req *Request) SetBody(body io.ReadSeekCloser, contentType string) error {
	var err error
	var size int64

	// get body size
	if body != nil {
		size, err = body.Seek(0, io.SeekEnd)
		if err != nil {
			return err
		}
	}

	// treat an empty stream as a nil body
	if size == 0 {
		body = nil
		req.Header.Del(HeaderContentLength)
	} else {
		_, err = body.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}

		req.Header.Set(HeaderContentLength, strconv.FormatInt(size, 10))

		req.GetBody = func() (io.ReadCloser, error) {
			_, err := body.Seek(0, io.SeekStart)
			return body, err
		}
	}

	req.Body = body
	req.ContentLength = size
	if contentType == "" {
		req.Header.Del(HeaderContentType)
	} else {
		req.Header.Set(HeaderContentType, contentType)
	}
	return nil
}

// JoinPaths concatenates multiple URL path segments into one path, inserting path separation characters as required.
func JoinPaths(root string, paths ...string) string {
	if len(paths) == 0 {
		return root
	}

	finalPath := path.Join(paths...)
	// path.Join will remove any trailing slashes.
	// if one was provided, preserve it.
	if strings.HasSuffix(paths[len(paths)-1], "/") && !strings.HasSuffix(finalPath, "/") {
		finalPath += "/"
	}

	return root + finalPath
}

// EncodeAsJSON calls encodes an object as JSON encoding with SetBody.
func (req *Request) EncodeAsJSON(obj any) error {
	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(obj); err != nil {
		return fmt.Errorf("error marshalling type %T: %s", obj, err)
	}

	return req.SetBody(nopCloser{bytes.NewReader(buffer.Bytes())}, ContentTypeAppJSON)
}

// extracted from io
type nopCloser struct {
	io.ReadSeeker
}

func (n nopCloser) Close() error {
	return nil
}
