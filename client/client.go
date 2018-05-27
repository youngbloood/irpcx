package client

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"irpcx"
	"log"
	"strings"
	"sync"

	"github.com/smallnest/rpcx/client"
)

var (
	etcdAddr []string
	once     sync.Once
)

func lazyInit(addr []string) {
	once.Do(
		func() {
			etcdAddr = append(etcdAddr, addr...)
		})
}

type iclient struct {
	discovery client.ServiceDiscovery
	cli       client.XClient
}

var rpcClient = make(map[string]*iclient)

func get(basePath, servicePath string) *iclient {
	if etcdAddr == nil || len(etcdAddr) == 0 {
		log.Fatalln("invoke func client.InitEtcdAddr(etcdAddr []string) to initialize etcd cluster")
	}
	basePath = strings.Trim(basePath, "/")
	basePath = "/" + basePath
	servicePath = strings.Trim(servicePath, "/")
	allPath := basePath + "/" + servicePath
	token := hash(allPath)

	cli, exist := rpcClient[token]
	if !exist {
		discovery := client.NewEtcdDiscovery(basePath, servicePath, etcdAddr, nil)
		xclient := client.NewXClient(servicePath, client.Failover, client.RandomSelect, discovery, client.DefaultOption)
		mc := new(iclient)
		mc.cli = xclient
		mc.discovery = discovery
		rpcClient[token] = mc
		return mc
	}
	return cli
}

func hash(str string) string {
	md := md5.New()
	md.Write([]byte(str)) // 需要加密的字符串为 str
	cipherStr := md.Sum(nil)
	return fmt.Sprintf("%s", hex.EncodeToString(cipherStr)) // 输出加密结果
}

func (mc *iclient) call(method string, args, reply interface{}) error {
	return mc.cli.Call(context.Background(), method, args, reply)
}
func (mc *iclient) gocall(method string, args, reply interface{}) (*client.Call, error) {
	return mc.cli.Go(context.Background(), method, args, reply, nil)
}

// InitEtcdAddr will initialize the etcd cluster
func InitEtcdAddr(etcdAddr []string) {
	lazyInit(etcdAddr)
}

// Call the "IPRCX" service and "Do" method by sync
func Call(req *irpcx.Request) (reply *irpcx.Response, err error) {
	reply = new(irpcx.Response)
	servicePath, method := getParam()
	err = get(req.BasePath, servicePath).call(method, req, reply)
	if err != nil {
		return nil, err
	}
	return
}

// Go will return response by async
func Go(req *irpcx.Request) (call *client.Call, reply *irpcx.Response, err error) {
	reply = new(irpcx.Response)
	servicePath, method := getParam()
	call, err = get(req.BasePath, servicePath).gocall(method, req, reply)
	if err != nil {
		return call, nil, err
	}
	return
}

// Mode mode
type Mode string

const (
	// ModeRelease mean release
	ModeRelease Mode = "release"
	// ModeDebug mean debug
	ModeDebug Mode = "debug"
)

var _mode Mode

// SetMode . set the client mode
func SetMode(mode Mode) {
	_mode = mode
}

// GetMode . return the client mode
func GetMode() Mode {
	return _mode
}

func getParam() (servicePath, method string) {
	if GetMode() == ModeRelease {
		return "IRPCX", "Do"
	}
	return "IRPCX_DEBUG", "Do"
}