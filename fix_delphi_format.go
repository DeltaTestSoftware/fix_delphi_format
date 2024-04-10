package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func main() {
	exitCode := 0

	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, `usage: %s one.pas two.pas three.pas [...]
provide at least one .pas file, they will all be fixed.`,
			filepath.Base(os.Args[0]))
		exitCode = 1
	}

	for _, path := range os.Args[1:] {
		if err := fix(path); err != nil {
			fmt.Fprintf(os.Stderr, "error for file '%s': %v\n", path, err)
			exitCode = 2
		}
	}

	os.Exit(exitCode)
}

func fix(path string) error {
	// Using this program with other tools, we might have a space at the end of
	// our path name. Trim it.
	path = strings.TrimRight(path, " \t\n")

	if !(strings.HasSuffix(strings.ToLower(path), ".pas") ||
		strings.HasSuffix(strings.ToLower(path), ".dpr")) {
		return errors.New("only .pas and .dpr files are supported")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	utf8bom := []byte{0xEF, 0xBB, 0xBF}

	var code string
	var isUTF8 bool

	if isASCII(data) {
		code = string(data)
	} else if bytes.HasPrefix(data, utf8bom) {
		code = string(data[len(utf8bom):])
		isUTF8 = true
	} else {
		return errors.New("only ASCII and UTF-8 (with BOM) files are supported")
	}

	if strings.Contains(code, "\n") && !strings.Contains(code, "\r\n") {
		return errors.New(`only files with \r\n as line breaks are supported`)
	}

	lines := strings.Split(code, "\r\n")
	lines = fixCommentIndentation(lines)
	lines = fixLocalVars(lines)

	var newData []byte
	if isUTF8 {
		newData = append(newData, utf8bom...)
	}
	newData = append(newData, []byte(strings.Join(lines, "\r\n"))...)

	newData = bytes.Replace(newData, []byte(" Default ("), []byte(" Default("), -1)
	newData = bytes.Replace(newData, []byte(" default ("), []byte(" Default("), -1)

	newData = bytes.Replace(newData, []byte(" low("), []byte(" Low("), -1)
	newData = bytes.Replace(newData, []byte(" high("), []byte(" High("), -1)
	newData = bytes.Replace(newData, []byte(" length("), []byte(" Length("), -1)
	newData = bytes.Replace(newData, []byte(" setlength("), []byte(" SetLength("), -1)

	return ioutil.WriteFile(path, newData, 0666)
}

func isASCII(data []byte) bool {
	for _, b := range data {
		if b >= 128 {
			return false
		}
	}
	return true
}

func isComment(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "//")
}

func indentationPrefix(s string) string {
	for i, r := range s {
		if !unicode.IsSpace(r) {
			return s[:i]
		}
	}
	return ""
}

func nextLineIndentation(line string) string {
	indent := indentationPrefix(line)
	clean := strings.ToLower(strings.TrimSpace(line))
	if clean == "end" || clean == "end;" || strings.HasPrefix(clean, "until ") {
		indent = "  " + indent
	}
	return indent
}

func findIndentation(lines []string) string {
	for _, line := range lines {
		if strings.TrimSpace(line) != "" && !isComment(line) {
			return nextLineIndentation(line)
		}
	}
	return ""
}

func lineIsIndentedVar(line string) bool {
	isVar := strings.ToLower(strings.TrimSpace(line)) == "var"
	isIndented := len(indentationPrefix(line)) > 0
	return isVar && isIndented
}

func isLocalVar(line string, nextLines []string) bool {
	if lineIsIndentedVar(line) {
		rest := ""
		for len(nextLines) > 0 && !strings.Contains(rest, ";") {
			rest += nextLines[0]
			nextLines = nextLines[1:]
		}
		return strings.Contains(rest, " := ")
	}
	return false
}

func addLocalVarPrefix(varLine, nextLine string) string {
	return varLine + " " + strings.TrimSpace(nextLine)
}

func fixCommentIndentation(lines []string) []string {
	var fixed []string
	for i, line := range lines {
		if i+1 < len(lines) && isComment(line) {
			line = findIndentation(lines[i+1:]) + strings.TrimSpace(line)
		}
		fixed = append(fixed, line)
	}
	return fixed
}

func fixLocalVars(lines []string) []string {
	var fixed []string
	for i := range lines {
		if i+1 < len(lines) && isLocalVar(lines[i], lines[i+1:]) {
			lines[i+1] = addLocalVarPrefix(lines[i], lines[i+1])
		} else {
			fixed = append(fixed, lines[i])
		}
	}
	return fixed
}
