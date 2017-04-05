package go_networking

import (
	"bytes"
	"encoding/base64"
	"git.cyberdust.com/radicalapp/go-networking/logger"
	"io/ioutil"
	"net/http"
	"time"
)

// Set up logging
var log = logger.Register("git.cyberdust.com/radicalapp/go-networking")

type HTTP_METHOD string

const (
	HTTP_METHOD_GET    HTTP_METHOD = "GET"
	HTTP_METHOD_POST   HTTP_METHOD = "POST"
	HTTP_METHOD_UPLOAD HTTP_METHOD = "UPLOAD"
)

type ConnectionState string

const (
	CONNECTION_STATE_DISCONNECTED ConnectionState = "Disconnected"
	CONNECTION_STATE_CONNECTING   ConnectionState = "Connecting"
	CONNECTION_STATE_CONNECTED    ConnectionState = "Connected"
	// MAYBE IN FUTURE...?
	//CONNECTION_STATE_RECONNECTING ConnectionState = 3
	//CONNECTION_STATE_FAILED ConnectionState = 4
)

// `Response` a struct that defines the response for a particular connection
type Response struct {
	Data     []byte
	Request  http.Request
	Response http.Response
}

func newResponse(data []byte, request http.Request, response http.Response) *Response {
	return &Response{Data: data, Request: request, Response: response}
}

func emptyResponse() *Response {
	return &Response{}
}

// Check if a `Response` is empty
func (r *Response) IsEmpty() bool {
	return r == emptyResponse()
}

// Connection States
type OnStarted func()
type OnReceived func(response Response)
type OnClosed func()
type OnError func(err error)
type OnProgress func(progress int)
type OnStateChanged func(state ConnectionState)

type basicAuthorization struct {
	username string
	password string
}

func (b *basicAuthorization) isEmpty() bool {
	return (len(b.username) == 0) && (len(b.password) == 0)
}

// Future implementation for offline usage
type IdealResponse interface {
	getIdealResponseTemplate() interface{}
}

// `Connection` a struct to help connect to the real world
type Connection struct {
	urlString        string
	params           Params
	timeoutInSeconds time.Duration
	method           HTTP_METHOD
	state            ConnectionState
	numberOfRetries  int
	headers          map[string]string

	// Basic auth variables
	basicAuthorization basicAuthorization

	// Functions to handle connection states
	OnStarted      OnStarted
	OnReceived     OnReceived
	OnClosed       OnClosed
	OnError        OnError
	OnProgress     OnProgress
	OnStateChanged OnStateChanged
	//
	//idealResponse      IdealResponse
	//guaranteeExecution bool
}

// NewConnection creates a new connection with the appropriate Dust auth headers.
func NewConnection(urlString string, params Params) *Connection {
	/*
		Creates a new connection with the appropriate Dust auth headers.

		Args:
		  url (string): The URL to make the request to. (includes full url)
		  params ([]string): List of key,value strings to add to the request.
	*/
	conn := Connection{urlString: urlString, params: params, timeoutInSeconds: 30, headers: map[string]string{}}

	return &conn
}

// Add basic authorization to the headers
func (c *Connection) SetBasicAuth(username string, password string) {
	log.Debug("Setting basic authentication.")
	c.basicAuthorization = basicAuthorization{username: username, password: password}
}

// Set method type. E.g.: GET, POST
func (c *Connection) SetMethod(method HTTP_METHOD) {
	log.Debug("Setting method to :", string(method))
	c.method = method
}

// Set timeout for the connection
func (c *Connection) SetTimeout(timeoutInSeconds time.Duration) {
	log.Debug("Setting timeout in seconds: ", timeoutInSeconds)
	c.timeoutInSeconds = timeoutInSeconds
}

// Add header, like Authorization or custom
func (c *Connection) PutHeader(key, value string) {
	log.Debug("Adding header with key, val: ", key, ", ", value)
	c.headers[key] = value
}

// Set number of retries for a HTTP request
func (c *Connection) SetNumberOfRetries(number int) {
	log.Debug("Setting number of retries: ", number)
	c.numberOfRetries = number
}

// Helper to send GET request
func (c *Connection) GET() {
	c.Get(nil)
}

