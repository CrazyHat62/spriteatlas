package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
)

func main() {
	line, _ := ReadAtlas("atiles.atlas")
	println(line)
}

func ReadAtlas(name string) (string, error) {

	var v []byte

	file, err := os.Open(name)
	if err != nil && err != io.EOF {
		os.Exit(1)
	}

	var str string

	defer file.Close()
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return "", err
		}
		//remove newline and carrige return if it exists from line
		line = bytes.ReplaceAll(line, []byte("\r"), []byte(""))
		line = bytes.ReplaceAll(line, []byte("\n"), []byte(""))
		//line = bytes.ReplaceAll(line, []byte(" "), []byte(""))
		v = line
		s := string(v)
		s = strings.Trim(s, " ")
		if s == "" || s[:1] == "#" || s[:2] == "//" {
			continue
		}
		str = s
		break
	}

	return str, err

}
