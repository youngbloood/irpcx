package main

import (
	"fmt"

	"github.com/youngbloood/irpcx"

	"github.com/youngbloood/irpcx/client"
)

func main() {

	basePath := "/irpcx/youngbloood"
	service := "Test"
	method := "Add"

	req := irpcx.NewRequest(basePath, service, method)

	client.InitEtcdAddr([]string{"127.0.0.1:2379"})

	// sync invoke rpc
	resp, err := client.Call(req)
	if err != nil {
		panic(err)
	}
	fmt.Println("resp=", string(resp.Body))

	// async invoke rpc
	call, resp, err := client.Go(req)
	if err != nil {
		panic(err)
	}
	<-call.Done
	fmt.Println("resp=", string(resp.Body))

}
