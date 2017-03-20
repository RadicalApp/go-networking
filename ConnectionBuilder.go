package go_networking

type ConnectionBuilder struct {
	connection Connection
}

func (cb *ConnectionBuilder) SetUrlString(urlString string) *ConnectionBuilder {
	cb.connection.urlString = urlString
	return cb
}

func (cb *ConnectionBuilder) SetParams(params Params) *ConnectionBuilder {
	cb.connection.params = params
	return cb
}

func (cb *ConnectionBuilder) SetMethod(method HTTP_METHOD) *ConnectionBuilder {
	cb.connection.method = method
	return cb
}

func (cb *ConnectionBuilder) Build() Connection {
	return cb.connection
}
