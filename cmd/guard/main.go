package main

import (
	"fmt"
	"guard/bridge/constant"
	"guard/cmd"
	"time"
)

func main() {
	t := time.Now()
	fmt.Println(t.Format(constant.TimeFormatlayout))
	fmt.Println(t.Format("2006-01-02 15:04:05 -0700"))
	fmt.Println(t.Format("2006-01-02 15:04:05 -0600"))
	cmd.Main()
}
