package main

import (
	"github.com/youngbloood/irpcx/server"
)

func main() {

	basePath := "/irpcx/youngbloood"

	serve := server.New(basePath, []string{"127.0.0.1:2379"})

	serve.Register(&Test{})

	serve.Start(":8888")

}
