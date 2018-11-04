package main

import (
	"fmt"
	"time"

	"github.com/youngbloood/irpcx"
)

// Test struct
type Test struct{}

// Add of Test
func (Test) Add(c *irpcx.Context) error {


	c.Req.Body=[]byte(`add test body`)
	time.Sleep(2 * time.Second)
	fmt.Println("doing Test.Add()")
	return c.SetResp("the resp.body")
	//return nil
}
