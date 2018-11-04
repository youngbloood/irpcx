package irpcx

import (
	"encoding/json"
	// "github.com/gin-gonic/gin/binding"
	// "strings"
)

// Context of irpcx
type Context struct {
	Req  *Request
	Resp *Response
}

// SetResp . set the ctx's response
func (c *Context) SetResp(obj interface{}) ( error){
	if c.Resp == nil {
		return nil
	}
	return 	c.Resp.Marshal(obj)
}

func (c *Context) Bind(obj interface{}) error {


	return json.Unmarshal(c.Req.Body,obj)

	// bind := c.getMIME()
	// return c.bind(obj, bind)
}

// func (c *Context) getMIME() binding.Binding {
// 	contentTypes := c.Req.Header.Get("Content-Type")
// 	contentType := strings.Split(contentTypes, ";")
// 	if len(contentTypes) > 0 {
// 		return binding.Default(c.Req.Method,contentType[0])
// 	}
// 	return binding.JSON
// }

// func (c *Context) bind(obj interface{}, b binding.Binding) error {
// 	return b.Bind(c.Req, obj)
// }

