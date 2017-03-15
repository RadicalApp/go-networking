package go_networking

type Params struct {
	stringParams map[string]string
	intParams    map[string]int
	fileParams   map[string]fileWrapper
}

func NewParams() Params {
	return Params{
		stringParams: map[string]string{},
		intParams:    map[string]int{},
		fileParams:   map[string]fileWrapper{},
	}
}

// PutString add a key/val pair to the RequestParams, where value is a string.
func (rp *Params) PutString(key, value string) {
	rp.stringParams[key] = value
}

// PutInt add a key/val pair to the RequestParams, where value is an int.
func (rp *Params) PutInt(key string, value int) {
	rp.intParams[key] = value
}

// PutFile add a key/val pair to the RequestParams, where value is a FileWrapper.
func (rp *Params) PutFile(key string, value fileWrapper) {
	rp.fileParams[key] = value
}
