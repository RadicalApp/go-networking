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

	connection.GET()
}

//func TestPostConnection(t *testing.T) {
//
//}

//func TestUploadConnection(t *testing.T) {
//
//}


//connection := networking.NewConnection(url).WithParams(params).WithCompletion(completion).WithOnReveived(received).GET()

