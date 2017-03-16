package go_networking

type ConnectionBuilder struct {

	connection *Connection

}

//func NewConnectionBuilder(urlString string, method)

func (cb *ConnectionBuilder) setUrlString(urlString string) {
	cb.connection.urlString = urlString
}

func (cb *ConnectionBuilder) build() *Connection {
	return cb.connection
}