package main

import (
	"fmt"
	"time"
)

func main() {
	s := time.Now().Weekday()
	fmt.Println(int(s))
}
