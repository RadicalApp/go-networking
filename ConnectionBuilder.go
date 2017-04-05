package go_networking

// ConnectionBuilder is used to help build a `Connection`.
type ConnectionBuilder struct {
	connection Connection
}

// Set url string for a HTTP request
func (cb *ConnectionBuilder) SetUrlString(urlString string) *ConnectionBuilder {
	cb.connection.urlString = urlString
	return cb
}

// Set `Params` for a HTTP request
func (cb *ConnectionBuilder) SetParams(params *Params) *ConnectionBuilder {
	cb.connection.params = *params
	return cb
}

// Set HTTP Method type for a HTTP request
func (cb *ConnectionBuilder) SetMethod(method HTTP_METHOD) *ConnectionBuilder {
	cb.connection.method = method
	return cb
}

// Build a `Connection`
func (cb *ConnectionBuilder) Build() Connection {
	return cb.connection
}
