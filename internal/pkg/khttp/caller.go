package khttp

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type HttpCaller interface {
	GET() ([]byte, error)
	POST() ([]byte, error)
	PUT() ([]byte, error)
	PATCH() ([]byte, error)
	DELETE() ([]byte, error)
	SetClient(client *http.Client)
}

type Caller struct {
	URL    string
	Body   interface{}
	Header map[string]string
	Client *http.Client
}

type value struct {
	URL    string
	Method string
	Body   interface{}
	Header map[string]string
	Client *http.Client
}

// New create the new http caller with parameter url, request body, and http header
func New(url string, body interface{}, header map[string]string) HttpCaller {
	return &Caller{
		URL:    url,
		Body:   body,
		Header: header,
	}
}

// SetClient set new http client
func (c *Caller) SetClient(client *http.Client) {
	c.Client = client
}

// GET func call server with method GET
func (c *Caller) GET() ([]byte, error) {
	return invoke(value{
		URL:    c.URL,
		Method: http.MethodGet,
		Body:   c.Body,
		Header: c.Header,
		Client: c.Client,
	})
}

// POST func call server with method POST
func (c *Caller) POST() ([]byte, error) {
	return invoke(value{
		URL:    c.URL,
		Method: http.MethodPost,
		Body:   c.Body,
		Header: c.Header,
		Client: c.Client,
	})
}

// PUT func call server with method PUT
func (c *Caller) PUT() ([]byte, error) {
	return invoke(value{
		URL:    c.URL,
		Method: http.MethodPut,
		Body:   c.Body,
		Header: c.Header,
		Client: c.Client,
	})
}

// PATCH func call server with method PATCH
func (c *Caller) PATCH() ([]byte, error) {
	return invoke(value{
		URL:    c.URL,
		Method: http.MethodPatch,
		Body:   c.Body,
		Header: c.Header,
		Client: c.Client,
	})
}

// DELETE func call server with method DELETE
func (c *Caller) DELETE() ([]byte, error) {
	return invoke(value{
		URL:    c.URL,
		Method: http.MethodDelete,
		Body:   c.Body,
		Header: c.Header,
		Client: c.Client,
	})
}

func invoke(v value) ([]byte, error) {
	if v.Client == nil {
		v.Client = &http.Client{}
	}

	body := bytes.NewBuffer(nil)
	if v.Body != nil {
		contentType, ok := v.Header["Content-Type"]
		if !ok {
			return nil, errors.New("Please provide Content-Type")
		}

		var b []byte
		var errMarshal error
		normalizeType := strings.ToLower(contentType)
		if strings.Contains(normalizeType, "json") {
			b, errMarshal = json.Marshal(v.Body)
			if errMarshal != nil {
				return nil, errMarshal
			}
		} else if strings.Contains(normalizeType, "xml") {
			b, errMarshal = xml.Marshal(v.Body)
			if errMarshal != nil {
				return nil, errMarshal
			}
		} else {
			return nil, errors.New("Content-Type of request need to be in JSON or XML format")
		}

		body = bytes.NewBuffer(b)
	}
	// make a new http request
	req, _ := http.NewRequest(v.Method, v.URL, body)

	// set the request header
	for k, v := range v.Header {
		req.Header.Set(k, v)
	}
	res, err := v.Client.Do(req)
	// cannot connect to server
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bytesContent, err := ioutil.ReadAll(res.Body)

	return bytesContent, err
}
