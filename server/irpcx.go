package server

import (
	"fmt"
	"context"
	"reflect"

	"github.com/youngbloood/irpcx"
)

type functioner func(*irpcx.Context) error

var functor functioner

// IRPCX struct
type IRPCX struct {
	store store
}

// Do will register as a rpc function in rpcx
func (r *IRPCX) Do(ctx context.Context, request *irpcx.Request, response *irpcx.Response)(err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[service internal error]: %v", r)
		}
	}()

	var fn *function
	fn, err = r.store.get(request.Service, request.Method)
	if err != nil {
		return err
	}

	var in []reflect.Value
	for i, v := range fn.In {
		var vt reflect.Value
		if v.Kind() != reflect.Ptr{
			vt = reflect.New(v).Elem()
		}else{
			vt = reflect.New(reflect.PtrTo(v).Elem().Elem())
		}
		
		switch i {
		case 1:
			ctxParam := new(irpcx.Context)
			ctxParam.Req = request
			ctxParam.Resp = response
			vt = reflect.ValueOf(ctxParam)
		}
		in = append(in, vt)
	}

	out := fn.Func.Call(in)

	var result interface{} = out[0].Interface()

	err, ok := result.(error)
	if ok && err != nil {
		return err
	}

	fmt.Println("req1111=",string(request.Body))

	fmt.Println("resp1111=",string(response.Body()))
	return nil
}
