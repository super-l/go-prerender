package utils

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type fileLib struct{}

var FileLib = fileLib{}

func (fileLib) CreateDir(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (fileLib) ReadFileContent(path string) (string, error) {
	var result string
	fd, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fd.Close()
	buff := bufio.NewReader(fd)
	for {
		data, _, eof := buff.ReadLine()
		if eof == io.EOF {
			break
		}
		var value = strings.TrimSpace(string(data))
		if value != "" {
			result = result + value
		}
	}
	return result, nil
}
