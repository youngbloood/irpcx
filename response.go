package irpcx

import (
	"encoding/json"
	"net/http"
)

// Response of irpcx
type Response struct {
	Body   []byte
	Data   interface{}
	Status int
	Writer http.ResponseWriter
}

// Unmarshal will decode the Response.Body to obj
func (resp *Response) Unmarshal(obj interface{}) error {
	return json.Unmarshal(resp.Body, obj)
}

// Marshal will encode the obj to Response.Body
func (resp *Response) Marshal(obj interface{}) (err error) {
	resp.Body, err = json.Marshal(obj)
	return
}
