package server

import (
	"fmt"
	"reflect"
)

// function's name and in parameter
type function struct {
	Func reflect.Value  // function name
	In   []reflect.Type // function in parameters
}

// method's string name: function
type fn map[string]*function

// store the service:function
type store map[string]fn

// check the service:method is exist in store
func (s store) isExist(service string, method reflect.Value) bool {
	funcs, exist := s[service]
	if !exist {
		return false
	}
	for i := range funcs {
		if funcs[i].Func == method {
			return true
		}
	}
	return false
}

// add the service:method in store
func (s store) set(service string, method reflect.Method, in []reflect.Type) {
	_, exist := s[service]
	if !exist {
		fnmap := make(fn)
		fnmap[method.Name] = &function{method.Func, in}
		s[service] = fnmap
	} else {
		s[service][method.Name] = &function{method.Func, in}
	}
}

// get function from store by service:method
func (s store) get(service, method string) (*function, error) {
	fnmap, exist := s[service]
	if !exist {
		return nil, fmt.Errorf("not found service : %s", service)
	}
	function, exist := fnmap[method]
	if !exist {
		return nil, fmt.Errorf("not found method : %s", method)
	}
	return function, nil
}
