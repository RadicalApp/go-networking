package tests

import (
	"fmt"
	networking "git.cyberdust.com/radicalapp/go-networking"
	"testing"
)

func TestGetConnection(t *testing.T) {
	urlString := "https://jsonplaceholder.typicode.com/posts"
	params := networking.NewParams()
	connection := networking.NewConnection(urlString, params)

	connection.OnReceived = func(response []byte) {
		fmt.Println("Response: !!! ", string(response))
	}
	connection.OnError = func (err error) {
		t.Error("Error in GET request for url: ", urlString)
		t.Fail()
	}

	connection.GET()
}

func TestPostConnection(t *testing.T) {
	urlString := "https://jsonplaceholder.typicode.com/posts"

	params := networking.NewParams()
	params.PutString("title", "foo")
	params.PutString("body", "bar")
	params.PutInt("userId", 1)

	connection := networking.NewConnection(urlString, params)
	connection.OnReceived = func (response []byte) {
		fmt.Println("Response: !!! ", string(response))
	}
	connection.OnError = func (err error) {
		t.Error("Error in POST request for url: ", urlString)
		t.Fail()
	}
	connection.POST()

}
