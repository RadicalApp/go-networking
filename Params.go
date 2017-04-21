package go_networking

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// `Params` a wrapper struct to hold URL/ File/ Post parameters for HTTP Connection
type Params struct {
	stringParams map[string]string
	intParams    map[string]int
	fileParams   map[string]FileWrapper
	jsonObjects  map[string]interface{}
}

// Constructor for `Params`
func NewParams() *Params {
	return &Params{
		stringParams: map[string]string{},
		intParams:    map[string]int{},
		fileParams:   map[string]FileWrapper{},
		jsonObjects:  make(map[string]interface{}),
	}
}

func QueryParams(stringParams map[string]string) *Params {
	return &Params{
		stringParams: stringParams,
		intParams:    map[string]int{},
		fileParams:   map[string]FileWrapper{},
		jsonObjects:  make(map[string]interface{}),
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

// Put an object which will be converted to JSON object while `POST` request.
func (p *Params) PutObject(key string, value interface{}) {
	p.jsonObjects[key] = value
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

func (p *Params) Encoded() string {
	data := url.Values{}
	for key, val := range p.jsonObjects {
		bytes123, _ := json.Marshal(val)
		data.Set(key, string(bytes123))
	}
	for key, val := range p.stringParams {
		data.Set(key, val)
	}
	for key, val := range p.intParams {
		data.Set(key, strconv.Itoa(val))
	}
	return data.Encode()
}
