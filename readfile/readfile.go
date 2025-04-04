package readfile

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

type LineType string

const (
	LineTypeAnsibleGroupVars     LineType = "ansible_group_vars"
	LineTypeAnsibleGroupChildren LineType = "ansible_group_children"
	LineTypeAnsibleHost          LineType = "ansible_host"
)

type Line struct {
	LineNum  int
	Line     string
	LineType LineType
}

func ReadFile(filepath string) ([]Line, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	blankLine, err := regexp.Compile(`^\s*$`)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(f)
	lines := []Line{}
	lineNum := 0
	for {
		lineNum++
		line, err := reader.ReadString('\n')

		// 防止最后一行没有换行符
		if err == io.EOF && len(line) == 0 {
			break
		}

		line = strings.TrimSuffix(line, "\n")
		if strings.HasPrefix(line, "#") {
			continue
		}
		if blankLine.Match([]byte(line)) {
			continue
		}
		lines = append(lines, Line{
			LineNum: lineNum,
			Line:    line,
		})
	}

	return lines, nil
}
