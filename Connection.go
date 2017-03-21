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

type ConnectionState int

const (
	CONNECTION_STATE_DISCONNECTED ConnectionState = 0
	CONNECTION_STATE_CONNECTING   ConnectionState = 1
	CONNECTION_STATE_CONNECTED    ConnectionState = 2
	// MAYBE IN FUTURE...?
	//CONNECTION_STATE_RECONNECTING ConnectionState = 3
	//CONNECTION_STATE_FAILED ConnectionState = 4
)

func (c ConnectionState) ToString() string {
	switch c {
	case CONNECTION_STATE_DISCONNECTED:
		return "Disconnected"
	case CONNECTION_STATE_CONNECTING:
		return "Connecting"
	case CONNECTION_STATE_CONNECTED:
		return "Connected"
	default:
		return "INVALID STATE"
	}
}

// Connection States
type OnStarted func()
type OnReceived func([]byte)
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

type IdealResponse interface {
	getIdealResponseTemplate() interface{}
}

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

func (c *Connection) SetBasicAuth(username string, password string) {
	log.Debug("Setting basic authentication.")
	c.basicAuthorization = basicAuthorization{username: username, password: password}
}

func (c *Connection) SetMethod(method HTTP_METHOD) {
	log.Debug("Setting method to :", string(method))
	c.method = method
}

func (c *Connection) SetTimeout(timeoutInSeconds time.Duration) {
	log.Debug("Setting timeout in seconds: ", timeoutInSeconds)
	c.timeoutInSeconds = timeoutInSeconds
}

func (c *Connection) PutHeader(key, value string) {
	log.Debug("Adding header with key, val: ", key, ", ", value)
	c.headers[key] = value
}

func (c *Connection) SetNumberOfRetries(number int) {
	log.Debug("Setting number of retries: ", number)
	c.numberOfRetries = number
}

func (c *Connection) GET() {
	c.Get(nil)
}

func (c *Connection) Get(completion func([]byte, error)) {
	log.Debug("GET")
	c.method = HTTP_METHOD_GET
	c.makeRequest(completion)
}

func (c *Connection) POST() {
	c.Post(nil)
}

func (c *Connection) Post(completion func([]byte, error)) {
	c.method = HTTP_METHOD_POST
	c.makeRequest(completion)
}

func (c *Connection) UPLOAD() {
	c.Upload(nil)
}

func (c *Connection) Upload(completion func([]byte, error)) {
	c.method = HTTP_METHOD_UPLOAD
	c.makeRequest(completion)
}

func (c *Connection) makeRequest(completion func([]byte, error)) {
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
	c.makeParams(req)
	if err != nil {
		if c.OnError != nil {
			c.OnError(err)
		}
		if completion != nil {
			completion(nil, err)
		}
	} else {
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

func (c *Connection) doRequest(req *http.Request, completion func([]byte, error)) {
	log.Debug("Doing request")
	client := &http.Client{
		Timeout: c.timeoutInSeconds * time.Second,
	}

	c.changeState(CONNECTION_STATE_CONNECTING, CONNECTION_STATE_CONNECTED)
	response, err := client.Do(req)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		c.processError(err, completion)
	} else {
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			c.processError(err, completion)
		} else {
			c.processResponse(contents, completion)
		}
	}
}

func (c *Connection) changeState(from, to ConnectionState) {
	// MAYBE GUARANTEE CORRECT STATE CHANGE, IN FUTURE
	log.Debug("Changing state from: ", from.ToString(), " to: ", to.ToString())
	c.state = to
	if c.OnStateChanged != nil {
		c.OnStateChanged(to)
	}
}

func (c *Connection) processError(err error, completion func([]byte, error)) {
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
			completion(nil, err)
		}
	}
}

func (c *Connection) processResponse(response []byte, completion func([]byte, error)) {
	log.Debug("Processing response")
	c.changeState(CONNECTION_STATE_CONNECTED, CONNECTION_STATE_DISCONNECTED)
	if c.OnReceived != nil {
		c.OnReceived(response)
	}
	if c.OnClosed != nil {
		c.OnClosed()
	}
	if completion != nil {
		completion(response, nil)
	}
}
