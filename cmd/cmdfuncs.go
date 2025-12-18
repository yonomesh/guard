package guardcmd

import (
	"fmt"
	"guard"
	"runtime/debug"
)

type moduleInfo struct {
	guardModuleID string
	golangModule  *debug.Module
	err           error
}

func cmdStart(fl Flags) (int, error) {
	configFlag := fl.String("config")
	pidfileFlag := fl.String("pidfile")
	profileFlag := fl.String("profile")
	fmt.Println("cmdStart testing func")

	fmt.Println(configFlag)
	fmt.Println(pidfileFlag)
	fmt.Println(profileFlag)

	return guard.ExitCodeSuccess, nil
}

func cmdTest(fl Flags) (int, error) {
	helloFlag := fl.String("hello")
	fmt.Println("test cmd")
	fmt.Println(helloFlag)
	return guard.ExitCodeSuccess, nil
}
