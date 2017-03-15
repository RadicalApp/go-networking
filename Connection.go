package go_networking

import (
	"errors"
	"git.cyberdust.com/radicalapp/go-networking/logger"
	"io/ioutil"
	"time"
)

// Set up logging
var log = logger.Register("git.cyberdust.com/radicalapp/go-networking")

//const POST string = "POST"
//const GET string = "GET"
//const UPLOAD string = "UPLOAD"

type Connection struct {
	url    string
	params Params

	headers map[string]string

	// Basic auth variables
	username    string
	password    string
}

// NewConnection creates a new connection with the appropriate Dust auth headers.
func NewConnection(url string, params Params) *Connection {
	/*
		Creates a new connection with the appropriate Dust auth headers.

		Args:
		  url (string): The URL to make the request to. (includes full url)
		  params ([]string): List of key,value strings to add to the request.
	*/
	conn := Connection{url: url, params: params, headers: map[string]string{}}

	return &conn
}

func (c *Connection) SetBasicAuth(username, password string) {
	log.Debug("Setting basic authentication.")
	c.username = username
	c.password = password
}

func (c *Connection) PutHeader(key, value string) {
	log.Debug("Adding header with key, val: ", key, ", ", value)
	c.headers[key] = value
}

func (c *Connection) GET() ([]byte, error) {
	log.Debug("GET Request")
}

func (c *Connection) POST() ([]byte, error) {
	log.Debug("POST Request")
	return nil, nil
}

func (c *Connection) UPLOAD() ([]byte, error) {
	log.Debug("UPLOAD Request")
	return nil, nil
}

func (c *Connection) send(method string) ([]byte, error) {
	return nil, nil
}


func (c *Connection) sendGet(method string) ([]byte, error) {
	return nil, nil
}

func (c *Connection) sendPost(method string) ([]byte, error) {
	return nil, nil
}


func (c *Connection) sendUpload(method string) ([]byte, error) {
	return nil, nil
}


func (c *Connection) makeRequest(method string) ([]byte, error) {
	/*
		MakeRequest makes an HTTP request to the provided URL.

		Args:
			method (str): GET or POST
		  	url (str): The address to make the request on.
			headers ([]str): List of headers to include with the request.
			params ([]str): List of parameters to include with the request.
	*/
	//log.Debug("Making request:", method)
	//log.Debug("  Headers:")
	//for i, value := range headers {
	//	if i%2 == 1 {
	//		log.Debug("   ", headers[i-1]+":", value)
	//	}
	//}
	//log.Debug("  Params:")
	//for i, value := range params {
	//	if i%2 == 1 {
	//		log.Debug("   ", params[i-1]+":", value)
	//	}
	//}
	//
	//client := &http.Client{
	//	Timeout: 30 * time.Second,
	//}
	//req, err := http.NewRequest(method, url, nil)
	//
	//if len(params)%2 != 0 {
	//	paramErr := errors.New("Must have key and value for params")
	//	log.Error(paramErr)
	//	d.ConnectionState = ConnectionStateDisconnected
	//	//d.OnError(paramErr)
	//	return nil, paramErr
	//}
	//
	//// Set the request parameters
	//values := req.URL.Query()
	//for i := 0; i < len(params); i += 2 {
	//	values.Add(params[i], params[i+1])
	//}
	//req.URL.RawQuery = values.Encode()
	//
	//// Set the headers of the request
	//for i := 0; i < len(headers); i += 2 {
	//	req.Header.Set(headers[i], headers[i+1])
	//}
	//
	//// Make the request
	//d.ConnectionState = ConnectionStateConnecting
	////d.OnStarted()
	//response, err := client.Do(req)
	//if err != nil {
	//	d.ConnectionState = ConnectionStateDisconnected
	//	return nil, err
	//}
	//d.ConnectionState = ConnectionStateConnected
	////d.OnReceived(u.GenerateUUID())
	//
	//// Read the contents of the response
	//defer response.Body.Close()
	//contents, err := ioutil.ReadAll(response.Body)
	//if err != nil {
	//	log.Error(err)
	//	//d.OnError(err)
	//	return nil, err
	//}
	//
	//log.Debug("With URL: ", url)
	//log.Debug("Recieved response code:", response.StatusCode)
	//log.Debug("Response:", string(contents))
	//
	//if response.StatusCode != 200 {
	//	statusErr := errors.New("Recieved non 200 status response.")
	//	//d.OnError(statusErr)
	//	return contents, statusErr
	//}
	//
	////d.OnClosed()

	return nil, nil
}
