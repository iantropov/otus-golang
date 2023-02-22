package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("ReadDir: %w", err)
	}

	env := make(Environment)
	for _, dirEntry := range dirEntries {
		env[dirEntry.Name()], err = evaluateEnvValue(dir, dirEntry)
		if err != nil {
			return nil, fmt.Errorf("ReadDir: %w", err)
		}
	}

	return env, nil
}

func evaluateEnvValue(dir string, dirEntry fs.DirEntry) (envValue EnvValue, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("evaluateEnvValue: %w", err)
		}
	}()

	fileInfo, err := dirEntry.Info()
	if err != nil {
		return
	}
	size := fileInfo.Size()
	if size == 0 {
		envValue.NeedRemove = true
		return
	}

	filePath := path.Join(dir, fileInfo.Name())
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	var firstLine []byte
	reader := bufio.NewReader(file)
	for {
		var prefix []byte
		var isPrefix bool
		prefix, isPrefix, err = reader.ReadLine()
		if err != nil {
			return
		}

		firstLine = append(firstLine, prefix...)
		if !isPrefix {
			break
		}
	}

	for i := 0; i < len(firstLine); i++ {
		if firstLine[i] == 0 {
			firstLine[i] = '\n'
		}
	}
	envValue.Value = strings.TrimRight(string(firstLine), " \t")
	return
}
