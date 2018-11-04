package server

import (
	"log"
	"net"
	"path"
	"reflect"
	"strings"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
)

// Mode of Server
type Mode string

const (
	// ModeRelease mean release
	ModeRelease Mode = "release"
	// ModeDebug mean debug
	ModeDebug Mode = "debug"
)

// Server :rpcx server
type Server struct {
	// store
	store store
	// rpcx.Server
	// you can define it by youself
	Serve *server.Server
	// the basepath in etcd
	basepath string
	// the etcd group address
	addr []string
	// mode: release/debug
	mode Mode
}

// New config the basepath and etcd addr
func New(basepath string, addr []string) *Server {
	s := server.NewServer()
	store := make(map[string]fn)

	basepath = strings.Trim(basepath, "/")
	basepath = "/" + basepath

	return &Server{store, s, basepath, addr, ModeDebug}
}

// SetMode is set the Server's mode (default: debug)
func (s *Server) SetMode(mode Mode) {
	if s != nil {
		s.mode = mode
	}
}

// GetMode return the Server Mode
func (s *Server) GetMode() Mode {
	return s.mode
}

// Register the obj's function
func (s *Server) Register(obj interface{}) {
	typ := reflect.TypeOf(obj)
	var name string
	if typ.Kind()!=reflect.Ptr{
		name = path.Base(typ.Name())
	}else{
		name = typ.Elem().Name()
	}

	
	s.RegisterName(name, obj)
}

// RegisterName the obj's function with define name
func (s *Server) RegisterName(name string, obj interface{}) {
	typ := reflect.TypeOf(obj)
	for i := 0; i < typ.NumMethod(); i++ {
		t := typ.Method(i)
		// only one parameter(the first is reciever parameter)
		if t.Type.NumIn() != 2 {
			continue
		}

		//fmtT := reflect.TypeOf(functor)

		// TODO: if the first parameter is not *irpcx.Context type ,then continue
		// for j := 0; j < t.Type.NumIn(); j++ {
		// 	if j == 1 {
		// 		t.Type.In(j).Elem().c
		// 		var inter interface{} = reflect.New(t.Type.In(j)).Interface()
		// 		switch d := inter.(type) {
		// 		default:
		// 			fmt.Println("d=", d)
		// 		}

		// 		if _, ok := inter.(*irpcx.Context); !ok {
		// 			continue
		// 		}
		// 	}
		// }

		// only allow one reply parameter
		if t.Type.NumOut() != 1 {
			continue
		}

		// TODO: the reply parameter must be error interface
		// if _, ok := reflect.New(t.Type.Out(i)).Interface().(error); !ok {
		// 	continue
		// }

		// check handler:function is exist
		if s.store.isExist(name, t.Func) {
			log.Fatalf("function : %s have exist !\n", typ.Method(i).Name)
		}

		// range the in parameters
		var in []reflect.Type
		for i := 0; i < t.Type.NumIn(); i++ {
			in = append(in, t.Type.In(i))
		}
		// store the handler:function
		s.store.set(name, typ.Method(i), in)

		if s.mode == ModeDebug {
			log.Printf("register hander : %s , function : %s\n", name, typ.Method(i).Name)
		}
	}
}

// Start the service and use default register (etcd register)
func (s *Server) Start(addr string) error {
	addr, err := initAddr(addr)
	if err != nil {
		return err
	}

	// add default etcd register
	r := &serverplugin.EtcdRegisterPlugin{
		ServiceAddress: "tcp@" + addr,
		EtcdServers:    s.addr,
		BasePath:       s.basepath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	err = r.Start()
	if err != nil {
		log.Fatal(err)
	}

	s.Serve.Plugins.Add(r)

	serviceName := s.getServiceName()

	s.Serve.RegisterName(serviceName, &IRPCX{s.store}, "")
	if s.mode == ModeDebug {
		log.Printf("Listening and serving TCP on : %s\n", addr)
	}

	if err := s.Serve.Serve("tcp", addr); err != nil {
		log.Fatalln(err)
	}
	return nil
}

// StartWithPlugin the service on addr with plugin
func (s *Server) StartWithPlugin(addr string, plugin server.Plugin) error {
	addr, err := initAddr(addr)
	if err != nil {
		return err
	}

	// add plugin
	s.Serve.Plugins.Add(plugin)

	serviceName := s.getServiceName()

	s.Serve.RegisterName(serviceName, &IRPCX{s.store}, "")
	if s.mode == ModeDebug {
		log.Printf("Listening and serving TCP on : %s\n", addr)
	}
	if err := s.Serve.Serve("tcp", addr); err != nil {
		log.Fatalln(err)
	}
	return nil

}

// init the addr (such as ":8080" to "127.0.0.1:8080")
func initAddr(addr string) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}
	if host == "" {
		host = "127.0.0.1"
	}
	return net.JoinHostPort(host, port), nil
}

func (s *Server) getServiceName() string {
	serviceName := "IRPCX"
	if s.mode == ModeDebug {
		serviceName = "IRPCX_DEBUG"
	}
	return serviceName
}

// func SetServer(ser *server.Server){
// 	ser.Plugins.
// }
