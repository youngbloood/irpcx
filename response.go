package irpcx

import (
	"encoding/json"
	"net/http"
)

// Response of irpcx
type Response struct {
	body   []byte
	data   interface{}
	Status int
	Writer http.ResponseWriter
}

// Unmarshal will decode the Response.Body to obj
func (resp *Response) Unmarshal(obj interface{}) error {
	return json.Unmarshal(resp.body, obj)
}

// Marshal will encode the obj to Response.Body
func (resp *Response) Marshal(obj interface{}) (err error) {
	resp.body, err = json.Marshal(obj)
	return
}

// Encode the data to body
func (resp *Response) Encode() (err error) {
	resp.body, err = json.Marshal(resp.data)
	return
}

// Body return resp's body
func (resp *Response) Body() []byte {
	return resp.body
}

// Data return resp's data
func (resp *Response) Data() interface{} {
	return resp.data
}
