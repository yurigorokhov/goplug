// Package goplug is a mutable HTTP request library
package goplug

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

//--- Types ---

// Plug object
type Plug struct {
	Uri     *url.URL
	headers map[string]string
}

// Result object
type Result struct {
	Response *http.Response
	Error    error
}

//--- Constructors ---

// Create new plug from URL
func NewFromUrl(uri *url.URL) *Plug {
	return &Plug{
		Uri:     uri,
		headers: make(map[string]string),
	}
}

// Create new plug from a string that represents a URI
func New(uristring string) (plug *Plug, err error) {
	uri, err := url.Parse(uristring)
	if err != nil {
		return nil, err
	}
	return &Plug{
		Uri:     uri,
		headers: make(map[string]string),
	}, nil
}

//--- Functions ---

// Use AtPath to set the full path of the url
func (p *Plug) AtPath(path string) *Plug {
	if !strings.HasSuffix(p.Uri.Path, "/") {
		p.Uri.Path = p.Uri.Path + "/"
	}
	p.Uri.Path = path
	return p
}

// Use At to add segments to the url
func (p *Plug) At(paths ...string) *Plug {
	if !strings.HasSuffix(p.Uri.Path, "/") {
		p.Uri.Path = p.Uri.Path + "/"
	}
	p.Uri.Path = p.Uri.Path + strings.Join(paths, "/")
	return p
}

// Use WithUserPassword to add an authorization header
func (p *Plug) WithUser(username string) *Plug {
	p.Uri.User = url.User(username)
	return p
}

// Use WithUserPassword to add an authorization header
func (p *Plug) WithUserPassword(username, password string) *Plug {
	p.Uri.User = url.UserPassword(username, password)
	return p
}

// Use With to add a query parameter to the url
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

// WithHeader sets the request header
func (p *Plug) WithHeader(name string, value string) *Plug {
	p.headers[name] = value
	return p
}

// Async HEAD request
func (p *Plug) Head() (response chan Result) {
	response = make(chan Result)
	go performRequest(p, func() (*http.Response, error) {
		if req, err := prepareRequest(p, "HEAD", nil); err != nil {
			return nil, err
		} else {
			return http.DefaultClient.Do(req)
		}
	}, response)
	return response
}

// Async GET request
func (p *Plug) Get() (response chan Result) {
	response = make(chan Result)
	go performRequest(p, func() (*http.Response, error) {
		if req, err := prepareRequest(p, "GET", nil); err != nil {
			return nil, err
		} else {
			return http.DefaultClient.Do(req)
		}
	}, response)
	return response
}

// String representation of the uri
func (p *Plug) String() string {
	return p.Uri.String()
}

// Async DELETE request with no body
func (p *Plug) Delete() (response chan Result) {
	response = make(chan Result)
	go performRequest(p, func() (*http.Response, error) {
		if req, err := prepareRequest(p, "DELETE", nil); err != nil {
			return nil, err
		} else {
			return http.DefaultClient.Do(req)
		}
	}, response)
	return response
}

// Async DELETE request with a body in the request
func (p *Plug) DeleteWithBody(reader io.Reader, contentType string) (response chan Result) {
	response = make(chan Result)
	go performRequest(p, func() (*http.Response, error) {
		if req, err := prepareRequest(p, "DELETE", &reader); err != nil {
			return nil, err
		} else {
			return http.DefaultClient.Do(req)
		}
	}, response)
	return response
}

// Async PUT request
func (p *Plug) Put(reader io.Reader, contentType string) (response chan Result) {
	response = make(chan Result)
	go performRequest(p, func() (*http.Response, error) {
		if req, err := prepareRequest(p, "PUT", &reader); err != nil {
			return nil, err
		} else {
			return http.DefaultClient.Do(req)
		}
	}, response)
	return response
}

// Post request
func (p *Plug) Post(reader io.Reader, contentType string) (response chan Result) {
	response = make(chan Result)
	go performRequest(p, func() (*http.Response, error) {
		if req, err := prepareRequest(p, "POST", &reader); err != nil {
			return nil, err
		} else {
			return http.DefaultClient.Do(req)
		}
	}, response)
	return response
}

// Clone this plug and return the copy
func (p *Plug) Clone() *Plug {

	var userInfo *url.Userinfo = nil
	if p.Uri.User != nil {
		if pass, ok := p.Uri.User.Password(); ok {
			userInfo = url.UserPassword(p.Uri.User.Username(), pass)
		} else {
			userInfo = url.User(p.Uri.User.Username())
		}
	}
	newHeaders := make(map[string]string)
	for k, v := range p.headers {
		newHeaders[k] = v
	}
	return &Plug{
		Uri: &url.URL{
			Scheme:   p.Uri.Scheme,
			Opaque:   p.Uri.Opaque,
			User:     userInfo,
			Host:     p.Uri.Host,
			Path:     p.Uri.Path,
			RawQuery: p.Uri.RawQuery,
			Fragment: p.Uri.Fragment,
		},
		headers: newHeaders,
	}
}

// prepareRequest creates a request object ased on the verb, uri and headers
func prepareRequest(p *Plug, verb string, r *io.Reader) (*http.Request, error) {
	if req, err := http.NewRequest(verb, p.Uri.String(), nil); err != nil {
		return nil, err
	} else {
		for k, v := range p.headers {
			req.Header.Add(k, v)
		}
		return req, nil
	}
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
