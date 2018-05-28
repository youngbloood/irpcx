package irpcx

// Context of irpcx
type Context struct {
	Req  *Request
	Resp *Response
}

// SetResp . set the ctx's response
func (c *Context) SetResp(obj interface{}) {
	if c.Resp == nil {
		return
	}
	c.Resp.Data = obj
}

// TODO: Bind()
// func (c *Context) Bind(obj interface{}) error {
// 	mime := c.getMIME()
// 	return c.bind(obj, mime)
// }

// func (c *Context) getMIME() binding.MIME {
// 	mimeStr := c.Req.Header.Get("Content-Type")
// 	mimeSlice := strings.Split(mimeStr, ";")
// 	if len(mimeSlice) > 0 {
// 		return binding.MIME(mimeSlice[0])
// 	}
// 	return binding.MIMEJSON
// }

// func (c *Context) bind(obj interface{}, b binding.Binding) error {
// 	return b.Bind(c.Req.Body, obj)
// }
