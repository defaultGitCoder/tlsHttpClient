package tlsHttpClient

import (
	"github.com/quotpw/tlsHttpClient/tlsHttpClient/cycletls"
)

type Response struct {
	Bytes      []byte
	Text       string
	json       map[string]any
	StatusCode int
	Headers    map[string]string
	Cookies    []cycletls.Cookie
}

func (r *Response) Json() map[string]any {
	if r.json == nil {
		r.json = map[string]any{}
		Unmarshal(r.Bytes, &r.json)
	}
	return r.json
}
func (r *Response) ToStruct(v interface{}) (interface{}, error) {
	err := Unmarshal(r.Bytes, &v)
	return v, err
}
