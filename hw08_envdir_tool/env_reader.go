package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	envs := make(Environment)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if strings.Contains(name, "=") {
			return nil, fmt.Errorf("invalid env name: %s", name)
		}

		path := filepath.Join(dir, name)
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		reader := bufio.NewReader(file)
		text, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("file reading error: %w", err)
		}
		file.Close()

		text = strings.TrimSuffix(text, "\n")
		text = strings.TrimRight(text, " \t")
		text = strings.ReplaceAll(text, "\x00", "\n")

		if len(text) == 0 && errors.Is(err, io.EOF) {
			envs[name] = EnvValue{NeedRemove: true}
		} else {
			envs[name] = EnvValue{Value: text, NeedRemove: false}
		}
	}
	return envs, nil
}
