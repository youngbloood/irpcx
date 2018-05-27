package server

import (
	"context"
	"reflect"

	"github.com/youngbloood/irpcx"
)

type functioner func(*irpcx.Context) error

// IRPCX struct
type IRPCX struct {
	store store
}

// Do will register as a rpc function in rpcx
func (r *IRPCX) Do(ctx context.Context, request *irpcx.Request, response *irpcx.Response) error {

	fn, err := r.store.get(request.Service, request.Method)
	if err != nil {
		return err
	}

	var in []reflect.Value
	for i, v := range fn.In {
		vt := reflect.New(v).Elem()
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

	return nil
}
