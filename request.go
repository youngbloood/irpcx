package irpcx

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Request fo irpcx
type Request struct {
	Body   []byte
	Path   string
	Host   string
	IP     string
	Header http.Header
	Params gin.Params
	// the base path in etcd
	BasePath string
	// the service path in etcd
	Service string
	// the method in etcd
	Method string
}

// NewRequest with service and function
func NewRequest(basePath, service, method string) *Request {

	basePath = "/" + strings.Trim(basePath, "/")
	service = strings.Trim(service, "/")

	host, _ := os.Hostname()
	return &Request{Host: host, BasePath: basePath, Service: service, Method: method}
}

// Marshal will encode the obj to Request.Body
func (req *Request) Marshal(obj interface{}) (err error) {
	req.Body, err = json.Marshal(obj)
	return
}
