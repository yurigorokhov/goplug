// Package plug is a mutable HTTP request library
package goplug

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

//--- Types ---

// mutable HTTP Plug object
type Plug struct {
	Uri *url.URL
}

type Result struct {
	Response *http.Response
	Error    error
}

//--- Constructors ---

// create new plug from URL
func NewFromUrl(uri *url.URL) *Plug {
	return &Plug{
		Uri: uri,
	}
}

// create new plug from a string that represents a URI
func New(uristring string) (plug *Plug, err error) {
	uri, err := url.Parse(uristring)
	if err != nil {
		return nil, err
	}
	return &Plug{
		Uri: uri,
	}, nil
}

//--- Functions ---

// Use AtPath to set the full path of the url
func (p *Plug) AtPath(path string) *Plug {
	p.Uri.Path = path
	return p
}

// Use At to add segments to the url
func (p *Plug) At(paths ...string) *Plug {
	p.Uri.Path = p.Uri.Path + strings.Join(paths, "/")
	return p
}

// User With to add a query parameter to the url

func (p *Plug) With(name string, value string) *Plug {
	q := p.Uri.Query()
	q.Set(name, value)
	p.Uri.RawQuery = q.Encode()
	return p
}

// Use WithParams to add a map of query parameters to the url
func (p *Plug) WithParams(params map[string]string) {
	q := p.Uri.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	p.Uri.RawQuery = q.Encode()
}

// GetParam returns the parameter with a certain name
func (p *Plug) GetParam(name string) string {
	return p.Uri.Query().Get(name)
}

// User Without to remove a query parameter
func (p *Plug) Without(name string) *Plug {
	q := p.Uri.Query()
	q.Del(name)
	p.Uri.RawQuery = q.Encode()
	return p
}

// Async GET request
func (p *Plug) Get() (response chan Result) {
	response = make(chan Result)
	go performRequest(p, func() (*http.Response, error) {
		return http.Get(p.Uri.String())
	}, response)
	return response
}

// String representation of the uri
func (p *Plug) String() string {
	return p.Uri.String()
}

// Async POST request
func (p *Plug) Post(reader io.Reader, contentType string) (response chan Result) {
	response = make(chan Result)
	go performRequest(p, func() (*http.Response, error) {
		return http.Post(p.Uri.String(), contentType, reader)
	}, response)
	return response
}

// Async POST with text/plain
func (p *Plug) PostPlainText(contents string) (response chan Result) {
	return p.Post(strings.NewReader(contents), "text/plain")
}

// Handles passing http response and error onto response channel
func performRequest(p *Plug, fetch func() (*http.Response, error), response chan Result) {
	if res, err := fetch(); err != nil {
		response <- Result{
			Response: nil,
			Error:    err}
	} else {
		response <- Result{
			Response: res,
			Error:    nil}
	}
	close(response)
}
