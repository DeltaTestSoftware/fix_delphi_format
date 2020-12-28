package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf(`usage: %s one.pas two.pas three.pas [...]
provide at least one .pas file, they will all be fixed.`, os.Args[0])
	}

	for _, path := range os.Args[1:] {
		if err := fix(path); err != nil {
			fmt.Printf("error for file '%s': %v\n", path, err)
		}
	}
}

func fix(path string) error {
	if !(strings.HasSuffix(strings.ToLower(path), ".pas") ||
		strings.HasSuffix(strings.ToLower(path), ".dpr")) {
		return errors.New("only .pas and .dpr files are supported")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if !isASCII(data) {
		return errors.New("only ASCII files are supported")
	}

	code := string(data)
	if strings.Contains(code, "\n") && !strings.Contains(code, "\r\n") {
		return errors.New(`only files with \r\n as line breaks are supported`)
	}

	isComment := func(line string) bool {
		return strings.HasPrefix(strings.TrimSpace(line), "//")
	}

	indentation := func(line string) string {
		for i := range line {
			if line[i] != ' ' {
				return line[:i]
			}
		}
		return ""
	}

	findIndentation := func(lines []string) string {
		for _, line := range lines {
			if strings.TrimSpace(line) != "" && !isComment(line) {
				return indentation(line)
			}
		}
		return ""
	}

	lines := strings.Split(code, "\r\n")
	var fixed []string
	for i, line := range lines {
		if isComment(line) {
			line = findIndentation(lines[i+1:]) + strings.TrimSpace(line)
		}
		fixed = append(fixed, line)
	}

	code = strings.Join(fixed, "\r\n")
	return ioutil.WriteFile(path, []byte(code), 0666)
}

func isASCII(data []byte) bool {
	for _, b := range data {
		if b >= 128 {
			return false
		}
	}
	return true
}
