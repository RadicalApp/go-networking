package go_networking

import (
	"strconv"
	"net/url"
)

type Params struct {
	stringParams map[string]string
	intParams    map[string]int
	fileParams   map[string]fileWrapper
}

func NewParams() Params {
	return Params{
		stringParams: map[string]string{},
		intParams:    map[string]int{},
		fileParams:   map[string]fileWrapper{},
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
func (p *Params) PutFile(key string, value fileWrapper) {
	p.fileParams[key] = value
}

func (p *Params) urlParameters() map[string]string {
	params := p.stringParams
	for key, val := range p.intParams {
		params[key] = strconv.Itoa(val)
	}
	return params
}

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
