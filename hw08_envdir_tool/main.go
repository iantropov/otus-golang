package main

import (
	"fmt"
	"os"
)

func main() {
	// if len(os.Args) < 3 {
	// 	fmt.Println("Please, provide more than 2 arguments")
	// 	return
	// }

	// envDir := os.Args[1]
	// cmdWithArgs := os.Args[2:]

	envDir := "./testdata/env"
	cmdWithArgs := []string{"/bin/bash", "./testdata/echo.sh", "arg1=1", "arg2=2"}

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	returnCode, err := RunCmd(cmdWithArgs, env)
	if err != nil {
		fmt.Println(err)
		return
	}
	os.Exit(returnCode)
}
