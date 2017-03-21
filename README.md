[![Latest Tag](https://img.shields.io/badge/tag-0.1.1-green.svg?style=flat)](https://github.com/RadicalApp/go-networking/releases/tag/0.1.1)
[![Twitter](https://img.shields.io/badge/twitter-@DustMessaging-blue.svg?style=flat)](https://twitter.com/dustmessaging)

# go-networking

`go-networking` is a networking library. It's built on top of go's [http](https://golang.org/pkg/net/http/) package and designed to create a common interface for simple RESTful services. 


## Installation

```go
import "github.com/RadicalApp/go-networking"
```


## Usage

Simplest way to use `go-networking` create an instance of `connection` and call `GET()`, `POST()` on it.

### Connection

#### GET

```go
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
```

### Parameters

#### Params

`Params` is a common interface for all types of parameters. `int`, `string` or `file`.

> Use `FileWrapperBuilder` class to initialize a file parameter using file url or file bytes.


## Credits

`go-networking` is owned and maintained by the [Radical App, LLC](https://usedust.com/)

`go-networking` was created in order to cater towards the development of [Dust App](https://usedust.com/) which uses RESTful service for all activities.

Main contributors are:
- [Rohit Kotian](rohit@usedust.com)
- [Elliot Sperling](elliot@usedust.com)

Thanks for [JSONPlaceholder](https://jsonplaceholder.typicode.com/) for having test APIs.

If you wish to contribute to the project please follow the [guidelines](CONTRIBUTING.md).


## License

`go-networking` is released under the MIT license. See [LICENSE](LICENSE) for details.


