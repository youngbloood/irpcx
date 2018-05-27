package main

import (
	"fmt"
	"irpcx"
	"time"
)

// Test struct
type Test struct{}

// Add of Test
func (Test) Add(c *irpcx.Context) error {
	time.Sleep(5 * time.Second)
	fmt.Println("doing Test.Add()")
	c.Resp.Body = []byte("the resp.body")
	return nil
}
