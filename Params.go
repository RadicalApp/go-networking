package go_networking

import (
	"net/url"
	"strconv"
)

// `Params` a wrapper struct to hold URL/ File/ Post parameters for HTTP Connection
type Params struct {
	stringParams map[string]string
	intParams    map[string]int
	fileParams   map[string]FileWrapper
}

// Constructor for `Params`
func NewParams() *Params {
	return &Params{
		stringParams: map[string]string{},
		intParams:    map[string]int{},
		fileParams:   map[string]FileWrapper{},
	}
}

func QueryParams(stringParams map[string]string) *Params {
	return &Params{
		stringParams: stringParams,
		intParams:    map[string]int{},
		fileParams:   map[string]FileWrapper{},
	}
}

// PutString add a key/val pair to the RequestParams, where value is a string.
func (p *Params) PutString(key, value string) {
	p.stringParams[key] = value
}

// PutInt add a key/val pair to the RequestParams, where value is an int.
func (p *Params) PutInt(key string, value int) {
	p.intParams[key] = value
}

// PutFile add a key/val pair to the RequestParams, where value is a FileWrapper.
func (p *Params) PutFile(key string, value FileWrapper) {
	p.fileParams[key] = value
}

// Internal method to get url parameters for HTTP request
func (p *Params) urlParameters() map[string]string {
	params := p.stringParams
	for key, val := range p.intParams {
		params[key] = strconv.Itoa(val)
	}
	return params
}

// Internal method to get url encoded values for HTTP request
func (p *Params) urlEncodeValues() string {
	values := url.Values{}
	for key, value := range p.stringParams {
		values.Add(key, value)
	}
	for key, value := range p.intParams {
		values.Add(key, strconv.Itoa(value))
	}
	return values.Encode()
}
