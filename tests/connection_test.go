package tests

import (
	"fmt"
	networking "github.com/RadicalApp/go-networking"
	"github.com/franela/goblin"
	"testing"
)

func TestGetConnection(t *testing.T) {
	urlString := "https://jsonplaceholder.typicode.com/posts"
	params := networking.NewParams()
	connection := networking.NewConnection(urlString, params)

	g := goblin.Goblin(t)
	g.Describe("Get dummy json", func() {
		g.It("Should get dummy json data", func(done goblin.Done) {
			connection.OnReceived = func(response []byte) {
				fmt.Println("Response: !!! ", string(response))
				done()
			}
			connection.OnError = func(err error) {
				t.Error("Error in GET request for url: ", urlString)
				t.Fail()
				done()
			}

			connection.GET()
		})
	})
}

func TestPostConnection(t *testing.T) {
	urlString := "https://jsonplaceholder.typicode.com/posts"
	params := networking.NewParams()
	params.PutString("title", "foo")
	params.PutString("body", "bar")
	params.PutInt("userId", 1)
	connection := networking.NewConnection(urlString, params)

	g := goblin.Goblin(t)
	g.Describe("Post dummy data", func() {
		g.It("Should post dummy data", func(done goblin.Done) {
			connection.OnReceived = func(response []byte) {
				fmt.Println("Response: !!! ", string(response))
				done()
			}

			connection.OnError = func(err error) {
				t.Error("Error in POST request for url: ", urlString)
				t.Fail()
				done()
			}

			connection.POST()
		})
	})
}