// Helper to send GET request with completion
func (c *Connection) Get(completion func(Response, error)) {
	log.Debug("GET")
	c.method = HTTP_METHOD_GET
	c.makeRequest(completion)
}

// Helper to send POST request
func (c *Connection) POST() {
	c.Post(nil)
}

// Helper to send POST request with completion
func (c *Connection) Post(completion func(Response, error)) {
	c.method = HTTP_METHOD_POST
	c.makeRequest(completion)
}

// Helper to send UPLOAD (type of POST request) request
// This should be used while POST'ing MULTIPART-FORM
func (c *Connection) UPLOAD() {
	c.Upload(nil)
}

// Helper to send UPLOAD (type of POST request) request with completion
// This should be used while POST'ing MULTIPART-FORM
func (c *Connection) Upload(completion func(Response, error)) {
	c.method = HTTP_METHOD_UPLOAD
	c.makeRequest(completion)
}

func (c *Connection) makeRequest(completion func(Response, error)) {
	//if c.guaranteeExecution {
	//	if c.idealResponse == nil {
	//		err := errors.New("Ideal response not configured")
	//		c.processError(err, completion)
	//		return
	//	}
	//}
	if c.OnStarted != nil {
		c.OnStarted()
	}
	c.changeState(CONNECTION_STATE_DISCONNECTED, CONNECTION_STATE_CONNECTING)

	body := c.makeBody()

	log.Debug("Method ", string(c.method))

	req, err := http.NewRequest(string(c.method), c.urlString, body)
	if err != nil {
		c.processError(err, completion)
	} else {
		c.makeParams(req)
		// Set the headers of the request
		if !c.basicAuthorization.isEmpty() {
			encoded := base64.StdEncoding.EncodeToString([]byte(c.basicAuthorization.username + ":" + c.basicAuthorization.password))
			c.PutHeader("Authorization", "Basic "+encoded)

		}
		for key, val := range c.headers {
			req.Header.Set(key, val)
		}
		go c.doRequest(req, completion)
	}
}

func (c *Connection) makeParams(req *http.Request) {
	if c.method != HTTP_METHOD_UPLOAD {
		req.URL.RawQuery = c.params.urlEncodeValues()
	}
}

func (c *Connection) makeBody() *bytes.Buffer {
	body := &bytes.Buffer{}
	if c.method == HTTP_METHOD_UPLOAD {
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
	}
	return body
}

func (c *Connection) doRequest(req *http.Request, completion func(Response, error)) {
	log.Debug("Doing request")
	client := &http.Client{
		Timeout: c.timeoutInSeconds * time.Second,
	}

	c.changeState(CONNECTION_STATE_CONNECTING, CONNECTION_STATE_CONNECTED)
	response, err := client.Do(req)
	if err != nil {
		c.processError(err, completion)
	} else {
		c.processResponse(req, response, completion)
	}
}

func (c *Connection) changeState(from, to ConnectionState) {
	// MAYBE GUARANTEE CORRECT STATE CHANGE, IN FUTURE
	log.Debug("Changing state from: ", string(from), " to: ", string(to))
	c.state = to
	if c.OnStateChanged != nil {
		c.OnStateChanged(to)
	}
}

func (c *Connection) processError(err error, completion func(Response, error)) {
	c.changeState(CONNECTION_STATE_CONNECTED, CONNECTION_STATE_DISCONNECTED)
	if c.numberOfRetries > 0 {
		c.numberOfRetries--
		c.makeRequest(completion)
	} else {
		log.Error(err)
		if c.OnError != nil {
			c.OnError(err)
		}
		if completion != nil {
			completion(*emptyResponse(), err)
		}
		if c.OnClosed != nil {
			c.OnClosed()
		}
	}
}

func (c *Connection) processResponse(req *http.Request, res *http.Response, completion func(Response, error)) {
	log.Debug("Processing response")
	c.changeState(CONNECTION_STATE_CONNECTED, CONNECTION_STATE_DISCONNECTED)

	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.processError(err, completion)
	} else {
		response := newResponse(contents, *req, *res)

		if c.OnReceived != nil {
			c.OnReceived(*response)
		}
		if completion != nil {
			completion(*response, nil)
		}
		if c.OnClosed != nil {
			c.OnClosed()
		}
	}
}
