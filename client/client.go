package client

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/smallnest/rpcx/client"
	"github.com/youngbloood/irpcx"
)

// Mode mode
type Mode string

const (
	// ModeRelease mean release
	ModeRelease Mode = "release"
	// ModeDebug mean debug
	ModeDebug Mode = "debug"
)

/**   MODE   */

var _mode Mode

// SetMode . set the client mode
func SetMode(mode Mode) {
	if len(irpcxCli.cli) > 0 {
		return
	}
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

var (
	etcdAddr []string
	once     sync.Once
)

func lazyInit(addr []string) {
	once.Do(
		func() {
			etcdAddr = append(etcdAddr, addr...)
			irpcxCli.cli = make(map[string]*iclient, 100)
		})
}

/** IPRC-Client */

type iclient struct {
	discovery client.ServiceDiscovery
	cli       client.XClient
}

var rpcClient = make(map[string]*iclient)

type irpcxClientMap map[string]*iclient

type irpcxClient struct {
	cli irpcxClientMap
	// Hash is the func to generate basePath+servicePath hash value
	hash func(string) string
}

var irpcxCli irpcxClient

func get(basePath, servicePath string) *iclient {
	if etcdAddr == nil || len(etcdAddr) == 0 {
		log.Fatalln("invoke func client.InitEtcdAddr(etcdAddr []string) to initialize etcd cluster")
	}
	basePath = strings.Trim(basePath, "/")
	basePath = "/" + basePath
	servicePath = strings.Trim(servicePath, "/")
	allPath := basePath + "/" + servicePath
	token := hashSelf(allPath)

	cli, exist := irpcxCli.cli[token]
	if !exist {
		return set(basePath, servicePath, token)
	}
	return cli
}
func set(basePath, servicePath, token string) *iclient {
	discovery := client.NewEtcdDiscovery(basePath, servicePath, etcdAddr, nil)
	xclient := client.NewXClient(servicePath, client.Failover, client.RandomSelect, discovery, client.DefaultOption)
	mc := new(iclient)
	mc.cli = xclient
	mc.discovery = discovery

	irpcxCli.cli[token] = mc
	rpcClient[token] = mc
	return mc
}

func del(basePath, servicePath string) {
	token := hashSelf(basePath + "/" + servicePath)
	delete(irpcxCli.cli, token)
}

// SetHashFunc . define youself hash func
func SetHashFunc(hash func(string) string) {
	if len(irpcxCli.cli) > 0 {
		return
	}
	irpcxCli.hash = hash
}

// SetXClient . define youself XClient
func SetXClient(basePath, servicePath string, xClient client.XClient) {
	token := hashSelf(basePath + "/" + servicePath)

	c := new(iclient)
	c.cli = xClient
	irpcxCli.cli[token] = c
}
func hash(str string) string {
	md := md5.New()
	md.Write([]byte(str)) // 需要加密的字符串为 str
	cipherStr := md.Sum(nil)
	return fmt.Sprintf("%s", hex.EncodeToString(cipherStr)) // 输出加密结果
}
func hashSelf(src string) string {
	if irpcxCli.hash != nil {
		return irpcxCli.hash(src)
	}
	return hash(src)
}

func (mc *iclient) call(method string, args, reply interface{}) error {
	return mc.cli.Call(context.Background(), method, args, reply)
}
func (mc *iclient) gocall(method string, args, reply interface{}) (*client.Call, error) {
	return mc.cli.Go(context.Background(), method, args, reply, nil)
}
func (mc *iclient) fork(method string, args, reply interface{}) error {
	return mc.cli.Fork(context.Background(), method, args, reply)
}
func (mc *iclient) broadcast(method string, args, reply interface{}) error {
	return mc.cli.Broadcast(context.Background(), method, args, reply)
}
func (mc *iclient) close() error {
	return mc.cli.Close()
}
func (mc *iclient) auth(src string) {
	mc.cli.Auth(src)
}

// InitEtcdAddr will initialize the etcd cluster
func InitEtcdAddr(etcdAddr []string) {
	SetMode(ModeDebug)
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

// Fork sends requests to all servers and Success once one server returns OK.
// FailMode and SelectMode are meanless for this method.
func Fork(req *irpcx.Request) (reply *irpcx.Response, err error) {
	reply = new(irpcx.Response)
	servicePath, method := getParam()
	err = get(req.BasePath, servicePath).fork(method, req, reply)
	if err != nil {
		return nil, err
	}
	return
}

// Broadcast sends requests to all servers and Success only when all servers return OK.
// FailMode and SelectMode are meanless for this method.
// Please set timeout to avoid hanging.
func Broadcast(req *irpcx.Request) (reply *irpcx.Response, err error) {

	reply = new(irpcx.Response)
	servicePath, method := getParam()
	err = get(req.BasePath, servicePath).broadcast(method, req, reply)
	if err != nil {
		return nil, err
	}
	return
}

// Close closes this client and its underlying connnections to services.
func Close(basePath string) error {
	servicePath, _ := getParam()
	if err := get(basePath, servicePath).close(); err != nil {
		return err
	}
	del(basePath, servicePath)
	return nil
}

// Auth sets s token for Authentication.
func Auth(basePath string, auth string) {
	servicePath, _ := getParam()
	get(basePath, servicePath).auth(auth)
}
