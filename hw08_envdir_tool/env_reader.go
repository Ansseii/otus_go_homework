package main

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var ErrorIncorrectFileName = errors.New(`filename mustn't contain "=" symbol'`)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := Environment{}
	handler := func(fd *os.File, fileName string) {
		scanner := bufio.NewScanner(fd)
		scanner.Scan()

		buf := bytes.ReplaceAll(scanner.Bytes(), []byte{0x00}, []byte("\n"))
		env[fileName] = EnvValue{
			Value:      strings.TrimRight(string(buf), " \t"),
			NeedRemove: len(buf) == 0,
		}
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if err := validateFileName(file.Name()); err != nil {
			return nil, err
		}
		if err := processFile(dir, file.Name(), handler); err != nil {
			return nil, err
		}
	}

	return env, nil
}

func validateFileName(fileName string) error {
	if strings.Contains(fileName, "=") {
		return ErrorIncorrectFileName
	}
	return nil
}

func processFile(root string, fileName string, dataHandler func(fd *os.File, fileName string)) (err error) {
	filePath := filepath.Join(root, fileName)
	var fd *os.File
	fd, err = os.Open(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := fd.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	dataHandler(fd, fileName)

	return err
}
