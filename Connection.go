package go_networking

import (
	"encoding/base64"
	"git.cyberdust.com/radicalapp/go-networking/logger"
	"io/ioutil"
	"time"
	"net/http"
	"io"
	"bytes"
)

// Set up logging
var log = logger.Register("git.cyberdust.com/radicalapp/go-networking")

const http_get_method string = "GET"
const http_post_method string = "POST"
const http_upload_method string = "UPLOAD"

type Completion func([]byte, error)

// Connection States
type OnStarted func()
type OnReceived func([]byte)
type OnClosed func()
type OnError func(err error)
type OnProgress func(progress int)

type basicAuthorization struct {
	username string
	password string
}

type Connection struct {
	url              string
	params           Params
	timeoutInSeconds time.Duration

	headers map[string]string

	// Basic auth variables
	basicAuthorization basicAuthorization

	// Functions to handle connection states
	OnStarted  OnStarted
	OnReceived OnReceived
	OnClosed   OnClosed
	OnError    OnError
	OnProgress OnProgress
}

// NewConnection creates a new connection with the appropriate Dust auth headers.
func NewConnection(url string, params Params) *Connection {
	/*
		Creates a new connection with the appropriate Dust auth headers.

		Args:
		  url (string): The URL to make the request to. (includes full url)
		  params ([]string): List of key,value strings to add to the request.
	*/
	conn := Connection{url: url, params: params, timeoutInSeconds: 30, headers: map[string]string{}}

	return &conn
}

func (c *Connection) SetBasicAuth(username string, password string) {
	log.Debug("Setting basic authentication.")
	c.basicAuthorization = basicAuthorization{username: username, password: password}
}

func (c *Connection) SetTimeout(timeoutInSeconds time.Duration) {
	log.Debug("Setting timeout in seconds: ", timeoutInSeconds)
	c.timeoutInSeconds = timeoutInSeconds
}

func (c *Connection) PutHeader(key, value string) {
	log.Debug("Adding header with key, val: ", key, ", ", value)
	c.headers[key] = value
}

func (c *Connection) GET(completion Completion) {
	c.makeRequest(http_get_method, completion)
}

func (c *Connection) POST(completion Completion) {
	c.makeRequest(http_post_method, completion)
}

func (c *Connection) UPLOAD(completion Completion) {
	c.makeRequest(http_upload_method, completion)
}

func (c *Connection) makeRequest(method string, completion Completion) {
	c.OnStarted()
	body := c.makeBody(method)

	req, err := http.NewRequest(method, c.url, body)
	c.makeParams(method, req)
	if err != nil {
		c.OnError(err)
		completion(nil, err)

	} else {
		// Set the headers of the request
		if &c.basicAuthorization != nil {
			encoded := base64.StdEncoding.EncodeToString([]byte(c.basicAuthorization.username + ":" + c.basicAuthorization.password))
			c.PutHeader("Authorization", "Bearer " + encoded)
		}
		for key, val := range c.headers {
			req.Header.Set(key, val)
		}
		go c.doRequest(req, completion)
	}
}

func (c *Connection) makeParams(method string, req *http.Request) {
	if method != http_upload_method {
		req.URL.RawQuery = c.params.urlEncodeValues()
	}
}

func (c *Connection) makeBody(method string) *bytes.Buffer {
	if method == http_upload_method {
		body := &bytes.Buffer{}
		writer := newByteWriter(body)
		defer writer.close()

		for key, val := range c.params.urlParameters() {
			_ = writer.w.WriteField(key, val)
		}

		isMultipart := false
		for key, val := range c.params.fileParams {
			if val.path == "" {
				_ = writer.writeByteField(key, val.data, val.name)
			} else {
				_ = writer.writeFileField(key, val.path, val.name)
			}
			isMultipart = true
		}

		if isMultipart {
			c.PutHeader("Content-Type", writer.w.FormDataContentType())
		}

		return body
	}
	return nil
}

func (c *Connection) doRequest(req *http.Request, completion Completion) ([]byte, error) {
	client := &http.Client{
		Timeout: c.timeoutInSeconds * time.Second,
	}

	response, err := client.Do(req)

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		c.OnError(err)
		completion(nil, err)
		return nil, err
	} else {
		c.OnReceived(contents)
		c.OnClosed()
		return contents, nil
	}
}
