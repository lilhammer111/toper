package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

func main() {
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("*/3 * * * * *", func() {
		fmt.Printf("per 3 seconds to do something,now is %s\n", time.Now().String())
	})
	if err != nil {
		fmt.Printf("what is the problem ? %s ", err)
	}
	c.Start()
	defer c.Stop()

	select {}

}
