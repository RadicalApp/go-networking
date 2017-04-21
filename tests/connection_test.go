package tests

import (
	networking "github.com/RadicalApp/go-networking"
	"github.com/franela/goblin"
	"testing"
	"encoding/json"
)

func TestGetConnection(t *testing.T) {
	urlString := "https://jsonplaceholder.typicode.com/posts"
	params := networking.NewParams()
	connection := networking.NewConnection(urlString, params)

	g := goblin.Goblin(t)
	g.Describe("Get dummy json", func() {
		g.It("Should get dummy json data", func(done goblin.Done) {
			connection.OnReceived = func(response networking.Response) {
				t.Log("Response: !!! ", string(response.Data))
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
			connection.OnReceived = func(response networking.Response) {
				t.Log("Response: !!! ", string(response.Data))
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

func TestPostConnection2(t *testing.T) {
	type UploadContent struct {
		text string
	}
	c := UploadContent{text: "hello, world"}
	uploadContentBytes, err := json.Marshal(c)
	if err != nil {
		t.Error("Failed to marshal")
		t.Fail()
	}

	urlString := "https://jsonplaceholder.typicode.com/posts"
	params := networking.NewParams()
	params.PutString("content", string(uploadContentBytes))
	connection := networking.NewConnection(urlString, params)
	g := goblin.Goblin(t)
	g.Describe("Post dummy data", func() {
		g.It("Should post dummy data", func(done goblin.Done) {
			connection.OnReceived = func(response networking.Response) {
				t.Log("Response: !!! ", string(response.Data))
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

func TestFailConnection(t *testing.T) {
	urlString := "https://foo.bar"
	params := networking.NewParams()
	connection := networking.NewConnection(urlString, params)

	g := goblin.Goblin(t)
	g.Describe("Hit fake url", func() {
		g.It("Should fail to connect to fake url", func(done goblin.Done) {
			connection.OnReceived = func(response networking.Response) {
				t.Error("Incorrectly passed...", urlString)
				t.Fail()
				done()
			}
			connection.OnError = func(err error) {
				t.Log("Successfully failed!")
				done()
			}
			connection.GET()
		})
	})
}
