package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmdWithArgs []string, env Environment) (returnCode int, err error) {
	cmd := exec.Command(cmdWithArgs[0], cmdWithArgs[1:]...)

	cmdEnvMap := buildEnvMap(os.Environ())
	for k, v := range env {
		if v.NeedRemove {
			delete(cmdEnvMap, k)
		} else {
			cmdEnvMap[k] = v.Value
		}
	}

	cmd.Env = compileEnvMap(cmdEnvMap)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return exiterr.ExitCode(), nil
		} else {
			return 0, fmt.Errorf("RunCmd: %w", err)
		}
	}

	return 0, nil
}

func buildEnvMap(envPairs []string) map[string]string {
	res := make(map[string]string)
	for _, pair := range envPairs {
		parts := strings.SplitN(pair, "=", 2)
		res[parts[0]] = parts[1]
	}
	return res
}

func compileEnvMap(env map[string]string) []string {
	res := make([]string, 0, len(env))
	for k, v := range env {
		res = append(res, fmt.Sprintf("%s=%s", k, v))
	}
	return res
}
