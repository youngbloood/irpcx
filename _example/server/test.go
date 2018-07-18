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
	time.Sleep(2 * time.Second)
	fmt.Println("doing Test.Add()")
	c.SetResp("the resp.body")
	return nil
}
